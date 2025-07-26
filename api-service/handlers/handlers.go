package handlers

import (
	"emtest/api-service/subscription"
	"fmt"
	"net/http"
	"strings"

	"emtest/api-service/db"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// @Description Error response object
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// @Description Success response object
type SuccessResponse struct {
	Message string `json:"message" example:"Subscription deleted successfully"`
}

// @description Success calculate object
type SuccessCostResponse struct {
	Total int `json:"total"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int       `json:"price,omitempty"`
	UserId      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	EndDate     *string    `json:"end_date,omitempty"`
}

// @Summary Create a new subscription
// @Description Создание новой подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param body body subscription.Subscription true "Body of the request"
// @Success 200 {object} subscription.Subscription
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/subscriptions [post]
func CreateSubscription(c *fiber.Ctx) error {

	var sub subscription.Subscription

	if err := c.BodyParser(&sub); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := validate.Struct(sub); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}

		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Message: strings.Join(errors, "; "),
		})
	}

	if err := subscription.ValidateDateFormat(sub.StartDate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid date format",
			Message: err.Error(),
		})
	}

	result := db.DB.Create(&sub)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to create subscription",
			Message: result.Error.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(sub)
}

// @Summary Get subscriptions
// @Description Получение подписки по её id. Если id не указано, возвращаются все
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string false "Subscription ID (UUID format)" Format(uuid)
// @Success 200 {array} subscription.Subscription "List of subscriptions"
// @Success 200 {object} subscription.Subscription "Single subscription when ID provided"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions [get]
func GetSubscriptions(c *fiber.Ctx) error {

	id := c.Query("id")

	if id == "" {
		var subs []subscription.Subscription

		result := db.DB.Find(&subs)
		if result.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to fetch subscriptions",
				Message: result.Error.Error(),
			})
		}

		return c.JSON(subs)
	}
	var sub subscription.Subscription

	result := db.DB.First(&sub, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Error: fmt.Sprintf("Sunscription %s not found", id),
		})
	}

	return c.JSON(sub)
}

// @Summary Update subscription
// @Description Обновление подписки по её id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string true "Subscription ID (UUID format)" Format(uuid)
// @Param subscription body UpdateSubscriptionRequest true "Subscription object with updated fields"
// @Success 200 {object} subscription.Subscription "Updated subscription"
// @Failure 400 {object} ErrorResponse "Bad request - missing ID or invalid body"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Router /api/v1/subscriptions [put]
func UpdateSubscription(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad request",
			Message: "Parameter 'id' is required",
		})
	}

	var updatedReq UpdateSubscriptionRequest

	if err := c.BodyParser(&updatedReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}
	updates := make(map[string]interface{})

	if updatedReq.ServiceName != nil {
		updates["service_name"] = *updatedReq.ServiceName
	}
	if updatedReq.Price != nil {
		updates["price"] = *updatedReq.Price
	}
	if updatedReq.UserId != nil {
		updates["user_id"] = *updatedReq.UserId
	}
	if updatedReq.StartDate != nil {
		if err := subscription.ValidateDateFormat(*updatedReq.StartDate); err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid date format",
				Message: err.Error(),
			})
		}
		updates["start_date"] = *updatedReq.StartDate
	}
	if updatedReq.EndDate != nil {
		if err := subscription.ValidateDateFormat(*updatedReq.EndDate); err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid date format",
				Message: err.Error(),
			})
		}
		updates["end_date"] = *updatedReq.EndDate
	}

	if len(updates) == 0 {
		var currentSub subscription.Subscription
		result := db.DB.First(&currentSub, "id = ?", id)
		if result.RowsAffected == 0 {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse{
				Error:   "Subscription not found",
				Message: fmt.Sprintf("Subscription %s not found", id),
			})
		}
		return c.Status(http.StatusOK).JSON(currentSub)
	}

	result := db.DB.Model(&subscription.Subscription{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Error:   "Subscription not found",
			Message: fmt.Sprintf("Subscription %s not found", id),
		})
	}

	var updatedSubscription subscription.Subscription
	db.DB.First(&updatedSubscription, "id = ?", id)
	return c.Status(fiber.StatusOK).JSON(updatedSubscription)
}

// @Summary Delete subscription
// @Description Удаление подписки по её id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string true "Subscription ID (UUID format)" Format(uuid) Example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} SuccessResponse "Success message"
// @Failure 400 {object} ErrorResponse "Bad request - missing ID"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions [delete]
func DeleteSubscription(c *fiber.Ctx) error {

	id := c.Query("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad request",
			Message: "Parameter 'id' is required",
		})
	}

	result := db.DB.Delete(&subscription.Subscription{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Error:   "Subscription not found",
			Message: fmt.Sprintf("Subscription %s not found", id),
		})
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal server error",
			Message: result.Error.Error(),
		})
	}

	return c.JSON(SuccessResponse{Message: "Subscription deleted successfully"})
}

// @Summary Calculate total cost of subscriptions
// @Description Подсчет стоимости подписки по фильтрам
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user ID (UUID format)" Format(uuid) Example(550e8400-e29b-41d4-a716-446655440000)
// @Param service_name query string false "Filter by service name"
// @Param start_date query string false "Filter by start date (MM-YYYY format)" Format(MM-YYYY)
// @Param end_date query string false "Filter by end date (MM-YYYY format)" Format(MM-YYYY)
// @Success 200 {object} SuccessCostResponse "Total cost calculation result"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions/calculate [get]
func CalculateTotalCost(c *fiber.Ctx) error {

	var totalCost int64
	filters := map[string]interface{}{
		"user_id":      c.Query("user_id"),
		"service_name": c.Query("service_name"),
		"start_date":   c.Query("start_date"),
		"end_date":     c.Query("end_date"),
	}

	query := db.DB.Model(&subscription.Subscription{})
	query = applyFilters(query, filters)

	row := query.Select("COALESCE(SUM(price), 0)").Row()
	if err := row.Scan(&totalCost); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to count",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessCostResponse{Total: int(totalCost)})
}

func applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {

	conditions := map[string]string{
		"user_id":      "user_id = ?",
		"service_name": "service_name = ?",
		"start_date":   "start_date >= ?",
		"end_date":     "start_date <= ?",
	}
	for field, value := range filters {
		if value != "" {
			if condition, exists := conditions[field]; exists {
				query = query.Where(condition, value)
			}
		}
	}

	return query
}
