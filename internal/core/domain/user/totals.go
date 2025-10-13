package user

type Totals struct {
	Workouts int
	Lifted   float64 // always stored in kg
	Time     int     // minutes
}

func NewTotals() Totals {
	return Totals{
		Workouts: 0,
		Lifted:   0,
		Time:     0,
	}
}

func (t *Totals) RecordWorkout(lifted WeightValue, duration Duration) {
	t.Workouts++
	t.Lifted += float64(lifted)
	t.Time += duration.Minutes()
}
