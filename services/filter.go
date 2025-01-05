package handlers

import (
	"log"
	"taas-api/models"
	"time"
)

func FilterValidSlots(availability []models.AvailableTimeSlots, cutoffTime time.Time) []models.AvailableTimeSlots {
	var filtered []models.AvailableTimeSlots

	for _, record := range availability {
		var validSlots []models.TimeRange
		date := record.AvailableDate

		for _, slot := range record.AvailableSlots {
			startTime, err := parseTimeRangeToTime(date, slot.StartTime)
			if err != nil {
				log.Printf("Error parsing slot %s for date %v: %v", slot.StartTime, date, err)
				continue
			}
			// Include only future slots
			if startTime.After(cutoffTime) {
				validSlots = append(validSlots, slot)
			}

		}

		// If there are valid slots, add to the filtered list
		if len(validSlots) > 0 {
			record.AvailableSlots = validSlots
			filtered = append(filtered, record)
		}
	}

	log.Printf("Filtered records: %v", filtered)
	return filtered
}

func parseTimeRangeToTime(date time.Time, timeStr string) (time.Time, error) {
	timePart, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	combinedTime := time.Date(date.Year(), date.Month(), date.Day(), timePart.Hour(), timePart.Minute(), 0, 0, date.Location())
	return combinedTime, nil
}
