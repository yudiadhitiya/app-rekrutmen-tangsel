// Session state (plain reactive store, no Pinia).

import { reactive } from 'vue'
import { api } from './api'

export const auth = reactive({
  user: null,
  loaded: false,
  loading: false
})

export function isAdmin() {
  return !!(auth.user && auth.user.role === 'admin')
}

// The full administrator is the account of level 3. It is the only one that
// may manage the master data, the settings, and delete a record. Accounts
// without a level fall back to the role, so older sessions keep working.
export function isFullAdmin() {
  if (!auth.user) return false
  if (auth.user.level) return auth.user.level >= 3
  return auth.user.role === 'admin'
}

// The label of the account, which follows its access level.
export function roleLabel() {
  if (!auth.user) return ''
  if (auth.user.level === 3) return 'Admin'
  if (auth.user.level === 2) return 'Admin Kecamatan'
  return auth.user.role === 'admin' ? 'Admin' : 'Petugas'
}

// A member only reaches the records it created itself, so an empty list
// means it has not entered anything yet.
export function isMemberOnly() {
  if (!auth.user) return false
  if (auth.user.level) return auth.user.level < 2
  return auth.user.role !== 'admin'
}

export async function loadUser(force = false) {
  if (auth.loaded && !force) return auth.user
  auth.loading = true
  try {
    const data = await api.me()
    auth.user = data && data.user ? data.user : null
  } catch (err) {
    auth.user = null
  } finally {
    auth.loaded = true
    auth.loading = false
  }
  return auth.user
}

export async function login(username, password) {
  const data = await api.login(username, password)
  auth.user = data && data.user ? data.user : null
  auth.loaded = true
  // The response also reports whether the account was migrated from the
  // previous application, so the interface can tell the user about it.
  return data
}

export async function logout() {
  try {
    await api.logout()
  } catch (err) {
    // Ignore logout failures - the local session is dropped anyway.
  }
  auth.user = null
  auth.loaded = true
}
