package handlers

import (
	"net/http"
	"time"

	"mywebsite.tv/name/cmd/model"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")

	var updateData model.Posts
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	existingPost, err := h.service.GetPostByID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get existing post",
		})
	}

	if updateData.Title == "" {
		updateData.Title = existingPost.Title
	}
	if updateData.Content == "" {
		updateData.Content = existingPost.Content
	}
	if updateData.Published {
		updateData.Published = true
	} else {
		updateData.Published = false
	}

	updateData.UpdatedAt = time.Now()
	updateData.CreatedAt = existingPost.CreatedAt

	updatedPost, err := h.service.UpdatePost(id, updateData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post",
		})
	}

	response := fiber.Map{
		"id":         updatedPost.ID,
		"title":      updatedPost.Title,
		"content":    updatedPost.Content,
		"published":  updatedPost.Published,
		"created_at": updatedPost.CreatedAt.Format("2006-01-02T15:04:05"),
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
