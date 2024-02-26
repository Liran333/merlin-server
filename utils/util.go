/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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

// LoadFromYaml reads a YAML file from the given path and unmarshals it into the provided interface.
func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

// Now returns the current Unix timestamp as an int64.
func Now() int64 {
	return time.Now().Unix()
}

// ToDate converts the given Unix timestamp to a date string formatted according to the layout variable.
func ToDate(n int64) string {
	if n == 0 {
		n = Now()
	}

	return time.Unix(n, 0).Format(layout)
}

// Time returns the current time formatted according to the timeLayout variable.
func Time() string {
	return time.Now().Format(timeLayout)
}

// DateAndTime returns a pair of strings representing the date and time for the given Unix timestamp.
func DateAndTime(n int64) (string, string) {
	if n <= 0 {
		return "", ""
	}

	t := time.Unix(n, 0)

	return t.Format(layout), t.Format(timeLayout)
}

// Expiry calculates the expiration Unix timestamp by adding the given duration in seconds to the current time.
func Expiry(expiry int64) int64 {
	return time.Now().Add(time.Second * time.Duration(expiry)).Unix()
}

// StrLen returns the length of a string in terms of runes (UTF-8 characters).
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
