package subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDateFormat(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		wantErr bool
	}{
		{"Valid date", "01-2025", false},
		{"Valid date", "11-2025", false},
		{"Invalid one digit date", "1-2025", true},
		{"Invalid date -- no dash", "012025", true},
		{"Invalid date -- have days", "01-12-2025", true},
		{"Invalid date -- invalid month - low", "00-2025", true},
		{"Invalid date -- invalid month - high", "14-2025", true},
		{"Invalid date -- invalid year - low", "02-1000", true},
		{"Invalid date -- invalid year - high", "02-20000", true},
		{"Invalid date -- invalid year - lett", "02-absd", true},
		{"Empty date", "", false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateDateFormat(testCase.dateStr)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
