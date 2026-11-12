import {defineConfig} from '@playwright/test';

export default defineConfig({
  testDir: './test/browser',
  testIgnore: '**/._*',
  retries: 0,
  timeout: 60000,
  workers: 2,
  reporter: 'list',
  use: {
    baseURL: 'https://storm-harbor.github.io/storm-harbor-control-plane/',
    trace: 'retain-on-failure',
    launchOptions: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH
      ? {executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH} : {},
  },
});
