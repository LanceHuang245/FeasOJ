<script setup>
import { onMounted, ref, computed, onUnmounted } from 'vue';
import { useI18n } from "vue-i18n";
import { token } from '../../utils/account';
import { showAlert } from '../../utils/alert';
import { useRoute, useRouter } from 'vue-router';
import { getCompetitionById, getCompetitionUsers, getCompetitionProblems, isInCompetition, quitCompetition } from '../../utils/api/competitions';
import { avatarServer } from '../../utils/axios';
import { MdPreview } from 'md-editor-v3';
import { getMdPreviewTheme } from '../../utils/theme';
import 'md-editor-v3/lib/preview.css';
import moment from "moment";
import { difficultyColor, difficultyLang } from '../../utils/dynamic_styles';

const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

const route = useRoute();
const router = useRouter();

const quitDialog = ref(false);

// 竞赛信息
const contestInfo = ref({});
const competitionId = route.params.cid;

// 参赛人员信息
const usersInfo = ref([]);

// 竞赛题目信息
const problems = ref([]);

// 显示题目状态
const compStatus = (status) => {
    switch (status) {
        case 0:
            return 'message.compenotstarted';
        case 1:
            return 'message.compeprogress';
        case 2:
            return 'message.compeover';
        default:
            return 'message.compenotstarted';
    }
}

const networkloading = ref(false);
const loading = ref(false);
const id = 'preview-only';
const previewTheme = ref(getMdPreviewTheme());

// 监听主题变化
const handleThemeChange = (event) => {
  previewTheme.value = event.detail.theme === 'dark' ? 'dark' : 'light';
};

const quitComp = async () => {
    networkloading.value = true;
    try {
        const response = await quitCompetition(competitionId);
        showAlert(response.data.message, "/competitions");
    } catch (error) {
        showAlert(error.response.data.message, "");
    } finally {
        networkloading.value = false;
    }
}

const isContestEnded = computed(() => {
    if (!contestInfo.value.end_at) return false;
    return moment().isSameOrAfter(moment(contestInfo.value.end_at));
});

onMounted(async () => {
    loading.value = true;
    if (userLoggedIn.value) {
        const response = await isInCompetition(competitionId);
        if (response.data.isIn) {
            try {
                // 获取该竞赛详细信息
                const resp = await getCompetitionById(competitionId);
                contestInfo.value = resp.data.data;

                // 获取该竞赛参赛人员列表
                const resp2 = await getCompetitionUsers(competitionId);
                usersInfo.value = resp2.data.data;

                // 获取该竞赛题目列表
                if (contestInfo.value.status != 0) {
                    const resp3 = await getCompetitionProblems(competitionId);
                    problems.value = resp3.data.data;
                }
            } catch (error) {
                showAlert(t("message.failed") + "!", "");
            } finally {
                loading.value = false;
            }
        } else {
            window.location = '#/competitions'
        }
    } else {
        window.location = '#/login'
    }
    
    // 监听主题变化
    window.addEventListener('theme-change', handleThemeChange);
})

onUnmounted(() => {
    // 清理事件监听器
    window.removeEventListener('theme-change', handleThemeChange);
});
</script>

