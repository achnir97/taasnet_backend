package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validate if the slot is in the correct format (HH:mm-HH:mm)
func IsValidSlotFormat(slot string) bool {
	parts := strings.Split(slot, "-")
	if len(parts) != 2 {
		return false
	}

	timeRegex := `^([01]?[0-9]|2[0-3]):[0-5][0-9]$`
	return regexp.MustCompile(timeRegex).MatchString(parts[0]) && regexp.MustCompile(timeRegex).MatchString(parts[1])
}

func ParseSlotToTime(date time.Time, slot string) (time.Time, error) {
	parts := strings.Split(slot, "-")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid slot format")
	}

	startParts := strings.Split(parts[0], ":")
	if len(startParts) != 2 {
		return time.Time{}, fmt.Errorf("invalid start time format")
	}

	hour, err := strconv.Atoi(startParts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %v", err)
	}

	minute, err := strconv.Atoi(startParts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %v", err)
	}

	// Combine the parsed hour and minute with the provided date
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location()), nil
}

// Preprocess slots to ensure they are valid and formatted correctly
func PreprocessSlots(slots []string) []string {
	var validSlots []string
	for _, slot := range slots {
		if IsValidSlotFormat(slot) {
			validSlots = append(validSlots, slot)
		}
	}
	return validSlots
}
