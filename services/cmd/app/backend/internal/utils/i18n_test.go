package utils

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestInitI18n(t *testing.T) {
	bundle := InitI18n()
	if bundle == nil {
		t.Fatal("InitI18n 返回了 nil bundle")
	}
}

func TestLocalizeExistingMessage(t *testing.T) {
	bundle := InitI18n()
	localizer := i18n.NewLocalizer(bundle, "en")

	_, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "loginSuccess",
		DefaultMessage: &i18n.Message{
			ID:    "loginSuccess",
			Other: "default login success",
		},
	})
	if err != nil {
		t.Fatalf("Localize 出错: %v", err)
	}
}

func TestLocalizeMessageWithTemplate(t *testing.T) {
	bundle := InitI18n()
	localizer := i18n.NewLocalizer(bundle, "en")

	_, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "problem_completed",
		TemplateData: map[string]interface{}{
			"PID": "123",
		},
		DefaultMessage: &i18n.Message{
			ID:    "problem_completed",
			Other: "default problem completed",
		},
	})
	if err != nil {
		t.Fatalf("Localize 出错: %v", err)
	}
}

func TestLocalizeFallbackMessage(t *testing.T) {
	bundle := InitI18n()
	localizer := i18n.NewLocalizer(bundle, "en")

	_, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "nonexistent_message",
		DefaultMessage: &i18n.Message{
			ID:    "nonexistent_message",
			Other: "fallback message",
		},
	})
	if err != nil {
		t.Fatalf("Localize 出错: %v", err)
	}
}
