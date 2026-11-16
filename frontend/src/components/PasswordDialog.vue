<script setup>
import { ref, watch } from 'vue'
import { api } from '../api'
import { toast } from '../toast'

// The dialog is opened by the header of the application, so every signed in
// account may change its own password.
const props = defineProps({
  open: { type: Boolean, default: false }
})
const emit = defineEmits(['close'])

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const errorMessage = ref('')
const submitting = ref(false)

watch(
  () => props.open,
  (open) => {
    if (open) {
      currentPassword.value = ''
      newPassword.value = ''
      confirmPassword.value = ''
      errorMessage.value = ''
    }
  }
)

async function handleSubmit() {
  errorMessage.value = ''
  if (!currentPassword.value) {
    errorMessage.value = 'Password lama wajib diisi.'
    return
  }
  if (!newPassword.value) {
    errorMessage.value = 'Password baru wajib diisi.'
    return
  }
  if (newPassword.value.length < 6) {
    errorMessage.value = 'Password baru minimal 6 karakter.'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Konfirmasi password tidak sama.'
    return
  }
  if (newPassword.value === currentPassword.value) {
    errorMessage.value = 'Password baru harus berbeda dari password lama.'
    return
  }

  submitting.value = true
  try {
    await api.password.change(currentPassword.value, newPassword.value)
    toast.success('Password berhasil diubah. Gunakan password baru saat masuk berikutnya.')
    emit('close')
  } catch (err) {
    errorMessage.value = err.message || 'Gagal mengubah password. Silakan coba lagi.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="modal-backdrop" @click.self="emit('close')">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Ubah password">
        <h3 class="modal__title">Ubah Password</h3>
        <p class="modal__message">Password baru berlaku juga di aplikasi sebelumnya.</p>

        <form @submit.prevent="handleSubmit">
          <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

          <div class="field">
            <label class="field__label" for="password-current">Password Lama</label>
            <input
              id="password-current"
              v-model="currentPassword"
              class="input"
              type="password"
              autocomplete="current-password"
              placeholder="Masukkan password lama"
              :disabled="submitting"
            />
          </div>

          <div class="field">
            <label class="field__label" for="password-new">Password Baru</label>
            <input
              id="password-new"
              v-model="newPassword"
              class="input"
              type="password"
              autocomplete="new-password"
              placeholder="Minimal 6 karakter"
              :disabled="submitting"
            />
          </div>

          <div class="field">
            <label class="field__label" for="password-confirm">Ulangi Password Baru</label>
            <input
              id="password-confirm"
              v-model="confirmPassword"
              class="input"
              type="password"
              autocomplete="new-password"
              placeholder="Ulangi password baru"
              :disabled="submitting"
            />
          </div>

          <div class="modal__actions">
            <button class="btn btn--secondary" type="button" :disabled="submitting" @click="emit('close')">
              Batal
            </button>
            <button class="btn btn--primary" type="submit" :disabled="submitting">
              <span v-if="submitting" class="spinner" aria-hidden="true"></span>
              {{ submitting ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
