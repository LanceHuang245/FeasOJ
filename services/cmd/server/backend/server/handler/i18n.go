package handler

import (
	"FeasOJ/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func GetLangFromHeader(c *gin.Context) language.Tag {
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		return language.English
	}
	tags, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	if len(tags) == 0 {
		return language.English
	}
	return tags[0]
}

// 获取国际化消息
func GetMessage(c *gin.Context, messageID string) string {
	langBundle := utils.InitI18n()
	tag := GetLangFromHeader(c)
	localizer := i18n.NewLocalizer(langBundle, tag.String())

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		localizer = i18n.NewLocalizer(langBundle, language.English.String())
		message, _ = localizer.Localize(&i18n.LocalizeConfig{
			MessageID: messageID,
		})
	}
	return message
}
