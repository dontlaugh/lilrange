package lilrange

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
	repr     string
}

func (t Range) Within(instant time.Time) bool {
	return instant.After(t.Start) && instant.Before(t.End)
}

func (r Range) Next() Range {
	return Range{
		r.Start.Add(24 * time.Hour),
		r.End.Add(24 * time.Hour),
		r.Duration,
		r.repr,
	}
}

func Parse(s string) (*Range, error) {
	var tr Range
	splitted := strings.Split(s, "-")
	if len(splitted) != 2 {
		return nil, fmt.Errorf("invalid time range value: %s", s)
	}
	startHr, startMin, err := extractAndValidate(splitted[0])
	if err != nil {
		return nil, err
	}
	endHr, endMin, err := extractAndValidate(splitted[1])
	if err != nil {
		return nil, err
	}
	// 0. determine duration
	dur, _ := CalculateDurationMinutes(startHr, startMin, endHr, endMin)
	// 1. determine if "today's start" has happened
	now := time.Now()
	nowYear, nowMonth, nowDate := now.Date()
	todaysStart := time.Date(nowYear, nowMonth, nowDate, startHr, startMin, 0, 0, time.UTC)
	endTime := todaysStart.Add(time.Minute * time.Duration(dur))

	if now.After(endTime) {
		// today's range has ended. Our default behavior is to return the next
		// range, e.g. tomorrow's
		tr.Start = todaysStart.Add(time.Hour * 24)
		tr.End = endTime.Add(time.Hour * 24)
	} else {
		// Today's range has not ended. We are either a) inside this range, or
		// b) it is in the future and has not started. Either way, assign the
		// computed times and let the caller decide what to do.
		tr.Start = todaysStart
		tr.End = endTime
	}
	tr.Duration = time.Minute * time.Duration(dur)
	tr.repr = s
	return &tr, nil
}

func CalculateDurationMinutes(startHr, startMin, endHr, endMin int) (int, bool) {
	if !between(startHr, 0, 23) || !between(endHr, 0, 23) || !between(startMin, 0, 59) || !between(endMin, 0, 59) {
		panic("invalid values for hours and/or minutes")
	}
	var crossesMidnight bool
	if startHr > endHr {
		crossesMidnight = true
	}
	var duration int
	if crossesMidnight {
		minUntilMidnight := ((24 - startHr) * 60) - startMin
		minAfterMidnight := (endHr * 60) + endMin
		duration = minUntilMidnight + minAfterMidnight
	} else {
		duration = ((endHr - startHr) * 60) - startMin + endMin
	}
	return duration, crossesMidnight
}

func validRune(r rune) bool {
	valid := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for _, x := range valid {
		if x == r {
			return true
		}
	}
	return false
}

func extractAndValidate(s string) (int, int, error) {
	runes := []rune(s)
	if len(runes) != 4 {
		return -1, -1, fmt.Errorf("invalid time component: %s", s)
	}
	for pos, char := range s {
		// TODO make this a comparison to range of ASCII bytes rather than another loop
		if !validRune(char) {
			return -1, -1, fmt.Errorf("invalid char at position %v: %c [%v]",
				pos, char, char)

		}
	}
	hr, min := string(runes[0:2]), string(runes[2:4])
	// We've asserted that our chars are ASCII 0-9, so Atoi should never fail.
	hrInt, err := strconv.Atoi(hr)
	if err != nil {
		panic("how is this possible")
	}
	minInt, err := strconv.Atoi(min)
	if err != nil {
		panic("universe is broken, sorry")
	}
	if !between(hrInt, 0, 23) {
		return -1, -1, errors.New("hour value must be between 0 and 23")
	}
	if !between(minInt, 0, 59) {
		return -1, -1, errors.New("minute value must be between 0 and 59")
	}
	return hrInt, minInt, nil
}

// x is between this and that
func between(x, this, that int) bool {
	return x >= this && x <= that
}
