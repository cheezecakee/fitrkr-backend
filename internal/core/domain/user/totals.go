package user

type Totals struct {
	Workouts int     `json:"workouts"`
	Lifted   float64 `json:"lifted"` // always stored in kg
	Time     int     `json:"time"`   // minutes
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
