/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package utils

import (
	"html/template"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// RePositiveInterger represents the regular expression pattern for matching positive integers.
	RePositiveInterger = "^[1-9]\\d*$"

	// RePositiveScientificNotation represents the regular expression pattern
	// for matching positive numbers in scientific notation.
	RePositiveScientificNotation = "^(\\d+(.{0}|.\\d+))[Ee]{1}([\\+|-]?\\d+)$"

	// RePositiveFloatPoint represents the regular expression pattern for matching positive floating-point numbers.
	RePositiveFloatPoint = "^(?:[1-9][0-9]*\\.[0-9]+|0\\.(?!0+$)[0-9]+)$"

	// ReURL represents the regular expression pattern for matching URLs.
	ReURL = "[\\w-]+(/[\\w-./?%&=]*)?"

	// ReFileName represents the regular expression pattern for matching valid file names.
	ReFileName = "^[a-zA-Z0-9-_\\.]+$"

	// ReChinesePhone represents the regular expression pattern for matching Chinese phone numbers.
	ReChinesePhone = "^1\\d{10}$"
)

// IsPositiveInterger checks if the given string is a positive integer.
func IsPositiveInterger(num string) bool {
	return isMatchRegex(RePositiveInterger, num)
}

// IsPositiveScientificNotation checks if the given string is a positive scientific notation number.
func IsPositiveScientificNotation(num string) bool {
	return isMatchRegex(RePositiveScientificNotation, num)
}

// IsPositiveFloatPoint checks if the given string is a positive floating-point number.
func IsPositiveFloatPoint(num string) bool {
	return isMatchRegex(RePositiveFloatPoint, num)
}

// IsSafeFileName checks if the given string is a safe file name.
func IsSafeFileName(name string) bool {
	return isMatchRegex(ReFileName, name)
}

// IsPath checks if the given string is a valid URL path.
func IsPath(url string) bool {
	return isMatchRegex(ReURL, url)
}

// IsChinesePhone checks if the given string is a Chinese phone number.
func IsChinesePhone(phone string) bool {
	return isMatchRegex(ReChinesePhone, phone)
}

// IsUrl checks if the given string is a valid URL.
func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsPictureName checks if the given picture name has an allowed extension.
func IsPictureName(pictureName string) bool {
	ext := filepath.Ext(pictureName)
	ext = strings.ToLower(ext)

	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	return allowed
}

// IsInt checks if the given string can be converted to an integer.
func IsInt(input string) bool {
	_, err := strconv.Atoi(input)
	return err == nil
}

// IsTxt checks if the given file name has a .txt extension.
func IsTxt(fileName string) bool {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	allowedExtensions := []string{".txt"}
	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	return allowed
}

func isMatchRegex(pattern string, v string) bool {
	matched, err := regexp.MatchString(pattern, v)
	if err != nil {
		return false
	}

	return matched
}

// XSSEscapeString escapes the input string for safe use in HTML content.
func XSSEscapeString(input string) (output string) {
	return template.HTMLEscapeString(input)
}
