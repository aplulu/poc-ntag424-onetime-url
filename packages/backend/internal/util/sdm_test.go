package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSDMShortMacAES(t *testing.T) {
	testCases := []struct {
		name    string
		uid     []byte
		counter []byte
		key     []byte
		expect  []byte
	}{
		{
			name: "Factory Key",
			// Factory Key
			// 043923627E7580000021F99802BF8FA8C446
			uid:     []byte{0x04, 0x39, 0x23, 0x62, 0x7e, 0x75, 0x80},
			counter: []byte{0x00, 0x00, 0x21}, // 33
			key:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expect:  []byte{0xf9, 0x98, 0x02, 0xbf, 0x8f, 0xa8, 0xc4, 0x46},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CalculateSDMShortMacAES(tc.key, tc.uid, tc.counter)

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestIsEmptyBytes(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		expect bool
	}{
		{
			name:   "zero length",
			input:  []byte{},
			expect: true,
		},
		{
			name:   "empty",
			input:  []byte{0x00},
			expect: true,
		},
		{
			name:   "not empty",
			input:  []byte{0x00, 0x01},
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isEmptyBytes(tc.input)

			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestToShortMAC(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		expect []byte
	}{
		{
			name:   "valid",
			input:  []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
			expect: []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := toShortMAC(tc.input)

			assert.Equal(t, tc.expect, actual)
		})
	}
}
