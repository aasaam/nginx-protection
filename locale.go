package main

import (
	"embed"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v2"
)

//go:embed locale/*
var localeEmbed embed.FS

var translateData map[string]map[string]string

func isFileLocaleExist(path, lang string) bool {
	if _, err := os.Stat(path + "/" + lang + ".yml"); err == nil {
		return true
	}
	return false
}

func loadFileLocale(path, lang string) map[string]string {
	file, err := os.ReadFile(path + "/" + lang + ".yml")
	if err != nil {
		panic(err)
	}

	var tr map[string]string
	err = yaml.Unmarshal(file, &tr)
	if err != nil {
		panic(err)
	}

	icuLang := language.Make(lang)
	p := message.NewPrinter(icuLang)

	tr["dir"] = getLanguageDirection(lang)
	tr["name"] = p.Sprintf("%v", display.Language(icuLang))

	return tr
}

func loadEmbedLocale(lang string) map[string]string {
	file, err := localeEmbed.ReadFile("locale/" + lang + ".yml")
	if err != nil {
		panic(err)
	}

	var tr map[string]string
	err = yaml.Unmarshal(file, &tr)
	if err != nil {
		panic(err)
	}

	icuLang := language.Make(lang)
	p := message.NewPrinter(icuLang)

	tr["dir"] = getLanguageDirection(lang)
	tr["name"] = p.Sprintf("%v", display.Language(icuLang))

	return tr
}

func loadLocales(config *config) {
	result := make(map[string]map[string]string)

	result["en"] = loadEmbedLocale("en")

	for _, lang := range config.supportedLanguages {
		if lang == "en" {
			continue
		}
		result[lang] = loadEmbedLocale(lang)
	}

	for lang, tdData := range result {
		if isFileLocaleExist(config.localePath, lang) {
			mergeData := loadFileLocale(config.localePath, "en")
			for key := range tdData {
				if newTranslate, ok := mergeData[key]; ok {
					result[lang][key] = newTranslate
				}
			}
		}
	}

	translateData = result
}

func languagesData(languages []string, currentLanguage string) map[string]string {
	result := make(map[string]string)
	for _, lang := range languages {
		if lang == currentLanguage {
			continue
		}
		icuLang := language.Make(lang)
		p := message.NewPrinter(icuLang)
		result[lang] = p.Sprintf("%v", display.Language(icuLang))
	}
	return result
}

func getLanguageDirection(lang string) string {
	if _, ok := rtlLanguagesMap[lang]; ok {
		return "rtl"
	}
	return "ltr"
}

func isSupportedLanguage(lang string) bool {
	return supportedLanguagesMap[lang]
}

func isSupportedLanguageConfig(lang string, supported []string) bool {
	for _, supportedLanguage := range supported {
		if lang == supportedLanguage {
			return true
		}
	}
	return false
}
