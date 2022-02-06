package main

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	c1 := newConfig("info", true, "en", "", "", "", "", "", "", "")
	c1.getLogger().Info().Msg("info")
	c1.getLogger().Trace().Msg("trace")

	c2 := newConfig("panic", false, "fa", "", "", "", "", "http://cdn.example.com", "", "")
	c2.getLogger().Info().Msg("info")

	c3 := newConfig("panic", false, "zz", "zz,cc", "", "", "", "http://cdn.example.com", "", "")
	c3.getLogger().Info().Msg("info")
}
