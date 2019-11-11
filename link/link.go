package link

import (
	"fmt"
	"strings"
	"time"

	"github.com/mzjp2/link-of-the-day/storage"
)

func normaliseTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// GetURL returns the URL scheduled for the day correspondoing to t
func GetURL(svc storage.Service, t time.Time) (string, error) {
	scheduledTime := normaliseTime(t)

	record, err := svc.LoadScheduled(scheduledTime)
	if err != nil {
		return "", fmt.Errorf("Could not get URL: %v", err)
	}

	if record == nil {
		return "/nothing-scheduled", nil
	}

	svc.UpdateCount(record.ID)

	return record.URL, nil
}

// SaveURL saves the URL and schedules it for the correct time
func SaveURL(svc storage.Service, url string, t time.Time) error {
	lastScheduled, err := svc.LoadLast()
	if err != nil {
		return fmt.Errorf("could not get last scheduled time: %v", err)
	}

	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	normTime := normaliseTime(t)
	if lastScheduled == nil {
		_, err := svc.Save(url, normTime, normTime)
		if err != nil {
			return fmt.Errorf("could not save url: %v", err)
		}
		return nil
	}

	lastScheduledTime := normaliseTime(lastScheduled.Scheduled)
	newScheduledTime := lastScheduledTime.Add(time.Hour * 24)

	_, err = svc.Save(url, normaliseTime(newScheduledTime), normTime)
	if err != nil {
		return fmt.Errorf("could not save url: %v", err)
	}
	return nil
}
