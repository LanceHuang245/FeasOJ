import { createI18n } from "vue-i18n";
import ar from "./locales/ar.js";
import en from "./locales/en.js";
import es from "./locales/es.js";
import fr from "./locales/fr.js";
import it from "./locales/it.js";
import ja from "./locales/ja.js";
import pt from "./locales/pt.js";
import ru from "./locales/ru.js";
import zh_CN from "./locales/zh_CN.js";
import zh_TW from "./locales/zh_TW.js";
import { DEFAULT_LOCALE } from "../utils/locale_constants";

const messages = {
  ar,
  en,
  es,
  fr,
  it,
  ja,
  pt,
  ru,
  zh_CN,
  zh_TW
};

export const i18n = createI18n({
  legacy: false,
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages,
})