<!-- 题目详细页 -->
<script setup>
import { ref, onMounted, computed, watch, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import { getPbDetails, uploadCode } from "../../utils/api/problems";
import { VAceEditor } from "vue3-ace-editor";
import { showAlert } from "../../utils/alert";
import { token } from "../../utils/account";
import { MdPreview } from "md-editor-v3";
import { useI18n } from "vue-i18n";
import { difficultyColor, difficultyLang } from "../../utils/dynamic_styles";
import { getMdPreviewTheme } from "../../utils/theme";
import {
  DEFAULT_LANGUAGE,
  getFileExtension,
  getAceLang,
  getLanguageOptions
} from "../../utils/language_constants";
import "md-editor-v3/lib/preview.css";

const { t } = useI18n();

const networkloading = ref(false);
const id = "preview-only";
const route = useRoute();
const loading = ref(true);
const problemInfo = ref({});
const content = ref("");
const lang = ref(DEFAULT_LANGUAGE);
const previewTheme = ref(getMdPreviewTheme());

// Ace Editor字体
const fontSize = ref(14);
const minFontSize = 10;
const maxFontSize = 36;
const stepFontSize = 1;

// Ace Editor主题
const theme = ref("chrome");

// 生成组件的列表
const languageOptions = computed(() => getLanguageOptions());

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

// 监听主题变化
const handleThemeChange = (event) => {
  previewTheme.value = event.detail.theme === "dark" ? "dark" : "light";
};

// 代码模板
const templates = {
  Java: `
// Please don't edit the 'Main' class name. 请不要编辑 'Main' 类名。
public class Main {
    public static void main(String[] args) {
        
    }
}`,
  "C++": `#include <iostream>
            
using namespace std;

int main() {
    
    return 0;
}`,
  Python: ``,
  Rust: `

fn main() {
    
}`,
  Php: `<?php

?>`,
  Pascal: `{ O2 Enabled / O2已启用 }
begin

end.`,
  Golang: `
package main

import "fmt"

func main() {

}`,
};

// 代码上传
const uploadContentAsFile = async () => {
  if (!content.value) {
    showAlert(t("message.formCheckfailed"), "");
    return;
  }
  const blob = new Blob([content.value], { type: "text/plain" });
  const codefile = new File([blob], `main.${getFileExtension(lang.value)}`, {
    type: "text/plain",
  });
  try {
    networkloading.value = true;
    const resp = await uploadCode(codefile, route.params.Pid);
    networkloading.value = false;
    showAlert(resp.data.message, "reload");
  } catch (error) {
    networkloading.value = false;
    showAlert(error.response.data.message, "");
  }
};

// 应用用户偏好语言
const applyPreferredLanguage = () => {
  const savedPreferredLanguage = localStorage.getItem('preferredLanguage') || DEFAULT_LANGUAGE;
  lang.value = savedPreferredLanguage;
  content.value = templates[savedPreferredLanguage];
};

onMounted(async () => {
  loading.value = true;
  if (userLoggedIn.value) {
    try {
      const problemId = route.params.Pid;
      const resp = await getPbDetails(problemId);
      problemInfo.value = resp.data.data;
    } catch (error) {
      showAlert(error.response.data.message, "");
    } finally {
      loading.value = false;
    }
  } else {
    window.location = "#/login";
  }

  // 应用用户偏好语言
  applyPreferredLanguage();

  watch(
    lang,
    (newLang) => {
      content.value = templates[newLang];
    },
    { immediate: true }
  );

  // 监听主题变化
  window.addEventListener("theme-change", handleThemeChange);
});

onUnmounted(() => {
  // 清理事件监听器
  window.removeEventListener("theme-change", handleThemeChange);
});
</script>

<template>
  <v-dialog v-model="networkloading" max-width="600px">
    <v-card rounded="xl">
      <div class="networkloading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
      </div>
    </v-card>
  </v-dialog>
  <div v-if="loading" class="loading">
    <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
  </div>
  <div v-else>
    <v-app-bar :elevation="2">
      <template v-slot:prepend>
        <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
      </template>
      <v-col class="align-left">
        <v-row style="align-items: center">
          <p style="font-size: 24px">{{ problemInfo.title }}</p>
          <div style="margin-left: 10px"></div>
          <v-chip :style="difficultyColor(problemInfo.difficulty)">
            {{ $t(difficultyLang(problemInfo.difficulty)) }}
          </v-chip>
        </v-row>
      </v-col>
    </v-app-bar>

    <v-container fluid class="problem-container">
      <v-row>
        <!-- 题目描述部分 -->
        <v-col cols="12" lg="6" class="problem-description-col">
          <v-card class="problem-card" elevation="2" rounded="lg">
            <v-card-text class="problem-content">
              <!-- 题目限制信息 -->
              <v-row class="mb-6">
                <v-col cols="6">
                  <v-card variant="outlined" class="limit-card" rounded="xl">
                    <v-card-text class="text-center pa-4">
                      <v-icon icon="mdi-clock-outline" size="24" color="primary" class="mb-2"></v-icon>
                      <div class="text-h6 font-weight-medium">
                        {{ problemInfo.time_limit }}s
                      </div>
                      <div class="text-caption text-medium-emphasis">
                        {{ $t("message.timeLimit") }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
                <v-col cols="6">
                  <v-card variant="outlined" class="limit-card" rounded="xl">
                    <v-card-text class="text-center pa-4">
                      <v-icon icon="mdi-memory" size="24" color="primary" class="mb-2"></v-icon>
                      <div class="text-h6 font-weight-medium">
                        {{ problemInfo.memory_limit }}MB
                      </div>
                      <div class="text-caption text-medium-emphasis">
                        {{ $t("message.memoryLimit") }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
              <!-- 题目内容 -->
              <div class="problem-description">
                <md-preview :modelValue="problemInfo.content" :editorId="id" class="md_preview" :theme="previewTheme" />
              </div>
              <div style="margin-bottom: 20px;"></div>
              <!-- 输入输出示例 -->
              <div class="example-section">
                <div class="example-item">
                  <div class="example-label">
                    <v-icon icon="mdi-arrow-down" size="20" class="me-1"></v-icon>
                    {{ $t("message.displayInputCase") }}
                  </div>
                  <v-card variant="outlined" class="example-card" rounded="md">
                    <v-card-text class="example-text">
                      <pre>{{ problemInfo.input }}</pre>
                    </v-card-text>
                  </v-card>
                </div>

                <div class="example-item mt-4">
                  <div class="example-label">
                    <v-icon icon="mdi-arrow-up" size="20" class="me-1"></v-icon>
                    {{ $t("message.displayOutputCase") }}
                  </div>
                  <v-card variant="outlined" class="example-card" rounded="md">
                    <v-card-text class="example-text">
                      <pre>{{ problemInfo.output }}</pre>
                    </v-card-text>
                  </v-card>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 代码编辑器部分 -->
        <v-col cols="12" lg="6" class="code-editor-col">
          <v-card class="editor-card" elevation="2" rounded="lg">
            <v-card-text class="editor-content">
              <v-alert :text="$t('message.a_code')" :title="$t('message.attention')" type="info" variant="tonal"
                border="start" closable align="left" class="mb-4"></v-alert>

              <!-- 编辑器控制面板 -->
              <v-card variant="outlined" class="control-panel" rounded="md">
                <v-card-text class="pa-4">
                  <v-row>
                    <v-col cols="12" sm="6">
                      <v-select :label="$t('message.lang')" v-model="lang" :items="languageOptions" variant="outlined"
                        density="compact" hide-details></v-select>
                    </v-col>
                    <v-col cols="12" sm="6">
                      <v-select :label="$t('message.editor_theme')" v-model="theme" :items="[
                        'ambiance',
                        'clouds',
                        'cobalt',
                        'chaos',
                        'crimson_editor',
                        'dawn',
                        'dracula',
                        'dreamweaver',
                        'chrome',
                        'github',
                        'terminal',
                        'monokai',
                        'mono_industrial',
                        'pastel_on_dark',
                        'sqlserver',
                        'solarized_light',
                        'solarized_dark',
                      ]" variant="outlined" density="compact" hide-details></v-select>
                    </v-col>
                  </v-row>

                  <v-row class="mt-3">
                    <v-col cols="12">
                      <div class="d-flex align-center">
                        <v-icon icon="mdi-format-size" class="me-2"></v-icon>
                        <v-slider v-model="fontSize" :min="minFontSize" :max="maxFontSize" :step="stepFontSize"
                          thumb-label hide-details class="flex-grow-1"></v-slider>
                        <span class="text-caption ml-2">{{ fontSize }}px</span>
                      </div>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>

              <!-- 代码编辑器 -->
              <div class="editor-wrapper mt-4">
                <v-ace-editor v-model:value="content" :theme="theme" :lang="getAceLang(lang)"
                  :style="`height: 650px; font-size: ${fontSize}px; border-radius: 8px;`" />
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <!-- 提交按钮 -->
      <div class="submit-section mt-4">
        <v-btn color="primary" size="large" rounded="xl" @click="uploadContentAsFile" :loading="networkloading"
          class="submit-btn">
          {{ $t("message.submit") }}
        </v-btn>
      </div>
    </v-container>
  </div>
</template>

<style scoped>
.problem-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.problem-description-col,
.code-editor-col {
  padding: 12px;
}

.problem-card,
.editor-card {
  height: calc(100vh - 140px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.problem-content,
.editor-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.limit-card {
  transition: all 0.3s ease;
  border: 1px solid var(--v-border-color);
}

.limit-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.problem-description {
  margin-top: 24px;
}

.md_preview {
  text-align: left;
  line-height: 1.6;
}

.control-panel {
  background: var(--v-surface-variant);
  border: 1px solid var(--v-border-color);
}

.editor-wrapper {
  border: 1px solid var(--v-border-color);
  border-radius: 8px;
  overflow: hidden;
}

.submit-section {
  display: flex;
  justify-content: center;
}

.submit-btn {
  min-width: 120px;
  font-weight: 600;
  text-transform: none;
  letter-spacing: 0.5px;
}

.example-section {
  padding: 8px 0;
}

.example-item {
  margin-bottom: 16px;
}

.example-label {
  display: flex;
  align-items: center;
  font-weight: 600;
  color: var(--v-primary-base);
  margin-bottom: 8px;
  font-size: 14px;
}

.example-card {
  background: var(--v-surface-variant);
  border: 1px solid var(--v-border-color);
}

.example-text {
  padding: 12px;
  font-size: 16px;
  text-align: left;
  background: var(--v-surface);
  border-radius: 4px;
}

.example-text pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
</style>
