package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// формируем правила повторения d, y, m, w
func NextDate(now time.Time, date string, repeat string) (string, error) {

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var nextDate time.Time

	parsedDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}

	if repeat == "" {
		return "", fmt.Errorf("неверный формат даты")
	}

	dateRepeat := strings.Split(repeat, " ")

	switch dateRepeat[0] {

	case "d":

		if len(dateRepeat) != 2 {
			return "", fmt.Errorf("не указано повторение")
		}

		days, err := strconv.Atoi(dateRepeat[1])
		if err != nil || days < 0 || days > 400 {
			return "", fmt.Errorf("указано неверное число (необходимо 1-400)")
		}

		if parsedDate.Equal(now) {
			nextDate = now
		} else {
			nextDate = parsedDate.AddDate(0, 0, days)
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(0, 0, days)
			}
		}
	case "y":
		nextDate = parsedDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	case "m":
		return "", fmt.Errorf("данная функци пока не доступна")
	case "w":
		return "", fmt.Errorf("данная функци пока не доступна")
	default:
		return "", fmt.Errorf("неверный формат повторения")
	}
	return nextDate.Format(DateFormat), nil
}
