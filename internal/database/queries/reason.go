package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) GetReportReasonById(id uint) (*models.ReportReason, error) {
	db := s.s.DB()
	var reason models.ReportReason
	if err := db.First(&reason, id).Error; err != nil {
		return nil, err
	}
	return &reason, nil
}

func (s *DatabaseService) GetReasons() ([]*models.ReportReason, error) {
	db := s.s.DB()
	var reasons []*models.ReportReason
	if err := db.Find(&reasons).Error; err != nil {
		return nil, err
	}
	return reasons, nil
}
