// +build windows

package jibberjabber

import (
	"errors"
	"syscall"
	"unsafe"
)

const LOCALE_NAME_MAX_LENGTH uint32 = 85

var SUPPORTED_LOCALES = map[uintptr]string{
	0x0407: "de-DE",
	0x0409: "en-US",
	0x0c0a: "es-ES", //or is it 0x040a
	0x040c: "fr-FR",
	0x0410: "it-IT",
	0x0411: "ja-JA",
	0x0412: "ko_KR",
	0x0416: "pt-BR",
	0x0419: "ru_RU",
	0x0804: "zh-CN",
	0x0c04: "zh-HK",
	0x0404: "zh-TW",
}

func getWindowsLocaleFrom(sysCall string) (string, error) {
	buffer := make([]uint16, LOCALE_NAME_MAX_LENGTH)

	dll := syscall.MustLoadDLL("kernel32")
	proc := dll.MustFindProc(sysCall)
	r, _, dllError := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(LOCALE_NAME_MAX_LENGTH))
	if r == 0 {
		return "", errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
	}

	return syscall.UTF16ToString(buffer), nil
}

func getAllWindowsLocaleFrom(sysCall string) (string, error) {
	dll, err := syscall.LoadDLL("kernel32")
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

	return SUPPORTED_LOCALES[locale], nil
}

func getWindowsLocale() (string, error) {
	dll, err := syscall.LoadDLL("kernel32")
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
		locale, err = getWindowsLocaleFrom("GetUserDefaultLocaleName")
		if err != nil {
			locale, err = getWindowsLocaleFrom("GetSystemDefaultLocaleName")
		}
		return locale, err
	}

	// prior Windows Vista
	locale, err = getAllWindowsLocaleFrom("GetUserDefaultLCID")
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
