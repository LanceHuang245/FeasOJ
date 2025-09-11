<script setup>
import { VApp, VMain } from 'vuetify/components';
import VNavigation from './components/VNavigation.vue'
import VFill from './components/VFill.vue';
import VAlert from './components/VAlert.vue';
import { RouterView } from 'vue-router';
import Notification from './components/Notification.vue';
import { useTheme } from 'vuetify';
import { onMounted, onUnmounted, watch } from 'vue';
import { watchSystemTheme } from './utils/theme';

const vuetifyTheme = useTheme();

// 初始化主题
const initTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'system';
  
  if (savedTheme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    vuetifyTheme.change(isDark ? 'dark' : 'light');
  } else {
    vuetifyTheme.change(savedTheme);
  }
};

// 监听系统主题变化
let cleanup = null;

// 发送主题变化事件
const emitThemeChange = (theme) => {
  window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
};

onMounted(() => {
  initTheme();
  cleanup = watchSystemTheme((theme) => {
    vuetifyTheme.change(theme);
    emitThemeChange(theme);
  });
});

onUnmounted(() => {
  if (cleanup) {
    cleanup();
  }
});

// 监听Vuetify主题变化
watch(() => vuetifyTheme.global.name.value, (newTheme) => {
  emitThemeChange(newTheme);
});
</script>

<template>
  <v-app>
    <VAlert />
    <VNavigation></VNavigation>
    <Notification />
    <v-main>
      <RouterView></RouterView>
      <div style="margin-top: 50px"></div>
      <!-- 备案组件 -->
      <!-- <VFill class="fill"></VFill> -->
    </v-main>
  </v-app>
</template>

<style>
.fill {
  text-align: center;
  display: inline-block;
}
</style>
