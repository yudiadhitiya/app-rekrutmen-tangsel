import { createRouter, createWebHistory } from 'vue-router'
import { auth, loadUser, isFullAdmin } from './auth'

import LoginView from './views/LoginView.vue'
import DashboardView from './views/DashboardView.vue'
import ListView from './views/ListView.vue'
import FormView from './views/FormView.vue'
import AktivitasView from './views/AktivitasView.vue'
import SettingsView from './views/SettingsView.vue'

const routes = [
  { path: '/login', name: 'login', component: LoginView, meta: { public: true, title: 'Masuk' } },
  { path: '/dashboard', name: 'dashboard', component: DashboardView, meta: { title: 'Dashboard' } },
  { path: '/', name: 'list', component: ListView, meta: { title: 'Data Rekrutmen' } },
  { path: '/tambah', name: 'tambah', component: FormView, meta: { title: 'Tambah Rekrutmen' } },
  { path: '/ubah/:id', name: 'ubah', component: FormView, meta: { title: 'Ubah Rekrutmen' } },
  { path: '/aktivitas', name: 'aktivitas', component: AktivitasView, meta: { title: 'Master Aktivitas', admin: true } },
  { path: '/pengaturan', name: 'pengaturan', component: SettingsView, meta: { title: 'Pengaturan', admin: true } },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach(async (to) => {
  if (to.meta.public) {
    // Already signed in users should not stay on the login screen.
    if (auth.user) return { path: '/dashboard' }
    return true
  }

  const user = await loadUser()
  if (!user) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.admin && !isFullAdmin()) {
    return { path: '/dashboard' }
  }
  return true
})

router.afterEach((to) => {
  const title = to.meta.title ? `${to.meta.title} · Rekrutmen Tangsel` : 'Rekrutmen Tangsel'
  document.title = title
})

export default router
