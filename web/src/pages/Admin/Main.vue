<script setup>
import { useRouter } from 'vue-router';
import { verifyUserInfo } from '../../utils/api/auth';
import { token, userName } from '../../utils/account'
import { ref, computed, onMounted } from 'vue';

const router = useRouter();
const loading = ref(true);
const userPrivilege = ref("")

// 计算属性来判断用户是否已经登录
const userLoggedIn = computed(() => !!token.value)

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
        loading.value = false;
    } catch (error) {
        window.location = '#/403';
    }
});
</script>

<template>
    <v-app-bar :elevation="2">
        <p style="font-size: 24px;margin-left: 20px;">{{ $t("message.management") }}</p>
    </v-app-bar>
    <div v-if="loading" class="loading">
        <v-progress-circular indeterminate color="primary" :width="12" :size="100"></v-progress-circular>
    </div>
    <div v-else>
        <VRow justify="center" class="my-5">
            <VCol cols="3">
                <VCard class="d-flex align-center justify-center bold" variant="outlined" height="15vh"
                    @click="router.push('/admin/user')">
                    {{ $t('message.usermanagement') }}
                </VCard>
            </VCol>
            <VCol cols="3">
                <VCard class="d-flex align-center justify-center bold" variant="outlined" height="15vh"
                    @click="router.push('/admin/problem')">
                    {{ $t('message.problemmanagement') }}
                </VCard>
            </VCol>
            <VCol cols="3">
                <VCard class="d-flex align-center justify-center bold" variant="outlined" height="15vh"
                    @click="router.push('/admin/competition')">
                    {{ $t('message.competitionmanagement') }}
                </VCard>
            </VCol>
        </VRow>
        <VRow justify="center" class="my-5">
            <VCol cols="3">
                <VCard class="d-flex align-center justify-center bold" variant="outlined" height="15vh"
                    @click="router.push('/admin/ipstat')">
                    {{ $t('message.ipstat') }}
                </VCard>
            </VCol>
            <VCol cols="3"></VCol>
            <VCol cols="3"></VCol>
        </VRow>
    </div>
</template>

<style scoped>
.loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.bold {
    font-weight: bold;
}
</style>