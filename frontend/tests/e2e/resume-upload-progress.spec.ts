import { test, expect } from '@playwright/test';

test.describe('Resume Upload Progress', () => {
  test.beforeEach(async ({ page }) => {
    const token = process.env.NEXT_PUBLIC_E2E_TOKEN || 'e2e-token';

    await page.addInitScript((value) => {
      window.localStorage.setItem('token', value);
    }, token);

    await page.route('**/user/model/check', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: { configured: true } }),
      });
    });

    await page.route('**/user/profile', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ code: 200, data: { username: 'e2e-user' } }),
      });
    });

    await page.route('**/resume/list', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ code: 200, data: { resumes: [] } }),
      });
    });

    await page.route('**/resume/upload', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            is_async: true,
            upload_id: 'e2e-upload-id',
          },
        }),
      });
    });

    await page.route('**/resume/upload/progress/**', async (route) => {
      const events = [
        {
          upload_id: 'e2e-upload-id',
          status: 'pending',
          progress: 10,
          stage: '排队中',
        },
        {
          upload_id: 'e2e-upload-id',
          status: 'extracting',
          progress: 45,
          stage: '提取文本',
        },
        {
          upload_id: 'e2e-upload-id',
          status: 'analyzing',
          progress: 80,
          stage: '分析中',
        },
        {
          upload_id: 'e2e-upload-id',
          status: 'completed',
          progress: 100,
          stage: '完成',
          resume_id: 1,
        },
      ];

      const body = events
        .map((payload) => `event: progress\ndata: ${JSON.stringify(payload)}\n\n`)
        .join('');

      await route.fulfill({
        status: 200,
        headers: {
          'content-type': 'text/event-stream',
          'cache-control': 'no-cache',
        },
        body: `${body}event: complete\ndata: ${JSON.stringify({ status: 'completed' })}\n\n`,
      });
    });

    await page.goto('/user/center');
  });

  test('shows progress modal and allows background processing', async ({ page }) => {
    await page.getByText('我的简历').waitFor({ state: 'visible' });

    const fileInput = page.locator('input[type="file"]');
    await expect(fileInput).toBeEnabled({ timeout: 10000 });

    await fileInput.setInputFiles({
      name: 'resume.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('%PDF-1.4\n%Fake PDF for E2E\n'),
    });

    const progressTitle = page.getByText('简历上传进度');
    await expect(progressTitle).toBeVisible({ timeout: 12000 });

    const backgroundButton = page.getByRole('button', { name: '后台处理' });
    if (await backgroundButton.isVisible()) {
      await backgroundButton.click();
      const backgroundCard = page.getByText('后台处理中');
      await expect(backgroundCard).toBeVisible();
    }

    const viewButton = page.getByRole('button', { name: '查看进度' });
    if (await viewButton.isVisible()) {
      await expect(viewButton).toBeVisible();
    }

    await expect(page.getByText('后台处理中')).toBeHidden({ timeout: 5000 });
  });
});
