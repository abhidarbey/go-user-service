package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
	"unsafe"
)

const (
	letterBytes   = "0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetMd5 function
func GetMd5(input string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateID function
func GenerateID(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
