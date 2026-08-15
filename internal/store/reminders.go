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
	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		return nil, err
	}
	now := time.Now().In(loc)
	currDate := now.Format("02-01-2006")
	tmr := now.AddDate(0, 0, 1).Format("02-01-2006")
	err = s.Db.Where("date = ?", currDate).Or("date = ?", tmr).Find(&reminders).Error
	return reminders, err
}

func (s *Store) GetReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder
	err := s.Db.Find(&reminders).Error
	return reminders, err
}

func (s *Store) UpdateNotification(id uint, channel string) error {
	err := s.Db.Model(&model.Savings{}).
		Where("id = ?", id).
		Update(channel, true).Error

	return err
}

func (s *Store) DeleteCompletedReminders() error {
	// Delete reminders where both discord and gcal are true (fully processed)
	err := s.Db.Where("discord = ? AND gcal = ?", true, true).
		Delete(&model.Reminder{}).Error
	return err
}
