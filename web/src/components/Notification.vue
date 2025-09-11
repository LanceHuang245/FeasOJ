<script setup>
import { useI18n } from 'vue-i18n';
import { useNotificationState, closeNotification } from '../utils/notification.js';
import { initSSE } from '../utils/notification.js';
import { userId, token } from '../utils/account.js';
import { onMounted, computed } from 'vue';

const { t } = useI18n();
const { snackbar, snackbarMessage } = useNotificationState();
const userLoggedIn = computed(() => !!token.value);

onMounted(() => {
  if (userLoggedIn.value) {
    initSSE(userId.value);
  }
});
const timeout = 4000;
</script>

<template>
    <v-snackbar v-model="snackbar" :timeout="timeout" bottom>
        {{ snackbarMessage }}
        <template v-slot:actions>
            <v-btn color="blue" @click="closeNotification()">
                {{ t('message.ok') }}
            </v-btn>
        </template>
    </v-snackbar>
</template>