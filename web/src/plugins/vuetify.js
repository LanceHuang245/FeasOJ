import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import '@mdi/font/css/materialdesignicons.css'

// 获取初始主题设置
const getInitialTheme = () => {
    const savedTheme = localStorage.getItem('theme') || 'system';
    
    if (savedTheme === 'system') {
        const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        return isDark ? 'dark' : 'light';
    }
    
    return savedTheme;
};

export default createVuetify({
  theme: {
    defaultTheme: getInitialTheme(),
  },
  icons: {
    defaultSet: 'mdi',
  },
})