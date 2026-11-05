// Promise based confirm dialog, driven by a single global modal component.

import { reactive } from 'vue'

export const confirmState = reactive({
  open: false,
  title: 'Konfirmasi',
  message: '',
  confirmText: 'Ya, lanjutkan',
  cancelText: 'Batal',
  tone: 'danger',
  resolve: null
})

export function confirmDialog(options = {}) {
  return new Promise((resolve) => {
    confirmState.title = options.title || 'Konfirmasi'
    confirmState.message = options.message || ''
    confirmState.confirmText = options.confirmText || 'Ya, lanjutkan'
    confirmState.cancelText = options.cancelText || 'Batal'
    confirmState.tone = options.tone || 'danger'
    confirmState.resolve = resolve
    confirmState.open = true
  })
}

export function resolveConfirm(result) {
  confirmState.open = false
  const resolve = confirmState.resolve
  confirmState.resolve = null
  if (resolve) resolve(result)
}
