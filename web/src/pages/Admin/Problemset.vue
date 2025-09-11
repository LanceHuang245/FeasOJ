<!-- 题目后台管理页 -->
<script setup>
import { ref, onMounted, computed, reactive, onUnmounted } from 'vue'
import { token, userName } from '../../utils/account'
import { getProblemAllInfoByAdmin, updateProblemInfo, deleteProblemAllInfo, getAllCompetitionsInfo, getAllProblemsAdmin } from '../../utils/api/admin';
import { verifyUserInfo } from '../../utils/api/auth';
import { showAlert } from '../../utils/alert';
import { MdEditor } from 'md-editor-v3';
import { useI18n } from 'vue-i18n';
import { getMdEditorTheme } from '../../utils/theme';
import 'md-editor-v3/lib/style.css';
import { difficultyLang } from '../../utils/dynamic_styles';

const { locale } = useI18n();
const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

const loading = ref(true)
const networkloading = ref(false)

const delDialog = ref(false)
const userPrivilege = ref("")
const problems = ref([])
const totalProblems = ref(0)
const competitionIds = ref([0])
const dialog = ref(false)
const editorTheme = ref(getMdEditorTheme());
const searchQuery = ref('');
const currentPage = ref(1);
const itemsPerPage = ref(50);
const batchDialog = ref(false)
const batchInput = ref('')

const testCaseHeaders = [
    { title: t('message.input'), value: 'input', align: 'center' },
    { title: t('message.output'), value: 'output', align: 'center' },
    { title: t('message.operation'), value: 'actions', align: 'center', sortable: false }
]

const difficultyOptions = [
    { value: '简单', label: t('message.easy') },
    { value: '中等', label: t('message.medium') },
    { value: '困难', label: t('message.hard') }
];

const problemFields = reactive({
    pid: null,
    title: "",
    content: "",
    difficulty: "",
    time_limit: "",
    memory_limit: "",
    input: "",
    output: "",
    contest_id: null,
    is_visible: true,
    test_cases: [{ input: '', output: '' }]
});

const headers = ref([
    { title: 'ID', value: 'Pid', align: 'center', sortable: false },
    { title: t('message.problem'), value: 'Title', align: 'center', sortable: false },
    { title: t('message.difficulty'), value: 'Difficulty', align: 'center', sortable: false },
    { title: t('message.contestid'), value: 'ContestId', align: 'center', sortable: false },
    { title: t('message.isvisible'), value: 'IsVisible', align: 'center', sortable: false },
    { title: t('message.operation'), value: 'actions', align: 'center', sortable: false },
])

const isCreate = ref(false)

// 过滤后的数据
const filteredProblems = computed(() => {
    if (!searchQuery.value) {
        return problems.value;
    }
    return problems.value.filter((problem) =>
        problem.Title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        problem.Pid.toString().includes(searchQuery.value)
    )
})

// 分页后的数据
const paginatedProblems = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value;
    const end = start + itemsPerPage.value;
    return filteredProblems.value.slice(start, end);
});

// 总页数
const totalPages = computed(() => {
    return Math.ceil(filteredProblems.value.length / itemsPerPage.value);
});

// 页面变化处理
const handlePageChange = (page) => {
    currentPage.value = page;
};

