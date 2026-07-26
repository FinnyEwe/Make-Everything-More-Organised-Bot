package store

import (
	"backend/internal/model"
	"time"
)

// commands:
// Remind description, time as a date
// remind list

func (s *Store) CreateReminder(description string, date string) {
	s.Db.Create(&model.Reminder{Description: description, Date: date})
}

func (s *Store) PollReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder
	currDate := time.Now().Format("02-01-2006")
	tmr := time.Now().AddDate(1, 0, 0).Format("02-01-2006")
	err := s.Db.Where("date = ?", currDate).Or("date = ?", tmr).Find(&reminders).Error
	return reminders, err
}

func (s *Store) GetReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder
    err := s.Db.Find(&reminders).Error
    return reminders, err
}



