package components

import (
	"fmt"
	"testing"
)

func TestGenerateTimeOptions(t *testing.T) {
	tests := []struct {
		name            string
		intervalMinutes int
		wantLength      int
		wantFirst       timeOption
		wantLast        timeOption
	}{
		{
			name:            "15-minute interval",
			intervalMinutes: 15,
			wantLength:      96,
			wantFirst:       timeOption{hour: 12, minute: 0, period: "AM"},
			wantLast:        timeOption{hour: 11, minute: 45, period: "PM"},
		},
		{
			name:            "30-minute interval",
			intervalMinutes: 30,
			wantLength:      48,
			wantFirst:       timeOption{hour: 12, minute: 0, period: "AM"},
			wantLast:        timeOption{hour: 11, minute: 30, period: "PM"},
		},
		{
			name:            "60-minute interval",
			intervalMinutes: 60,
			wantLength:      24,
			wantFirst:       timeOption{hour: 12, minute: 0, period: "AM"},
			wantLast:        timeOption{hour: 11, minute: 0, period: "PM"},
		},
		{
			name:            "1-minute interval",
			intervalMinutes: 1,
			wantLength:      1440,
			wantFirst:       timeOption{hour: 12, minute: 0, period: "AM"},
			wantLast:        timeOption{hour: 11, minute: 59, period: "PM"},
		},
		{
			name:            "5-minute interval",
			intervalMinutes: 5,
			wantLength:      288,
			wantFirst:       timeOption{hour: 12, minute: 0, period: "AM"},
			wantLast:        timeOption{hour: 11, minute: 55, period: "PM"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateTimeOptions(tt.intervalMinutes)

			if len(got) != tt.wantLength {
				t.Errorf("generateTimeOptions() length = %v, want %v", len(got), tt.wantLength)
			}

			if len(got) > 0 {
				if got[0] != tt.wantFirst {
					t.Errorf("generateTimeOptions() first = %v, want %v", got[0], tt.wantFirst)
				}

				if got[len(got)-1] != tt.wantLast {
					t.Errorf("generateTimeOptions() last = %v, want %v", got[len(got)-1], tt.wantLast)
				}
			}
		})
	}
}

func TestTimeOptionFormatting(t *testing.T) {
	tests := []struct {
		name       string
		timeOption timeOption
		wantHour   string
		wantMinute string
	}{
		{
			name:       "Midnight",
			timeOption: timeOption{hour: 12, minute: 0, period: "AM"},
			wantHour:   "12",
			wantMinute: "00",
		},
		{
			name:       "Noon",
			timeOption: timeOption{hour: 12, minute: 0, period: "PM"},
			wantHour:   "12",
			wantMinute: "00",
		},
		{
			name:       "Single Digit Hour",
			timeOption: timeOption{hour: 1, minute: 30, period: "AM"},
			wantHour:   "01",
			wantMinute: "30",
		},
		{
			name:       "Double Digit Hour",
			timeOption: timeOption{hour: 11, minute: 59, period: "PM"},
			wantHour:   "11",
			wantMinute: "59",
		},
		{
			name:       "Single Digit Minute",
			timeOption: timeOption{hour: 5, minute: 5, period: "AM"},
			wantHour:   "05",
			wantMinute: "05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHour := fmt.Sprintf("%02d", tt.timeOption.hour)
			gotMinute := fmt.Sprintf("%02d", tt.timeOption.minute)

			if gotHour != tt.wantHour {
				t.Errorf("Hour formatting: got %v, want %v", gotHour, tt.wantHour)
			}

			if gotMinute != tt.wantMinute {
				t.Errorf("Minute formatting: got %v, want %v", gotMinute, tt.wantMinute)
			}
		})
	}
}
