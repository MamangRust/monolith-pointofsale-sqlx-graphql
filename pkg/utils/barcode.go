package utils

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"math/rand"
)

func GenerateBarcode(name string) string {
	var initials []rune
	for _, word := range strings.Fields(name) {
		if len(word) > 0 {
			initials = append(initials, unicode.ToUpper(rune(word[0])))
		}
	}
	if len(initials) == 0 {
		initials = []rune{'P', 'R', 'D'}
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	numbers := r.Intn(900000) + 100000

	return string(initials) + strconv.Itoa(numbers)
}
