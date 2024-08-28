package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mywebsite.tv/name/cmd/database"
	"mywebsite.tv/name/cmd/handlers"
	"mywebsite.tv/name/cmd/model"
	"mywebsite.tv/name/cmd/repository"
	"mywebsite.tv/name/cmd/service"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Data struct {
	Posts     []model.Posts `json:"posts"`
	Count     int           `json:"count"`
	Limit     int           `json:"limit"`
	Page      int           `json:"page"`
	TotalPage int           `json:"total_page"`
	Published string        `json:"published"`
}

func main() {
	fiberApp := fiber.New()
	db := database.Connect()
	repo := repository.NewRepo(db)
	srv := service.NewService(repo)
	handler := handlers.NewHandler(srv)

	api := fiberApp.Group("/api/v1")
	api.Get("/posts", handler.GetPosts)
	api.Get("/posts/:id", handler.GetPostID)
	api.Post("/posts", handler.CreatePosts)
	api.Put("/posts/:id", handler.UpdatePost)
	api.Delete("/posts/:id", handler.DeletePost)

	go func() {
		if err := fiberApp.Listen(":8081"); err != nil {
			fmt.Println("Error starting Fiber server:", err.Error())
		}
	}()

	echoApp := echo.New()
	echoApp.Use(middleware.Logger())
	echoApp.Renderer = NewTemplates()
	echoApp.Static("/css", "css")

	echoApp.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"ActiveTab": "posts",
		})
	})

	echoApp.GET("/posts", func(c echo.Context) error {
		published := c.QueryParam("published")
		limit := c.QueryParam("limit")
		page := c.QueryParam("page")

		queryParams := fmt.Sprintf("published=%s&limit=%s&page=%s", published, limit, page)
		if queryParams != "" {
			queryParams = "?" + queryParams
		}

		resp, err := http.Get("http://localhost:8081/api/v1/posts" + queryParams)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch posts from API")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.Logger().Errorf("Failed to fetch posts, status code: %d", resp.StatusCode)
			return c.String(http.StatusInternalServerError, "Failed to fetch posts from API")
		}

		var data Data
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			c.Logger().Error("Failed to parse posts data:", err)
			return c.String(http.StatusInternalServerError, "Failed to parse posts data")
		}

		paginationData := map[string]interface{}{
			"Posts":           data.Posts,
			"PrevPage":        getPage(data.Page-1, data.TotalPage),
			"NextPage":        getPage(data.Page+1, data.TotalPage),
			"PublishedFilter": published,
		}
		fmt.Println("sssss :", paginationData)
		return c.Render(http.StatusOK, "post-content", paginationData)
	})

	echoApp.GET("/posts/:id", func(c echo.Context) error {
		id := c.Param("id")

		apiURL := fmt.Sprintf("http://localhost:8081/api/v1/posts/%s", id)

		resp, err := http.Get(apiURL)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch post")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c.String(http.StatusInternalServerError, "Failed to fetch post data")
		}

		var post model.Posts
		if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
			c.Logger().Error("Failed to decode post data:", err)
			return c.String(http.StatusInternalServerError, "Failed to decode post data")
		}

		return c.Render(http.StatusOK, "editPostPopup", post)
	})

	echoApp.PUT("/posts/:id", func(c echo.Context) error {
		id := c.Param("id")
		title := c.FormValue("title")
		content := c.FormValue("content")
		published := c.FormValue("published") == "true"
		createdAt := c.FormValue("createdAt")
		if title == "" || content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and content are required"})
		}

		updateData := model.Posts{
			ID:        id,
			Title:     title,
			Content:   content,
			Published: published,
		}

		fmt.Println(updateData)
		c.Logger().Infof("Received update data: %v", updateData)

		updateDataJSON, err := json.Marshal(updateData)
		if err != nil {
			c.Logger().Error("Failed to encode JSON:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encode JSON"})
		}

		apiURL := fmt.Sprintf("http://localhost:8081/api/v1/posts/%s", id)
		req, err := http.NewRequest(http.MethodPut, apiURL, bytes.NewReader(updateDataJSON))
		if err != nil {
			c.Logger().Error("Failed to create PUT request:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create PUT request"})
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.Logger().Error("Failed to send PUT request:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send PUT request"})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.Logger().Errorf("Failed to update post, status code: %d", resp.StatusCode)
			return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to update post"})
		}

		updatedPost := map[string]interface{}{
			"ID":        id,
			"Title":     updateData.Title,
			"Content":   updateData.Content,
			"Published": updateData.Published,
			"CreatedAt": createdAt,
		}
		if err := json.NewDecoder(resp.Body).Decode(&updatedPost); err != nil {
			c.Logger().Error("Failed to decode post data:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode post data"})
		}

		c.Logger().Infof("Successfully updated post: %v", updatedPost)
		return c.Render(http.StatusOK, "editPosts", updatedPost)
	})

	echoApp.PUT("/posts/publish/:id", func(c echo.Context) error {
		id := c.Param("id")

		published := c.FormValue("published") == "true"

		updateData := model.Posts{
			Published: published,
		}

		c.Logger().Infof("Received update data: %v", updateData)

		updateDataJSON, err := json.Marshal(updateData)
		if err != nil {
			c.Logger().Error("Failed to encode JSON:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encode JSON"})
		}

		apiURL := fmt.Sprintf("http://localhost:8081/api/v1/posts/%s", id)
		req, err := http.NewRequest(http.MethodPut, apiURL, bytes.NewReader(updateDataJSON))
		if err != nil {
			c.Logger().Error("Failed to create PUT request:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create PUT request"})
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.Logger().Error("Failed to send PUT request:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send PUT request"})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.Logger().Errorf("Failed to update post, status code: %d", resp.StatusCode)
			return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to update post"})
		}

		return c.NoContent(http.StatusOK)
	})

	echoApp.DELETE("/posts/:id", func(c echo.Context) error {
		id := c.Param("id")
		apiURL := fmt.Sprintf("http://localhost:8081/api/v1/posts/%s", id)

		req, err := http.NewRequest("DELETE", apiURL, nil)
		if err != nil {
			c.Logger().Error("Failed to create DELETE request:", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.Logger().Error("Failed to send DELETE request:", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return c.NoContent(resp.StatusCode)
		}

		return c.NoContent(http.StatusNoContent)
	})

	echoApp.POST("/posts", func(c echo.Context) error {
		title := c.FormValue("title")
		content := c.FormValue("content")
		published := c.FormValue("published") == "true"

		if title == "" || content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and content are required"})
		}

		newPost := model.Posts{
			Title:     title,
			Content:   content,
			Published: published,
		}

		newPostJSON, err := json.Marshal(newPost)
		if err != nil {
			c.Logger().Error("Failed to encode JSON:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encode JSON"})
		}

		apiURL := "http://localhost:8081/api/v1/posts"
		resp, err := http.Post(apiURL, "application/json", bytes.NewReader(newPostJSON))
		if err != nil {
			c.Logger().Error("Failed to send POST request:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send POST request"})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			c.Logger().Error("Failed to create post, status code:", resp.StatusCode)
			return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to create post"})
		}

		return c.NoContent(http.StatusCreated)
	})

	echoApp.Logger.Fatal(echoApp.Start(":8080"))

	if err := echoApp.Start(":8080"); err != nil {
		echoApp.Logger.Fatal(err)
	}
}

func getPage(page, totalPage int) string {
	if page > 0 && page <= totalPage {
		return fmt.Sprintf("%d", page)
	}
	return ""
}
