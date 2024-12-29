package services

import (
	"log"
	"strings"
	"time"

	"taas-api/models"
	"taas-api/utils"
)

func FilterValidSlots(availability []models.AvailableTimeSlots, cutoffTime time.Time) []models.AvailableTimeSlots {
	var filtered []models.AvailableTimeSlots

	for _, record := range availability {
		var validSlots []string
		date := record.AvailableDate

		// Split slots into individual times
		slots := strings.Split(record.AvailableSlots, ",")

		for _, slot := range slots {
			startTime, err := utils.ParseSlotToTime(date, slot)
			if err != nil {
				log.Printf("Error parsing slot %s for date %v: %v", slot, date, err)
				continue
			}

			// Include only future slots
			if startTime.After(cutoffTime) {
				validSlots = append(validSlots, slot)
			}
		}

		// If there are valid slots, add to the filtered list
		if len(validSlots) > 0 {
			record.AvailableSlots = strings.Join(validSlots, ",")
			filtered = append(filtered, record)
		}
	}

	log.Printf("Filtered records: %v", filtered)
	return filtered
}
