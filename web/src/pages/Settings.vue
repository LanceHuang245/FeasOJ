<script setup>
import { ref, onMounted, watch, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useTheme } from 'vuetify';
import { ojAuthor, ojVersion } from '../utils/oj_constants';
import {
  getLanguageOptionsWithIcon,
  DEFAULT_LANGUAGE
} from '../utils/language_constants';
import {
  getLocaleOptions,
  DEFAULT_LOCALE
} from '../utils/locale_constants';

const { t, locale } = useI18n();
const vuetifyTheme = useTheme();

// 使用共享的多语言常量
const langs = ref(getLocaleOptions());

const currentLocale = ref(locale.value);

const theme = ref('system');
const preferredLanguage = ref(DEFAULT_LANGUAGE);
const showAboutDialog = ref(false);

const themes = [
  { title: t('message.light'), value: 'light', icon: 'mdi-white-balance-sunny' },
  { title: t('message.dark'), value: 'dark', icon: 'mdi-weather-night' },
  { title: t('message.system'), value: 'system', icon: 'mdi-monitor' }
];

// 使用共享的语言常量（编程语言名称不需要翻译）
const languages = getLanguageOptionsWithIcon();

const applyThemeHandler = (selectedTheme) => {
  if (selectedTheme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    vuetifyTheme.change(isDark ? 'dark' : 'light');
  } else {
    vuetifyTheme.change(selectedTheme);
  }
};

const applyLanguageHandler = (selectedLanguage) => {
  preferredLanguage.value = selectedLanguage;
};

const openAboutDialog = () => {
  showAboutDialog.value = true;
};

const handleSystemThemeChange = (e) => {
  if (theme.value === 'system') {
    applyThemeHandler('system');
  }
};

watch(currentLocale, (val) => {
  locale.value = val;
  localStorage.setItem('language', val);
});
watch(locale, (val) => {
  currentLocale.value = val;
});

onMounted(() => {
  theme.value = localStorage.getItem('theme') || 'system';
  locale.value = localStorage.getItem('language') || DEFAULT_LOCALE;
  preferredLanguage.value = localStorage.getItem('preferredLanguage') || DEFAULT_LANGUAGE;
  currentLocale.value = locale.value;
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', handleSystemThemeChange);
});

onUnmounted(() => {
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', handleSystemThemeChange);
});

watch(theme, (val) => {
  localStorage.setItem('theme', val);
  applyThemeHandler(val);
});

watch(preferredLanguage, (val) => {
  localStorage.setItem('preferredLanguage', val);
});

watch(locale, (val) => {
  localStorage.setItem('language', val);
  themes[0].title = t('message.light');
  themes[1].title = t('message.dark');
  themes[2].title = t('message.system');
  // 编程语言名称不需要翻译，无需动态更新
});
</script>

