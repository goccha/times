package gauge

import (
	"fmt"
	"testing"
	"time"
)

func TestCalc(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T17:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-01T18:00:00+09:00")
	rec := New(begin, end)
	times := rec.Split(18, 0, 0, 0, time.Local)
	if len(times) != 1 {
		t.Errorf("expected=1, actual=%d", len(times))
		return
	}
	if times[0].Date() != "2020-04-01" {
		t.Errorf("expected=2020-04-01, actual=%s", times[0].Date())
		return
	}
}

func TestCalc2(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T17:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-01T19:00:00+09:00")
	rec := New(begin, end)
	times := rec.Split(18, 0, 0, 0, time.Local)
	if len(times) != 2 {
		t.Errorf("[length] expected=2, actual=%d", len(times))
		return
	}
	if times[0].Date() != "2020-04-01" {
		t.Errorf("expected=2020-04-01, actual=%s", times[0].Date())
		return
	}
	base, _ := time.Parse(time.RFC3339, "2020-04-01T18:00:00+09:00")
	if times[0].end != base {
		t.Errorf("[0] expected=%v, actual=%v", base, times[0].end)
		return
	}

	if times[1].Date() != "2020-04-02" {
		t.Errorf("[1] expected=2020-04-02, actual=%s", times[1].Date())
		return
	}
	if times[1].begin != base {
		t.Errorf("[1].Start expected=%v, actual=%v", base, times[1].begin)
		return
	}
	if times[1].end != end {
		t.Errorf("[1].end expected=%v, actual=%v", end, times[1].end)
		return
	}
}

func TestCalc3(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T23:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-05T07:00:00+09:00")
	rec := New(begin, end)
	times := rec.Split(18, 0, 0, 0, time.Local)
	if len(times) != 4 {
		t.Errorf("expected=4, actual=%d", len(times))
		return
	}
	for _, v := range times {
		fmt.Printf("%v\n", v)
	}
}

func TestCalc4(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T17:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-01T17:00:00+09:00")
	rec := New(begin, end)
	times := rec.Split(18, 0, 0, 0, time.Local)
	if len(times) != 0 {
		t.Errorf("expected=0, actual=%d", len(times))
		return
	}
}

func TestTimeGauge_Hours(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T22:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-02T06:15:00+09:00")
	rec := New(begin, end)
	hours := rec.Hours()
	if hours != 8.25 {
		t.Errorf("expected=8.25, actual=%v", hours)
	}
}

func TestTimeGauge_Rounds(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T22:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-02T06:15:12.012345678+09:00")
	rec := New(begin, end)
	hour, minute, sec := rec.Rounds()
	if !(hour == 8 && minute == 15 && sec == 12) {
		t.Errorf("expected=8,15,12,actual=%d,%d,%d", hour, minute, sec)
	}
}

func TestTimeGauge_RoundAll(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T22:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-02T06:15:12.012345678+09:00")
	rec := New(begin, end)
	hour, minute, sec, mill, micro, nano := rec.RoundAll()
	if !(hour == 8 && minute == 15 && sec == 12 && mill == 12 && micro == 345 && nano == 678) {
		t.Errorf("expected=8,15,12,12,345,678,actual=%d,%d,%d,%d,%d,%d", hour, minute, sec, mill, micro, nano)
	}
}

func TestTimeGauge_Format(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T22:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-02T06:15:12+09:00")
	rec := New(begin, end)
	actual := fmt.Sprintf("%h時間%m分%s秒", rec, rec, rec)
	expected := "8時間15分12秒"
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	}

	begin, _ = time.Parse(time.RFC3339, "2020-04-01T22:00:00+09:00")
	end, _ = time.Parse(time.RFC3339, "2020-04-02T06:15:12.023123456+09:00")
	rec = New(begin, end)
	actual = fmt.Sprintf("%h時間%m分%s秒.%S%M%n", rec, rec, rec, rec, rec, rec)
	expected = "8時間15分12秒.023123456"
	if actual != expected {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	}
}

func TestTimeGauge_Overlap(t *testing.T) {
	begin, _ := time.Parse(time.RFC3339, "2020-04-01T23:00:00+09:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-02T06:15:00+09:00")
	rec := New(begin, end)
	begin, _ = time.Parse(time.RFC3339, "2020-04-01T23:01:00+09:00")
	end, _ = time.Parse(time.RFC3339, "2020-04-02T06:16:00+09:00")
	if !rec.Overlap(begin, end) {
		t.Error("expected=true, actual=false")
	}
	begin, _ = time.Parse(time.RFC3339, "2020-04-01T22:59:00+09:00")
	end, _ = time.Parse(time.RFC3339, "2020-04-02T06:14:00+09:00")
	if !rec.Overlap(begin, end) {
		t.Error("expected=true, actual=false")
	}
	begin, _ = time.Parse(time.RFC3339, "2020-04-01T22:59:00+09:00")
	end, _ = time.Parse(time.RFC3339, "2020-04-02T06:16:00+09:00")
	if !rec.Overlap(begin, end) {
		t.Error("expected=true, actual=false")
	}
	begin, _ = time.Parse(time.RFC3339, "2020-04-01T22:50:00+09:00")
	end, _ = time.Parse(time.RFC3339, "2020-04-01T22:59:00+09:00")
	if rec.Overlap(begin, end) {
		t.Error("expected=false, actual=true")
	}
}
