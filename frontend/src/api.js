// Thin fetch wrapper around the backend REST API (see API_CONTRACT.md).
// Base path is proxied by Vite in dev to http://localhost:8080.

const API_BASE = '/api'

export class ApiError extends Error {
  constructor(message, status = 0, fields = {}) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.fields = fields || {}
  }
}

const unauthorizedHandlers = new Set()

// Register a callback fired whenever the API answers 401.
export function onUnauthorized(handler) {
  unauthorizedHandlers.add(handler)
  return () => unauthorizedHandlers.delete(handler)
}

function buildUrl(path, query) {
  let url = API_BASE + path
  if (query) {
    const params = new URLSearchParams()
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') params.set(key, value)
    })
    const qs = params.toString()
    if (qs) url += (url.includes('?') ? '&' : '?') + qs
  }
  return url
}

function defaultMessage(status) {
  if (status === 401) return 'Sesi Anda telah berakhir. Silakan masuk kembali.'
  if (status === 403) return 'Anda tidak memiliki hak akses untuk tindakan ini.'
  if (status === 404) return 'Data tidak ditemukan.'
  if (status === 409) return 'Data sudah digunakan / bentrok dengan data lain.'
  if (status >= 500) return 'Terjadi kesalahan pada server. Silakan coba lagi.'
  return 'Permintaan gagal diproses.'
}

export async function request(path, { method = 'GET', body, query } = {}) {
  const options = {
    method,
    credentials: 'include',
    headers: {}
  }
  if (body !== undefined) {
    options.headers['Content-Type'] = 'application/json'
    options.body = JSON.stringify(body)
  }

  let response
  try {
    response = await fetch(buildUrl(path, query), options)
  } catch (err) {
    throw new ApiError('Tidak dapat terhubung ke server. Periksa koneksi Anda.', 0, {})
  }

  const raw = await response.text()
  let data = null
  if (raw) {
    try {
      data = JSON.parse(raw)
    } catch (err) {
      data = null
    }
  }

  if (!response.ok) {
    if (response.status === 401) {
      unauthorizedHandlers.forEach((handler) => handler())
    }
    const message = (data && data.error) || defaultMessage(response.status)
    const fields = (data && data.fields) || {}
    throw new ApiError(message, response.status, fields)
  }

  return data
}

export const api = {
  login: (username, password) => request('/login', { method: 'POST', body: { username, password } }),
  logout: () => request('/logout', { method: 'POST' }),
  me: () => request('/me'),

  wilayah: () => request('/wilayah'),

  aktivitas: {
    list: () => request('/aktivitas'),
    create: (nama) => request('/aktivitas', { method: 'POST', body: { nama } }),
    update: (id, payload) => request(`/aktivitas/${encodeURIComponent(id)}`, { method: 'PUT', body: payload }),
    remove: (id) => request(`/aktivitas/${encodeURIComponent(id)}`, { method: 'DELETE' })
  },

  rekrutmen: {
    list: (query) => request('/rekrutmen', { query }),
    get: (id) => request(`/rekrutmen/${encodeURIComponent(id)}`),
    create: (body) => request('/rekrutmen', { method: 'POST', body }),
    update: (id, body) => request(`/rekrutmen/${encodeURIComponent(id)}`, { method: 'PUT', body }),
    remove: (id) => request(`/rekrutmen/${encodeURIComponent(id)}`, { method: 'DELETE' }),
    checkWa: (no) => request('/rekrutmen/check-wa', { query: { no } })
  },

  settings: {
    get: () => request('/settings'),
    update: (provider) => request('/settings', { method: 'PUT', body: { provider } })
  }
}

export default api
