// +build windows

package jibberjabber

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TODO: read LOCALE_NAME_MAX_LENGTH from Windows, instead of hard-coding '85'
const LOCALE_NAME_MAX_LENGTH uint32 = 85

var SUPPORTED_LANGUAGE_PACKS = map[uintptr]string{
	0x0401: "ar-SA",
	0x0402: "bg-BG",
	0x0804: "zh-CN",
	0x0404: "zh-TW",
	0x041a: "hr-HR",
	0x0405: "cs-CZ",
	0x0406: "da-DK",
	0x0413: "nl-NL",
	0x0409: "en-US",
	0x0809: "en-GB",
	0x0425: "et-EE",
	0x040b: "fi-FI",
	0x0c0c: "fr-CA",
	0x040c: "fr-FR",
	0x0407: "de-DE",
	0x0408: "el-GR",
	0x040d: "he-IL",
	0x040e: "hu-HU",
	0x0410: "it-IT",
	0x0411: "ja-JP",
	0x0412: "ko-KR",
	0x0426: "lv-LV",
	0x0427: "lt-LT",
	0x0414: "nb-NO",
	0x0415: "pl-PL",
	0x0416: "pt-BR",
	0x0816: "pt-PT",
	0x0418: "ro-RO",
	0x0419: "ru-RU",
	0x241A: "sr-Latn-RS",
	0x041b: "sk-SK",
	0x0424: "sl-SI",
	0x080a: "es-MX",
	0x0c0a: "es-ES",
	0x041d: "sv-SE",
	0x041e: "th-TH",
	0x041f: "tr-TR",
	0x0422: "uk-UA",
}

func getWindowsLocaleFrom(sysCall string) (string, error) {
	buffer := make([]uint16, LOCALE_NAME_MAX_LENGTH)

	dll, err := windows.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("could not find kernel32 dll: " + err.Error())
	}

	proc, err := dll.FindProc(sysCall)
	if err != nil {
		return "", err
	}

	r, _, dllError := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(LOCALE_NAME_MAX_LENGTH))
	if r == 0 {
		return "", errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
	}

	return windows.UTF16ToString(buffer), nil
}

func getAllWindowsLocaleFrom(sysCall string) (string, error) {
	dll, err := windows.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("could not find kernel32 dll: " + err.Error())
	}

	proc, err := dll.FindProc(sysCall)
	if err != nil {
		return "", err
	}

	locale, _, dllError := proc.Call()
	if locale == 0 {
		return "", errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
	}

	return SUPPORTED_LANGUAGE_PACKS[locale], nil
}

func getWindowsLocale() (string, error) {
	dll, err := windows.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("could not find kernel32 dll: " + err.Error())
	}

	proc, err := dll.FindProc("GetVersion")
	if err != nil {
		return "", err
	}

	v, _, _ := proc.Call()
	windowsVersion := byte(v)

	// is of version Windows Vista or greater
	if windowsVersion >= 6 {
		locale, err := getWindowsLocaleFrom("GetUserDefaultLocaleName")
		if err != nil {
			locale, err = getWindowsLocaleFrom("GetSystemDefaultLocaleName")
		}
		return locale, err
	}

	// prior Windows Vista
	locale, err := getAllWindowsLocaleFrom("GetUserDefaultLCID")
	if err != nil {
		locale, err = getAllWindowsLocaleFrom("GetSystemDefaultLCID")
	}
	return locale, err
}

// DetectIETF detects and returns the IETF language tag of Windows.
func DetectIETF() (string, error) {
	return getWindowsLocale()
}

// DetectLanguage detects the IETF language tag of Windows
// and returns the first half of the string, before the `_`.
func DetectLanguage() (string, error) {
	locale, err := getWindowsLocale()
	if err != nil {
		return "", err
	}
	language, _ := splitLocale(locale)
	return language, nil
}

// DetectLanguageTag detects the IETF language tag of Windows
// and returns a fitting language tag.
func DetectLanguageTag() (language.Tag, error) {
	locale, err := getWindowsLocale()
	if err != nil {
		return language.Und, err
	}
	return language.Parse(locale)
}

// DetectTerritory detects the IETF language tag of Windows
// and returns the second half of the string, after the `_`.
func DetectTerritory() (string, error) {
	locale, err := getWindowsLocale()
	if err != nil {
		return "", nil
	}
	_, territory := splitLocale(locale)
	return territory, nil
}
