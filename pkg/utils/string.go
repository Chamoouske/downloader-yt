package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

func SanitizeFilename(name string) string {
	if name == "" {
		return "file"
	}

	ext := filepath.Ext(name)
	base := name[:len(name)-len(ext)]

	n := norm.NFKD.String(base)
	var buf []rune
	for _, r := range n {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		buf = append(buf, r)
	}
	s := string(buf)

	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")

	reInvalid := regexp.MustCompile(`[^\p{L}\p{N}\._ -]+`)
	s = reInvalid.ReplaceAllString(s, "_")

	reMulti := regexp.MustCompile(`[_\-.]{2,}`)
	s = reMulti.ReplaceAllString(s, "_")

	s = strings.Trim(s, "._- ")

	if s == "" {
		s = "file"
	}

	lower := strings.ToLower(s)
	reserved := map[string]struct{}{
		"con": {}, "prn": {}, "aux": {}, "nul": {},
		"com1": {}, "com2": {}, "com3": {}, "com4": {}, "com5": {}, "com6": {}, "com7": {}, "com8": {}, "com9": {},
		"lpt1": {}, "lpt2": {}, "lpt3": {}, "lpt4": {}, "lpt5": {}, "lpt6": {}, "lpt7": {}, "lpt8": {}, "lpt9": {},
	}
	if _, ok := reserved[lower]; ok {
		s = s + "_"
	}

	const maxRunes = 255
	if utf8.RuneCountInString(s) > maxRunes {
		rs := []rune(s)
		s = string(rs[:maxRunes])
		s = strings.TrimRight(s, "._- ")
		if s == "" {
			s = "file"
		}
	}

	return s + ext
}

func GetEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
