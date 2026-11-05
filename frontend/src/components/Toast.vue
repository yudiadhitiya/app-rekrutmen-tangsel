<script setup>
import { toasts, dismissToast } from '../toast'

const TITLES = {
  success: 'Berhasil',
  error: 'Gagal',
  warning: 'Perhatian',
  info: 'Informasi'
}
</script>

<template>
  <div class="toast-stack" role="status" aria-live="polite">
    <div v-for="item in toasts" :key="item.id" class="toast" :class="`toast--${item.type}`">
      <div class="toast__body">
        <div class="toast__title">{{ TITLES[item.type] || 'Informasi' }}</div>
        <div>{{ item.message }}</div>
      </div>
      <button class="toast__close" type="button" aria-label="Tutup notifikasi" @click="dismissToast(item.id)">
        &times;
      </button>
    </div>
  </div>
</template>
