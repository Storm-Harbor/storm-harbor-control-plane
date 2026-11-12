import {expect, test} from '@playwright/test';
import {resolve, sep} from 'node:path';

const preferenceKey = 'dr-docs-primary-navigation-collapsed';
const primary = page => page.locator('.md-sidebar--primary');
const menu = page => page.locator('[data-dr-nav-toggle]');

// Serve the real built artifact at its configured origin and Pages prefix.
// Material's instant navigation follows canonical URLs, so a localhost origin
// would exercise full reloads instead of the deployed navigation behavior.
test.beforeEach(async ({context, baseURL}) => {
  await context.route(`${baseURL}**`, async route => {
    const url = new URL(route.request().url());
    const path = url.pathname.slice(new URL(baseURL).pathname.length);
    const root = resolve('site');
    const file = resolve(root, decodeURIComponent(path + (path.endsWith('/') || !path ? 'index.html' : '')));
    expect(file.startsWith(root + sep)).toBe(true);
    await route.fulfill({path: file});
  });
});

async function expanded(page, value) {
  await expect(menu(page)).toHaveAttribute('aria-expanded', String(value));
  await expect.poll(() => primary(page).evaluate(el => el.inert)).toBe(!value);
  if (value) {
    await expect(primary(page)).toBeVisible();
    await expect.poll(() => primary(page).evaluate(el => el.getBoundingClientRect().x >= 0)).toBe(true);
  }
}

test('desktop home opens the hidden sidebar with keyboard and restores preference', async ({page}) => {
  await page.setViewportSize({width: 1440, height: 900});
  await page.goto('');
  await expanded(page, false);
  await expect(menu(page)).toHaveAccessibleName('Open navigation');
  await menu(page).focus();
  await page.keyboard.press('Enter');
  await expanded(page, true);
  await expect(primary(page).getByRole('link', {name: 'Product', exact: true})).toBeVisible();
  await page.keyboard.press('Space');
  await expanded(page, false);
  await expect.poll(() => page.evaluate(key => localStorage.getItem(key), preferenceKey)).toBe('true');
  await page.reload();
  await expanded(page, false);
  await menu(page).click();
  await expanded(page, true);
  await page.reload();
  await expanded(page, true);
});

test('instant navigation keeps a single working control and accurate page state', async ({page}) => {
  await page.setViewportSize({width: 1440, height: 900});
  await page.goto('');
  await menu(page).click();
  await expanded(page, true);
  await page.evaluate(() => { window.navigationSentinel = 'same-document'; });
  await page.getByRole('link', {name: 'Get started', exact: true}).click();
  await expect(page).toHaveURL(/\/getting-started\/$/);
  await expect(page.locator('h1')).toContainText('Getting started');
  expect(await page.evaluate(() => window.navigationSentinel)).toBe('same-document');
  await expect(menu(page)).toHaveCount(1);
  await expanded(page, true);
  await menu(page).click();
  await expanded(page, false);
  await menu(page).click();
  await expanded(page, true);
});

for (const path of ['', 'getting-started/']) {
  test(`mobile drawer supports open, Escape, overlay and link selection at ${path || 'home'}`, async ({page}) => {
    await page.setViewportSize({width: 390, height: 844});
    await page.goto(path);
    await expanded(page, false);
    await menu(page).click();
    await expanded(page, true);
    await page.keyboard.press('Escape');
    await expanded(page, false);
    await expect(menu(page)).toBeFocused();
    await menu(page).click();
    await expanded(page, true);
    await page.mouse.click(380, 650);
    await expanded(page, false);
    await menu(page).click();
    await expanded(page, true);
    await primary(page).getByRole('link', {name: path ? 'YAML structure' : 'Product', exact: true}).click();
    await expect(page).toHaveURL(path ? /\/getting-started\/yaml\/$/ : /\/product\/overview\/$/);
    await expanded(page, false);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  });
}

test('resize separates mobile drawer state from desktop preference', async ({page}) => {
  await page.setViewportSize({width: 1440, height: 900});
  await page.goto('getting-started/');
  await expanded(page, true);
  await menu(page).click();
  await expanded(page, false);
  await page.setViewportSize({width: 768, height: 900});
  await expanded(page, false);
  await menu(page).click();
  await expanded(page, true);
  await page.setViewportSize({width: 1440, height: 900});
  await expanded(page, false);
  await menu(page).click();
  await expanded(page, true);
  await page.keyboard.press('Escape');
  await expanded(page, false);
});

test('navigation works when local storage is unavailable', async ({page}) => {
  await page.addInitScript(() => {
    Object.defineProperty(window, 'localStorage', {get() { throw new DOMException('Storage unavailable', 'SecurityError'); }});
  });
  await page.setViewportSize({width: 1440, height: 900});
  await page.goto('');
  await expanded(page, false);
  await menu(page).click();
  await expanded(page, true);
  await menu(page).click();
  await expanded(page, false);
});

test('public documentation has no corporate attribution or corporate-site links', async ({page}) => {
  for (const path of ['', 'about/', 'community/sponsorship/', 'architecture/repository-boundaries/']) {
    await page.goto(path);
    await expect(page.locator('body')).not.toContainText(/iot\.?erax/i);
    await expect(page.locator('a[href*="ioterax.com"]')).toHaveCount(0);
  }
});
