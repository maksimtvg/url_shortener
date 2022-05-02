// Package generator.
//max unique Uris is (2^62 - basePad), equals 4.6e18 unique uris
//keyChars is alphabet for base62 encoding
package generator

import "sync"

var (
	keyChars          = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	keyCharsLen int64 = 62
)

// number padding in order to not get small uris in the beginning
const basePad = 100_000_000

// Generator declares generator interface for generating unique uris
type Generator interface {
	GenerateUri() (string, error)
}

// UriGenerator is for generating unique uris
type UriGenerator struct {
	index int64
	mx    *sync.Mutex
}

// NewUriGenerator constructor
func NewUriGenerator(n int64) *UriGenerator {
	return &UriGenerator{
		index: n,
		mx:    &sync.Mutex{},
	}
}

// GenerateUri locks index generates unique uri and unlock index. Concurrently save
func (u *UriGenerator) GenerateUri() (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	n := u.index + basePad
	shortenedUrl := make([]byte, 20)
	i := len(shortenedUrl)

	for n > 0 && i >= 0 {
		i--
		reminder := n % keyCharsLen
		n = (n - reminder) / keyCharsLen
		shortenedUrl[i] = keyChars[reminder]
	}

	u.index++
	return string(shortenedUrl[i:]), nil
}
