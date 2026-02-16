import { defineConfig, devices } from '@playwright/test';

const port = process.env.E2E_PORT || '4173';
const baseURL = process.env.NEXT_PUBLIC_E2E_BASE_URL || `http://localhost:${port}`;

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false,
  retries: 0,
  timeout: 120000,
  webServer: {
    command: `NEXT_PUBLIC_E2E=true npm run dev -- --port ${port}`,
    url: baseURL,
    reuseExistingServer: true,
    timeout: 120000,
  },
  use: {
    baseURL,
    headless: true,
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  reporter: [['list']],
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
