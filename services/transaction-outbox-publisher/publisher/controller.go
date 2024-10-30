package publisher

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	TriggerPublisher(ctx *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(s Service) Controller {
	return &controller{service: s}
}

// TriggerPublisher godoc
// @Summary Triggers the event publisher
// @Description Manually triggers the publisher to retrieve events and publish them
// @Tags Publisher
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int "eventCount"
// @Failure 500 {object} map[string]string "error"
// @Router /trigger-publisher [post]
func (c *controller) TriggerPublisher(ctx *fiber.Ctx) error {
	eventCount, err := c.service.TriggerPublisher()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": fmt.Errorf("failed to trigger publisher: %w", err).Error()})
	}

	return ctx.JSON(fiber.Map{"eventCount": eventCount})
}
