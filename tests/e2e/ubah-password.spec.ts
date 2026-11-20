import { test, expect } from '@playwright/test'

// Checks of the change of the own password, against a running server that
// serves the built frontend.
//
// The accounts used here are test accounts, so the password of the built-in
// accounts stays untouched for the other checks.
const BASE = process.env.E2E_BASE_URL || 'http://localhost:8100'

const ACCOUNTS = [
  { username: 'pwpetugas', password: 'pwpetugas-anggota-2024', label: 'anggota' },
  { username: 'pwkec', password: 'pwkec-kecamatan-2024', label: 'admin kecamatan' },
  { username: 'pwadmin', password: 'pwadmin-sekolah-2024', label: 'admin penuh' }
]

async function login(page, username: string, password: string) {
  await page.goto(`${BASE}/login`)
  await page.locator('#login-username').fill(username)
  await page.locator('#login-password').fill(password)
  await page.getByRole('button', { name: 'Masuk' }).click()
  await page.waitForURL((url) => !url.pathname.includes('/login'))
  await page.waitForLoadState('networkidle')
}

async function openPasswordDialog(page) {
  await page.getByRole('button', { name: 'Password' }).click()
  await expect(page.getByRole('dialog', { name: 'Ubah password' })).toBeVisible()
}

test.describe.configure({ mode: 'serial' })

test('Tombol Password tersedia untuk semua level akun', async ({ page }) => {
  for (const account of ACCOUNTS) {
    await login(page, account.username, account.password)
    await expect(page.getByRole('button', { name: 'Password' })).toBeVisible()
    await page.goto(`${BASE}/login`)
  }
})

test('Dialog mengubah password menampilkan tiga field', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await expect(page.locator('#password-current')).toBeVisible()
  await expect(page.locator('#password-new')).toBeVisible()
  await expect(page.locator('#password-confirm')).toBeVisible()

  // The password is never shown while it is typed.
  await expect(page.locator('#password-current')).toHaveAttribute('type', 'password')
  await expect(page.locator('#password-new')).toHaveAttribute('type', 'password')
})

test('Konfirmasi password yang tidak sama ditolak di sisi klien', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await page.locator('#password-current').fill('pwpetugas-anggota-2024')
  await page.locator('#password-new').fill('baru-sekarang-2024')
  await page.locator('#password-confirm').fill('beda-sekarang-2024')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(page.locator('.modal .alert--error')).toContainText('Konfirmasi password tidak sama.')
  await expect(page.getByRole('dialog', { name: 'Ubah password' })).toBeVisible()
})

test('Password lama yang salah ditolak oleh server', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await page.locator('#password-current').fill('password-salah')
  await page.locator('#password-new').fill('baru-sekarang-2024')
  await page.locator('#password-confirm').fill('baru-sekarang-2024')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(page.locator('.modal .alert--error')).toContainText('password lama salah')
})

test('Password baru terlalu pendek ditolak di sisi klien', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await page.locator('#password-current').fill('pwpetugas-anggota-2024')
  await page.locator('#password-new').fill('abc')
  await page.locator('#password-confirm').fill('abc')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(page.locator('.modal .alert--error')).toContainText('minimal 12 karakter')
})

test('Dialog bisa ditutup dan isiannya dikosongkan', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await page.locator('#password-current').fill('pwpetugas-anggota-2024')
  await page.locator('#password-new').fill('pwpetugas456')
  await page.getByRole('button', { name: 'Batal' }).click()
  await expect(page.getByRole('dialog', { name: 'Ubah password' })).toHaveCount(0)

  await openPasswordDialog(page)
  await expect(page.locator('#password-current')).toHaveValue('')
  await expect(page.locator('#password-new')).toHaveValue('')
  await expect(page.locator('#password-confirm')).toHaveValue('')
})

test('Password berhasil diganti dan dipakai untuk masuk lagi', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas-anggota-2024')
  await openPasswordDialog(page)

  await page.locator('#password-current').fill('pwpetugas-anggota-2024')
  await page.locator('#password-new').fill('pwpetugas456')
  await page.locator('#password-confirm').fill('pwpetugas456')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(page.locator('.toast--success')).toContainText('Password berhasil diubah')
  await expect(page.getByRole('dialog', { name: 'Ubah password' })).toHaveCount(0)

  // The session of the account stays usable.
  await page.reload()
  await expect(page.getByRole('button', { name: 'Password' })).toBeVisible()

  // The new password works, the old one does not.
  await page.getByRole('button', { name: 'Keluar' }).click()
  await page.waitForURL((url) => url.pathname.includes('/login'))
  await page.locator('#login-username').fill('pwpetugas')
  await page.locator('#login-password').fill('pwpetugas-anggota-2024')
  await page.getByRole('button', { name: 'Masuk' }).click()
  await expect(page.locator('#login-username')).toBeVisible()
  await expect(page.locator('.alert--error')).toContainText('username atau password salah')

  // The test account is restored, so a repeated run starts from the same
  // state.
  await page.locator('#login-username').fill('pwpetugas')
  await page.locator('#login-password').fill('pwpetugas456')
  await page.getByRole('button', { name: 'Masuk' }).click()
  await page.waitForURL((url) => !url.pathname.includes('/login'))
  await openPasswordDialog(page)
  await page.locator('#password-current').fill('pwpetugas456')
  await page.locator('#password-new').fill('pwpetugas-anggota-2024')
  await page.locator('#password-confirm').fill('pwpetugas-anggota-2024')
  await page.getByRole('button', { name: 'Simpan' }).click()
  await expect(page.locator('.toast--success')).toContainText('Password berhasil diubah')
})
