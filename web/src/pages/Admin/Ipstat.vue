<script setup>
import { onMounted, ref, computed } from 'vue';
import { getIpStat } from '../../utils/api/admin';
import { verifyUserInfo } from '../../utils/api/auth';
import { showNotification } from '../../utils/notification';
import { token, userName } from '../../utils/account'
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import moment from 'moment';

const { t } = useI18n();
const router = useRouter();

// 响应式数据
const loading = ref(true);
const networkloading = ref(false);
const userPrivilege = ref("");
const ipStats = ref([]);
const searchQuery = ref('');
const currentPage = ref(1);
const itemsPerPage = ref(50);

// 表格头部定义
const headers = ref([
    { title: t('message.ipAddress'), value: 'ip', align: 'center', sortable: false },
    { title: t('message.visitCount'), value: 'visit_count', align: 'center', sortable: false },
    { title: t('message.lastVisitTime'), value: 'last_visit', align: 'center', sortable: false },
]);

// 过滤后的数据
const filteredIpStats = computed(() => {
    if (!searchQuery.value) {
        return ipStats.value;
    }
    return ipStats.value.filter((item) =>
        item.IP.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

// 分页后的数据
const paginatedIpStats = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value;
    const end = start + itemsPerPage.value;
    return filteredIpStats.value.slice(start, end);
});

// 总页数
const totalPages = computed(() => {
    return Math.ceil(filteredIpStats.value.length / itemsPerPage.value);
});

// 获取IP统计数据
const fetchData = async () => {
    loading.value = true;
    try {
        const response = await getIpStat();
        ipStats.value = response.data.ipStatistics || [];
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    } finally {
        loading.value = false;
    }
};

// 格式化时间
const formatTime = (time) => {
    return moment(time).format('YYYY-MM-DD HH:mm:ss');
};

// 页面变化处理
const handlePageChange = (page) => {
    currentPage.value = page;
};

// 复制IP到剪贴板
const copyIp = async (ip) => {
    try {
        await navigator.clipboard.writeText(ip);
        showNotification(t('message.copied'));
    } catch (e) {
        showNotification(t('message.copyFailed'));
    }
};

// 初始化
onMounted(async () => {
    loading.value = true;
    try {
        if (!token.value) {
            window.location = "#/login";
            window.location.reload();
            return;
        }
        const userInfoResponse = await verifyUserInfo(userName.value, token.value);
        userPrivilege.value = userInfoResponse.data.info.role;
        if (userPrivilege.value !== 1) {
            window.location = '#/403';
            return;
        }
        await fetchData();
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
            <p style="font-size: 24px;">{{ t('message.ipstat') }}</p>
        </v-col>
    </v-app-bar>

    <!-- 加载对话框 -->
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
                <!-- IP统计卡片 -->
                <v-card class="ipstat-card" rounded="xl" elevation="2">
                    <!-- 搜索栏 -->
                    <v-card-text class="pa-6 pb-0">
                        <v-row align="center">
                            <v-col cols="12">
                                <v-text-field v-model="searchQuery" :placeholder="t('message.searchIpAddress')"
                                    prepend-inner-icon="mdi-magnify" variant="outlined" density="comfortable"
                                    hide-details clearable @update:model-value="currentPage = 1"></v-text-field>
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <!-- 数据表格 -->
                    <v-card-text class="pa-0">
                        <v-data-table :headers="headers" :items="paginatedIpStats" :loading="loading"
                            :loading-text="t('message.loading')" :no-data-text="t('message.nodata')"
                            class="ipstat-table" density="comfortable" hover>
                            <template v-slot:item="{ item }">
                                <tr class="ipstat-table-row">
                                    <td class="text-center pa-4 font-weight-medium">
                                        <v-chip color="primary" variant="elevated" class="ip-chip font-weight-medium"
                                            @click.stop="copyIp(item.IP)">
                                            <v-icon start small>mdi-content-copy</v-icon>
                                            {{ item.IP }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip color="secondary" variant="tonal" size="small"
                                            class="font-weight-medium">
                                            {{ item.VisitCount }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <span class="text-body-2 text-medium-emphasis">
                                            {{ formatTime(item.LastVisit) }}
                                        </span>
                                    </td>
                                </tr>
                            </template>
                        </v-data-table>
                    </v-card-text>

                    <!-- 分页 -->
                    <v-card-actions class="pa-6 pt-0" v-if="totalPages > 1">
                        <v-spacer></v-spacer>
                        <v-pagination v-model="currentPage" :length="totalPages" :total-visible="7" rounded="circle"
                            @update:model-value="handlePageChange"></v-pagination>
                        <v-spacer></v-spacer>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<style scoped>
.ipstat-card {
    border-radius: 16px !important;
}

.ipstat-table {
    border-radius: 0 0 16px 16px;
}

.ipstat-table-row {
    transition: background-color 0.2s ease;
}

.ipstat-table-row:hover {
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

.ip-chip {
    transition: box-shadow 0.2s, background 0.2s;
    cursor: pointer;
    user-select: all;
}

.ip-chip:active {
    box-shadow: 0 0 0 4px rgba(var(--v-theme-primary), 0.12);
    background: rgba(var(--v-theme-primary), 0.08);
}
</style>