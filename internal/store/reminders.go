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
	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		return err
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var reminders []model.Reminder
	if err := s.Db.Find(&reminders).Error; err != nil {
		return err
	}

	for _, reminder := range reminders {
		if reminder.Discord && reminder.Gcal {
			if err := s.Db.Delete(&reminder).Error; err != nil {
				return err
			}
			continue
		}

		day, err := time.ParseInLocation("02-01-2006", reminder.Date, loc)
		if err != nil {
			continue
		}
		if day.Before(today) {
			if err := s.Db.Delete(&reminder).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) DeleteRemindersAtEOD() error {
	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		return err
	}
	now := time.Now().In(loc)
	// 12:00 AM is the start of a new day; the day that just ended is yesterday.
	ended := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -1)

	var reminders []model.Reminder
	if err := s.Db.Find(&reminders).Error; err != nil {
		return err
	}

	for _, reminder := range reminders {
		day, err := time.ParseInLocation("02-01-2006", reminder.Date, loc)
		if err != nil {
			continue
		}
		if !ended.Before(day) {
			if err := s.Db.Delete(&reminder).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
