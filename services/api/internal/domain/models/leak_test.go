package models

import (
	"testing"
)

func TestLeak_Validate(t *testing.T) {
	tests := []struct {
		name      string
		leak      Leak
		wantErr   error
	}{
		{
			name: "valid leak",
			leak: Leak{
				Amount:     100,
				Confidence: 50,
			},
			wantErr: nil,
		},
		{
			name: "invalid amount (zero)",
			leak: Leak{
				Amount:     0,
				Confidence: 50,
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "invalid amount (negative)",
			leak: Leak{
				Amount:     -10,
				Confidence: 50,
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "invalid confidence (negative)",
			leak: Leak{
				Amount:     100,
				Confidence: -1,
			},
			wantErr: ErrInvalidConfidence,
		},
		{
			name: "invalid confidence (over 100)",
			leak: Leak{
				Amount:     100,
				Confidence: 101,
			},
			wantErr: ErrInvalidConfidence,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.leak.Validate()
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}