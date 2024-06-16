package http

import "testing"

func TestIsHex(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "valid hex",
			input:  "0123456789abcdefABCDEF",
			expect: true,
		},
		{
			name:   "invalid hex",
			input:  "ghijklmnopqrstuvwxyz",
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isHex(tc.input)

			if actual != tc.expect {
				t.Errorf("expected: %v, got: %v", tc.expect, actual)
			}
		})
	}
}
