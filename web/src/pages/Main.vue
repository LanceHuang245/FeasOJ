<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { getAnnouncement, getNotification } from '../utils/api/users';
import { getDailyProblem } from '../utils/api/problems';
import { getAllCompetitions } from '../utils/api/competitions';
import { useTheme } from 'vuetify';
import { useI18n } from "vue-i18n";
import { difficultyColor, difficultyLang } from "../utils/dynamic_styles";
import { MdPreview } from "md-editor-v3";
import "md-editor-v3/lib/preview.css";
import { getMdPreviewTheme } from '../utils/theme';

const { t } = useI18n();
const announcement = ref('');
const notice = ref('');
const dailyProblem = ref(null);
const competitions = ref([]);
const vuetifyTheme = useTheme();
const previewTheme = ref(getMdPreviewTheme());

const getCompetitionStatus = (status) => {
  switch (status) {
        case 0:
            return 'message.compenotstarted';
        case 1:
            return 'message.compeprogress';
        case 2:
            return 'message.compeover';
        default:
            return 'message.compenotstarted';
    }
};

// 监听主题变化
const handleThemeChange = (event) => {
  const theme = event.detail.theme === 'dark' ? 'dark' : 'light';
  previewTheme.value = theme;
};

onMounted(async () => {
  announcement.value = await getAnnouncement();
  notice.value = await getNotification();
  try {
    const res = await getDailyProblem();
    if (res && res.data) {
      dailyProblem.value = res.data.data;
    }
  } catch (error) {
    console.error("Failed to fetch daily problem:", error);
  }

  try {
    const res = await getAllCompetitions();
    if (res && res.data) {
      competitions.value = res.data.data.slice(0, 3);
    }
  } catch (error) {
    console.error("Failed to fetch competitions:", error);
  }

  // 监听主题变化
  window.addEventListener('theme-change', handleThemeChange);
});

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
  <v-container fluid class="main-container pa-4 pa-md-8">
    <v-row>
      <v-col cols="12" class="d-flex justify-center text-center mb-6">
        <div class="d-flex align-center ga-4">
          <v-img src="/logo.png" width="80px" height="80px" class="logo-image flex-shrink-0"
            :class="{ 'dark-mode': vuetifyTheme.global.name.value === 'dark' }"></v-img>
          <h1 class="text-h3 font-weight-medium">FeasOJ</h1>
        </div>
      </v-col>
      <v-col cols="12" md="8">
        <v-row>
          <v-col cols="12">
            <v-card class="card-hover" elevation="2" rounded="lg">
              <v-card-title class="d-flex align-center">
                <v-icon start icon="mdi-bullhorn-variant-outline"></v-icon>
                {{ $t('message.announcement') }}
              </v-card-title>
              <v-divider></v-divider>
              <v-card-text>
                <MdPreview :editorId="'announcement-preview'" :modelValue="announcement" :theme="previewTheme" />
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12">
            <v-card class="card-hover" elevation="2" rounded="lg">
              <v-card-title class="d-flex align-center">
                <v-icon start icon="mdi-bell-ring-outline"></v-icon>
                {{ $t('message.notice') }}
              </v-card-title>
              <v-divider></v-divider>
              <v-card-text>
                <MdPreview :editorId="'notice-preview'" :modelValue="notice" :theme="previewTheme" />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="12" md="4">
        <v-row>
          <v-col cols="12">
            <v-card v-if="dailyProblem" class="card-hover" elevation="2" rounded="lg"
              :to="`/problemset/${dailyProblem.id}`">
              <v-card-title class="d-flex align-center">
                <v-icon start icon="mdi-lightbulb-on-outline"></v-icon>
                 {{ $t("message.problemOfTheDay") }}
              </v-card-title>
              <v-divider></v-divider>
              <v-card-text>
                <p class="font-weight-bold">{{ dailyProblem.title }}</p>
                <v-chip :style="difficultyColor(dailyProblem.difficulty)" size="small" class="mt-2">
                  {{ $t(difficultyLang(dailyProblem.difficulty)) }}
                </v-chip>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12">
            <v-card v-if="competitions.length > 0" class="card-hover" elevation="2" rounded="lg">
              <v-card-title class="d-flex align-center">
                <v-icon start icon="mdi-trophy-outline"></v-icon>
                {{ $t("message.competition") }}
              </v-card-title>
              <v-divider></v-divider>
              <v-list lines="two" class="bg-transparent">
                <template v-for="(comp, index) in competitions" :key="comp.id">
                  <v-list-item
                    :title="comp.title"
                    :subtitle="`${t('message.status')}: ${t(getCompetitionStatus(comp.status))}`"
                    :to="`/competitions`"
                  >
                  </v-list-item>
                  <v-divider v-if="index < competitions.length - 1" inset></v-divider>
                </template>
              </v-list>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.main-container {
  min-height: 100vh;
}

.logo-image {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.logo-image:hover {
  transform: scale(1.05);
}

.logo-image.dark-mode {
  filter: brightness(0) invert(1);
}

.v-card.card-hover {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.card-hover:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2) !important;
}
</style>
