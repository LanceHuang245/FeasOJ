<!-- 讨论帖子详情页 -->
<script setup>
import { ref, onMounted, computed, onUnmounted } from "vue";
import { getDisDetails, deleteDiscussion, getComments, deleteComment, addComment } from "../../utils/api/discussions";
import { avatarServer } from "../../utils/axios";
import { useRoute, useRouter } from "vue-router";
import { showAlert } from "../../utils/alert";
import { token, userName } from "../../utils/account";
import { MdPreview, MdEditor } from "md-editor-v3";
import { useI18n } from "vue-i18n";
import moment from "moment";
import "md-editor-v3/lib/style.css";
import "md-editor-v3/lib/preview.css";
import { getMdEditorTheme, getMdPreviewTheme } from '../../utils/theme';

// 用于获取当前语言
const { locale } = useI18n();
const { t } = useI18n();

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value);

// 用于获取路由的参数
const route = useRoute();
const Did = route.params.Did;

// 用于路由处理
const router = useRouter();

const loading = ref(true);
const discussionInfos = ref({});
const id = "preview-only";
const page = ref(1);
const comments = ref([]);
const commentContent = ref("");
const expandedComments = ref({});
const checkDialog = ref(false);
const isDelComment = ref(false);
const deleteCommentID = ref("");

const editorTheme = ref(getMdEditorTheme());
const previewTheme = ref(getMdPreviewTheme());

// 监听主题变化
const handleThemeChange = (event) => {
    const theme = event.detail.theme === 'dark' ? 'dark' : 'light';
    editorTheme.value = theme;
    previewTheme.value = theme;
};

// 用作分页
const itemsPerPage = 5;
const totalPages = computed(() => {
    return Math.ceil(comments.value.length / itemsPerPage);
});

