package v1

func FilterByTime(set []*MeterSession, model *MeterSession) ([]*MeterSession, error) {
	if model == nil {
		return set, nil
	}
	filter := Window(model.StartTime, model.EndTime).Filter()
	filtered := []*MeterSession{}
	for _, item := range set {
		if filter(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

type SessionFilterFn func(*MeterSession) bool

func (fn SessionFilterFn) And(fn2 SessionFilterFn) SessionFilterFn {
	return func(m *MeterSession) bool { return fn(m) && fn2(m) }
}

type TimeWindow interface {
	Start() int64
	End() int64
	Filter() SessionFilterFn
}

type window struct {
	start int64
	end   int64
}

var _ TimeWindow = (*window)(nil)

func WindowFromModel(model *MeterSession) TimeWindow {
	if model == nil {
		return Window(0, 0)
	}
	return Window(model.StartTime, model.EndTime)
}

func Window(start int64, end int64) TimeWindow {
	w := &window{start: start, end: end}
	if w.start == 0 && w.end == 0 {
		return w
	}
	if w.start > w.end {
		w.start, w.end = w.end, w.start
	}
	if w.start == 0 {
		w.start = w.end
	}
	if w.end == 0 {
		w.end = w.start
	}
	return w
}

func (w *window) Start() int64 {
	if w != nil {
		return w.start
	}
	return 0
}

func (w *window) End() int64 {
	if w != nil {
		return w.end
	}
	return 0
}

func (w *window) Filter() SessionFilterFn {
	if w != nil && w.start != 0 && w.end != 0 {
		return sessionIntersects(w)
	}
	return passthrough
}

func sessionIntersects(window TimeWindow) SessionFilterFn {
	return func(item *MeterSession) bool {
		// filter out invalid records
		if item == nil || item.StartTime == 0 {
			return false
		}
		// record.start < window.end && (record.end == 0 || record.end > window.start)
		return item.StartTime < window.End() && (item.EndTime == 0 || item.EndTime > window.Start())
	}
}

func passthrough(i *MeterSession) bool { return i != nil && i.StartTime != 0 }
