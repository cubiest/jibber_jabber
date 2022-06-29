package jibberjabber

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var (
	ErrLangDetectFail          = errors.New("could not detect Language")
	ErrLangFallbackUndefined   = errors.New("no fallback language defined")
	ErrLangFallbackUnsupported = errors.New("defined fallback language is not supported")
	ErrLangUnsupported         = errors.New("language not supported")
	ErrLangParse               = errors.New("language identifier cannot be parsed")
)

func splitLocale(locale string) (string, string) {
	formattedLocale := strings.Split(locale, ".")[0]
	formattedLocale = strings.Replace(formattedLocale, "-", "_", -1)

	pieces := strings.Split(formattedLocale, "_")
	language := pieces[0]
	territory := ""
	if len(pieces) > 1 {
		territory = strings.Split(formattedLocale, "_")[1]
	}
	return language, territory
}

/**
 * languageServer
 */

type languageServer struct {
	supportedLanguages map[language.Tag]string // the string can be used to link to a localization file for that language
	fallbackLanguage   language.Tag
}

var (
	languageServerSingletonOnce sync.Once
	languageServerInstance      *languageServer
	languageServerMutex         = &sync.Mutex{}
)

func LanguageServer() *languageServer {
	languageServerSingletonOnce.Do(func() {
		if languageServerInstance == nil {
			languageServerInstance = new(languageServer)
		}
	})
	return languageServerInstance
}

// SetSupportedLanguages defines the supported languages checked against in other funcs.
// The values (type `string`) can be used to link to a localization file for that language.
func (server *languageServer) SetSupportedLanguages(supported map[language.Tag]string) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	server.supportedLanguages = supported
}

// GetSupportedLanguages returns the supported languages.
func (server *languageServer) GetSupportedLanguages() map[language.Tag]string {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	return server.supportedLanguages
}

// SetFallbackLanguage defines the language used as a fallback language Tag if any other func returns no valid value.
func (server *languageServer) SetFallbackLanguage(fallback language.Tag) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	server.fallbackLanguage = fallback
}

// GetFallbackLanguage returns the language fallback.
func (server *languageServer) GetFallbackLanguage() language.Tag {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	return server.fallbackLanguage
}

// DetectSupportedLanguage returns the language tag detected from the system.
// If it's not supported, it returns the fallback.
// Returns ErrLangParse, if library cannot detect language or parse value given from your operating system.
// Returns ErrLangFallbackUndefined, if fallback is undefined.
// Returns ErrLangFallbackUnsupported, if fallaback is defined but unsupported.
// If you want to check for jibberjabber errors, call `jibberjabber.IsError()`.
func (server *languageServer) DetectSupportedLanguage() (language.Tag, error) {

	tag, err := DetectLanguageTag()
	if err != nil {
		return language.Und, fmt.Errorf("%v: %w", ErrLangParse.Error(), err)
	}

	if server.LanguageTagIsSupported(tag) {
		return tag, nil
	}

	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	fallbackTag := server.fallbackLanguage
	if fallbackTag == language.Und {
		return language.Und, ErrLangFallbackUndefined
	} else if _, supported := server.supportedLanguages[fallbackTag]; !supported {
		return language.Und, ErrLangFallbackUnsupported
	} else {
		return fallbackTag, nil
	}
}

// ListSupportedLanguages returns the language tags in a language.Tag slice.
func (server *languageServer) ListSupportedLanguages() []language.Tag {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	supportedLangTags := make([]language.Tag, 0, len(server.supportedLanguages))

	for tag := range server.supportedLanguages {
		supportedLangTags = append(supportedLangTags, tag)
	}

	return supportedLangTags
}

// ListSupportedLanguagesAsStrings returns the language tags in a slice of string representation of the language tags.
func (server *languageServer) ListSupportedLanguagesAsStrings() []string {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	supportedLangs := make([]string, 0, len(server.supportedLanguages))

	for tag := range server.supportedLanguages {
		supportedLangs = append(supportedLangs, tag.String())
	}

	return supportedLangs
}

// ListSupportedLanguagesAsStringsSorted returns the language tags in a slice of string representation of the language tags, alphabetically sorted.
func (server *languageServer) ListSupportedLanguagesAsStringsSorted() []string {
	supportedLangs := server.ListSupportedLanguagesAsStrings()
	sort.Strings(supportedLangs)
	return supportedLangs
}

// ListSupportedLanguagesForDisplay returns the language tags in a slice of human readable strings.
func (server *languageServer) ListSupportedLanguagesForDisplay() []string {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	supportedLangs := make([]string, 0, len(server.supportedLanguages))

	for tag := range server.supportedLanguages {
		supportedLangs = append(supportedLangs, display.Self.Name(tag))
	}

	return supportedLangs
}

// ListSupportedLanguagesForDisplaySorted returns the language tags in a string slice, alphabetically sorted.
func (server *languageServer) ListSupportedLanguagesForDisplaySorted() []string {
	supportedLangs := server.ListSupportedLanguagesForDisplay()
	sort.Strings(supportedLangs)
	return supportedLangs
}

