<script setup>
// 跳转至beian.miit.gov.cn
const goToICP = () => {
  window.open("https://beian.miit.gov.cn/");
};
</script>

<!-- TODO: 若需要将站点部署于中国境内服务器，请修改此组件，并在App.vue中启用该组件！ -->
 
<template>
    <v-btn variant="plain" rounded="xl" style="font-size: 12px" @click="goToICP">闽ICP备**********</v-btn>
</template>