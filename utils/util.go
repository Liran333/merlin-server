package utils

import (
	"net/url"
	"os"
	"time"
	"unicode/utf8"

	"sigs.k8s.io/yaml"
)

const (
	layout     = "2006-01-02"
	timeLayout = "2006-01-02 15:04:05"
)

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func Now() int64 {
	return time.Now().Unix()
}

func ToDate(n int64) string {
	if n == 0 {
		n = Now()
	}

	return time.Unix(n, 0).Format(layout)
}

func Time() string {
	return time.Now().Format(timeLayout)
}

func DateAndTime(n int64) (string, string) {
	if n <= 0 {
		return "", ""
	}

	t := time.Unix(n, 0)

	return t.Format(layout), t.Format(timeLayout)
}

func Expiry(expiry int64) int64 {
	return time.Now().Add(time.Second * time.Duration(expiry)).Unix()
}

func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

// ExtractDomain extract hostname in URL
func ExtractDomain(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	domain := parsedURL.Hostname()

	return domain, nil
}
