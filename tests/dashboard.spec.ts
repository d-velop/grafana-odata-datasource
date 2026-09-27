import { expect, test } from '@grafana/plugin-e2e';

test('provisioned test dashboard should render all panels without errors', async ({
  gotoDashboardPage,
  readProvisionedDashboard,
}) => {
  const dashboard = await readProvisionedDashboard({ fileName: 'test_dashboard.json' });
  const dashboardPage = await gotoDashboardPage(dashboard);
  await dashboardPage.waitForPanelsQueriesToComplete({ scrollAll: true });

  await expect(dashboardPage).toHavePanelErrors(0);
  for (const title of ['Temperatures table', 'Temperatures table (mocked)']) {
    await expect(dashboardPage.getPanelByTitle(title).data.first()).toBeVisible();
  }
});
