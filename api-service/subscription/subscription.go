package subscription

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required,gt=0"`
	UserId      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func ValidateDateFormat(dateStr string) error {
	if dateStr == "" {
		return nil
	}

	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid date format: expected 'MM-YYYY'")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid month: %v", err)
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month: must be between 1 and 12")
	}

	if year < 1900 || year > 9999 {
		return fmt.Errorf("invalid year: must be between 1900 and 9999")
	}

	_, err = time.Parse("01-2006", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	return nil
}
