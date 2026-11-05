<script setup>
import { confirmState, resolveConfirm } from '../confirm'
</script>

<template>
  <Teleport to="body">
    <div v-if="confirmState.open" class="modal-backdrop" @click.self="resolveConfirm(false)">
      <div class="modal" role="dialog" aria-modal="true" :aria-label="confirmState.title">
        <h3 class="modal__title">{{ confirmState.title }}</h3>
        <p class="modal__message">{{ confirmState.message }}</p>
        <div class="modal__actions">
          <button class="btn btn--secondary" type="button" @click="resolveConfirm(false)">
            {{ confirmState.cancelText }}
          </button>
          <button
            class="btn"
            :class="confirmState.tone === 'danger' ? 'btn--danger' : 'btn--primary'"
            type="button"
            @click="resolveConfirm(true)"
          >
            {{ confirmState.confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
