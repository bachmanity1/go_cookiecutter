package util

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	uuid "github.com/lithammer/shortuuid/v3"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type lockedSource struct {
	lk   sync.Mutex
	seed rand.Source
}

var lSrc = &lockedSource{
	seed: rand.NewSource(time.Now().UnixNano()),
}

func (l *lockedSource) Intn(n int) (rnum int) {
	l.lk.Lock()
	rnum = rand.New(l.seed).Intn(n)
	l.lk.Unlock()
	return
}

func (l *lockedSource) Int63() (rnum int64) {
	l.lk.Lock()
	rnum = rand.New(l.seed).Int63()
	l.lk.Unlock()
	return
}

func (l *lockedSource) Float64() (rnum float64) {
	l.lk.Lock()
	rnum = rand.New(l.seed).Float64()
	l.lk.Unlock()
	return
}

func (l *lockedSource) NormFloat64() (rnum float64) {
	l.lk.Lock()
	rnum = rand.New(l.seed).NormFloat64()
	l.lk.Unlock()
	return
}

// NewID ...
func NewID() (id string) {
	rnum := lSrc.Intn(10000)
	id = fmt.Sprintf("%s%04d", strings.Replace(time.Now().Truncate(time.Millisecond).Format("20060102150405.00"), ".", "", -1), rnum)
	return
}

// RandString : make Random String
func RandString(n int) string {
	if n > 0 {
		return uuid.New()[:n]
	}
	return uuid.New()
}

// ConvertBase ...
func ConvertBase(n, base int) (s string) {
	if n == 0 {
		return "0"
	}
	marks := "0123456789ABCDEF"
	for n > 0 {
		s = fmt.Sprintf("%c%s", marks[n%base], s)
		n = n / base
	}
	return s
}

// RandNum ...
func RandNum(max int) int {
	return lSrc.Intn(max)
}

// RandLog ...
func RandLog() float64 {
	return lSrc.NormFloat64()
}

// Uniform ...
func Uniform() float64 {
	return float64(RandNum(100)) / 1000
}

// IntRayleighCDF ...
func IntRayleighCDF() int {
	return int(RayleighCDF())
}

// RayleighCDF ...
func RayleighCDF() float64 {
	rnum := lSrc.Float64()
	return math.Sqrt(-2 * math.Log(float64(1)-rnum))
}

// ArrayToString ...
func ArrayToString(a []uint, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
