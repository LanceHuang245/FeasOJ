<script setup>
import { dialog, dialogMessage } from '../utils/alert.js';
</script>

<template>
    <v-dialog v-model="dialog" persistent max-width="290">
        <v-card rounded="xl">
            <v-card-title class="text-h5">{{$t('message.notify')}}</v-card-title>
            <v-card-text>{{ dialogMessage }}</v-card-text>
            <v-card-actions>
                <v-btn color="blue darken-1" @click="dialog = false" rounded="xl">{{$t('message.ok')}}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<style scoped></style>