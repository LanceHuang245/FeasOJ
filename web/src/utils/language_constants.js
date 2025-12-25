// 统一的编程语言配置（单一数据源）
export const SUPPORTED_LANGUAGES = [
  { value: 'C++', label: 'C++', ext: 'cpp', aceLang: 'c_cpp', icon: 'mdi-language-cpp' },
  { value: 'Java', label: 'Java', ext: 'java', aceLang: 'java', icon: 'mdi-language-java' },
  { value: 'Pascal', label: 'Pascal', ext: 'pas', aceLang: 'pascal', icon: 'mdi-language-pascal' },
  { value: 'Python', label: 'Python', ext: 'py', aceLang: 'python', icon: 'mdi-language-python' },
  { value: 'Php', label: 'PHP', ext: 'php', aceLang: 'php', icon: 'mdi-language-php' },
  { value: 'Rust', label: 'Rust', ext: 'rs', aceLang: 'rust', icon: 'mdi-language-rust' },
  { value: 'Golang', label: 'Go', ext: 'go', aceLang: 'golang', icon: 'mdi-language-go' },
];

// 默认编程语言
export const DEFAULT_LANGUAGE = 'C++';

// 根据语言值获取文件扩展名
export const getFileExtension = (langValue) => {
  const lang = SUPPORTED_LANGUAGES.find(l => l.value === langValue);
  return lang ? lang.ext : 'txt';
};

// 根据语言值获取 Ace Editor 的语言标识符
export const getAceLang = (langValue) => {
  const lang = SUPPORTED_LANGUAGES.find(l => l.value === langValue);
  return lang ? lang.aceLang : 'text';
};

// 获取语言选项列表
export const getLanguageOptions = () => {
  return SUPPORTED_LANGUAGES.map(lang => lang.value);
};

// 获取带图标的语言选项列表
export const getLanguageOptionsWithIcon = () => {
  return SUPPORTED_LANGUAGES.map(lang => ({
    title: lang.label,
    value: lang.value,
    icon: lang.icon
  }));
};
