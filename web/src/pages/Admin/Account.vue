<script setup>
import { onMounted, ref, computed } from 'vue';
import { getAllUsersInfo, promoteUser, demoteUser, banUser, unbanUser } from '../../utils/api/admin';
import { verifyUserInfo } from '../../utils/api/auth';
import { showAlert } from '../../utils/alert';
import { token, userName } from '../../utils/account'
import { useI18n } from 'vue-i18n';
import moment from 'moment';

const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)
const loading = ref(true)
const userPrivilege = ref("")
const users = ref([])
const totalUsers = ref(0)
const searchQuery = ref('')
const networkloading = ref(false);
const currentPage = ref(1);
const itemsPerPage = ref(50);

const getMenus = (item) => {
    let filteredMenus = []
    if (!item.isBan) {
        filteredMenus.push({ title: t("message.banUser"), icon: 'mdi-account-off' })
    } else {
        filteredMenus.push({ title: t("message.unbanUser"), icon: 'mdi-account-check' })
    }
    if (item.role === 1) {
        filteredMenus.push({ title: t("message.demoteUser"), icon: 'mdi-account-minus' })
    } else {
        filteredMenus.push({ title: t("message.promoteUser"), icon: 'mdi-account-plus' })
    }
    return filteredMenus
}

// 简化公共处理逻辑
const handleAction = async (action, item) => {
    networkloading.value = true;
    try {
        const resp = await action(item.uid)
        networkloading.value = false;
        showAlert(resp.data.message, "reload")
    } catch (error) {
        networkloading.value = false;
        showAlert(error.response.data.message, "")
    }
}

const handleMenuClick = (menu, item) => {
    const actions = {
        [t("message.banUser")]: banUser,
        [t("message.unbanUser")]: unbanUser,
        [t("message.demoteUser")]: demoteUser,
        [t("message.promoteUser")]: promoteUser
    };

    const action = actions[menu.title];
    if (action) {
        handleAction(action, item);
    }
}

const headers = ref([
    { title: t("message.username"), value: 'username', align: 'center', sortable: false },
    { title: t("message.email"), value: 'email', align: 'center', sortable: false },
    { title: t("message.role"), value: 'role', align: 'center', sortable: false },
    { title: t("message.timeOfRegister"), value: 'registerTime', align: 'center', sortable: false },
    { title: t("message.status"), value: 'status', align: 'center', sortable: false },
    { title: t("message.operation"), value: 'actions', align: 'center', sortable: false },
])

// 通过用户名搜索
const filteredUsers = computed(() => {
    if (!searchQuery.value) {
        return users.value;
    }
    return users.value.filter((user) =>
        user.username.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
})

// 分页后的数据
const paginatedUsers = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value;
    const end = start + itemsPerPage.value;
    return filteredUsers.value.slice(start, end);
});

// 总页数
const totalPages = computed(() => {
    return Math.ceil(filteredUsers.value.length / itemsPerPage.value);
});

// 页面变化处理
const handlePageChange = (page) => {
    currentPage.value = page;
};

// 获取数据
const fetchData = async () => {
    loading.value = true;
    try {
        const usersInfoResp = await getAllUsersInfo()
        users.value = usersInfoResp.data.data;
        totalUsers.value = usersInfoResp.data.data.length;
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    } finally {
        loading.value = false;
    }
}

onMounted(async () => {
    loading.value = true;
    try {
        if (!userLoggedIn.value) {
            window.location = "#/login";
            window.location.reload();
            return;
        }
        const userInfoResponse = await verifyUserInfo(userName.value, token.value);
        userPrivilege.value = userInfoResponse.data.data.role;
        if (userPrivilege.value !== 1) {
            window.location = '#/403';
            return;
        }
        await fetchData();
        loading.value = false;
    } catch (error) {
        window.location = '#/403';
    }
});
</script>

