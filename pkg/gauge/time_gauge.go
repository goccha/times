package gauge

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// New 指定期間の時間計測機を生成する
func New(begin, end time.Time) *TimeGauge {
	tg := &TimeGauge{
		begin: begin,
		end:   end,
	}
	_ = tg.Duration()
	return tg
}

// TimeGauge 時間計測機
type TimeGauge struct {
	date     string
	begin    time.Time
	end      time.Time
	duration *time.Duration
}

// Date 日付を返す
func (t *TimeGauge) Date() string {
	if t.date == "" {
		t.date = t.begin.Format("2006-01-02")
	}
	return t.date
}

// Begin 開始日時を返す
func (t *TimeGauge) Begin() time.Time {
	return t.begin
}

// End 終了日時を返す
func (t *TimeGauge) End() time.Time {
	return t.end
}

// Duration 期間を返す
func (t *TimeGauge) Duration() time.Duration {
	if t.duration == nil {
		d := t.end.Sub(t.begin)
		t.duration = &d
	}
	return *t.duration
}

// Seconds returns the duration as a floating point number of seconds.
func (t *TimeGauge) Seconds() float64 {
	return t.Duration().Seconds()
}

// Minutes returns the duration as a floating point number of minutes.
func (t *TimeGauge) Minutes() float64 {
	return t.Duration().Minutes()
}

// Hours returns the duration as a floating point number of hours.
func (t *TimeGauge) Hours() float64 {
	return t.Duration().Hours()
}

// Days returns the duration as a floating point number of days(24h).
func (t *TimeGauge) Days() float64 {
	return t.Duration().Hours() / 24
}

// Weeks returns the duration as a floating point number of weeks(168h).
func (t *TimeGauge) Weeks() float64 {
	return t.Duration().Hours() / 168
}

// Months returns the duration as a floating point number of months(720h).
func (t *TimeGauge) Months() float64 {
	return t.Duration().Hours() / 720
}

// Years returns the duration as a floating point number of years(8760h).
func (t *TimeGauge) Years() float64 {
	return t.Duration().Hours() / 8760
}

// Rounds returns the number of hours, minutes, and seconds.
func (t *TimeGauge) Rounds() (hour int, minute int, second int) {
	d := t.Duration().Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return int(h), int(m), int(s)
}

// RoundAll returns the number of hours, minutes, seconds, milliseconds, and nanoseconds.
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

// Format 文字列に変換する
func (t *TimeGauge) Format(s fmt.State, verb rune) {
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

// Split 期間を基準時刻で分割する
func (t *TimeGauge) Split(hour, min, sec, ns int, loc *time.Location) []TimeGauge {
	b := baseTime{
		hour: hour,
		min:  min,
		sec:  sec,
		ns:   ns,
		loc:  loc,
	}
	timeRange := t.end.Sub(t.begin)
	if timeRange < 0 {
		return []TimeGauge{}
	}
	m := make(map[string]*TimeGauge)
	split(b, t.begin, timeRange, m)
	times := make([]TimeGauge, 0, len(m))
	for _, v := range m {
		times = append(times, *v)
	}
	return times
}

// baseTime 基準時刻
type baseTime struct {
	hour int
	min  int
	sec  int
	ns   int
	loc  *time.Location
}

// Time 指定日の基準時刻(time.Time)を返す
func (t baseTime) Time(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), t.hour, t.min, t.sec, t.ns, t.loc)
}

// split 開始日時からの経過時間を基準時刻で日付毎に分割する
func split(b baseTime, tm time.Time, timeRange time.Duration, times map[string]*TimeGauge) {
	if timeRange <= 0 {
		return
	}
	base := b.Time(tm)
	diff := base.Sub(tm)
	if diff > 0 { // 開始が基準時刻よりも前
		key := tm.Format("2006-01-02")
		if timeRange < diff { // 基準時刻よりも前に終了している場合
			times[key] = &TimeGauge{
				date:  key,
				begin: tm,
				end:   tm.Add(timeRange),
			}
			return
		} else {
			times[key] = &TimeGauge{
				date:  key,
				begin: tm,
				end:   base,
			}
			timeRange -= diff
			split(b, base, timeRange, times)
			return
		}
	} else {
		base = base.AddDate(0, 0, 1) // 基準日時を翌日にする
		key := base.Format("2006-01-02")
		diff = base.Sub(tm)   // 基準日時までの時間を算出
		if timeRange < diff { // 基準日時以前に開始している場合
			if v, ok := times[key]; ok {
				v.end = tm.Add(timeRange)
			} else {
				times[key] = &TimeGauge{
					date:  key,
					begin: tm,
					end:   tm.Add(timeRange),
				}
			}
			return
		} else { // 基準日時以降に開始している場合
			if v, ok := times[key]; ok {
				v.end = base
			} else {
				times[key] = &TimeGauge{
					date:  key,
					begin: tm,
					end:   base,
				}
			}
			timeRange -= diff
			split(b, base, timeRange, times)
			return
		}
	}
}

// Overlap 指定した期間と重複しているかどうか
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

// Contains 指定した日時が期間内に含まれるかどうか
func (t *TimeGauge) Contains(tm time.Time) bool {
	return t.begin.Before(tm) && t.end.After(tm)
}
