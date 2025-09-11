<!-- 重置密码页 -->
<script setup>
import { ref, reactive } from 'vue';
import { updatePassword, getCaptchaCode } from '../../utils/api/auth';
import { regex, rules } from '../../utils/rules'
import { showAlert } from '../../utils/alert';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const forms = reactive({
    email: '',
    vcode: '',
    password: '',
    confirmPassword: ''
});

const isButtonDisabled = ref(false);
const countdown = ref(60);

const nextStep = async () => {
    if (forms.email === "" || forms.password === "" || forms.confirmPassword === "" || forms.vcode === "") {
        showAlert(t("message.formCheckfailed") + "!", "");
        return;
    }
    if (forms.password.length < 8) {
        showAlert(t("message.formRuleCheckfailed") + "!", "");
        return;
    }
    if (forms.password !== forms.confirmPassword) {
        showAlert(t("message.formRuleCheckfailed") + "!", "");
        return;
    }
    try {
        const response = await updatePassword(forms.email, forms.vcode, forms.password);
        localStorage.clear();
        showAlert(response.data.message, "/login");
        return;
    } catch (error) {
        showAlert(error.response.data.message, "");
        return;
    }
}

const getCaptcha = async () => {
    if (!regex.test(forms.email)) {
        showAlert(t("message.formRuleCheckfailed") + "!", "");
        return;
    }
    if (!forms.email) {
        showAlert(t("message.formCheckfailed") + "!", "");
        return;
    }
    try {
        const response = await getCaptchaCode(forms.email, "false");
        showAlert(response.data.message, "");
        if (isButtonDisabled.value) {
            return;
        }
        isButtonDisabled.value = true;
        let timer = setInterval(() => {
            if (countdown.value > 0) {
                countdown.value--;
            } else {
                clearInterval(timer);
                isButtonDisabled.value = false;
                countdown.value = 60; // 重置倒计时
            }
        }, 1000);
    } catch (error) {
        showAlert(error.response.data.message, "");
    }
}
</script>

<template>
    <v-app-bar :elevation="2">
        <template v-slot:prepend>
            <v-btn icon="mdi-chevron-left" size="x-large" @click="$router.back"></v-btn>
            <v-col class="align-left">
                <p style="font-size: 24px;">{{ t("message.resetpwd") }}</p>
            </v-col>
        </template>
    </v-app-bar>
    <v-sheet class="constrainsheet" rounded="xl" :elevation="10">
        <v-form fast-fail width="400px" class="mx-auto" @submit.prevent="nextStep" style="margin: 20px;">
            <v-text-field v-model="forms.email" :rules="[rules.userEmail.required, rules.userEmail.email]" rounded="xl"
                variant="solo-filled" :label="$t('message.email')" />
            <v-text-field v-model="forms.vcode" :rules="[rules.vcode.required]" rounded="xl" variant="solo-filled"
                :label="$t('message.vCode')">
                <template v-slot:append>
                    <v-btn @click="getCaptcha" :disabled="isButtonDisabled" size="25" icon>
                        <v-icon v-if="!isButtonDisabled" icon="mdi-email" />
                        <span v-if="isButtonDisabled">{{ countdown }}</span>
                    </v-btn>
                </template>
            </v-text-field>
            <v-text-field v-model="forms.password" :rules="[rules.password.required, rules.password.minLength]"
                rounded="xl" variant="solo-filled" type="password" :label="$t('message.password')" />
            <v-text-field v-model="forms.confirmPassword"
                :rules="[rules.confirmPassword.required, rules.confirmPassword.minLength]" rounded="xl"
                variant="solo-filled" type="password" :label="$t('message.confirmPwd')" />
            <v-btn type="submit" color="primary" rounded="xl">{{ $t('message.submit') }}</v-btn>
        </v-form>
    </v-sheet>
</template>