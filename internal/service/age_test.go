package service

import (
	"testing"
	"time"
)

// TestCalculateAge tests the age calculation logic with various scenarios
func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		now      time.Time
		expected int
	}{
		// Test 1: Birthday today - should be full age
		{
			name:     "birthday_today",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC),
			expected: 35,
		},
		// Test 2: Birthday tomorrow - still previous age
		{
			name:     "birthday_tomorrow",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 5, 9, 12, 0, 0, 0, time.UTC),
			expected: 34,
		},
		// Test 3: Birthday passed this month
		{
			name:     "birthday_passed_month",
			dob:      time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 5, 20, 12, 0, 0, 0, time.UTC),
			expected: 35,
		},
		// Test 4: Same month, day passed
		{
			name:     "same_month_day_passed",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
			expected: 35,
		},
		// Test 5: Same month, before birthday
		{
			name:     "same_month_before_birthday",
			dob:      time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC),
			expected: 34,
		},
		// Test 6: Future DOB (invalid, returns 0)
		{
			name:     "future_dob",
			dob:      time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		// Test 7: Leap year birthday (Feb 29)
		{
			name:     "leap_year_birthday",
			dob:      time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
		// Test 8: Edge case - Dec 31 to Jan 1
		{
			name:     "year_boundary",
			dob:      time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 34,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAge(tt.dob, tt.now)
			if got != tt.expected {
				t.Errorf("calculateAge(%v, %v) = %d; want %d",
					tt.dob.Format("2006-01-02"), tt.now.Format("2006-01-02"),
					got, tt.expected)
			}
		})
	}
}
