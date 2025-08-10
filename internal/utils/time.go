package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeParser struct {
	timezone *time.Location
}

func NewTimeParser(timezone string) (*TimeParser, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}
	return &TimeParser{timezone: loc}, nil
}

func (tp *TimeParser) ParseRelativeTime(input string) (time.Time, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	now := time.Now().In(tp.timezone)

	if strings.Contains(input, "after") || strings.Contains(input, "بعد") {
		return tp.parseAfterTime(input, now)
	}

	return tp.parseAbsoluteTime(input, now)
}

func (tp *TimeParser) parseAfterTime(input string, now time.Time) (time.Time, error) {
	input = strings.ReplaceAll(input, "after", "")
	input = strings.ReplaceAll(input, "بعد", "")
	input = strings.TrimSpace(input)

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid relative time format")
	}

	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number: %s", parts[0])
	}

	unit := strings.ToLower(parts[1])
	var duration time.Duration

	switch {
	case strings.Contains(unit, "minute") || strings.Contains(unit, "دقيقة"):
		duration = time.Duration(num) * time.Minute
	case strings.Contains(unit, "hour") || strings.Contains(unit, "ساعة"):
		duration = time.Duration(num) * time.Hour
	case strings.Contains(unit, "day") || strings.Contains(unit, "يوم"):
		duration = time.Duration(num) * 24 * time.Hour
	case strings.Contains(unit, "week") || strings.Contains(unit, "أسبوع"):
		duration = time.Duration(num) * 7 * 24 * time.Hour
	case strings.Contains(unit, "month") || strings.Contains(unit, "شهر"):
		duration = time.Duration(num) * 30 * 24 * time.Hour
	case strings.Contains(unit, "year") || strings.Contains(unit, "سنة"):
		duration = time.Duration(num) * 365 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return now.Add(duration), nil
}

func (tp *TimeParser) parseAbsoluteTime(input string, now time.Time) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04",
		"15:04 02/01/2006",
		"15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, input, tp.timezone); err == nil {
			// If only time is provided, assume it's for today or tomorrow
			if format == "15:04" {
				today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, tp.timezone)
				if today.Before(now) {
					// If the time has passed today, schedule for tomorrow
					today = today.Add(24 * time.Hour)
				}
				return today, nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", input)
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

func GetNextRecurrenceTime(baseTime time.Time, recurrenceType string, count int) time.Time {
	switch recurrenceType {
	case "daily":
		return baseTime.AddDate(0, 0, count)
	case "weekly":
		return baseTime.AddDate(0, 0, count*7)
	case "monthly":
		return baseTime.AddDate(0, count, 0)
	case "yearly":
		return baseTime.AddDate(count, 0, 0)
	default:
		return baseTime
	}
}
