package database

import (
	"database/sql"
	"errors"
	"synthori/ediary/m/models"
	"time"
)

type DiaryDatabase interface {
	GetDayByID(userID int, id int) (models.DiaryEntries, error)
	GetAllWeek(userID int, date time.Time) ([]models.DiaryEntries, error)
}

type SQLiteDiaryDatabase struct {
	db *sql.DB
}

func NewSQLiteDiaryDatabase(db *sql.DB) *SQLiteDiaryDatabase {
	return &SQLiteDiaryDatabase{db: db}
}

var _ DiaryDatabase = (*SQLiteDiaryDatabase)(nil)

func (repo *SQLiteDiaryDatabase) GetDayByID(userID int, id int) (models.DiaryEntries, error) {
	var diaryEntry models.DiaryEntries

	err := repo.db.QueryRow(`select id, user_id, date, subject, lesson_topic, is_done, notes, created_at, updated_at
		from diary_entries
		where id=? AND user_id=?`,
		id, userID).Scan(
		&diaryEntry.ID,
		&diaryEntry.UserID,
		&diaryEntry.Date,
		&diaryEntry.Subject,
		&diaryEntry.LessonTopic,
		&diaryEntry.IsDone,
		&diaryEntry.Notes,
		&diaryEntry.CreatedAt,
		&diaryEntry.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.DiaryEntries{}, sql.ErrNoRows
	}
	if err != nil {
		return models.DiaryEntries{}, err
	}

	// user.CreatedAt, err = time.Parse(time.UnixDate, CreatedAtStr)
	println(diaryEntry.ID)
	return diaryEntry, nil
}

func (repo *SQLiteDiaryDatabase) GetAllWeek(userID int, date time.Time) ([]models.DiaryEntries, error) {
	date = date.In(date.Location())
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start := time.Date(date.Year(), date.Month(), date.Day()-weekday+1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 7)

	rows, err := repo.db.Query(`
    SELECT id, user_id, date, subject, lesson_topic, is_done, notes, created_at, updated_at
    FROM diary_entries
    WHERE date >= ? AND date < ? AND user_id = ?
    ORDER BY date, id`, start, end, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	week := make([]models.DiaryEntries, 0)
	for rows.Next() {
		var entry models.DiaryEntries
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Date,
			&entry.Subject,
			&entry.LessonTopic,
			&entry.IsDone,
			&entry.Notes,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		week = append(week, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return week, nil
}
