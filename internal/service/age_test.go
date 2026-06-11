package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name string
		dob  time.Time
		now  time.Time
		want int
	}{
		{
			name: "birthday passed",
			dob:  time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			want: 36,
		},
		{
			name: "birthday upcoming",
			dob:  time.Date(1990, 12, 10, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			want: 35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateAge(tt.dob, tt.now); got != tt.want {
				t.Fatalf("CalculateAge() = %d, want %d", got, tt.want)
			}
		})
	}
}
