package services

import (
	"synthori/ediary/m/database"
	"synthori/ediary/m/models"
	"time"
)

type Diary interface {
	GetDayByID(userID int, dayID int) (models.DiaryEntries, error)
	GetAllWeek(userID int, date time.Time) ([]models.DiaryEntries, error)
}

var _ Diary = (*DiaryService)(nil)

type DiaryService struct {
	repo database.SQLiteDiaryDatabase
}

func NewDiaryService(repo database.SQLiteDiaryDatabase) *DiaryService {
	return &DiaryService{repo: repo}
}

func (d *DiaryService) GetDayByID(userID int, dayID int) (models.DiaryEntries, error) {
	day, err := d.repo.GetDayByID(userID, dayID)
	if err != nil {
		return models.DiaryEntries{}, err
	}

	return day, nil
}

func (d *DiaryService) GetAllWeek(userID int, date time.Time) ([]models.DiaryEntries, error) {
	week, err := d.repo.GetAllWeek(userID, date)
	if err != nil {
		return []models.DiaryEntries{}, err
	}

	return week, nil
}
