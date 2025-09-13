<script setup>
import { ref, onMounted, computed, reactive, watch, onUnmounted } from 'vue'
import { token, userName } from '../../utils/account'
import { useI18n } from 'vue-i18n';
import { getAllCompetitionsInfo, getCompetitionInfoByIDAdmin, deleteCompetition, updateComInfo, caculateComScore, getScores } from '../../utils/api/admin';
import { verifyUserInfo } from '../../utils/api/auth';
import { showAlert } from '../../utils/alert';
import { MdEditor } from "md-editor-v3";
import { getMdEditorTheme } from '../../utils/theme';
import moment from 'moment';
import "md-editor-v3/lib/style.css";
import { difficultyLang } from '../../utils/dynamic_styles';

const { locale } = useI18n();
const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

const loading = ref(true)
const scoreloading = ref(false)
const networkloading = ref(false)
const delDialog = ref(false)
const scoreDialog = ref(false)

const id = "preview-only";
const editorTheme = ref(getMdEditorTheme());

const scores = ref([])
const totalScores = ref(0)
const itemsPerPage = ref(10)
const page = ref(1)

const cid = ref()
const userPrivilege = ref("")
const competitions = ref([])
const totalCompetitions = ref(0)
const dialog = ref(false)
const isCreate = ref(false)
const showPwd = ref(false)
const formatStartDate = ref("")
const formatEndDate = ref("")
const currentPage = ref(1);
const competitionItemsPerPage = ref(50);

const competitionFields = reactive({
    id: null,
    title: "",
    subtitle: "",
    difficulty: 0,
    password: "",
    crypto: false,
    is_visible: true,
    start_at: "",
    end_at: "",
    announcement: ""
});

const headers = ref([
    { title: 'ID', value: 'id', align: 'center', sortable: false },
    { title: t('message.competition'), value: 'title', align: 'center', sortable: false },
    { title: t('message.difficulty'), value: 'difficulty', align: 'center', sortable: false },
    { title: t('message.status'), value: 'status', align: 'center', sortable: false },
    { title: t('message.isvisible'), value: 'is_visible', align: 'center', sortable: false },
    { title: t("message.operation"), value: 'actions', align: 'center', sortable: false }
])

const difficultyOptions = [
    { value: 0, label: t('message.easy') },
    { value: 1, label: t('message.medium') },
    { value: 2, label: t('message.hard') }
];

const scoreHeaders = ref([
    { title: t('message.username'), value: 'Username', align: 'center' },
    { title: 'Score', value: 'Score', align: 'center' },
])

// 分页后的数据
const paginatedCompetitions = computed(() => {
    const start = (currentPage.value - 1) * competitionItemsPerPage.value;
    const end = start + competitionItemsPerPage.value;
    return competitions.value.slice(start, end);
});

// 总页数
const totalPages = computed(() => {
    return Math.ceil(competitions.value.length / competitionItemsPerPage.value);
});

// 页面变化处理
const handlePageChange = (page) => {
    currentPage.value = page;
};

// 获取竞赛状态文本
const getStatusText = (status) => {
    switch (status) {
        case 0: return t('message.compenotstarted');
        case 1: return t('message.compeprogress');
        case 2: return t('message.compeover');
        default: return t('message.unknown');
    }
}

// 获取竞赛状态颜色
const getStatusColor = (status) => {
    switch (status) {
        case 0: return 'warning';
        case 1: return 'success';
        case 2: return 'error';
        default: return 'grey';
    }
}

