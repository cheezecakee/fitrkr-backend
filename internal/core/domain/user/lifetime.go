package user

type Totals struct {
	workouts int
	volume   float64 // always stored in kg
	time     int     // minutes
}

func NewTotals() Totals {
	return Totals{
		workouts: 0,
		volume:   0,
		time:     0,
	}
}

func (t Totals) AddWorkout() Totals {
	return Totals{
		workouts: t.workouts + 1,
		volume:   t.volume,
		time:     t.time,
	}
}

func (t Totals) AddVolume(weight Weight) Totals {
	return Totals{
		workouts: t.workouts,
		volume:   t.volume + weight.ToKg(),
		time:     t.time,
	}
}

func (t Totals) AddTime(minutes int) Totals {
	return Totals{
		workouts: t.workouts,
		volume:   t.volume,
		time:     t.time + minutes,
	}
}

func (t Totals) Workouts() int {
	return t.workouts
}

func (t Totals) Volume() float64 {
	return t.volume
}

func (t Totals) TimeMinutes() int {
	return t.time
}

func (t Totals) RecordWorkout(volume Weight, minutes int) Totals {
	return Totals{
		workouts: t.workouts + 1,
		volume:   t.volume + volume.ToKg(),
		time:     t.time + minutes,
	}
}
