<!-- 题库页 -->
<script setup>
import { ref, onMounted, computed } from 'vue';
import { getAllProblems } from '../../utils/api/problems';
import { showAlert } from '../../utils/alert';
import { token } from "../../utils/account";
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { difficultyColor, difficultyLang } from '../../utils/dynamic_styles';

const router = useRouter();
const { t } = useI18n();

const headers = ref([
  { title: t('message.problemId'), value: 'Pid', align: 'center' },
  { title: t('message.problem'), value: 'Title', align: 'center' },
  { title: t('message.difficulty'), value: 'Difficulty', align: 'center' }
])
const problems = ref([])
const totalProblems = ref(0)
const loading = ref(true)
const searchQuery = ref('')
const showSearch = ref(false) // 控制搜索栏显示状态

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

// 搜索功能
const filteredProblems = computed(() => {
  return problems.value.filter((problem) =>
    problem.Title.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 切换搜索栏显示状态
const toggleSearch = () => {
  showSearch.value = !showSearch.value
  if (!showSearch.value) {
    searchQuery.value = '' // 关闭搜索时清空搜索内容
  }
}

// 从后端获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await getAllProblems();
    problems.value = response.data.problems
    totalProblems.value = problems.value.length
  } catch (error) {
    showAlert(t("message.failed") + "!", "")
  } finally {
    loading.value = false
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
})
</script>

<template>
  <v-container fluid class="pa-0">
    <v-app-bar :elevation="2">
      <p style="font-size: 24px;margin-left: 20px;">{{ t('message.problemset') }}</p>
      <v-spacer></v-spacer>
      <v-btn icon @click="toggleSearch" variant="text">
        <v-icon>mdi-magnify</v-icon>
      </v-btn>
    </v-app-bar>
    <!-- 搜索栏 -->
    <v-expand-transition>
      <v-row v-show="showSearch" class="mb-6">
        <v-col cols="12">
          <v-card elevation="0" rounded="sm">
            <v-card-text class="pa-0">
              <v-text-field 
                v-model="searchQuery" 
                :placeholder="$t('message.searchProblem')" 
                variant="solo-filled"
                density="comfortable" 
                prepend-inner-icon="mdi-magnify" 
                rounded="sm"
                hide-details 
                @keyup.enter="showSearch = false"
              ></v-text-field>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-expand-transition>

    <!-- 数据表格 -->
    <v-row>
      <v-col cols="12">
        <v-card class="problem-table-card" elevation="0" rounded="lg">
          <v-data-table :headers="headers" :items="filteredProblems" :items-length="totalProblems" :loading="loading"
            :loading-text="$t('message.loading')"
            :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')" class="problem-table"
            density="comfortable" hover rounded="lg" elevation="0">
            <!-- 自定义空状态模板 -->
            <template v-slot:no-data>
              <div class="empty-state-container">
                <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-file-document-outline</v-icon>
                <h3 class="text-h6 text-grey-darken-1 mb-2">
                  {{ !userLoggedIn ? t('message.nologin') : (searchQuery ? t('message.noSearchResults') : t('message.nodata')) }}
                </h3>
                <p v-if="userLoggedIn" class="text-body-2 text-grey-lighten-1 mb-4">
                  {{ searchQuery ? t('message.tryDifferentKeywords') : t('message.noProblemsAvailable') }}
                </p>
              </div>
            </template>
            <template v-slot:item.Pid="{ item }">
              <v-chip variant="tonal" size="small" color="primary">
                {{ item.Pid }}
              </v-chip>
            </template>
            <template v-slot:item.Title="{ item }">
              <v-btn @click="router.push({ path: `/problemset/${item.Pid}` })" variant="text" color="primary">
                {{ item.Title }}
              </v-btn>
            </template>
            <template v-slot:item.Difficulty="{ item }">
              <v-chip :style="difficultyColor(item.Difficulty)">
                {{ t(difficultyLang(item.Difficulty)) }}
              </v-chip>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.problem-table-card {
  background: rgba(var(--v-theme-surface), 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(var(--v-theme-outline), 0.08);
}

.problem-table {
  --v-data-table-header-background: rgba(var(--v-theme-surface), 0.95);
  --v-data-table-row-hover-background: rgba(var(--v-theme-primary), 0.04);
}

.problem-table :deep(.v-btn) {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  justify-content: flex-start;
  padding: 8px 12px;
  border-radius: 8px;
}

.problem-table :deep(.v-btn:hover) {
  background: rgba(var(--v-theme-primary), 0.08);
  transform: translateX(4px);
}

.empty-state-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  min-height: 300px;
  background: rgba(var(--v-theme-surface), 0.5);
  border-radius: 12px;
  margin: 20px;
}

.empty-state-container .v-icon {
  opacity: 0.6;
  transition: all 0.3s ease;
}

.empty-state-container:hover .v-icon {
  opacity: 0.8;
  transform: scale(1.05);
}

.empty-state-container h3 {
  font-weight: 500;
  line-height: 1.4;
}

.empty-state-container p {
  max-width: 400px;
  line-height: 1.6;
}
</style>