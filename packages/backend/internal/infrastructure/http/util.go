package http

import "strings"

func isHex(s string) bool {
	sl := strings.ToUpper(s)
	for _, r := range sl {
		if r < '0' || '9' < r && r < 'A' || 'F' < r {
			return false
		}
	}
	return true
}
