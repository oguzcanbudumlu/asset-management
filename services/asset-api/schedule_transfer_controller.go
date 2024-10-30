package main

//
//import (
//	"asset-management/internal/schedule"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"strconv"
//	"time"
//)
//
//type ScheduleTransactionController struct {
//	service schedule.ScheduleTransactionService
//}
//
//func NewScheduleTransactionController(service schedule.ScheduleTransactionService) *ScheduleTransactionController {
//	return &ScheduleTransactionController{service: service}
//}
//
//// CreateScheduleTransaction godoc
//// @Summary      Create a new scheduled transaction
//// @Description  Schedules a new transaction to be executed at a specified future time
//// @Tags         ScheduleTransaction
//// @Accept       json
//// @Produce      json
//// @Param        transaction body Request true "Schedule Transfer request payload"
//// @Success      201  {object}  map[string]int "Created transaction ID"  example: {"transaction_id": 123}
//// @Failure      400  {object}  map[string]string "Invalid request payload or scheduled time format" example: {"error": "Invalid scheduled time format"}
//// @Failure      500  {object}  map[string]string "Failed to create scheduled transaction" example: {"error": "Failed to create scheduled transaction"}
//// @Router       /scheduled-transaction [post]
//func (c *ScheduleTransactionController) CreateScheduleTransaction(ctx *fiber.Ctx) error {
//	var req Request
//	if err := ctx.BodyParser(&req); err != nil {
//		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
//	}
//
//	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
//	if err != nil {
//		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid scheduled time format"})
//	}
//
//	id, err := c.service.Create(req.From, req.To, req.Network, req.Amount, scheduledTime)
//	if err != nil {
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"transaction_id": id})
//}
//
//// GetNextMinuteTransactions godoc
//// @Summary Get transactions scheduled for the next minute
//// @Description Retrieve all transactions scheduled for the upcoming minute
//// @Tags ScheduleTransaction
//// @Accept  json
//// @Produce  json
//// @Success 200 {array} ScheduleTransaction
//// @Failure 500 {object} map[string]interface{}
//// @Router /scheduled-transaction/next-minute [get]
//func (c *ScheduleTransactionController) GetNextMinuteTransactions(ctx *fiber.Ctx) error {
//	transactions, err := c.service.GetNextMinuteTransactions()
//	if err != nil {
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transactions"})
//	}
//	return ctx.JSON(transactions)
//}
//
//// Process godoc
//// @Summary Process a scheduled transaction
//// @Description Processes a scheduled transaction by its ID
//// @Tags Transactions
//// @Param id path int true "Transaction ID"
//// @Success 200 {object} map[string]string "message": "Transaction processed successfully"
//// @Failure 400 {object} map[string]string "error": "Invalid transaction ID"
//// @Failure 500 {object} map[string]string "error": "Failed to process transaction"
//// @Router /scheduled-transaction/{id}/process [post]
//func (c *ScheduleTransactionController) Process(ctx *fiber.Ctx) error {
//	transactionIDParam := ctx.Params("id")
//	transactionID, err := strconv.ParseInt(transactionIDParam, 10, 64)
//	if err != nil {
//		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Invalid transaction ID",
//		})
//	}
//
//	if err := c.service.Process(transactionID); err != nil {
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": fmt.Errorf("failed to process transaction: %w", err).Error(),
//		})
//	}
//
//	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
//		"message": "Transaction processed successfully",
//	})
//}
