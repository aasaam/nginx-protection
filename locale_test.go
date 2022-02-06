package main

import (
	"testing"
)

func TestLocaleFunction(t *testing.T) {
	langs1 := []string{"fa", "en"}
	languagesData(langs1, "en")
	langs2 := []string{"fa", "en"}
	languagesData(langs2, "ar")

	if isSupportedLangauge("en") != true {
		t.Errorf("en is valid")
	}

	if isSupportedLangauge("11") != false {
		t.Errorf("11 is in valid")
	}

	if getLanguageDirection("fa") != "rtl" {
		t.Errorf("fa dir is rtl")
	}

	if getLanguageDirection("en") != "ltr" {
		t.Errorf("en dir is ltr")
	}
}

func TestLocaleData1(t *testing.T) {
	c1 := newConfig("error", false, "en", "en,fa", "", "", "", "", "", "test/locale")
	loadLocales(c1)

	c2 := newConfig("error", false, "fa", "en,fa", "", "", "", "", "", "")
	loadLocales(c2)
}
