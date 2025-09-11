<script setup>
const back = () => {
  window.location = '#/'
};

const logout = () => {
  localStorage.clear();
  window.location = '#/login';
  location.reload();
};
</script>

<template>
  <div style="margin: 10%"></div>
  <v-card class="mx-auto" max-width="50%" min-width="50%" rounded="xl" elevation="10">
    <v-empty-state headline="403" :title="$t('message.nopermission')"></v-empty-state>
    <v-btn color="primary" variant="text" rounded="xl" @click="back">{{ $t('message.back') }}</v-btn>
    <v-btn color="primary" variant="text" rounded="xl" @click="logout">{{ $t('message.logout') }}</v-btn>
    <div style="margin: 10px"></div>
  </v-card>
</template>