package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.service.GetPostByID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get post",
		})
	}
	err = h.service.DeletePost(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
