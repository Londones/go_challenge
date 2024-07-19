package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateReportedMessage(report *models.ReportedMessage) (*models.ReportedMessage, error) {
	db := s.s.DB()

	if err := db.Create(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

func (s *DatabaseService) GetUnHandledReportedMessages() ([]*models.ReportedMessage, error) {
	db := s.s.DB()
	var reports []*models.ReportedMessage
	if err := db.Where("is_handled = ?", false).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *DatabaseService) GetHandledReportedMessages() ([]*models.ReportedMessage, error) {
	db := s.s.DB()
	var reports []*models.ReportedMessage
	if err := db.Where("is_handled = ?", true).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *DatabaseService) HandleReportedMessage(id uint) error {
	db := s.s.DB()

	var report models.ReportedMessage
	if err := db.Where("id = ?", id).First(&report).Error; err != nil {
		return err
	}

	report.IsHandled = true

	if err := db.Save(&report).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetReportedMessages() ([]*models.ReportedMessage, error) {
	db := s.s.DB()
	var reports []*models.ReportedMessage
	if err := db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}