const paginatedComments = computed(() => {
    const start = (page.value - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    return comments.value.slice(start, end);
});

const onPageChange = (newPage) => {
    page.value = newPage;
};

// 添加评论
const addComments = async (content) => {
    loading.value = true;
    try {
        await addComment(Did, content);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        showAlert(error.response.data.message, "");
    } finally {
        loading.value = false;
    }
};

// 删除评论
const deleteCommentByID = async (commentID) => {
    loading.value = true;
    try {

        await deleteComment(commentID);
        showAlert(t("message.success") + "!", "reload");
    } catch (error) {
        window.location = "#/403";
    } finally {
        loading.value = false;
    }
};

// 删除讨论
const deleteDis = async () => {
    loading.value = true;
    try {
        await deleteDiscussion(Did);
        showAlert(t("message.success") + "!", "/discussion");
    } catch (error) {
        window.location = "#/403";
    } finally {
        loading.value = false;
    }
};

// 确认删除讨论
const showDelDialog = (isComment, commentID) => {
    checkDialog.value = true;
    if (isComment && commentID != null) {
        isDelComment.value = true;
        deleteCommentID.value = commentID;
    } else {
        isDelComment.value = false;
    }
}

// 确认删除
const deleteChecker = async () => {
    if (isDelComment.value) {
        await deleteCommentByID(deleteCommentID.value);
    } else {
        await deleteDis();
    }
    checkDialog.value = false;
}

onMounted(async () => {
    loading.value = true;
    if (userLoggedIn.value) {
        try {
            const response = await getDisDetails(Did);
            // 获取评论
            const commentsResp = await getComments(Did);
            comments.value = commentsResp.data.data;
            discussionInfos.value = response.data.data;
        } catch (error) {
            showAlert(t("message.failed") + "!", "/discussion");
        } finally {
            loading.value = false;
        }
    } else {
        window.location = "#/login";
        window.location.reload();
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
        <v-dialog v-model="checkDialog" persistent max-width="290">
            <v-card rounded="xl" elevation="2">
                <v-card-title class="text-h5">{{ $t('message.notify') }}</v-card-title>
                <v-card-text>{{ t('message.suredel') }}</v-card-text>
                <v-card-actions>
                    <v-btn variant="elevated" color="primary" @click="deleteChecker" rounded="xl">{{ $t('message.yes')
                    }}</v-btn>
                    <v-btn color="primary" @click="checkDialog = false" rounded="xl">{{ $t('message.cancel') }}</v-btn>
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
            <v-row style="align-items: center">
                <div style="margin-left: 30px"></div>
                <v-avatar size="50" color="surface-variant">
                    <v-img :src="avatarServer + discussionInfos.avatar" cover></v-img>
                </v-avatar>
                <v-btn variant="text" class="font-weight-black"
                    @click="router.push({ path: `/profile/${discussionInfos.username}` })">
                    {{ discussionInfos.username }}
                </v-btn>
            </v-row>
            <template v-if="discussionInfos.username === userName" v-slot:append>
                <v-btn icon="mdi-delete" size="x-large" @click="showDelDialog(false, null)"></v-btn>
            </template>
        </v-app-bar>
        <div style="margin-top: 30px"></div>
        <v-card class="mx-auto" width="75%" rounded="lg" elevation="2" style="display: grid;">
            <template v-slot:title>
                <span class="font-weight-black">{{ discussionInfos.title }}</span>
            </template>
            <MdPreview :editorId="id" :modelValue="discussionInfos.content" :theme="previewTheme" />
            <v-card-subtitle style="justify-self: right">
                <p style="font-size: 12px;margin: 10px">
                    {{ moment(discussionInfos.created_at).format("YYYY-MM-DD HH:mm") }}
                </p>
            </v-card-subtitle>
        </v-card>
        <div style="margin-top: 50px"></div>
        <v-card class="mx-auto" width="75%" rounded="lg" elevation="2">
            <template v-slot:title>
                <span class="font-weight-black">{{ $t("message.comments") }}</span>
            </template>
            <div style="max-height: 300px;">
                <md-editor v-model="commentContent" :editorId="id" :toolbars="[]" :noUploadImg="true" :preview="false"
                    :footers="[]" :language="locale === 'zh_CN' ? 'zh-CN' : 'en-US'" :theme="editorTheme" />
            </div>
            <div style="margin: 10px">
                <v-btn color="primary" rounded="xl" @click="addComments(commentContent)">{{ $t("message.submit")
                }}</v-btn>
            </div>
            <v-divider></v-divider>
            <v-list>
                <v-list-item v-for="comment in paginatedComments" :key="comment.cid">
                    <v-list-item>
                        <v-list-item-title class="username-avatar">
                            <v-avatar size="44" color="surface-variant">
                                <v-img :src="avatarServer + comment.avatar" alt="Avatar" cover></v-img>
                            </v-avatar>
                            <div style="margin: 5px">
                                <div class="username" @click="router.push({ path: `/profile/${comment.username}` })">
                                    {{ comment.username }}
                                </div>
                                <div class="timeline">
                                    {{ moment(comment.created_at).format("MM-DD HH:mm") }}
                                </div>
                            </div>
                        </v-list-item-title>
                        <v-chip color="warning" variant="tonal" class="ml-2"
                            v-if="comment.profanity && !expandedComments[comment.cid]" style="color: darkblue;"
                            @click="expandedComments[comment.cid] = true">
                            {{ $t("message.profanity_expand") }}
                        </v-chip>
                        <v-list-item v-else class="comment-content">
                            <md-preview :modelValue="comment.content" :theme="previewTheme" />
                        </v-list-item>
                    </v-list-item>
                    <div class="buttons">
                        <v-btn v-if="comment.username === userName" rounded="xl" variant="text" color="primary"
                            @click="showDelDialog(true, comment.id)">{{ $t("message.delete") }}</v-btn>
                    </div>
                    <v-divider></v-divider>
                </v-list-item>
            </v-list>
            <v-pagination v-model="page" :length="totalPages" @input="onPageChange" rounded="xl" />
        </v-card>
    </div>
</template>

<style scoped>
.buttons {
    display: flex;
    justify-content: flex-end;
    margin-top: 5px;
    margin-bottom: 10px;
}

.loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.username {
    font-weight: bold;
    font-size: 1em;
    cursor: pointer;
}

.comment-content {
    max-width: fit-content;
    max-height: 350px;
    text-align: left;
    margin-left: 0px;
}

.timeline {
    text-align: left;
    font-size: 0.8em;
    color: #999;
}

.username-avatar {
    text-align: left;
    display: flex;
    align-items: center;
}
</style>
