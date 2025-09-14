<!-- 讨论帖子列表页 -->
<script setup>
import { useRouter } from 'vue-router';
import { ref, onMounted, computed, watch } from 'vue';
import { getAllDis } from '../../utils/api/discussions';
import { showAlert } from '../../utils/alert';
import { token } from '../../utils/account';
import { useI18n } from 'vue-i18n';
import moment from 'moment';

const { t } = useI18n();

const headers = ref([
  { title: t('message.discussion'), value: 'Title', align: 'center' },
  { title: t('message.discussionOwner'), value: 'Owner', align: 'center' },
  { title: t('message.when'), value: 'Description', align: 'center' }
])

const router = useRouter()

// 用作分页
const itemsPerPage = ref(12)
const page = ref(1)
const discuss = ref([])
const discussCount = ref(0)
const loading = ref(true)

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

// 从后端获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await getAllDis(page.value, itemsPerPage.value);
    discuss.value = response.data.data ? response.data.data : [];
    discussCount.value = response.data.total ? response.data.total : 0;
  } catch (error) {
    showAlert(t("message.failed") + "!", '/discussion')
  } finally {
    loading.value = false
  }
}

// 监听分页变化
watch(page, async () => {
  await fetchData()
})

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
    <p style="font-size: 24px;margin-left: 20px;">{{ t('message.discussion') }}</p>
  </v-app-bar>
  <v-card class="mx-auto my-8" width="85%" elevation="2" rounded="lg">
    <v-data-table-server :headers="headers" :items="discuss" :items-length="discussCount" :loading="loading"
      :loading-text="$t('message.loading')" @update="fetchData" :hide-default-footer="true"
      :items-per-page="itemsPerPage" v-model:page="page"
      :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')">
      <template v-slot:no-data>
        <div class="empty-state-container">
          <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-forum-outline</v-icon>
          <h3 class="text-h6 text-grey-darken-1 mb-2">
            {{ !userLoggedIn ? t('message.nologin') : t('message.nodata') }}
          </h3>
          <p v-if="userLoggedIn" class="text-body-2 text-grey-lighten-1 mb-4">
            {{ t('message.noDiscussionsAvailable') }}
          </p>
          <v-btn v-if="userLoggedIn" @click="$router.push('/discussion/create')" color="primary" variant="tonal" class="mt-2">
            {{ t('message.createDiscussion') }}
          </v-btn>
        </div>
      </template>
      <template v-slot:item="{ item }">
        <tr>
          <td>
            <v-btn @click="router.push({ path: `/discussion/${item.id}` })" variant="text" color="primary" class="disc-btn">{{ item.title
              }}</v-btn>
          </td>
          <td><v-chip class="font-weight-medium author-chip" variant="outlined" color="primary" @click="router.push({ path: `/profile/${item.username}` })"
            >{{ item.username }}</v-chip></td>
          <td class="text-body-2 text-medium-emphasis">{{ moment(item.created_at).format('YYYY-MM-DD HH:mm') }}</td>
        </tr>
      </template>
      <template v-slot:bottom>
        <v-pagination v-model="page" :length="Math.ceil(discussCount / itemsPerPage)" rounded="xl" />
      </template>
    </v-data-table-server>
  </v-card>
  <div class="fab">
    <v-fab fixed icon="mdi-plus" size="64" color="primary" elevation="8" v-if="userLoggedIn && !loading"
      @click="$router.push('/discussion/create')"></v-fab>
  </div>
</template>

<style scoped>
.fab {
  position: fixed;
  bottom: 80px;
  right: 180px;
}

.disc-btn {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  justify-content: flex-start;
  padding: 8px 12px;
  border-radius: 8px;
}

.disc-btn:hover {
  background: rgba(var(--v-theme-primary), 0.08);
  transform: translateX(4px);
}

.author-chip {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.author-chip:hover {
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