<template>
  <v-app-bar :elevation="2">
    <p style="font-size: 24px;margin-left: 20px;">{{ t('message.settings') }}</p>
  </v-app-bar>
  <v-container fluid class="settings-container">
    <v-list class="settings-list">
      <v-list-item class="settings-list-item">
        <template v-slot:prepend>
          <v-icon icon="mdi-palette" size="24"></v-icon>
        </template>
        <v-list-item-title class="settings-item-title">{{ t('message.theme') }}</v-list-item-title>
        <template v-slot:append>
          <div class="select-wrapper" @click.stop>
            <v-select v-model="theme" :items="themes" item-title="title" item-value="value" variant="outlined"
              density="compact" hide-details @update:model-value="applyThemeHandler">
              <template v-slot:item="{ props, item }">
                <v-list-item v-bind="props" class="theme-select-item">
                  <template v-slot:prepend>
                    <v-icon :icon="item.raw.icon" size="20"></v-icon>
                  </template>
                </v-list-item>
              </template>
            </v-select>
          </div>
        </template>
      </v-list-item>
      <!-- 偏好编程语言 -->
      <v-list-item class="settings-list-item">
        <template v-slot:prepend>
          <v-icon icon="mdi-code-braces" size="24"></v-icon>
        </template>
        <v-list-item-title class="settings-item-title">{{ t('message.preferredLanguage') }}</v-list-item-title>
        <template v-slot:append>
          <div class="select-wrapper" @click.stop>
            <v-select v-model="preferredLanguage" :items="languages" item-title="title" item-value="value"
              variant="outlined" density="compact" hide-details @update:model-value="applyLanguageHandler">
              <template v-slot:item="{ props, item }">
                <v-list-item v-bind="props" class="language-select-item">
                </v-list-item>
              </template>
            </v-select>
          </div>
        </template>
      </v-list-item>
      <!-- 语言 -->
      <v-list-item class="settings-list-item">
        <template v-slot:prepend>
          <v-icon icon="mdi-translate" size="24"></v-icon>
        </template>
        <v-list-item-title class="settings-item-title">{{ t('message.lang') }}</v-list-item-title>
        <template v-slot:append>
          <div class="select-wrapper" @click.stop>
            <v-select v-model="currentLocale" :items="langs" item-title="title" item-value="value" variant="outlined"
              density="compact">
              <template v-slot:item="{ props, item }">
                <v-list-item v-bind="props" class="language-select-item">
                </v-list-item>
              </template>
            </v-select>
          </div>
        </template>
      </v-list-item>

      <v-list-item class="settings-list-item" @click="openAboutDialog">
        <template v-slot:prepend>
          <v-icon icon="mdi-information" size="24"></v-icon>
        </template>
        <v-list-item-title class="settings-item-title">{{ t('message.about') }}</v-list-item-title>
        <template v-slot:append>
          <v-icon icon="mdi-chevron-right" size="20"></v-icon>
        </template>
      </v-list-item>
    </v-list>
  </v-container>

  <!-- 关于对话框 -->
  <v-dialog v-model="showAboutDialog" max-width="400" persistent>
    <v-card rounded="xl" class="about-dialog">
      <v-card-title class="about-dialog-title">
        <v-icon icon="mdi-information" color="primary" size="24" class="mr-3"></v-icon>
        {{ t('message.aboutFeasOJ') }}
      </v-card-title>
      <v-card-text class="about-dialog-content">
        <div class="about-info">
          <div class="info-item">
            <span class="info-label">{{ t('message.version') }}：</span>
            <span class="info-value">{{ ojVersion }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('message.developer') }}：</span>
            <span class="info-value">{{ ojAuthor }}</span>
          </div>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="tonal" rounded="xl" @click="showAboutDialog = false" class="about-close-btn">
          {{ t('message.ok') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.settings-container {
  margin: 0px;
  padding: 0px;
  display: flex;
  width: 100%;
  overflow-y: auto;
}

.settings-list {
  background: transparent;
  padding: 0;
  width: 100%;
  height: 100%;
}

.settings-list-item {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  min-height: 64px;
  height: 64px;
  width: 100%;
  cursor: pointer;
  align-items: center;
  backdrop-filter: blur(20px);
}

.settings-item-title {
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 0.15px;
  text-align: left;
  margin-left: 0;
}

.select-wrapper {
  pointer-events: auto;
}

.theme-select-item,
.language-select-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  transition: all 0.2s ease;
}

.theme-select-item:hover,
.language-select-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

/* 关于对话框样式 */
.about-dialog {
  border: 1px solid rgba(var(--v-theme-outline), 0.08);
  background: rgba(var(--v-theme-surface), 0.98);
  backdrop-filter: blur(20px);
}

.about-dialog-title {
  font-size: 20px;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
  padding: 24px 24px 16px 24px;
  display: flex;
  align-items: center;
}

.about-dialog-content {
  padding: 0 24px 24px 24px;
}

.about-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-radius: 12px;
}

.info-label {
  font-weight: 500;
  font-size: 14px;
}

.info-value {
  font-weight: 600;
  font-size: 14px;
}

.about-dialog-actions {
  padding: 16px 24px 24px 24px;
  justify-content: center;
}
</style>