package i18n

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	messages     map[string]map[string]string
	messagesOnce sync.Once
)

// Init loads embedded flat locale files: locales/{lang}.json -> { errorCode: message }.
func Init() {
	messagesOnce.Do(func() {
		messages = make(map[string]map[string]string)
		_ = fs.WalkDir(localeFS, "locales", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
				return err
			}

			lang := strings.TrimSuffix(filepath.Base(path), ".json")
			data, readErr := localeFS.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			var locale map[string]string
			if unmarshalErr := json.Unmarshal(data, &locale); unmarshalErr != nil {
				return unmarshalErr
			}
			messages[lang] = locale
			return nil
		})
	})
}

func localize(lang, messageKey string, _ map[string]interface{}) string {
	Init()
	if messages == nil {
		return messageKey
	}

	for _, candidate := range languageCandidates(lang) {
		if locale, ok := messages[candidate]; ok {
			if msg, ok := locale[messageKey]; ok && msg != "" {
				return msg
			}
		}
	}
	return messageKey
}

func languageCandidates(lang string) []string {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return []string{"en"}
	}

	out := []string{lang}
	if primary, _, ok := strings.Cut(lang, "-"); ok && primary != lang {
		out = append(out, primary)
	}
	if lang != "en" {
		out = append(out, "en")
	}
	return out
}
