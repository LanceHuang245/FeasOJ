<!-- 个人信息页 -->
<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getUserSubmitRecords } from '../utils/api/submit_records.js';
import { uploadAvatar, updateSynopsis } from '../utils/api/users.js';
import { getResultStyle, getResultChipColor } from '../utils/dynamic_styles.js';
import { avatarServer } from '../utils/axios.js';
import { verifyUserInfo, getUserInfo } from '../utils/api/auth.js';
import { showAlert } from '../utils/alert.js';
import { userName, token } from '../utils/account.js';
import { useI18n } from 'vue-i18n';
import { MdPreview } from "md-editor-v3";
import { showNotification } from '../utils/notification.js';
import moment from 'moment';
import Heatmap from '../components/Profile/Heatmap.vue';
import 'md-editor-v3/lib/preview.css';
import { getMdPreviewTheme } from '../utils/theme.js';

const { t } = useI18n();

const userInfo = ref({});
const showCropper = ref(false);
const route = useRoute();
const router = useRouter();
const currentUsername = ref(route.params.Username);
const loading = ref(false);
const userSubmitRecords = ref([]);
const userSubmitRecordsLength = ref(0);
const networkloading = ref(false);
const dialog = ref(false);
const synopsis = ref('');
const id = 'preview-only';
const codeDialog = ref(false);
const previewTheme = ref(getMdPreviewTheme());

// 展示代码
const currentCode = ref('');
const headers = ref([
  { title: t('message.problemId'), value: 'Pid', align: 'center' },
  { title: t('message.result'), value: 'Result', align: 'center' },
  { title: t('message.lang'), value: 'Language', align: 'center' },
  { title: t('message.when'), value: 'Time', align: 'center' },
]);

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

// 监听主题变化
const handleThemeChange = (event) => {
  previewTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

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
    codeDialog.value = true;
  }
}

// 登出
const logout = () => {
  localStorage.clear();
  location.reload();
  location = '#/login';
};

// 上传头像至服务器
const uploadAvat = async (cropper) => {
  try {
    const canvas = cropper.getCroppedCanvas();
    const file = await new Promise((resolve) => canvas.toBlob(resolve));

    // 获取原始文件的名称和类型
    const imageData = cropper.getImageData();
    const originalFile = imageData.src ? imageData.src.split('/').pop() : 'avatar.png';
    const fileName = originalFile.split('?')[0];
    const fileType = file.type || 'image/png';

    // 创建一个新的文件对象，保留原始文件名和类型
    const newFile = new File([file], fileName, { type: fileType });
    networkloading.value = true;
    const resp = await uploadAvatar(newFile);
    networkloading.value = false;
    showAlert(resp.data.message, "reload");
  } catch (error) {
    networkloading.value = false;
    showAlert(error.response.data.message, "");
  }
};

// 获取用户提交信息
const fetchSubmitData = async () => {
  try {
    const submitResponse = await getUserSubmitRecords(currentUsername.value);
    userSubmitRecords.value = submitResponse.data.data;
    userSubmitRecordsLength.value = userSubmitRecords.value.length;
  } catch (error) {
    showAlert(t("message.failed") + "!", "");
  }
};

// 更新用户简介
const updateSyn = async () => {
  try {
    networkloading.value = true;
    const response = await updateSynopsis(synopsis.value);
    networkloading.value = false;
    showAlert(response.data.message, "reload");
  } catch (error) {
    networkloading.value = false;
    showAlert(error.response.data.message, "");
  }
};

// 检验获取用户信息
const verifyAndFetchUserInfo = async () => {
  loading.value = true;
  try {
    if (!userLoggedIn.value) {
      await router.push({ path: '/login' });
      return;
    }
    await verifyUserInfo(userName.value, token.value);
    const userInfoResponse = await getUserInfo(currentUsername.value);
    userInfo.value = userInfoResponse.data.data;
    synopsis.value = userInfo.value.synopsis;
    await fetchSubmitData();
  } catch (error) {
    await router.push({ path: '/403' });
  } finally {
    loading.value = false;
  }
};

// 监听路由参数变化
watch(() => route.params.Username, (newUsername) => {
  currentUsername.value = newUsername;
  verifyAndFetchUserInfo();
}, { immediate: true });

