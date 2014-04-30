package subchal

import (
    "time"
    "strings"
    "strconv"
)

// Converts a schedule time (e.g. "11:36:14") to a time.Time object.
//
// The time zone will be UTC. The date will be the reference date
// (January 2, 2006). This function handles times past midnight (e.g. 25:04:00)
// by dumbly subtracting 24 from the hour.
func ParseTime(timeStr string) (t time.Time, err error) {
    timeParts := strings.Split(timeStr, ":")
    hour, err := strconv.ParseInt(timeParts[0], 10, 8)
    if err != nil {
        return time.Time{}, err
    }
    
    if hour >= 24 {
        timeStr = strings.Join([]string{
            strconv.Itoa(int(hour - 24)),
            timeParts[1],
            timeParts[2],
        }, ":")
    }

    t, err = time.Parse("2006-01-02 15:04:05 UTC",
                        strings.Join([]string{"2006-01-02", timeStr, "UTC"}, " "))
    return
}
