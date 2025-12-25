// 支持的界面语言列表
export const SUPPORTED_LOCALES = [
  { title: "بالعربية", value: "ar", file: "ar" },
  { title: "English", value: "en", file: "en" },
  { title: "Español", value: "es", file: "es" },
  { title: "Français", value: "fr", file: "fr" },
  { title: "Italiano", value: "it", file: "it" },
  { title: "日本語", value: "ja", file: "ja" },
  { title: "Português", value: "pt", file: "pt" },
  { title: "Русский", value: "ru", file: "ru" },
  { title: "简体中文", value: "zh_CN", file: "zh_CN" },
  { title: "繁體中文", value: "zh_TW", file: "zh_TW" },
];

// 默认界面语言
export const DEFAULT_LOCALE = 'en';

// 获取语言选项列表（用于 v-select）
export const getLocaleOptions = () => {
  return SUPPORTED_LOCALES.map(locale => ({
    title: locale.title,
    value: locale.value
  }));
};

