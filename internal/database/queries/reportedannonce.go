package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateReportedAnnonce(report *models.ReportedAnnonce) (*models.ReportedAnnonce, error) {
	db := s.s.DB()

	if err := db.Create(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

func (s *DatabaseService) GetReportedAnnonceById(id uint) (*models.ReportedAnnonce, error) {
	db := s.s.DB()
	var report models.ReportedAnnonce
	if err := db.First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *DatabaseService) GetUnHandledReportedAnnonces() ([]*models.ReportedAnnonce, error) {
	db := s.s.DB()
	var reports []*models.ReportedAnnonce
	if err := db.Where("is_handled = ?", false).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *DatabaseService) GetHandledReportedAnnonces() ([]*models.ReportedAnnonce, error) {
	db := s.s.DB()
	var reports []*models.ReportedAnnonce
	if err := db.Where("is_handled = ?", true).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *DatabaseService) HandleReportedAnnonce(id uint) error {
	db := s.s.DB()

	var report models.ReportedAnnonce
	if err := db.Where("id = ?", id).First(&report).Error; err != nil {
		return err
	}

	report.IsHandled = true

	if err := db.Save(&report).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetReportedAnnonces() ([]*models.ReportedAnnonce, error) {
	db := s.s.DB()
	var reports []*models.ReportedAnnonce
	if err := db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}
