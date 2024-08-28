package handlers

import (
	"github.com/gofiber/fiber/v2"
	"mywebsite.tv/name/cmd/model"
)

func (h *Handler) CreatePosts(c *fiber.Ctx) error {
	var post model.Posts
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if post.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	// Assuming CreatePosts in service expects *model.Posts
	createdPost, err := h.service.CreatePosts(post.Title, post.Content, post.Published)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := fiber.Map{
		"id":         createdPost.ID,
		"title":      createdPost.Title,
		"content":    createdPost.Content,
		"published":  createdPost.Published,
		"created_at": createdPost.CreatedAt.Format("2006-01-02T15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}
