<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { verifyUserInfo } from "../utils/api/auth";
import { token, userName } from "../utils/account";
import { showAlert } from "../utils/alert";
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();

const privilege = ref("");
// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

// 路由实例
const router = useRouter();
const route = useRoute();

// 计算当前路由路径，用于高亮
const currentPath = computed(() => route.path);

// 判断是否为当前页面
const isCurrentPage = (path) => {
  if (path === '/') {
    return currentPath.value === '/';
  }
  return currentPath.value.startsWith(path);
};

// 根据用户登录状态修改导航目的
const navigate = async () => {
  if (userLoggedIn.value) {
    await router.push("/profile/" + userName.value);
  } else {
    await router.push("/login");
  }
};

// 校验Token是否过期
const isTokenExpired = () => {
  const expirationTime = localStorage.getItem('tokenExpiration');
  if (!expirationTime) {
    return true;
  }
  return Date.now() > expirationTime;
}

// 校验Token后获取用户信息，若privilege为1则表明是管理员
onMounted(async () => {
  if (userLoggedIn.value) {
    if (isTokenExpired()) {
      showAlert(t("message.tokenExpired") + "!", "#/login");
      localStorage.clear();
      userLoggedIn.value = false;
      return;
    }
    try {
      const resp = await verifyUserInfo(userName.value, token.value);
      privilege.value = resp.data.info.role;
    } catch (error) {
      showAlert(t("message.tokenCheckfailed") + "!", "reload");
      localStorage.clear();
      userLoggedIn.value = false;
    }
  }
});
</script>

<template>
  <v-navigation-drawer :width="200" permanent location="left" :elevation="2" class="modern-nav-drawer">
    <v-list nav class="nav-list">
      <!-- 主导航项目 -->
      <div class="nav-section">
        <v-list-item rounded="lg" prepend-icon="mdi-home-outline" value="HOME" @click="router.push('/')" color="primary"
          class="nav-item" :class="{ 'nav-item-active': isCurrentPage('/') }" :active="isCurrentPage('/')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.mainpage') }}</span>
          </template>
        </v-list-item>

        <v-list-item rounded="lg" prepend-icon="mdi-book-open-page-variant-outline" value="PROBLEMSET"
          @click="router.push('/problemset')" color="primary" class="nav-item"
          :class="{ 'nav-item-active': isCurrentPage('/problemset') }" :active="isCurrentPage('/problemset')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.problemset') }}</span>
          </template>
        </v-list-item>

        <v-list-item rounded="lg" prepend-icon="mdi-trophy-outline" value="COMPETITION" color="primary"
          @click="router.push('/competitions')" class="nav-item"
          :class="{ 'nav-item-active': isCurrentPage('/competitions') }" :active="isCurrentPage('/competitions')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.competition') }}</span>
          </template>
        </v-list-item>

        <v-list-item rounded="lg" prepend-icon="mdi-forum-outline" value="DISCUSS" color="primary"
          @click="router.push('/discussion')" class="nav-item"
          :class="{ 'nav-item-active': isCurrentPage('/discussion') }" :active="isCurrentPage('/discussion')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.discussion') }}</span>
          </template>
        </v-list-item>

        <v-list-item rounded="lg" prepend-icon="mdi-podium-gold" value="RANK" @click="router.push('/rank')"
          color="primary" class="nav-item" :class="{ 'nav-item-active': isCurrentPage('/rank') }"
          :active="isCurrentPage('/rank')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.rank') }}</span>
          </template>
        </v-list-item>

        <v-list-item rounded="lg" prepend-icon="mdi-clipboard-list-outline" value="STATUS"
          @click="router.push('/status')" color="primary" class="nav-item"
          :class="{ 'nav-item-active': isCurrentPage('/status') }" :active="isCurrentPage('/status')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.status') }}</span>
          </template>
        </v-list-item>
      </div>

      <!-- 分隔线 -->
      <v-divider class="nav-divider"></v-divider>

      <!-- 底部导航项目 -->
      <div class="nav-section bottom-section">
        <!-- 管理员入口 -->
        <v-list-item v-if="privilege === 1" rounded="lg" prepend-icon="mdi-shield-crown-outline"
          @click="router.push('/admin')" value="admin" base-color="primary" class="nav-item admin-item"
          :class="{ 'nav-item-active': isCurrentPage('/admin') }" :active="isCurrentPage('/admin')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.management') }}</span>
          </template>
        </v-list-item>

        <!-- 设置 -->
        <v-list-item rounded="lg" prepend-icon="mdi-cog-outline" value="SETTINGS" base-color="primary" class="nav-item"
          @click="router.push('/settings')" :class="{ 'nav-item-active': isCurrentPage('/settings') }"
          :active="isCurrentPage('/settings')">
          <template v-slot:title>
            <span class="nav-title">{{ $t('message.settings') }}</span>
          </template>
        </v-list-item>

        <!-- 用户资料/登录 -->
        <v-list-item rounded="lg" :prepend-icon="userLoggedIn ? 'mdi-account-circle-outline' : 'mdi-login'"
          @click="navigate" value="PROFILE" base-color="primary" class="nav-item"
          :class="{ 'nav-item-active': isCurrentPage('/profile') }" :active="isCurrentPage('/profile')">
          <template v-slot:title>
            <span class="nav-title">{{ userLoggedIn ? userName : $t('message.login') }}</span>
          </template>
        </v-list-item>
      </div>
    </v-list>
  </v-navigation-drawer>
</template>

<style scoped>
.modern-nav-drawer {
  border-right: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.nav-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px 8px;
  gap: 1px;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.bottom-section {
  margin-top: auto;
  padding-top: 12px;
}

.nav-item {
  height: 44px;
  margin: 1px 0;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 10px !important;
  position: relative;
}

.nav-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
  transform: translateX(2px);
}

.nav-item:active {
  transform: scale(0.98);
}

.nav-item-active {
  background-color: rgba(var(--v-theme-primary), 0.12) !important;
  color: rgb(var(--v-theme-primary)) !important;
  font-weight: 600;
}

.nav-item-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 24px;
  background-color: rgb(var(--v-theme-primary));
  border-radius: 0 2px 2px 0;
}

.nav-item-active:hover {
  background-color: rgba(var(--v-theme-primary), 0.16) !important;
  transform: translateX(2px);
}

.nav-title {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-divider {
  margin: 16px 0;
  opacity: 0.12;
}

.admin-item {
  background-color: rgba(var(--v-theme-warning), 0.08);
}

.admin-item:hover {
  background-color: rgba(var(--v-theme-warning), 0.12);
}

.admin-item.nav-item-active {
  background-color: rgba(var(--v-theme-warning), 0.16) !important;
  color: rgb(var(--v-theme-warning)) !important;
}

.admin-item.nav-item-active::before {
  background-color: rgb(var(--v-theme-warning));
}

.profile-item {
  background-color: rgba(var(--v-theme-primary), 0.04);
}

.profile-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.12);
}

.lang-item {
  padding: 12px 16px;
  transition: background-color 0.2s ease;
}

.lang-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.lang-title {
  font-size: 14px;
  font-weight: 500;
  color: rgba(var(--v-theme-on-surface), 0.87);
}
</style>
