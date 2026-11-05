import { createRouter, createWebHistory } from 'vue-router'
import { auth, loadUser, isAdmin } from './auth'

import LoginView from './views/LoginView.vue'
import ListView from './views/ListView.vue'
import FormView from './views/FormView.vue'
import AktivitasView from './views/AktivitasView.vue'
import SettingsView from './views/SettingsView.vue'

const routes = [
  { path: '/login', name: 'login', component: LoginView, meta: { public: true, title: 'Masuk' } },
  { path: '/', name: 'list', component: ListView, meta: { title: 'Data Rekrutmen' } },
  { path: '/tambah', name: 'tambah', component: FormView, meta: { title: 'Tambah Rekrutmen' } },
  { path: '/ubah/:id', name: 'ubah', component: FormView, meta: { title: 'Ubah Rekrutmen' } },
  { path: '/aktivitas', name: 'aktivitas', component: AktivitasView, meta: { title: 'Master Aktivitas' } },
  { path: '/pengaturan', name: 'pengaturan', component: SettingsView, meta: { title: 'Pengaturan', admin: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach(async (to) => {
  if (to.meta.public) {
    // Already signed in users should not stay on the login screen.
    if (auth.user) return { path: '/' }
    return true
  }

  const user = await loadUser()
  if (!user) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.admin && !isAdmin()) {
    return { path: '/' }
  }
  return true
})

router.afterEach((to) => {
  const title = to.meta.title ? `${to.meta.title} · Rekrutmen Tangsel` : 'Rekrutmen Tangsel'
  document.title = title
})

export default router
