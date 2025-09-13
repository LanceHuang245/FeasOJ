<script setup>
import { token } from "../../utils/account";
import { ref, onMounted, computed } from "vue";
import { getAllCompetitions, isInCompetition, joinCompWithPwd, joinCompetition } from "../../utils/api/competitions";
import { useI18n } from "vue-i18n";
import moment from "moment";
import { showAlert } from "../../utils/alert";
import { useRouter } from 'vue-router';
import { difficultyColor, difficultyLang } from "../../utils/dynamic_styles";

const router = useRouter();
const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

const withPwdDialog = ref(false);
const noPwdDialog = ref(false);
const password = ref("");
const selectedId = ref();
const networkloading = ref(false);
const competitions = ref([]);
const loading = ref(false);

// 加入带密码的竞赛
const joinCompetitionWithPwd = async (competitionId) => {
    if (withPwdDialog.value) {
        if (password.value === "") {
            showAlert(t("message.passwordRequired"), '');
        } else {
            try {
                networkloading.value = true;
                const resp = await joinCompWithPwd(competitionId, password.value);
                networkloading.value = false;
                showAlert(resp.data.message, 'reload')
            } catch (error) {
                networkloading.value = false;
                showAlert(error.response.data.message, '')
            }
        }
    }
}

// 加入竞赛
const joinComp = async (competitionId) => {
    if (noPwdDialog.value) {
        try {
            networkloading.value = true;
            const resp = await joinCompetition(competitionId);
            networkloading.value = false;
            showAlert(resp.data.message, 'reload')
        } catch (error) {
            networkloading.value = false;
            showAlert(error.response.data.message, '')
        }
    }
}

// 选择竞赛并弹出对话框 
const selectCompetition = async (contest) => {
    selectedId.value = contest.contest_id;

    // 检查用户是否在该竞赛中
    try {
        networkloading.value = true;
        const resp = await isInCompetition(selectedId.value);
        networkloading.value = false;
        if (resp.data.isIn) {
            await router.push({path: `/competitions/${selectedId.value}`})
        } else {
            contest.have_password ? withPwdDialog.value = true : noPwdDialog.value = true;
        }
    } catch (error) {
        showAlert(error.response.data.message, '')
    }

}

// 关闭对话框
const closeDialog = () => {
    withPwdDialog.value = false;
    noPwdDialog.value = false;
    password.value = "";
    selectedId.value = "";
}

onMounted(async () => {
    if (!userLoggedIn.value) {
        loading.value = false;
        setTimeout(() => {
            window.location = "#/login";
            window.location.reload();
        }, 2000);
        return;
    }
    loading.value = true;
    const response = await getAllCompetitions();
    competitions.value = response.data.data;
    loading.value = false;
});
</script>

<template>
    <v-dialog v-model="networkloading" max-width="500px">
        <v-card rounded=xl>
            <div class="networkloading">
                <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
            </div>
        </v-card>
    </v-dialog>
    <template>
        <v-dialog v-model="withPwdDialog" max-width="500px">
            <v-card rounded=xl>
                <v-card-title class="text-h5">{{ t('message.joinCompetition') }}</v-card-title>
                <v-card-subtitle>{{ t('message.followRules') }}</v-card-subtitle>
                <v-card-text>
                    <v-text-field v-model="password" :label="t('message.password')" type="password" rounded="xl"
                        variant="solo-filled"></v-text-field>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" rounded="xl" text @click="closeDialog">{{ t("message.cancel")
                        }}</v-btn>
                    <v-btn color="primary" rounded="xl" @click="joinCompetitionWithPwd(selectedId)">{{ t("message.ok")
                        }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </template>
    <template>
        <v-dialog v-model="noPwdDialog" max-width="500px">
            <v-card rounded=xl>
                <v-card-title class="text-h5">{{ t('message.joinCompetition') }}</v-card-title>
                <v-card-subtitle>{{ t('message.followRules') }}</v-card-subtitle>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" rounded="xl" text @click="closeDialog">{{ t("message.cancel")
                        }}</v-btn>
                    <v-btn color="primary" rounded="xl" @click="joinComp(selectedId)">{{ t("message.ok") }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </template>
    <div v-if="loading" class="loading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
    </div>
    <div v-else>
        <v-app-bar :elevation="2">
            <p style="font-size: 24px;margin-left: 20px;">{{ t("message.competition") }}</p>
        </v-app-bar>
        <div v-if="userLoggedIn">
            <div v-if="competitions.length === 0" class="empty-state">
                <v-icon size="120" color="grey-lighten-1" class="empty-icon">mdi-trophy-outline</v-icon>
                <h2 class="empty-title">{{ t("message.noCompetitions") }}</h2>
                <p class="empty-description">{{ t("message.noCompetitionsDesc") }}</p>
            </div>
            <v-container>
                <v-row>
                    <v-col v-for="contest in competitions" :key="contest.contest_id" cols="12" md="4">
                        <v-card rounded="xl" elevation="2" style="display: grid;" class="comp-card">
                            <v-card-title style="font-weight: bold;font-size:28px;justify-self: left;">{{ contest.title
                                }}</v-card-title>
                            <v-card-subtitle style="justify-self: left;">{{ contest.subtitle }}</v-card-subtitle>
                            <template v-slot:append>
                                <v-chip :style="difficultyColor(contest.difficulty)">
                                    {{ t(difficultyLang(contest.difficulty)) }}
                                </v-chip>
                            </template>
                            <v-card-text style="display: grid;">
                                <p style="justify-self: left;">
                                    {{ moment(contest.start_at).format("YYYY/MM/DD HH:mm") }} -
                                    {{ moment(contest.end_at).format("YYYY/MM/DD HH:mm") }}
                                </p>
                            </v-card-text>
                            <template v-slot:actions>
                                <v-btn color="primary" rounded="xl" append-icon="mdi-chevron-right"
                                    @click="selectCompetition(contest)">{{ $t("message.enter")
                                    }}</v-btn>
                            </template>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
        <div v-else>
            <div class="empty-state">
                <v-icon size="120" color="grey-lighten-1" class="empty-icon">mdi-login</v-icon>
                <h2 class="empty-title">{{ t("message.nologin") }}</h2>
            </div>
        </div>
    </div>
</template>

<style scoped>
.loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.comp-card {
  transition: all 0.3s ease;
  border: 1px solid var(--v-border-color);
}

.comp-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    padding: 40px 20px;
    text-align: center;
}

.empty-icon {
    margin-bottom: 24px;
    opacity: 0.6;
}

.empty-title {
    font-size: 28px;
    font-weight: 600;
    color: var(--v-on-surface-variant);
    margin-bottom: 16px;
    line-height: 1.2;
}

.empty-description {
    font-size: 16px;
    color: var(--v-on-surface-variant);
    opacity: 0.8;
    max-width: 400px;
    line-height: 1.5;
}

@media (max-width: 600px) {
    .empty-state {
        min-height: 50vh;
        padding: 20px 16px;
    }
    
    .empty-icon {
        font-size: 80px !important;
        margin-bottom: 20px;
    }
    
    .empty-title {
        font-size: 24px;
        margin-bottom: 12px;
    }
    
    .empty-description {
        font-size: 14px;
    }
}
</style>
