// +build darwin freebsd linux netbsd openbsd

package jibber_jabber

import (
	"errors"
	"os"
	"strings"
)

func getLangFromEnv() (locale string) {
	envs := []string{"LC_MESSAGES", "LC_ALL", "LANG"}

	for _, env := range envs {
		locale = os.Getenv(env)
		if len(locale) > 0 {
			return
		}
	}
	return
}

func getUnixLocale() (locale string, err error) {
	locale = getLangFromEnv()
	if len(locale) <= 0 {
		err = errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE)
	}
	return
}

// DetectIETF detects and returns the IETF language tag of UNIX systems, like Linux and macOS.
// If a territory is defined, the returned value will be in the format of `[language]-[territory]`,
// e.g. `en-GB`.
func DetectIETF() (locale string, err error) {
	locale, err := getUnixLocale()
	if err == nil {
		language, territory := splitLocale(locale)
		locale = language
		if len(territory) > 0 {
			locale = strings.Join([]string{language, territory}, "-")
		}
	}
	return
}

// DetectLanguage detects the IETF language tag of UNIX systems, like Linux and macOS,
// and returns the first half of the string, before the `_`.
func DetectLanguage() (language string, err error) {
	locale, err := getUnixLocale()
	if err == nil {
		language, _ = splitLocale(locale)
	}
	return
}

// DetectTerritory detects the IETF language tag of UNIX systems, like Linux and macOS,
// and returns the second half of the string, after the `_`.
func DetectTerritory() (territory string, err error) {
	locale, err := getUnixLocale()
	if err == nil {
		_, territory = splitLocale(locale)
	}
	return
}
