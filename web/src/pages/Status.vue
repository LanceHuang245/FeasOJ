<!-- 状态页 -->
<script setup>
import { getSubmitRecords } from '../utils/api/submit_records.js';
import { getResultStyle, getResultChipColor } from '../utils/dynamic_styles.js';
import { onMounted, computed, ref, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { showAlert } from '../utils/alert.js';
import { token } from "../utils/account.js";
import { useI18n } from 'vue-i18n';
import { MdPreview } from "md-editor-v3";
import { showNotification } from '../utils/notification.js';
import { getMdPreviewTheme } from '../utils/theme.js';
import 'md-editor-v3/lib/preview.css';
import moment from 'moment';

const { t } = useI18n();

const router = useRouter();

const headers = ref([
  { title: t('message.problemId'), value: 'Pid', align: 'center', sortable: false },
  { title: t('message.username'), value: 'Username', align: 'center', sortable: false },
  { title: t('message.result'), value: 'Result', align: 'center', sortable: false },
  { title: t('message.lang'), value: 'Language', align: 'center', sortable: false },
  { title: t('message.when'), value: 'Time', align: 'center', sortable: false },
])
const submitrecords = ref([])
const submitRecordsLength = ref(0)
const loading = ref(true)

const id = 'preview-only';
const dialog = ref(false);
const previewTheme = ref(getMdPreviewTheme());

// 展示代码
const currentCode = ref('');

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

// 监听主题变化
const handleThemeChange = (event) => {
  previewTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

// 从后端获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await getSubmitRecords()
    submitrecords.value = response.data.submitrecords
    submitRecordsLength.value = submitrecords.value.length
  } catch (error) {
    showAlert(t("message.failed") + "!", "")
  } finally {
    loading.value = false
  }
}

// 格式化代码
const formatAsFencedCode = (code, lang = '') => {
  return `\`\`\`${lang}
${code}
\`\`\``;
}

// 弹出对话框
const showCode = (code, lang) => {
  if (code === '') {
    showNotification(t('message.cannotViewCode'))
  } else {
    currentCode.value = formatAsFencedCode(code, lang) || '';
    dialog.value = true;
  }
}

// 初始化数据
onMounted(async () => {
  if (!userLoggedIn.value) {
    loading.value = false;
    setTimeout(() => {
      window.location = '#/login'
      window.location.reload()
    }, 2000);
    return;
  }
  await fetchData()

  // 监听主题变化
  window.addEventListener('theme-change', handleThemeChange);
})

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
  <v-app-bar :elevation="2">
    <v-col class="align-left-solo">
      <p style="font-size: 24px;">{{ $t('message.status') }}</p>
      <p style="font-size: 12px;">{{ $t('message.realtimeUpdate') }}</p>
    </v-col>
  </v-app-bar>

  <v-container fluid class="pa-6">
    <v-row justify="center">
      <v-col cols="12" lg="11" xl="10">
        <!-- 状态卡片 -->
        <v-card class="status-card" rounded="xl" elevation="2">
          <!-- 数据表格 -->
          <v-card-text class="pa-0">
            <v-data-table-server :headers="headers" :items="submitrecords" :items-length="submitRecordsLength"
              :loading="loading" :loading-text="$t('message.loading')" @update="fetchData" :hide-default-footer="true"
              :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')" class="status-table"
              density="comfortable" hover>
              <template v-slot:item="{ item }">
                <tr class="status-table-row">
                  <td class="text-center pa-4">
                    <v-btn @click="router.push({ path: `/problem/${item.Pid}` })" variant="text" color="primary"
                      class="font-weight-medium" size="small" :ripple="false">
                      {{ item.Pid }}
                    </v-btn>
                  </td>
                  <td class="text-center pa-4">
                    <v-btn @click="router.push({ path: `/profile/${item.Username}` })" variant="text" color="primary"
                      class="font-weight-medium" size="small" :ripple="false">
                      {{ item.Username }}
                    </v-btn>
                  </td>
                  <td v-if="item.Result === 'Running...'" class="text-center pa-4">
                    <v-progress-circular indeterminate color="primary" size="24" width="2"></v-progress-circular>
                  </td>
                  <td v-else :style="getResultStyle(item.Result)" @click="showCode(item.Code, item.Language)"
                    class="text-center pa-4 result-cell">
                    <v-chip :color="getResultChipColor(item.Result)" variant="tonal" size="small"
                      class="font-weight-medium">
                      {{ item.Result }}
                    </v-chip>
                  </td>
                  <td class="text-center pa-4">
                    <v-chip color="secondary" variant="outlined" size="small" class="font-weight-medium">
                      {{ item.Language }}
                    </v-chip>
                  </td>
                  <td class="text-center pa-4">
                    <span class="text-body-2 text-medium-emphasis">
                      {{ moment(item.Time).format('YYYY-MM-DD HH:mm') }}
                    </span>
                  </td>
                </tr>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <!-- 查看代码弹窗 -->
  <v-dialog v-model="dialog" max-width="800px" persistent transition="dialog-bottom-transition">
    <v-card rounded="lg" elevation="24">
      <v-card-title class="d-flex align-center pa-6">
        <v-icon icon="mdi-code-braces" size="24" class="me-3" color="primary"></v-icon>
        <span class="text-h6 font-weight-medium">{{ t('message.codePreview') }}</span>
        <v-spacer></v-spacer>
        <v-btn icon="mdi-close" variant="text" size="small" @click="dialog = false"></v-btn>
      </v-card-title>

      <v-card-text class="pa-0">
        <v-divider></v-divider>
        <div class="code-preview-container">
          <MdPreview v-if="currentCode" :id="id" :modelValue="currentCode" :theme="previewTheme" class="code-preview" />
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.status-card {
  border-radius: 16px !important;
}

.status-table {
  border-radius: 0 0 16px 16px;
}

.status-table-row {
  transition: background-color 0.2s ease;
}

.status-table-row:hover {
  background-color: rgba(var(--v-theme-primary), 0.04) !important;
}

.result-cell {
  cursor: pointer;
  transition: all 0.2s ease;
}

.code-preview-container {
  max-height: 70vh;
  overflow-y: auto;
  padding: 16px;
}

.code-preview {
  border-radius: 8px;
}

.align-left-solo {
  text-align: left;
  margin-left: 10px;
}

/* 自定义滚动条 */
.code-preview-container::-webkit-scrollbar {
  width: 8px;
}

.code-preview-container::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.code-preview-container::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.code-preview-container::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}
</style>