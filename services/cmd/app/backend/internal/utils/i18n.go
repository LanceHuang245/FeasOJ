package utils

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var LocaleFS embed.FS

// 初始化i18n文件
func InitI18n() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	bundle.LoadMessageFileFS(LocaleFS, "locales/en.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/zh_CN.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/zh_TW.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/ja.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/ru.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/es.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/fr.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/it.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/ar.json")
	bundle.LoadMessageFileFS(LocaleFS, "locales/pt.json")

	return bundle
}
