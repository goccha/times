package weeks

import (
	"time"
)

func Times(t time.Time) []time.Time {
	w := t.Weekday() // 曜日
	var sunday time.Time
	if w == time.Sunday {
		sunday = t
	} else {
		sunday = t.AddDate(0, 0, -int(w))
	}
	week := make([]time.Time, 0, 7)
	week = append(week, sunday)
	for i := 1; i < 7; i++ {
		week = append(week, sunday.AddDate(0, 0, i))
	}
	return week
}

func ISOTimes(t time.Time) []time.Time {
	w := t.Weekday() // 曜日
	var monday time.Time
	if w == time.Monday {
		monday = t
	} else if w == time.Sunday {
		monday = t.AddDate(0, 0, -6)
	} else {
		monday = t.AddDate(0, 0, -int(w-1))
	}
	week := make([]time.Time, 0, 7)
	week = append(week, monday)
	for i := 1; i < 7; i++ {
		week = append(week, monday.AddDate(0, 0, i))
	}
	return week
}

func DayStrings(t time.Time) []string {
	w := t.Weekday() // 曜日
	var sunday time.Time
	if w == time.Sunday {
		sunday = t
	} else {
		sunday = t.AddDate(0, 0, -int(w))
	}
	week := make([]string, 0, 7)
	week = append(week, sunday.Format("2006-01-02"))
	for i := 1; i < 7; i++ {
		week = append(week, sunday.AddDate(0, 0, i).Format("2006-01-02"))
	}
	return week
}

func ISODayStrings(t time.Time) []string {
	w := t.Weekday() // 曜日
	var monday time.Time
	if w == time.Monday {
		monday = t
	} else if w == time.Sunday {
		monday = t.AddDate(0, 0, -6)
	} else {
		monday = t.AddDate(0, 0, -int(w-1))
	}
	week := make([]string, 0, 7)
	week = append(week, monday.Format("2006-01-02"))
	for i := 1; i < 7; i++ {
		week = append(week, monday.AddDate(0, 0, i).Format("2006-01-02"))
	}
	return week
}

func WeekOfMonth(t time.Time) int {
	y, w := t.ISOWeek()
	if t.Weekday() == time.Sunday {
		w++
	}
	fom := time.Date(y, t.Month(), 1, 0, 0, 0, 0, time.Local)
	_, fw := fom.ISOWeek()
	if fom.Weekday() == time.Sunday {
		fw++
	}
	return w - fw + 1
}

func ISOWeekOfMonth(t time.Time) int {
	y, w := t.ISOWeek()
	fom := time.Date(y, t.Month(), 1, 0, 0, 0, 0, time.Local)
	_, fw := fom.ISOWeek()
	return w - fw + 1
}

func Add(t time.Time, w int) time.Time {
	return t.AddDate(0, 0, w*7)
}

func Same(t1 time.Time, t2 time.Time) bool {
	last, next := t1, t2
	if t2.Before(t1) {
		last = next
		next = t1
	}
	w := next.Weekday() // 曜日
	var sunday time.Time
	if w == time.Sunday {
		sunday = next
	} else {
		sunday = next.AddDate(0, 0, -int(w))
	}
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, sunday.Location())
	return sunday.Equal(last) || sunday.Before(last)
}

func ISOSame(t1 time.Time, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	if y1 == y2 {
		return w1 == w2
	}
	last, next := t1, t2
	if t2.Before(t1) {
		last = next
		next = t1
	}
	w := next.Weekday() // 曜日
	var monday time.Time
	switch w {
	case time.Monday:
		monday = next
	case time.Sunday:
		w = 8
		fallthrough
	default:
		w--
		monday = next.AddDate(0, 0, -int(w))
	}
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
	return monday.Equal(last) || monday.Before(last)
}