// ListSupportedLanguagesSorted returns the language tags + their strings sorted alphabetically by string.
// Use the elements for the first return value as key for the second return value.
func (server *languageServer) ListSupportedLanguagesSorted() ([]string, map[string]language.Tag) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	supportedLangs := make([]string, 0, len(server.supportedLanguages))
	supportedLangTags := make(map[string]language.Tag)

	for tag := range server.supportedLanguages {
		name := display.Self.Name(tag)
		supportedLangs = append(supportedLangs, name)
		supportedLangTags[name] = tag
	}

	sort.Strings(supportedLangs)

	return supportedLangs, supportedLangTags
}

// LanguageIsSupported returns true if the given BCP 47 string is in the list of supported languages.
// Returns ErrLangParse, if any parsing issue occured.
// If you want to check for jibberjabber errors, call `jibberjabber.IsError()`.
func (server *languageServer) LanguageIsSupported(bcp string) (bool, error) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	lang, parseErr := language.Parse(bcp)
	if parseErr != nil {
		return false, fmt.Errorf("%v: %w", ErrLangParse.Error(), parseErr)
	}

	_, supported := server.supportedLanguages[lang]

	return supported, nil
}

// LanguageTagIsSupported returns true if the given language tag is in the list of supported languages.
func (server *languageServer) LanguageTagIsSupported(lang language.Tag) bool {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	_, supported := server.supportedLanguages[lang]

	return supported
}

// StringToLanguageTag returns language tag for given BCP 47 string.
// Returns ErrLangParse, if parsing fails.
// If you want to check for jibberjabber errors, call `jibberjabber.IsError()`.
func (server *languageServer) StringToLanguageTag(bcp string) (language.Tag, error) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	lang, parseErr := language.Parse(bcp)
	if parseErr != nil {
		return language.Und, fmt.Errorf("%v: %w", ErrLangParse.Error(), parseErr)
	}
	return lang, nil
}

// StringToSupportedLanguageTag returns language tag for given BCP 47 string.
// Returns specified fallback, if language is not supported or parsing to language.Tag fails.
// Returns ErrLangUnsupported, if language could be parsed, but is not supported.
// Returns ErrLangFallbackUndefined, if ErrLangUnsupported and fallback is undefined.
// Returns ErrLangFallbackUnsupported, if ErrLangUnsupported and fallaback is defined but unsupported.
func (server *languageServer) StringToSupportedLanguageTag(bcp string) (language.Tag, error) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	var err error

	lang, parseErr := language.Parse(bcp)
	if parseErr == nil {
		if _, supported := server.supportedLanguages[lang]; supported {
			return lang, nil
		} else {
			err = ErrLangUnsupported
		}
	}

	lang = server.fallbackLanguage

	if _, supported := server.supportedLanguages[lang]; !supported {
		if server.fallbackLanguage == language.Und {
			err = ErrLangFallbackUndefined
		} else {
			err = ErrLangFallbackUnsupported
		}
	}

	return lang, err
}

// GetSupportedLanguageValue returns the value for given BCP 47 string.
// If parsing fails or the language is not supported, it will use the fallback's value.
// Returns ErrLangUnsupported, if language could be parsed, but is not supported.
// Returns ErrLangFallbackUndefined, if ErrLangUnsupported and fallback is undefined.
// Returns ErrLangFallbackUnsupported, if ErrLangUnsupported and fallaback is defined but unsupported.
func (server *languageServer) GetSupportedLanguageValue(bcp string) (string, error) {

	tag, err := server.StringToSupportedLanguageTag(bcp)
	if errors.Is(err, ErrLangFallbackUndefined) || errors.Is(err, ErrLangFallbackUnsupported) {
		return "", err
	}

	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	return server.supportedLanguages[tag], err
}

// GetSupportedLanguageValueByTag returns the value of the requested language tag.
// If the language is not supported, it will use the fallback's value.
// Returns ErrLangUnsupported, if language could be parsed, but is not supported.
// Returns ErrLangFallbackUndefined, if ErrLangUnsupported and fallback is undefined.
// Returns ErrLangFallbackUnsupported, if ErrLangUnsupported and fallaback is defined but unsupported.
func (server *languageServer) GetSupportedLanguageValueByTag(lang language.Tag) (string, error) {
	languageServerMutex.Lock()
	defer languageServerMutex.Unlock()

	var err error

	value, supported := server.supportedLanguages[lang]
	if !supported {
		err = ErrLangUnsupported
		value, supported = server.supportedLanguages[server.fallbackLanguage]
		if !supported {
			if server.fallbackLanguage == language.Und {
				err = ErrLangFallbackUndefined
			} else {
				err = ErrLangFallbackUnsupported
			}
		}
	}

	return value, err
}

// IsError checks an error you received from one of jibberjabber's funcs for a jibberjabber error like `ErrLangDetectFail`.
// Reason you cannot use e.g. `errors.Is()`: currently, golang does not allow native chain-wrapping errors. Therefore, `errors.Unwrap()`, `errors.Is()` & Co. won't return `true` for jibberjabber errors.
func IsError(err error, jjError error) bool {
	return strings.HasPrefix(err.Error(), jjError.Error())
}
