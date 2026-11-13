<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { auth, isFullAdmin, roleLabel, logout } from './auth'
import Toast from './components/Toast.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'

const route = useRoute()
const router = useRouter()

const isPublicPage = computed(() => !!route.meta.public)
const admin = computed(() => isFullAdmin())

const navLinks = computed(() => {
  const links = [
    { to: '/', label: 'Data Rekrutmen' },
    { to: '/tambah', label: 'Tambah' },
    { to: '/aktivitas', label: 'Aktivitas' }
  ]
  if (admin.value) links.push({ to: '/pengaturan', label: 'Pengaturan' })
  return links
})

const userRoleLabel = computed(() => roleLabel())

async function handleLogout() {
  await logout()
  router.replace({ path: '/login' })
}
</script>

<template>
  <div class="app-shell">
    <header v-if="!isPublicPage" class="app-bar">
      <div class="app-bar__inner">
        <RouterLink to="/" class="app-bar__brand">
          <span class="app-bar__logo">RT</span>
          <span>Rekrutmen Tangsel</span>
        </RouterLink>

        <nav class="app-nav" aria-label="Navigasi utama">
          <RouterLink v-for="link in navLinks" :key="link.to" :to="link.to" class="app-nav__link">
            {{ link.label }}
          </RouterLink>
        </nav>

        <div class="app-bar__user">
          <div class="app-bar__identity">
            <span class="app-bar__name">{{ auth.user ? auth.user.nama || auth.user.username : '' }}</span>
            <span class="app-bar__role">{{ userRoleLabel }}</span>
          </div>
          <button class="btn btn--secondary btn--sm" type="button" @click="handleLogout">Keluar</button>
        </div>
      </div>
    </header>

    <main class="app-main" :class="{ 'app-main--full': isPublicPage }">
      <RouterView />
    </main>

    <Toast />
    <ConfirmDialog />
  </div>
</template>

<style scoped>
.app-main--full {
  max-width: none;
  padding: 0;
}
</style>
