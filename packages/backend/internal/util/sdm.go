package util

import (
	"bytes"
	"crypto/aes"
	"fmt"

	"github.com/aead/cmac"
)

func CalculateSDMShortMacAES(key []byte, uid []byte, counter []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("key must be 16 bytes")
	}

	if len(uid) != 7 {
		return nil, fmt.Errorf("uid must be 7 bytes")
	}

	if len(counter) != 3 {
		return nil, fmt.Errorf("counter must be 3 bytes")
	}

	b := make([]byte, 16)
	b[0] = 0x3c
	b[1] = 0xc3
	b[2] = 0x00
	b[3] = 0x01
	b[4] = 0x00
	b[5] = 0x80

	idx := 6

	if !isEmptyBytes(uid) {
		copy(b[6:13], uid)
		idx = 13
	}

	if !isEmptyBytes(counter) {
		// little endianにする
		counter[0], counter[2] = counter[2], counter[0]
		copy(b[idx:idx+3], counter[:])
		idx += 3
	}

	for idx < 16 {
		b[idx] = 0
		idx++
	}

	sk, err := calculateCMAC(key, b)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate session key: %w", err)
	}

	cm, err := calculateCMAC(sk, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CMAC: %w", err)
	}

	return toShortMAC(cm), nil
}

func calculateCMAC(key []byte, b []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	h, err := cmac.NewWithTagSize(ci, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to create CMAC: %w", err)
	}

	if b != nil {
		if _, err := h.Write(b); err != nil {
			return nil, fmt.Errorf("failed to write to CMAC: %w", err)
		}
	}

	return h.Sum(nil), nil
}

func toShortMAC(mac []byte) []byte {
	return []byte{mac[1], mac[3], mac[5], mac[7], mac[9], mac[11], mac[13], mac[15]}
}

func isEmptyBytes(b []byte) bool {
	return bytes.Equal(b, make([]byte, len(b)))
}
