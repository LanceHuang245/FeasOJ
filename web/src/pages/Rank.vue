<script setup>
import { onMounted, computed, ref } from 'vue';
import { getRanking } from '../utils/api/users';
import { token } from '../utils/account';
import { showAlert } from '../utils/alert';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { avatarServer } from '../utils/axios';

const { t } = useI18n();

const router = useRouter();
const loading = ref(true);
const rankingData = ref([])
const rankingDataLength = ref(0)

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

const headers = ref([
  { title: t('message.rank'), value: 'rank', align: 'center', sortable: false, },
  { title: t('message.username'), value: 'username', align: 'center', sortable: false },
  { title: t('message.synopsis'), value: 'synopsis', align: 'center', sortable: false },
  { title: 'Score', value: 'score', align: 'center', sortable: false, },
])

// 从后端获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await getRanking()
    // 为每个用户添加排名
    rankingData.value = response.data.data.map((user, index) => ({
      ...user,
      rank: index + 1
    }));
    rankingDataLength.value = response.data.data.length;
  } catch (error) {
    showAlert(t("message.failed") + "!", "")
  } finally {
    loading.value = false
  }
}

// 获取排名徽章颜色
const getRankBadgeColor = (rank) => {
  if (rank === 1) return 'amber';
  if (rank === 2) return 'grey-lighten-1';
  if (rank === 3) return 'orange-darken-1';
  return 'primary';
}

// 获取排名徽章图标
const getRankIcon = (rank) => {
  if (rank === 1) return 'mdi-trophy';
  if (rank === 2) return 'mdi-medal';
  if (rank === 3) return 'mdi-medal-outline';
  return null;
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
  <v-app-bar :elevation="2">
    <v-col class="align-left-solo">
      <p style="font-size: 24px;">{{ $t('message.rank') }}</p>
      <p style="font-size: 12px;">{{ $t('message.refresh5m') }}</p>
    </v-col>
  </v-app-bar>

  <v-container class="pa-6">
    <v-card rounded="xl" elevation="4">
      <v-data-table-server :headers="headers" :items="rankingData" :items-length="rankingDataLength" :loading="loading"
        :loading-text="$t('message.loading')" @update="fetchData" :hide-default-footer="true"
        :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')" class="ranking-table"
        density="comfortable">
        <template v-slot:item="{ item, index }">
          <tr :class="{ 'top-three': item.rank <= 3 }">
            <!-- 排名列 -->
            <td class="text-center pa-4">
              <div class="d-flex align-center justify-center">
                <v-avatar :color="getRankBadgeColor(item.rank)" size="36" class="me-2" v-if="item.rank <= 3">
                  <v-icon :icon="getRankIcon(item.rank)" color="white" size="20"></v-icon>
                </v-avatar>
                <v-chip :color="getRankBadgeColor(item.rank)" variant="flat" size="small" class="font-weight-bold"
                  v-else>
                  #{{ item.rank }}
                </v-chip>
              </div>
            </td>

            <!-- 用户名列 -->
            <td class="pa-4">
              <div class="d-flex align-center">
                <v-avatar size="40" color="grey-lighten-1" class="me-3">
                  <v-img :src="avatarServer + item.avatar" cover alt="Avatar"></v-img>
                </v-avatar>
                <v-btn @click="router.push({ path: `/profile/${item.username}` })" variant="text"
                  class="text-body-1 font-weight-medium text-primary">
                  {{ item.username }}
                </v-btn>
              </div>
            </td>

            <!-- 简介列 -->
            <td class="pa-4">
              <span class="text-body-2 text-medium-emphasis">{{ item.synopsis || '-' }}</span>
            </td>

            <!-- 分数列 -->
            <td class="text-center pa-4">
              <v-chip :color="item.rank <= 3 ? 'success' : 'primary'" variant="flat" size="small"
                class="font-weight-bold">
                {{ item.score }}
              </v-chip>
            </td>
          </tr>
        </template>
      </v-data-table-server>
    </v-card>
  </v-container>
</template>

<style scoped>
.ranking-table {
  background: transparent;
}

.ranking-table :deep(.v-data-table__wrapper) {
  background: white;
  border-radius: 8px;
  margin: 0 16px 16px 16px;
}

.ranking-table :deep(.v-data-table-header) {
  background: #f8f9fa;
  border-radius: 8px 8px 0 0;
}

.ranking-table :deep(.v-data-table-header th) {
  font-weight: 600;
  color: #37474f;
  border-bottom: 2px solid #e0e0e0;
}

.ranking-table :deep(.v-data-table__tbody tr) {
  transition: all 0.2s ease;
}

.ranking-table :deep(.v-data-table__tbody tr:hover) {
  background-color: #f5f5f5;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ranking-table :deep(.v-data-table__tbody tr.top-three) {
  background: linear-gradient(90deg, rgba(255, 215, 0, 0.1) 0%, rgba(255, 255, 255, 1) 100%);
}

.ranking-table :deep(.v-data-table__tbody tr.top-three:hover) {
  background: linear-gradient(90deg, rgba(255, 215, 0, 0.15) 0%, rgba(255, 255, 255, 1) 100%);
}

.ranking-table :deep(.v-data-table__tbody tr:nth-child(1)) {
  background: linear-gradient(90deg, rgba(255, 215, 0, 0.2) 0%, rgba(255, 255, 255, 1) 100%);
}

.ranking-table :deep(.v-data-table__tbody tr:nth-child(2)) {
  background: linear-gradient(90deg, rgba(192, 192, 192, 0.15) 0%, rgba(255, 255, 255, 1) 100%);
}

.ranking-table :deep(.v-data-table__tbody tr:nth-child(3)) {
  background: linear-gradient(90deg, rgba(205, 127, 50, 0.15) 0%, rgba(255, 255, 255, 1) 100%);
}

.align-left-solo {
  text-align: left;
  margin-left: 10px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .ranking-table :deep(.v-data-table__tbody td) {
    padding: 8px 4px !important;
  }

  .ranking-table :deep(.v-data-table__tbody .v-avatar) {
    width: 32px !important;
    height: 32px !important;
  }
}
</style>