// 监听主题变化
const handleThemeChange = (event) => {
    editorTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

// 添加测试样例
const addTestCase = () => {
    problemFields.test_cases.push({ input: '', output: '' });
};

// 字段检查
const validateFields = () => {
    for (const key in problemFields) {
        if (problemFields[key] === "" || (Array.isArray(problemFields[key]) && problemFields[key].length === 0)) {
            showAlert(t("message.formCheckfailed") + "!", "");
            return false;
        }
    }
    if (problemFields.test_cases.some(testCase => testCase.output === "")) {
        showAlert(t("message.formCheckfailed") + "!", "");
        return false;
    }
    return true;
};

// 数据传至后端
const save = async () => {
    if (!validateFields()) return;
    networkloading.value = true;
    const problemData = { ...problemFields };
    try {
        await updateProblemInfo(problemData);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        showAlert(t("message.failed") + "!", "reload");
    } finally {
        networkloading.value = false;
    }
    dialog.value = false;
};

// 从后端获取数据
const fetchData = async () => {
    loading.value = true
    try {
        const response = await getAllProblemsAdmin();
        problems.value = response.data.problems
        totalProblems.value = problems.value.length
        problemFields.pid = totalProblems.value + 1
    } catch (error) {
        showAlert(t("message.failed") + "!", "")
    } finally {
        loading.value = false
    }
}

// 创建题目
const createProblem = async () => {
    networkloading.value = true;
    isCreate.value = true;
    dialog.value = true;
    // 重置内容
    problemFields.pid = totalProblems.value + 1;
    problemFields.title = "";
    problemFields.content = "";
    problemFields.difficulty = "简单";
    problemFields.time_limit = "";
    problemFields.memory_limit = "";
    problemFields.input = "";
    problemFields.output = "";
    problemFields.contest_id = 0;
    problemFields.is_visible = true;
    problemFields.test_cases = [{ input: '', output: '' }];
    const compResp = await getAllCompetitionsInfo();
    competitionIds.value = [0, ...compResp.data.contests.map(contest => contest.contest_id)];
    networkloading.value = false;
}

// 编辑题目显示
const goToEditProblem = async (pid) => {
    isCreate.value = false;
    dialog.value = true;
    networkloading.value = true;
    const problemResp = await getProblemAllInfoByAdmin(pid);
    Object.assign(problemFields, problemResp.data.problemInfo);
    const compResp = await getAllCompetitionsInfo();
    competitionIds.value = [0, ...compResp.data.contests.map(contest => contest.contest_id)];
    networkloading.value = false;
};

// 删除题目
const delProblem = async () => {
    delDialog.value = false;
    networkloading.value = true;
    try {
        await deleteProblemAllInfo(problemFields.pid);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        showAlert(t("message.failed") + "!", "");
    } finally {
        networkloading.value = false;
        dialog.value = false;
    }
}
// 批量导入处理
const handleBatchImport = () => {
    const lines = batchInput.value.split(/\r?\n/).map(l => l.trim()).filter(l => l !== '')
    const newCases = []
    for (let i = 0; i < lines.length; i += 2) {
        if (lines[i] !== undefined && lines[i + 1] !== undefined) {
            newCases.push({ input: lines[i], output: lines[i + 1] })
        }
    }
    if (newCases.length > 0) {
        problemFields.test_cases = newCases
    }
    batchDialog.value = false
    batchInput.value = ''
}
// 删除指定测试用例
const removeTestCaseAt = (index) => {
    problemFields.test_cases.splice(index, 1)
}

onMounted(async () => {
    loading.value = true;
    try {
        if (!userLoggedIn.value) {
            window.location = "#/login";
            return;
        }
        const userInfoResponse = await verifyUserInfo(userName.value, token.value);
        if (userInfoResponse.status !== 200) {
            window.location = '#/403';
            return;
        }
        userPrivilege.value = userInfoResponse.data.info.role;
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
        <v-dialog v-model="delDialog" persistent max-width="290">
            <v-card rounded="xl">
                <v-card-title class="text-h5">{{ $t('message.notify') }}</v-card-title>
                <v-card-text>{{ t('message.suredel') }}</v-card-text>
                <v-card-actions>
                    <v-btn variant="elevated" color="primary" @click="delProblem" rounded="xl">{{ $t('message.yes')
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
            <p style="font-size: 24px;">{{ t('message.problemmanagement') }}</p>
        </v-col>
    </v-app-bar>

    <v-container fluid class="pa-6">
        <v-row justify="center">
            <v-col cols="12" lg="11" xl="10">
                <!-- 题目管理卡片 -->
                <v-card class="problemset-card" rounded="xl" elevation="2">
                    <!-- 搜索栏 -->
                    <v-card-text class="pa-6 pb-0">
                        <v-row align="center">
                            <v-col cols="12">
                                <v-text-field v-model="searchQuery" :placeholder="$t('message.searchProblem')"
                                    prepend-inner-icon="mdi-magnify" variant="outlined" density="comfortable"
                                    hide-details clearable @update:model-value="currentPage = 1"></v-text-field>
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <!-- 数据表格 -->
                    <v-card-text class="pa-0">
                        <v-data-table :headers="headers" :items="paginatedProblems" :loading="loading"
                            :loading-text="$t('message.loading')"
                            :no-data-text="!userLoggedIn ? $t('message.nologin') : $t('message.nodata')"
                            class="problemset-table" density="comfortable" hover>
                            <template v-slot:item="{ item }">
                                <tr class="problemset-table-row">
                                    <td class="text-center pa-4 font-weight-medium">
                                        {{ item.Pid }}
                                    </td>
                                    <td class="text-center pa-4 font-weight-medium">
                                        {{ item.Title }}
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip
                                            :color="item.Difficulty === '简单' ? 'success' : item.Difficulty === '中等' ? 'warning' : 'error'"
                                            variant="tonal" size="small" class="font-weight-medium">
                                            {{ $t(difficultyLang(item.Difficulty)) }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <span class="text-body-2 text-medium-emphasis">
                                            {{ item.ContestID }}
                                        </span>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-chip :color="item.IsVisible ? 'success' : 'error'" variant="tonal"
                                            size="small" class="font-weight-medium">
                                            {{ item.IsVisible ? $t('message.visible') : $t('message.invisible') }}
                                        </v-chip>
                                    </td>
                                    <td class="text-center pa-4">
                                        <v-btn @click="goToEditProblem(item.Pid)" variant="text" icon="mdi-pencil"
                                            size="small" color="primary"></v-btn>
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

    <v-dialog v-model="dialog" max-width="800px">
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
                        <v-text-field label="ID" v-model="problemFields.pid" variant="solo-filled"
                            readonly></v-text-field>
                        <!-- 题目名称 -->
                        <v-text-field :label="$t('message.title')" v-model="problemFields.title"
                            variant="solo-filled"></v-text-field>
                        <!-- 题目描述 -->
                        <md-editor v-model="problemFields.content" :footers="[]" :noUploadImg="true" :tabWidth="4"
                            :language="locale === 'zh_CN' ? 'zh-CN' : 'en-US'" :theme="editorTheme" />
                        <div style="margin-top: 20px;"></div>
                        <!-- 难易程度 -->
                        <v-select :items="difficultyOptions" item-title="label" item-value="value"
                            :label="$t('message.difficulty')" v-model="problemFields.difficulty"
                            variant="solo-filled" />
                        <!-- 所属竞赛ID及是否可见 -->
                        <v-row class="limitRow">
                            <v-select :items="competitionIds" :label="$t('message.contestid')"
                                v-model.number="problemFields.contest_id" variant="solo-filled"></v-select>
                            <div style="margin-inline: 30px;"></div>
                            <v-switch v-model="problemFields.is_visible" :label="$t('message.isvisible')"
                                color="primary" inset></v-switch>
                        </v-row>
                        <!-- 时间、内存限制 -->
                        <v-row class="limitRow">
                            <v-text-field :label="$t('message.timeLimit') + '(S)'" v-model="problemFields.time_limit"
                                variant="solo-filled"></v-text-field>
                            <div style="margin-inline: 30px;"></div>
                            <v-text-field :label="$t('message.memoryLimit') + '(MB)'"
                                v-model="problemFields.memory_limit" variant="solo-filled"></v-text-field>
                        </v-row>
                        <!-- 显示给用户的输入输出样例 -->
                        <v-row class="limitRow">
                            <v-text-field :label="$t('message.displayInputCase')" v-model="problemFields.input"
                                variant="solo-filled"></v-text-field>
                            <div style="margin-inline: 30px;"></div>
                            <v-text-field :label="$t('message.displayOutputCase')" v-model="problemFields.output"
                                variant="solo-filled"></v-text-field>
                        </v-row>
                    </v-form>
                    <!-- 输入输出测试样例表格与批量导入 -->
                    <v-row>
                        <v-col cols="12">
                            <v-data-table :headers="testCaseHeaders" :items="problemFields.test_cases" class="mb-2"
                                hide-default-footer>
                                <template v-slot:item.input="{ item, index }">
                                    <v-text-field v-model="item.input" variant="solo-filled" density="compact"
                                        hide-details></v-text-field>
                                </template>
                                <template v-slot:item.output="{ item, index }">
                                    <v-text-field v-model="item.output" variant="solo-filled" density="compact"
                                        hide-details></v-text-field>
                                </template>
                                <template v-slot:item.actions="{ index }">
                                    <v-btn icon="mdi-delete" color="error" @click="removeTestCaseAt(index)"
                                        size="small"></v-btn>
                                </template>
                            </v-data-table>
                            <v-btn @click="batchDialog = true" color="primary" rounded="xl" style="margin-right: 10px;">
                                {{ $t('message.batchImportTestCases') }}
                            </v-btn>
                            <v-btn @click="addTestCase" color="primary" rounded="xl">{{ $t('message.addTestCase')
                                }}</v-btn>
                        </v-col>
                    </v-row>
                    <!-- 批量导入对话框 -->
                    <v-dialog v-model="batchDialog" max-width="600px">
                        <v-card>
                            <v-card-title>{{ $t('message.batchImportTestCases') }}</v-card-title>
                            <v-card-text>
                                <div style="margin-bottom: 10px; color: #888; font-size: 14px;">
                                    {{ $t('message.batchImportTip') }}
                                </div>
                                <v-textarea v-model="batchInput" rows="10" auto-grow
                                    :placeholder="$t('message.batchImportPlaceholder')" />
                            </v-card-text>
                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn color="primary" rounded="xl" @click="handleBatchImport">{{ $t('message.save')
                                    }}</v-btn>
                                <v-btn rounded="xl" @click="batchDialog = false">{{ $t('message.cancel') }}</v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
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
            @click="createProblem"></v-fab>
    </div>
</template>

<style scoped>
.problemset-card {
    border-radius: 16px !important;
}

.problemset-table {
    border-radius: 0 0 16px 16px;
}

.problemset-table-row {
    transition: background-color 0.2s ease;
}

.problemset-table-row:hover {
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

.fab {
    position: fixed;
    bottom: 80px;
    right: 180px;
}

.limitRow {
    display: flex;
    justify-content: space-between;
    margin: 5px 0;
}
</style>