<template>
    <v-app-bar :elevation="2">
        <template v-slot:prepend>
            <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
        </template>
        <v-col class="align-left">
            <p style="font-size: 24px;">{{ t('message.usermanagement') }}</p>
        </v-col>
    </v-app-bar>
    
    <v-dialog v-model="networkloading" max-width="500px" persistent>
        <v-card rounded="xl">
            <div class="networkloading">
                <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
            </div>
        </v-card>
    </v-dialog>
    
    <v-container fluid class="pa-6">
        <v-row justify="center">
            <v-col cols="12" lg="11" xl="10">
                <!-- 用户管理卡片 -->
                <v-card class="account-card" rounded="xl" elevation="2">
                    <!-- 搜索栏 -->
                    <v-card-text class="pa-6 pb-0">
                        <v-row align="center">
                            <v-col cols="12">
                                <v-text-field
                                    v-model="searchQuery"
                                    :placeholder="$t('message.searchUser')"
                                    prepend-inner-icon="mdi-magnify"
                                    variant="outlined"
                                    density="comfortable"
                                    hide-details
                                    clearable
                                    @update:model-value="currentPage = 1"
                                ></v-text-field>
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <!-- 数据表格 -->
                    <v-card-text class="pa-0">
                        <v-data-table
                            :headers="headers"
                            :items="paginatedUsers"
                            :loading="loading"
                            :loading-text="$t('message.loading')"
                            :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')"
                            class="account-table"
                            density="comfortable"
                            hover
                        >
                            <template v-slot:item="{ item }">
                                <tr class="account-table-row">
                                    <td class="text-center pa-4 font-weight-medium">
                                        {{ item.username }}
                                    </td>
                                    <td class="text-center pa-4">
                                        <span class="text-body-2 text-medium-emphasis">
                                            {{ item.email }}
                                        </span>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip
                                            :color="item.role === 1 ? 'error' : 'secondary'"
                                            variant="tonal"
                                            size="small"
                                            class="font-weight-medium"
                                        >
                                            {{ item.role === 1 ? $t('message.admin') : $t('message.regularUser') }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <span class="text-body-2 text-medium-emphasis">
                                            {{ moment(item.create_at).format('YYYY-MM-DD HH:mm') }}
                                        </span>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip
                                            :color="item.isBan ? 'error' : 'success'"
                                            variant="tonal"
                                            size="small"
                                            class="font-weight-medium"
                                        >
                                            {{ item.isBan ? $t('message.isBan') : $t('message.normal') }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-menu>
                                            <template v-slot:activator="{ props }">
                                                <v-btn v-bind="props" variant="text" icon="mdi-dots-horizontal"></v-btn>
                                            </template>
                                            <v-list rounded="xl">
                                                <v-list-item v-for="(menu, index) in getMenus(item)" :key="index" :value="index"
                                                    @click="handleMenuClick(menu, item)">
                                                    <template v-slot:default="{ active, toggle }">
                                                        <div class="d-flex align-center">
                                                            <v-icon :icon="menu.icon" class="me-2"></v-icon>
                                                            <v-list-item-title>{{ menu.title }}</v-list-item-title>
                                                        </div>
                                                    </template>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </td>
                                </tr>
                            </template>
                        </v-data-table>
                    </v-card-text>

                    <!-- 分页 -->
                    <v-card-actions class="pa-6 pt-0" v-if="totalPages > 1">
                        <v-spacer></v-spacer>
                        <v-pagination
                            v-model="currentPage"
                            :length="totalPages"
                            :total-visible="7"
                            rounded="circle"
                            @update:model-value="handlePageChange"
                        ></v-pagination>
                        <v-spacer></v-spacer>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<style scoped>
.account-card {
    border-radius: 16px !important;
}

.account-table {
    border-radius: 0 0 16px 16px;
}

.account-table-row {
    transition: background-color 0.2s ease;
}

.account-table-row:hover {
    background-color: rgba(var(--v-theme-primary), 0.04) !important;
}

.networkloading {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px;
}

.align-left {
    text-align: left;
    margin-left: 10px;
}
</style>