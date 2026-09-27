import { expect, test as setup } from '@playwright/test';

setup('sign in to Grafana via Keycloak', async ({ page }) => {
  await page.goto('/login');
  await page.getByRole('link', { name: /Keycloak-OAuth/ }).click();
  await page.locator('#username').fill('user');
  await page.locator('#password').fill('user');
  await page.getByRole('button', { name: 'Sign In' }).click();
  await page.waitForURL('/');

  const user = await page.request.get('/api/user');
  await expect(user).toBeOK();
  expect(await user.json()).toMatchObject({ login: 'user@test.localhost', isExternal: true });

  await page.context().storageState({ path: 'playwright/.auth/keycloak-user.json' });
});
