<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { getAnnouncement, getNotification } from '../utils/api/users';
import { useTheme } from 'vuetify';

const panel = ref([0, 1]);
const announcement = ref('');
const notice = ref('');
const vuetifyTheme = useTheme();

// 监听主题变化
const handleThemeChange = (event) => {
  // 主题变化时会触发重新渲染
};

onMounted(async () => {
  announcement.value = await getAnnouncement();
  notice.value = await getNotification();
  
  // 监听主题变化
  window.addEventListener('theme-change', handleThemeChange);
});

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
  <v-container fluid class="main-container">
    <div class="title-section">
      <v-img src="logo.png" width="150px" height="150px" class="logo-image" :class="{ 'logo-hover': true, 'dark-mode': vuetifyTheme.global.name.value === 'dark' }"></v-img>
      <h1 class="main-title">FeasOJ</h1>
    </div>

    <v-card rounded="xl" class="content-card" elevation="8" :class="{ 'card-hover': true }">
      <v-expansion-panels v-model="panel" variant="accordion">
        <v-expansion-panel :title="$t('message.announcement')" :text="announcement" rounded="xl">
        </v-expansion-panel>

        <v-expansion-panel :title="$t('message.notice')" :text="notice" rounded="xl">
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card>
  </v-container>
</template>

<style scoped>
.main-container {
  min-height: 100vh;
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 48px;
}

.title-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 16px;
}

.logo-image {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.logo-hover:hover {
  transform: scale(1.05);
  filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.15));
}

/* 深色模式下的logo样式 */
.logo-image.dark-mode {
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1)) brightness(0) invert(1);
}

.logo-hover.dark-mode:hover {
  transform: scale(1.05);
  filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.15)) brightness(0) invert(1);
}

.main-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0;
}

.content-card {
  width: 100%;
  max-width: 800px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(var(--v-theme-outline), 0.12);
}

.card-hover:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}
</style>
