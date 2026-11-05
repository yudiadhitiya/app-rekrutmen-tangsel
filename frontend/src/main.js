import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { onUnauthorized } from './api'
import './style.css'

const app = createApp(App)
app.use(router)

// Any 401 from the API sends the user back to the login page while keeping
// the originally intended destination in the `redirect` query.
onUnauthorized(() => {
  const current = router.currentRoute.value
  if (current.path === '/login') return
  router.replace({ path: '/login', query: { redirect: current.fullPath } })
})

app.mount('#app')
