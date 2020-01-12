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

func DetectLanguage() (language string, err error) {
	locale, err := getUnixLocale()
	if err == nil {
		language, _ = splitLocale(locale)
	}
	return
}

func DetectTerritory() (territory string, err error) {
	locale, err := getUnixLocale()
	if err == nil {
		_, territory = splitLocale(locale)
	}
	return
}
