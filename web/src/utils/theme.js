import { useTheme } from 'vuetify';

// 主题管理工具
export const initTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'system';
  
  if (savedTheme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return isDark ? 'dark' : 'light';
  }
  
  return savedTheme;
};

// 监听系统主题变化
export const watchSystemTheme = (callback) => {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  
  const handleChange = (e) => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'system') {
      const isDark = e.matches;
      callback(isDark ? 'dark' : 'light');
    }
  };
  
  mediaQuery.addEventListener('change', handleChange);
  
  // 返回清理函数
  return () => mediaQuery.removeEventListener('change', handleChange);
};

// 获取当前主题
export const getCurrentTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'system';
  
  if (savedTheme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return isDark ? 'dark' : 'light';
  }
  
  return savedTheme;
};

// 获取MdEditor主题配置
export const getMdEditorTheme = () => {
  const currentTheme = getCurrentTheme();
  return currentTheme === 'dark' ? 'dark' : 'light';
};

// 获取MdPreview主题配置
export const getMdPreviewTheme = () => {
  const currentTheme = getCurrentTheme();
  return currentTheme === 'dark' ? 'dark' : 'light';
}; 