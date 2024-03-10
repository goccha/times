package weeks

import (
	"time"
)

// Times 引数の日付が含まれる週の日曜日から土曜日までの日付を返す
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

// ISOTimes 引数の日付が含まれる週の月曜日から日曜日までの日付を返す
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

// DayStrings 引数の日付が含まれる週の日曜日から土曜日までの日付を文字列で返す
func DayStrings(t time.Time, layout ...string) []string {
	var format string
	if len(layout) == 0 {
		format = time.DateOnly
	} else {
		format = layout[0]
	}
	w := t.Weekday() // 曜日
	var sunday time.Time
	if w == time.Sunday {
		sunday = t
	} else {
		sunday = t.AddDate(0, 0, -int(w))
	}
	week := make([]string, 0, 7)
	week = append(week, sunday.Format(format))
	for i := 1; i < 7; i++ {
		week = append(week, sunday.AddDate(0, 0, i).Format(format))
	}
	return week
}

// ISODayStrings 引数の日付が含まれる週の月曜日から日曜日までの日付を文字列で返す
func ISODayStrings(t time.Time, layout ...string) []string {
	var format string
	if len(layout) == 0 {
		format = time.DateOnly
	} else {
		format = layout[0]
	}
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
	week = append(week, monday.Format(format))
	for i := 1; i < 7; i++ {
		week = append(week, monday.AddDate(0, 0, i).Format(format))
	}
	return week
}

// WeekOfMonth 引数の日付が含まれる週（日曜開始）がその月の何週目かを返す
func WeekOfMonth(t time.Time) int {
	y, w := t.ISOWeek()
	if t.Weekday() == time.Sunday {
		w++
	}
	fom := time.Date(y, t.Month(), 1, 0, 0, 0, 0, t.Location())
	_, fw := fom.ISOWeek()
	if fom.Weekday() == time.Sunday {
		fw++
	}
	return w - fw + 1
}

// ISOWeekOfMonth 引数の日付が含まれる週(月曜開始)がその月の何週目かを返す
func ISOWeekOfMonth(t time.Time) int {
	y, w := t.ISOWeek()
	fom := time.Date(y, t.Month(), 1, 0, 0, 0, 0, t.Location())
	_, fw := fom.ISOWeek()
	return w - fw + 1
}

// Add 引数の日付に週数を加算した日付を返す
func Add(t time.Time, w int) time.Time {
	return t.AddDate(0, 0, w*7)
}

// Same 引数の日付が同じ週（日曜開始）に含まれるかを返す
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

// ISOSame 引数の日付が同じ週（月曜開始）に含まれるかを返す
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
