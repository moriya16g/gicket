package i18n

import (
	"fmt"
	"os"
	"strings"
)

// Lang は現在の言語設定
var Lang = "en"

func init() {
	// 環境変数 GICKET_LANG > LANG で判定
	if v := os.Getenv("GICKET_LANG"); v != "" {
		SetLang(v)
		return
	}
	if v := os.Getenv("LANG"); v != "" {
		if strings.HasPrefix(strings.ToLower(v), "ja") {
			Lang = "ja"
		}
	}
}

// SetLang は言語を設定する ("en" or "ja")
func SetLang(lang string) {
	l := strings.ToLower(strings.TrimSpace(lang))
	if strings.HasPrefix(l, "ja") {
		Lang = "ja"
	} else {
		Lang = "en"
	}
}

// T はメッセージキーから現在の言語の文字列を返す
func T(key string) string {
	if msgs, ok := messages[Lang]; ok {
		if s, ok := msgs[key]; ok {
			return s
		}
	}
	// fallback to English
	if s, ok := messages["en"][key]; ok {
		return s
	}
	return key
}

// Tf はフォーマット付きメッセージを返す
func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

var messages = map[string]map[string]string{
	"en": messagesEN,
	"ja": messagesJA,
}
