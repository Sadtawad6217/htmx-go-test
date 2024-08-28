package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	defaultLimit := 10
	defaultPage := 1
	defaultPublished := true

	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	page, err := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	if err != nil || page <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page parameter",
		})
	}

	offset := (page - 1) * limit
	searchTitle := c.Query("title", "")

	publishedStr := c.Query("published")
	var published bool
	if publishedStr == "" {
		published = defaultPublished
	} else {
		published = strings.ToLower(publishedStr) == "true"
	}

	// Fetch posts for the current page
	articles, err := h.service.GetPostAll(limit, offset, searchTitle, published)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the total count of posts that match the search criteria
	totalArticles, err := h.service.GetTotalPostCount(searchTitle, published)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(totalArticles) / float64(limit)))

	response := fiber.Map{
		"posts":      articles,
		"count":      totalArticles, // Total count of posts
		"limit":      limit,
		"page":       page,
		"total_page": totalPages,
	}

	return c.JSON(response)
}