<template>
    <template>
        <v-dialog v-model="networkloading" max-width="600px">
            <v-card rounded="xl">
                <div class="networkloading">
                    <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
                </div>
            </v-card>
        </v-dialog>
    </template>
    <template>
        <v-dialog v-model="quitDialog" persistent max-width="290">
            <v-card rounded="xl">
                <v-card-title class="text-h5">{{ $t('message.notify') }}</v-card-title>
                <v-card-text>{{ t('message.surequit') }}</v-card-text>
                <v-card-actions>
                    <v-btn variant="elevated" color="primary" @click="quitComp" rounded="xl">{{ $t('message.yes')
                        }}</v-btn>
                    <v-btn color="primary" @click="quitDialog = false" rounded="xl">{{ $t('message.cancel') }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </template>
    <div v-if="loading" class="loading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
    </div>
    <div v-else>
        <v-app-bar :elevation="2">
            <template v-slot:prepend>
                <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
            </template>
            <v-col class="align-left">
                <p class="font-weight-black" style="font-size: 24px;">{{ contestInfo.title }}</p>
                <p>{{ contestInfo.subtitle }}</p>
            </v-col>
            <template v-slot:append>
                <v-btn v-if="!isContestEnded" color="primary" variant="flat" rounded="xl" @click="quitDialog = true">{{ t('message.quit')
                    }}</v-btn>
            </template>
        </v-app-bar>
        <v-container>
            <v-col>
                <v-row>
                    <v-card rounded="xl" width="100%" elevation="3">
                        <template v-slot:title>
                            <span class="font-weight-black">{{ t(compStatus(contestInfo.status)) }}</span>
                        </template>
                        <v-card-text class="bg-surface-light pt-4">
                            {{ moment(contestInfo.start_at).format("YYYY/MM/DD HH:mm") }} -
                            {{ moment(contestInfo.end_at).format("YYYY/MM/DD HH:mm") }}
                        </v-card-text>
                    </v-card>
                </v-row>
                <div style="height: 50px;"></div>
                <v-row>
                    <v-card rounded="xl" width="100%" elevation="3">
                        <template v-slot:title>
                            <span class="font-weight-black">{{ t('message.announcement') }}</span>
                        </template>
                        <md-preview v-if="contestInfo.announcement" :modelValue="contestInfo.announcement"
                            :editorId="id" class="md_preview" :theme="previewTheme" />
                    </v-card>
                </v-row>
                <div style="height: 50px;"></div>
                <v-row>
                    <v-col cols="8" style="padding: 0px 20px 0px 0px;">
                        <v-card rounded="xl" elevation="3" height="500px">
                            <template v-slot:title>
                                <span class="font-weight-black">{{ t("message.problem") }}</span>
                            </template>
                            <v-list v-if="problems.length > 0" style="max-height: 450px; overflow-y: auto;">
                                <v-list-item v-for="p in problems" :key="p.id">
                                    <v-list-item-title>
                                        <v-btn variant="text" color="primary" block
                                            @click="router.push({ path: `/problemset/${p.id}` })">
                                            {{ p.title }}
                                        </v-btn>
                                    </v-list-item-title>
                                    <template v-slot:append>
                                        <div style="margin-left: 10px;"></div>
                                        <v-chip :style="difficultyColor(p.difficulty)">
                                            {{ t(difficultyLang(p.difficulty)) }}
                                        </v-chip>
                                    </template>
                                </v-list-item>
                            </v-list>
                            <div v-else class="text-center">
                                <p>{{ t('message.unavailable') }}</p>
                            </div>
                        </v-card>
                    </v-col>
                    <v-col cols="4" style="padding: 0;">
                        <v-card rounded="xl" elevation="3" height="500px">
                            <template v-slot:title>
                                <span class="font-weight-black">{{ t('message.participants') }}</span>
                            </template>
                            <v-list style="max-height: 450px; overflow-y: auto;">
                                <v-list-item v-for="user in usersInfo" :key="user.Uid">
                                    <v-list-item-title class="username-avatar">
                                        <v-avatar size="40" color="surface-variant">
                                            <v-img :src="avatarServer + user.avatar" alt="Avatar" cover></v-img>
                                        </v-avatar>
                                        <div style="margin: 5px">
                                            <div class="username"
                                                @click="router.push({ path: `/profile/${user.username}` })">
                                                {{ user.username }}
                                            </div>
                                            <div class="timeline">
                                                {{ moment(user.join_date).format("MM-DD HH:mm") }}
                                            </div>
                                        </div>
                                    </v-list-item-title>
                                    <v-divider></v-divider>
                                </v-list-item>
                            </v-list>
                        </v-card>
                    </v-col>
                </v-row>
            </v-col>
        </v-container>
    </div>
</template>

<style scoped>
.loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.username-avatar {
    text-align: left;
    display: flex;
    align-items: center;
}

.timeline {
    text-align: left;
    font-size: 0.8em;
    color: #999;
}

.username {
    font-size: 0.9em;
    cursor: pointer;
}

.md_preview {
    text-align: left;
    margin: 10px;
}

.text-center {
    text-align: center;
}
</style>
