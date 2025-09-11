<!-- 创建讨论页 -->
<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { addDiscussion } from '../../utils/api/discussions.js';
import { showAlert } from '../../utils/alert.js';
import { token } from "../../utils/account.js";
import { MdEditor } from 'md-editor-v3';
import { useI18n } from 'vue-i18n';
import { getMdEditorTheme } from '../../utils/theme.js';
import 'md-editor-v3/lib/style.css';

const { t } = useI18n();
const { locale } = useI18n();
const title = ref('');
const content = ref('');
const loading = ref(false);
const editorTheme = ref(getMdEditorTheme());

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

// 监听主题变化
const handleThemeChange = (event) => {
  editorTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

// 监听MdEditor主题变化
const handleMdEditorThemeChange = (event) => {
  editorTheme.value = event.detail.theme;
};

const submit = async () => {
    if (!title.value || !content.value) {
        showAlert(t("message.formCheckfailed") + "!", "");
        return;
    }
    loading.value = true;
    try {
        const response = await addDiscussion(title.value, content.value)
        showAlert(response.data.message, "/discussion");
    } catch (error) {
        showAlert(error.response.data.message, "");
    } finally{
        loading.value = false;
    }
};

onMounted(async () => {
    loading.value = true;
    if (!userLoggedIn.value) {
        window.location = "#/login";
    } else {
        loading.value = false;
    }
    
    // 监听主题变化
    window.addEventListener('theme-change', handleThemeChange);
    // 监听MdEditor主题变化
    window.addEventListener('md-editor-theme-change', handleMdEditorThemeChange);
});

onUnmounted(() => {
    // 清理事件监听器
    window.removeEventListener('theme-change', handleThemeChange);
    window.removeEventListener('md-editor-theme-change', handleMdEditorThemeChange);
});
</script>

<template>
    <div v-if="loading" class="loading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
    </div>
    <div v-else>
        <v-app-bar :elevation="2">
            <template v-slot:prepend>
                <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
            </template>
            <v-row style="align-items: center;">
                <div style="margin-left: 15px;"></div>
                <p class="font-weight-black">{{ $t('message.createDiscussion') }}</p>
            </v-row>
        </v-app-bar>
        <div style="margin-top: 60px"></div>
        <div class="form-align">
            <v-form style="min-width: 50%" @submit.prevent="submit">
                <v-text-field v-model="title" :label="$t('message.title')" rounded="xl"
                    variant="solo-filled"></v-text-field>
                <MdEditor v-model="content" :noUploadImg="true" :footers="[]"
                    :language="locale === 'zh_CN' ? 'zh-CN' : 'en-US'" :theme="editorTheme" />
                <div style="margin-top: 30px;"></div>
                <v-btn rounded="xl" type="submit" color="primary">{{ $t('message.submit') }}</v-btn>
            </v-form>
        </div>
    </div>
</template>

<style scoped>
.loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.form-align {
    display: flex;
    justify-content: center;
    align-items: center;
}
</style>