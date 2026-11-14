import { test, expect } from '@playwright/test'
import path from 'path'

// A first check of the browser setup: Chromium must be able to open a page,
// read its title, and take a screenshot.
test('Chromium membuka example.com dan mengambil tangkapan layar', async ({ page }) => {
  const response = await page.goto('https://example.com', { waitUntil: 'domcontentloaded' })

  expect(response).not.toBeNull()
  expect(response!.status()).toBe(200)

  await expect(page).toHaveTitle(/Example Domain/)

  const screenshot = path.join(__dirname, 'browser-check.png')
  await page.screenshot({ path: screenshot, fullPage: true })
})