onMounted(() => {
  // 监听主题变化
  window.addEventListener('theme-change', handleThemeChange);
});

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
  <v-dialog v-model="networkloading" max-width="500px">
    <v-card rounded=xl>
      <div class="networkloading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
      </div>
    </v-card>
  </v-dialog>
  <div v-if="loading" class="loading">
    <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
  </div>
  <div v-else>
    <div style="margin: 10%"></div>
    <v-card class="mx-auto" max-width="60%" min-width="60%" rounded="xl" elevation="10">
      <div style="margin: 10px"></div>
      <div class="avatar-container">
        <v-avatar size="120" color="surface-variant">
          <v-img :src="avatarServer + userInfo.avatar" cover alt="Avatar"></v-img>
        </v-avatar>
        <v-btn v-if="currentUsername === userName" icon="mdi-pencil" size="30" @click="showCropper = true"
          class="edit-btn"></v-btn>
      </div>
      <v-card-text>
        <p class="text-h4 font-weight-black">{{ userInfo.username }}</p>
        <div style="margin: 10px"></div>
        <v-btn variant="text" class="text-medium-emphasis" @click="currentUsername === userName ? dialog = true : ''">{{
          !userInfo.synopsis ? $t('message.nosynopsis') : userInfo.synopsis }}</v-btn>
        <div style="margin: 20px"></div>
        <v-row style="justify-content: space-between;margin-inline: 5px">
          <v-text-field label="Email" :model-value="userInfo.email" readonly rounded="xl"
            variant="solo-filled"></v-text-field>
          <div style="margin: 5px;"></div>
          <v-text-field label="Score" :model-value="userInfo.score" readonly rounded="xl"
            variant="solo-filled"></v-text-field>
        </v-row>
      </v-card-text>
      <v-card-actions style="justify-content: end;">
        <v-btn v-if="currentUsername === userName" color="primary" variant="text" rounded="xl"
          style="margin-right: 10px;" @click="$router.push('/reset')">{{ $t('message.resetpwd') }}</v-btn>
        <v-btn v-if="currentUsername === userName" color="primary" variant="text" rounded="xl"
          style="margin-right: 10px;" @click="logout">{{
            $t('message.logout') }}</v-btn>
      </v-card-actions>
    </v-card>
    <v-card class="mx-auto" max-width="60%" min-width="60%" rounded="xl" elevation="10" style="margin-top: 30px;">
      <v-card-text>
        <p class="text-h4 font-weight-black">{{ $t('message.submissions') }}</p>
        <Heatmap :user-submit-records="userSubmitRecords" />
      </v-card-text>
    </v-card>
    <div style="margin: 30px"></div>
    <v-card class="mx-auto" max-width="60%" min-width="60%" rounded="xl" elevation="10">
      <v-card-text class="pa-0">
        <v-data-table-server :headers="headers" :items="userSubmitRecords" :items-length="userSubmitRecordsLength"
          :loading="loading" :loading-text="$t('message.loading')" @update="fetchSubmitData" :hide-default-footer="true"
          :no-data-text="$t('message.nodata')" class="profile-table" density="comfortable" hover>
          <template v-slot:item="{ item }">
            <tr class="profile-table-row">
              <td class="text-center pa-4">
                <v-btn @click="router.push({ path: `/problem/${item.Pid}` })" variant="text" color="primary"
                  class="font-weight-medium" size="small" :ripple="false">
                  {{ item.ProblemId }}
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
  </div>
  <avatar-cropper v-model="showCropper" :labels="{ submit: '上传头像', cancel: $t('message.cancel') }"
    :upload-handler="uploadAvat" />
  <v-dialog v-model="dialog" max-width="600px">
    <v-card rounded="xl">
      <div v-if="networkloading" class="loading" style="min-width: 600px;min-height: 500px;">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
      </div>
      <div v-else>
        <v-card-title>
          <span class="headline">{{ $t('message.edit') }}</span>
        </v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field :label="$t('message.synopsis')" v-model="synopsis" variant="solo-filled"></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="dialog = false" rounded="xl">{{ $t('message.cancel')
          }}</v-btn>
          <v-btn color="primary" @click="updateSyn" rounded="xl">{{ $t('message.save') }}</v-btn>
        </v-card-actions>
      </div>
    </v-card>
  </v-dialog>
  <!-- 查看代码弹窗 -->
  <v-dialog v-model="codeDialog" max-width="800px" persistent transition="dialog-bottom-transition">
    <v-card rounded="lg" elevation="24">
      <v-card-title class="d-flex align-center pa-6">
        <v-icon icon="mdi-code-braces" size="24" class="me-3" color="primary"></v-icon>
        <span class="text-h6 font-weight-medium">{{ t('message.codePreview') }}</span>
        <v-spacer></v-spacer>
        <v-btn icon="mdi-close" variant="text" size="small" @click="codeDialog = false"></v-btn>
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
.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.avatar-container {
  position: relative;
  display: inline-block;
}

.edit-btn {
  position: absolute;
  bottom: 0;
  right: 0;
}

.profile-table {
  border-radius: 0 0 16px 16px;
}

.profile-table-row {
  transition: background-color 0.2s ease;
}

.profile-table-row:hover {
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