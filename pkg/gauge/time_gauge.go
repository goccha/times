package gauge

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

func New(begin, end time.Time) *TimeGauge {
	tg := &TimeGauge{
		begin: begin,
		end:   end,
	}
	_ = tg.Duration()
	return tg
}

type TimeGauge struct {
	date     string
	begin    time.Time
	end      time.Time
	duration *time.Duration
}

func (t *TimeGauge) Date() string {
	return t.date
}
func (t *TimeGauge) Begin() time.Time {
	return t.begin
}
func (t *TimeGauge) End() time.Time {
	return t.end
}
func (t *TimeGauge) Duration() time.Duration {
	if t.duration == nil {
		d := t.end.Sub(t.begin)
		t.duration = &d
	}
	return *t.duration
}
func (t *TimeGauge) Seconds() float64 {
	return t.Duration().Seconds()
}
func (t *TimeGauge) Hours() float64 {
	return t.Duration().Hours()
}

func (t *TimeGauge) Rounds() (hour int, minute int, second int) {
	d := t.Duration().Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return int(h), int(m), int(s)
}

func (t *TimeGauge) RoundAll() (hour int, minute int, second int, milli int, micro int, nano int) {
	d := t.Duration()
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ml := d / time.Millisecond
	d -= ml * time.Millisecond
	mi := d / time.Microsecond
	d -= mi * time.Microsecond
	n := d / time.Nanosecond
	return int(h), int(m), int(s), int(ml), int(mi), int(n)
}

func (t TimeGauge) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = io.WriteString(s, t.Duration().String())
	case 'h':
		d := t.Duration().Round(time.Hour)
		h := d / time.Hour
		_, _ = io.WriteString(s, strconv.Itoa(int(h)))
	case 'm':
		d := t.Duration().Round(time.Minute)
		h := d / time.Hour
		d -= h * time.Hour
		m := d / time.Minute
		_, _ = io.WriteString(s, strconv.Itoa(int(m)))
	case 's':
		_, _, sec := t.Rounds()
		_, _ = io.WriteString(s, strconv.Itoa(int(sec)))
	case 'S':
		d := t.Duration().Round(time.Millisecond)
		h := d / time.Hour
		d -= h * time.Hour
		m := d / time.Minute
		d -= m * time.Minute
		sec := d / time.Second
		d -= sec * time.Second
		msec := d / time.Millisecond
		_, _ = io.WriteString(s, fmt.Sprintf("%03d", msec))
	case 'M':
		d := t.Duration().Round(time.Microsecond)
		h := d / time.Hour
		d -= h * time.Hour
		m := d / time.Minute
		d -= m * time.Minute
		sec := d / time.Second
		d -= sec * time.Second
		ml := d / time.Millisecond
		d -= ml * time.Millisecond
		mi := d / time.Microsecond
		_, _ = io.WriteString(s, fmt.Sprintf("%03d", mi))
	case 'n':
		_, _, _, _, _, n := t.RoundAll()
		_, _ = io.WriteString(s, fmt.Sprintf("%03d", n))
	}
}
func (t *TimeGauge) Split(hour, min, sec, nsec int, loc *time.Location) []TimeGauge {
	b := baseTime{
		hour: hour,
		min:  min,
		sec:  sec,
		nsec: nsec,
		loc:  loc,
	}
	sleepTime := t.end.Sub(t.begin)
	if sleepTime < 0 {
		return []TimeGauge{}
	}
	m := make(map[string]*TimeGauge)
	calc(b, t.begin, sleepTime, m)
	times := make([]TimeGauge, 0, len(m))
	for _, v := range m {
		times = append(times, *v)
	}
	return times
}

type baseTime struct {
	hour int
	min  int
	sec  int
	nsec int
	loc  *time.Location
}

func (t *baseTime) Time(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), t.hour, t.min, t.sec, t.nsec, t.loc)
}

func calc(b baseTime, tm time.Time, sleepTime time.Duration, times map[string]*TimeGauge) {
	if sleepTime <= 0 {
		return
	}
	base := b.Time(tm)
	diff := base.Sub(tm)
	if diff > 0 { // 基準時刻よりも前に就寝した場合
		key := tm.Format("2006-01-02")
		if sleepTime < diff { // 基準時刻よりも前に起床している場合
			times[key] = &TimeGauge{
				date:  key,
				begin: tm,
				end:   tm.Add(sleepTime),
			}
			return
		} else {
			times[key] = &TimeGauge{
				date:  key,
				begin: tm,
				end:   base,
			}
			sleepTime -= diff
			calc(b, base, sleepTime, times)
			return
		}
	} else {
		base = base.AddDate(0, 0, 1) // 基準日時を翌日にする
		key := base.Format("2006-01-02")
		diff := base.Sub(tm)  // 基準日時までの時間を算出
		if sleepTime < diff { // 基準日時以前に起床している場合
			if v, ok := times[key]; ok {
				v.end = tm.Add(sleepTime)
			} else {
				times[key] = &TimeGauge{
					date:  key,
					begin: tm,
					end:   tm.Add(sleepTime),
				}
			}
			return
		} else { // 基準日時以降に起床している場合
			if v, ok := times[key]; ok {
				v.end = base
			} else {
				times[key] = &TimeGauge{
					date:  key,
					begin: tm,
					end:   base,
				}
			}
			sleepTime -= diff
			calc(b, base, sleepTime, times)
			return
		}
	}
}

func (t *TimeGauge) Overlap(start time.Time, end time.Time) bool {
	if t.begin.Before(start) && t.end.After(start) {
		return true
	}
	if t.begin.Before(end) && t.end.After(end) {
		return true
	}
	if start.Before(t.begin) && end.After(t.begin) {
		return true
	}
	if start.Before(t.end) && end.After(t.end) {
		return true
	}
	return false
}
