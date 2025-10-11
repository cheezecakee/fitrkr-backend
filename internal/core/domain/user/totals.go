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

func (t *Totals) RecordWorkout(lifted Weight, duration Duration) {
	t.Workouts++
	t.Lifted += lifted.Value
	t.Time += duration.Minutes()
}
