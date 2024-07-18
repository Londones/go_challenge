package fixtures

import (
	"fmt"

	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
)

func NewReasonFixture(reason string) *models.ReportReason {
	return &models.ReportReason{
		Reason: reason,
	}
}

func CreateReasons(db *gorm.DB) (*models.ReportReason, error) {
	reasons := []string{"spam", "inappropriateContent", "illegalContent", "harassment", "other"}
	for _, r := range reasons {
		reason := NewReasonFixture(r)
		if err := db.Create(reason).Error; err != nil {
			return nil, fmt.Errorf("failed to create reason: %w", err)
		}
	}
	return nil, nil
}