// 监听主题变化
const handleThemeChange = (event) => {
    editorTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

// 字段检查
const validateFields = () => {
    for (const key in competitionFields) {
        if (key === "password" && !competitionFields.crypto || key === "announcement") {
            continue;
        }
        if (competitionFields[key] === "" || (Array.isArray(competitionFields[key]) && competitionFields[key].length === 0)) {
            showAlert(t("message.formCheckfailed") + "!", "");
            return false;
        }
    }
    return true;
};

// 删除竞赛
const delCompetition = async () => {
    networkloading.value = true;
    try {
        await deleteCompetition(competitionFields.id);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    } finally {
        networkloading.value = false;
    }
    dialog.value = false;
}

// 清除密码
const clearPassword = () => {
    if (!competitionFields.crypto) {
        competitionFields.password = '';
    }
}

// 创建竞赛
const createCompetition = async () => {
    isCreate.value = true;
    dialog.value = true;
    // 重置内容
    competitionFields.id = totalCompetitions.value + 1;
    competitionFields.title = "";
    competitionFields.subtitle = "";
    competitionFields.difficulty = 0;
    competitionFields.password = "";
    competitionFields.crypto = false;
    competitionFields.is_visible = true;
    competitionFields.start_at = "";
    competitionFields.end_at = "";
    competitionFields.announcement = "";
    formatStartDate.value = "";
    formatEndDate.value = "";
}

// 数据传至后端
const save = async () => {
    try {
        competitionFields.start_at = new Date(formatStartDate.value).toISOString();
        competitionFields.end_at = new Date(formatEndDate.value).toISOString();
    } catch (error) {
        showAlert(t('message.formRuleCheckfailed') + "!", "");
        return;
    }
    if (!validateFields()) return;
    networkloading.value = true;
    const comData = { ...competitionFields };
    try {
        await updateComInfo(comData);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        showAlert(t("message.failed") + "!", "reload");
    } finally {
        networkloading.value = false;
    }
    dialog.value = false;
};

// 编辑/创建竞赛
const goToEditCompetition = async (id) => {
    isCreate.value = false;
    dialog.value = true;
    networkloading.value = true;
    const comResp = await getCompetitionInfoByIDAdmin(id);
    networkloading.value = false;

    Object.assign(competitionFields, comResp.data.data);
    formatStartDate.value = moment(competitionFields.start_at).format('YYYY-MM-DD HH:mm');
    formatEndDate.value = moment(competitionFields.end_at).format('YYYY-MM-DD HH:mm');
}

// 从后端获取数据
const fetchData = async () => {
    loading.value = true
    try {
        const response = await getAllCompetitionsInfo();
        competitions.value = response.data.data
        totalCompetitions.value = competitions.value.length
        competitionFields.id = totalCompetitions.value + 1
    } catch (error) {
        loading.value = false
        showAlert(t("message.failed") + "!", "")
    } finally {
        loading.value = false
    }
}

// 竞赛计分
const calcCompetition = async (competitionId) => {
    try {
        await caculateComScore(competitionId);
        showAlert(t("message.calculating"), "");
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    }
}

// 查看竞赛得分情况
const getScoreBoard = async (competitionId) => {
    scoreloading.value = true;
    scoreDialog.value = true;
    cid.value = competitionId;
    try {
        const response = await getScores(competitionId, page.value, itemsPerPage.value);
        scores.value = response.data.users;
        totalScores.value = response.data.total;
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    } finally {
        scoreloading.value = false;
    }
}

// 监听分页变化
watch(page, async () => {
    await getScoreBoard(cid.value)
})

onMounted(async () => {
    loading.value = true;
    try {
        if (!userLoggedIn.value) {
            window.location = "#/login";
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

    // 监听主题变化
    window.addEventListener('theme-change', handleThemeChange);
});

onUnmounted(() => {
    // 清理事件监听器
    window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
    <template>
        <v-dialog v-model="scoreDialog" persistent max-width="1200px">
            <v-card rounded="xl">
                <v-data-table-server :headers="scoreHeaders" :items="scores" :items-length="totalScores"
                    :loading="scoreloading" :loading-text="$t('message.loading')" :hide-default-footer="true"
                    :items-per-page="itemsPerPage" v-model:page="page" @update="getScoreBoard(cid.value)">
                    <template v-slot:bottom>
                        <v-pagination v-model="page" :length="Math.ceil(totalScores / itemsPerPage)" rounded="xl" />
                    </template>
                </v-data-table-server>
                <v-card-actions>
                    <v-btn color="primary" @click="scoreDialog = false" rounded="xl">{{ $t('message.ok') }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </template>
    <template>
        <v-dialog v-model="delDialog" persistent max-width="290">
            <v-card rounded="xl">
                <v-card-title class="text-h5">{{ $t('message.notify') }}</v-card-title>
                <v-card-text>{{ t('message.suredel') }}</v-card-text>
                <v-card-actions>
                    <v-btn variant="elevated" color="primary" @click="delCompetition" rounded="xl">{{ $t('message.yes')
                    }}</v-btn>
                    <v-btn color="primary" @click="delDialog = false" rounded="xl">{{ $t('message.cancel') }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </template>

    <v-app-bar :elevation="2">
        <template v-slot:prepend>
            <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
        </template>
        <v-col class="align-left">
            <p style="font-size: 24px;">{{ t('message.competitionmanagement') }}</p>
        </v-col>
    </v-app-bar>

    <v-container fluid class="pa-6">
        <v-row justify="center">
            <v-col cols="12" lg="11" xl="10">
                <!-- 竞赛管理卡片 -->
                <v-card class="competitions-card" rounded="xl" elevation="2">
                    <!-- 数据表格 -->
                    <v-card-text class="pa-0">
                        <v-data-table :headers="headers" :items="paginatedCompetitions" :loading="loading"
                            :loading-text="$t('message.loading')"
                            :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')"
                            class="competitions-table" density="comfortable" hover>
                            <template v-slot:item="{ item }">
                                <tr class="competitions-table-row">
                                    <td class="text-center pa-4 font-weight-medium">
                                        {{ item.id }}
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-btn @click="goToEditCompetition(item.id)" variant="text"
                                            color="primary" class="font-weight-medium">
                                            {{ item.title }}
                                        </v-btn>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip
                                            :color="item.difficulty === 0 ? 'success' : item.difficulty === 1 ? 'warning' : 'error'"
                                            variant="tonal" size="small" class="font-weight-medium">
                                            {{ $t(difficultyLang(item.difficulty)) }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip :color="getStatusColor(item.status)" variant="tonal" size="small"
                                            class="font-weight-medium">
                                            {{ getStatusText(item.status) }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip :color="item.is_visible ? 'success' : 'error'" variant="tonal"
                                            size="small" class="font-weight-medium">
                                            {{ item.is_visible ? $t('message.visible') : $t('message.invisible') }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-menu>
                                            <template v-slot:activator="{ props }">
                                                <v-btn v-bind="props" variant="text" icon="mdi-dots-horizontal"></v-btn>
                                            </template>
                                            <v-list rounded="xl">
                                                <v-list-item :disabled="(item.status != 2) || (item.scored != false)"
                                                    @click="calcCompetition(item.id)">
                                                    <template v-slot:default="{ active, toggle }">
                                                        <div class="d-flex align-center">
                                                            <v-icon icon="mdi-calculator" class="me-2"></v-icon>
                                                            <v-list-item-title>{{ t('message.scoring')
                                                                }}</v-list-item-title>
                                                        </div>
                                                    </template>
                                                </v-list-item>
                                                <v-list-item @click="getScoreBoard(item.id)"
                                                    :disabled="item.scored != true">
                                                    <template v-slot:default="{ active, toggle }">
                                                        <div class="d-flex align-center">
                                                            <v-icon icon="mdi-note-text" class="me-2"></v-icon>
                                                            <v-list-item-title>{{ t("message.viewScores")
                                                                }}</v-list-item-title>
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
                        <v-pagination v-model="currentPage" :length="totalPages" :total-visible="7" rounded="circle"
                            @update:model-value="handlePageChange"></v-pagination>
                        <v-spacer></v-spacer>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
    </v-container>

    <v-dialog v-model="dialog" max-width="1200px">
        <v-card>
            <div v-if="networkloading" class="networkloading">
                <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
            </div>
            <div v-else>
                <v-card-title>
                    <span class="headline">{{ $t('message.edit') }}</span>
                </v-card-title>
                <v-card-text>
                    <v-form>
                        <v-text-field label="ID" v-model="competitionFields.id" variant="solo-filled"
                            readonly></v-text-field>
                        <!-- 竞赛名称 -->
                        <v-text-field :label="$t('message.title')" v-model="competitionFields.title"
                            variant="solo-filled"></v-text-field>
                        <!-- 竞赛副标题 -->
                        <v-text-field :label="$t('message.subtitle')" v-model="competitionFields.subtitle"
                            variant="solo-filled"></v-text-field>
                        <v-row class="limitRow">
                            <!-- 难易程度 -->
                            <v-select :items="difficultyOptions" item-title="label" item-value="value"
                                :label="$t('message.difficulty')" v-model="competitionFields.difficulty"
                                variant="solo-filled" />
                            <div style="margin-inline: 30px;"></div>
                            <!-- 是否可见 -->
                            <v-switch v-model="competitionFields.is_visible" :label="$t('message.isvisible')"
                                color="primary" inset></v-switch>
                            <div style="margin-inline: 30px;"></div>
                            <!-- 启用密码 -->
                            <v-switch v-model="competitionFields.crypto" :label="$t('message.withapwd')"
                                color="primary" inset @change="clearPassword"></v-switch>
                        </v-row>
                        <v-text-field v-if="competitionFields.crypto" v-model="competitionFields.password"
                            :append-icon="showPwd ? 'mdi-eye' : 'mdi-eye-off'" variant="solo-filled"
                            :type="showPwd ? 'text' : 'password'" :label="$t('message.password')" counter
                            @click:append="showPwd = !showPwd"></v-text-field>
                        <!-- 竞赛开始时间日期 -->
                        <v-row class="limitRow">
                            <v-text-field :label="$t('message.start_date') + '(YYYY-MM-DD HH:mm)'"
                                v-model="formatStartDate" variant="solo-filled"></v-text-field>
                            <div style="margin-inline: 30px;"></div>
                            <v-text-field :label="$t('message.end_date') + '(YYYY-MM-DD HH:mm)'" v-model="formatEndDate"
                                variant="solo-filled"></v-text-field>
                        </v-row>
                        <md-editor v-model="competitionFields.announcement" :editorId="id" :noUploadImg="true"
                            :footers="[]" :language="locale === 'zh_CN' ? 'zh-CN' : 'en-US'" :theme="editorTheme" />
                    </v-form>
                    <div style="position: fixed; bottom: 16px; right: 16px; z-index: 1000;">
                        <v-btn @click="dialog = false" rounded="xl" style="margin-right: 10px;">{{
                            $t('message.cancel') }}</v-btn>
                        <v-btn v-if="!isCreate" color="red-lighten-1" @click="delDialog = true" rounded="xl"
                            style="margin-right: 10px;">{{ $t('message.delete')
                            }}</v-btn>
                        <v-btn color="primary" @click="save" rounded="xl" style="margin-right: 10px;">{{
                            $t('message.save')
                            }}</v-btn>
                    </div>
                </v-card-text>
            </div>
        </v-card>
    </v-dialog>

    <div class="fab">
        <v-fab fixed icon="mdi-plus" size="64" color="primary" elevation="10" v-if="!loading"
            @click="createCompetition"></v-fab>
    </div>
</template>

<style scoped>
.competitions-card {
    border-radius: 16px !important;
}

.competitions-table {
    border-radius: 0 0 16px 16px;
}

.competitions-table-row {
    transition: background-color 0.2s ease;
}

.competitions-table-row:hover {
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

.limitRow {
    display: flex;
    justify-content: space-between;
    margin: 5px 0;
}

.fab {
    position: fixed;
    bottom: 80px;
    right: 180px;
}
</style>