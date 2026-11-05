// Tiny toast bus - no external toast library.

import { reactive } from 'vue'

let nextId = 1

export const toasts = reactive([])

export function showToast(message, type = 'success', timeout = 4000) {
  const id = nextId++
  toasts.push({ id, message, type })
  if (timeout > 0) {
    setTimeout(() => dismissToast(id), timeout)
  }
  return id
}

export function dismissToast(id) {
  const index = toasts.findIndex((toast) => toast.id === id)
  if (index !== -1) toasts.splice(index, 1)
}

export const toast = {
  success: (message) => showToast(message, 'success'),
  error: (message) => showToast(message, 'error', 6000),
  warning: (message) => showToast(message, 'warning', 5000),
  info: (message) => showToast(message, 'info')
}

export default toast
