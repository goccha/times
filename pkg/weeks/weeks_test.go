package weeks

import (
	"testing"
	"time"
)

func TestTimes(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	times := Times(tm)
	if len(times) != 7 {
		t.Errorf("expected=7, actual=%d", len(times))
		return
	}
	expected := "2020-03-29T10:00:00+09:00"
	actual := times[0].Format(time.RFC3339)
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
	expected = "2020-04-04T10:00:00+09:00"
	actual = times[6].Format(time.RFC3339)
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
}

func TestISOTimes(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	times := ISOTimes(tm)
	if len(times) != 7 {
		t.Errorf("expected=7, actual=%d", len(times))
		return
	}
	expected := "2020-03-30T10:00:00+09:00"
	actual := times[0].Format(time.RFC3339)
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
	expected = "2020-04-05T10:00:00+09:00"
	actual = times[6].Format(time.RFC3339)
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
}

func TestDayStrings(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	times := DayStrings(tm)
	if len(times) != 7 {
		t.Errorf("expected=7, actual=%d", len(times))
		return
	}
	expected := "2020-03-29"
	actual := times[0]
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
	expected = "2020-04-04"
	actual = times[6]
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
}

func TestISODayStrings(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	times := ISODayStrings(tm)
	if len(times) != 7 {
		t.Errorf("expected=7, actual=%d", len(times))
		return
	}
	expected := "2020-03-30"
	actual := times[0]
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
	expected = "2020-04-05"
	actual = times[6]
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
}

func TestWeekOfMonth(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	actual := WeekOfMonth(tm)
	if actual != 1 {
		t.Errorf("expected=1, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-04-05T10:00:00+09:00")
	actual = WeekOfMonth(tm)
	if actual != 2 {
		t.Errorf("expected=2, actual=%d", actual)
	}

	tm, _ = time.Parse(time.RFC3339, "2020-03-31T10:00:00+09:00")
	actual = WeekOfMonth(tm)
	if actual != 5 {
		t.Errorf("expected=5, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-03-29T10:00:00+09:00")
	actual = WeekOfMonth(tm)
	if actual != 5 {
		t.Errorf("expected=5, actual=%d", actual)
	}
}

func TestISOWeekOfMonth(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	actual := ISOWeekOfMonth(tm)
	if actual != 1 {
		t.Errorf("expected=1, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-04-05T10:00:00+09:00")
	actual = ISOWeekOfMonth(tm)
	if actual != 1 {
		t.Errorf("expected=1, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-04-06T10:00:00+09:00")
	actual = ISOWeekOfMonth(tm)
	if actual != 2 {
		t.Errorf("expected=2, actual=%d", actual)
	}

	tm, _ = time.Parse(time.RFC3339, "2020-03-01T10:00:00+09:00")
	actual = ISOWeekOfMonth(tm)
	if actual != 1 {
		t.Errorf("expected=1, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-03-02T10:00:00+09:00")
	actual = ISOWeekOfMonth(tm)
	if actual != 2 {
		t.Errorf("expected=1, actual=%d", actual)
	}
	tm, _ = time.Parse(time.RFC3339, "2020-03-31T10:00:00+09:00")
	actual = ISOWeekOfMonth(tm)
	if actual != 6 {
		t.Errorf("expected=1, actual=%d", actual)
	}
}

func TestAdd(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-04-01T10:00:00+09:00")
	w := Add(tm, 1)
	expected := "2020-04-08T10:00:00+09:00"
	actual := w.Format(time.RFC3339)
	if expected != actual {
		t.Errorf("expected=%s, actual=%s", expected, actual)
		return
	}
}

func TestSame(t *testing.T) {
	tm1, _ := time.Parse(time.RFC3339, "2020-03-29T10:00:00+09:00")
	tm2, _ := time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	if !Same(tm1, tm2) {
		t.Error("expected=true, actual=false")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2020-04-05T10:00:00+09:00")
	if Same(tm1, tm2) {
		t.Error("expected=false, actual=true")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2019-12-29T10:00:00+09:00")
	if !Same(tm1, tm2) {
		t.Error("expected=true, actual=false")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2020-04-10T10:00:00+09:00")
	if Same(tm1, tm2) {
		t.Error("expected=false, actual=true")
	}
}

func TestISOSame(t *testing.T) {
	tm1, _ := time.Parse(time.RFC3339, "2020-03-29T10:00:00+09:00")
	tm2, _ := time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	if ISOSame(tm1, tm2) {
		t.Error("expected=false, actual=true")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2020-04-05T10:00:00+09:00")
	if !ISOSame(tm1, tm2) {
		t.Error("expected=true, actual=false")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2019-12-29T10:00:00+09:00")
	if ISOSame(tm1, tm2) {
		t.Error("expected=false, actual=true")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-04-04T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2020-04-10T10:00:00+09:00")
	if ISOSame(tm1, tm2) {
		t.Error("expected=false, actual=true")
	}
	tm1, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00+09:00")
	tm2, _ = time.Parse(time.RFC3339, "2019-12-30T10:00:00+09:00")
	if !ISOSame(tm1, tm2) {
		t.Error("expected=true, actual=false")
	}
}
