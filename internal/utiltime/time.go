package utiltime

import (
	"time"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
)

const (
	dateLayout                     = "2006-01-02"
	safaricomTimeLayout            = "2006-01-02 15:04:05"
	timeLayout                     = "2006/01/02 15:04:05"
	timeandtimezonelayout          = "2006/01/02 15:04:05-07:00"
)

func ParseTime(
	timeString string,
) (time.Time, error) {

	parsedTime, err := time.ParseInLocation(timeLayout, timeString, time.UTC)
	if err != nil {

		parsedTime, err = time.ParseInLocation(timeandtimezonelayout, timeString, time.UTC)
		if err != nil {
			return parsedTime, apperr.NewErrorWithType(
				err,
				apperr.Internal,
			)
		}
	}

	return parsedTime, nil
}

func FormatDate(timeToFormat time.Time) string {
	return timeToFormat.In(time.UTC).Format(dateLayout)
}

func FormatTime(
	timeToFormat time.Time,
) string {
	return timeToFormat.In(time.UTC).Format(timeLayout)
}

func FormatDateTime(
	timeToFormat time.Time,
) string {
	return timeToFormat.In(time.UTC).Format(safaricomTimeLayout)
}
