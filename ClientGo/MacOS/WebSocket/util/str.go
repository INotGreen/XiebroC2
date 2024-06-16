package util

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func RandInt(min, max int) int {
	// if we get nonsense values
	// give them random int anyway
	if min > max ||
		min < 0 ||
		max < 0 {
		min = RandInt(0, 100)
		max = min + RandInt(0, 100)
	}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		log.Println("cannot seed math/rand package with cryptographically secure random number generator")
		log.Println("falling back to math/rand with time seed")
		return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	return min + rand.Intn(max-min)
}

func SplitString(input string) []string {
	var result []string
	regex := regexp.MustCompile(`(\".*?\"|\S+)`)
	matches := regex.FindAllString(input, -1)

	for _, match := range matches {
		result = append(result, strings.Trim(match, "\""))
	}

	return result
}
