package internal

import (
	"fmt"
	"os"
	"strconv"
)

// Those suggest value are all set according to
// https://github.com/Trojan-Qt5/shadowsocks-org/issues/44#issuecomment-281021054
// Due to this package contains various internal implementation so const named with DefaultBR prefix
const (
	DefaultSFCapacity = 1e6
	// FalsePositiveRate
	DefaultSFFPR  = 1e-6
	DefaultSFSlot = 10
)

const EnvironmentPrefix = "SHADOWSOCKS_"

// A shared instance used for checking salt repeat
var saltfilter *BloomRing

func init() {
	var (
		finalCapacity = DefaultSFCapacity
		finalFPR      = DefaultSFFPR
		finalSlot     = float64(DefaultSFSlot)
	)
	for _, opt := range []struct {
		ENVName string
		Target  *float64
	}{
		{
			ENVName: "CAPACITY",
			Target:  &finalCapacity,
		},
		{
			ENVName: "FPR",
			Target:  &finalFPR,
		},
		{
			ENVName: "SLOT",
			Target:  &finalSlot,
		},
	} {
		envKey := EnvironmentPrefix + "SF_" + opt.ENVName
		env := os.Getenv(envKey)
		if env != "" {
			p, err := strconv.ParseFloat(env, 64)
			if err != nil {
				panic(fmt.Sprintf("Invalid envrionment `%s` setting in saltfilter: %s", envKey, env))
			}
			*opt.Target = p
		}
	}
	// Support disable saltfilter by given a negative capacity
	if finalCapacity <= 0 {
		return
	}
	saltfilter = NewBloomRing(int(finalSlot), int(finalCapacity), finalFPR)
}

// TestSalt returns true if salt is repeated
func TestSalt(b []byte) bool {
	// If nil means feature disabled, return false to bypass salt repeat detection
	if saltfilter == nil {
		return false
	}
	return saltfilter.Test(b)
}

// AddSalt salt to filter
func AddSalt(b []byte) {
	// If nil means feature disabled
	if saltfilter == nil {
		return
	}
	saltfilter.Add(b)
}