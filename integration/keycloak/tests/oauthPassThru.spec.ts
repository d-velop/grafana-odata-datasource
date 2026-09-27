import { expect, test } from '@grafana/plugin-e2e';

test('health check of OData-Mock should succeed with forwarded Keycloak token', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml', name: 'OData-Mock' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await expect(configPage.saveAndTest()).toBeOK();
});

test('mocked panels should show data with forwarded Keycloak token', async ({
  gotoDashboardPage,
  readProvisionedDashboard,
}) => {
  const dashboard = await readProvisionedDashboard({ fileName: 'test_dashboard.json' });
  const dashboardPage = await gotoDashboardPage(dashboard);
  await dashboardPage.waitForPanelsQueriesToComplete({ scrollAll: true });

  await expect(dashboardPage).toHavePanelErrors(0);
  await expect(dashboardPage.getPanelByTitle('Temperatures table (mocked)').data.first()).toBeVisible();
});

test('health check of OData-Mock should fail without Keycloak token', async ({
  playwright,
  readProvisionedDataSource,
  baseURL,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml', name: 'OData-Mock' });
  const localAdmin = await playwright.request.newContext({ baseURL });
  await expect(await localAdmin.post('/login', { data: { user: 'admin', password: 'admin' } })).toBeOK();

  const health = await localAdmin.get(`/api/datasources/uid/${ds.uid}/health`);
  expect(await health.json()).toMatchObject({ status: 'ERROR', message: expect.stringContaining('401') });
  await localAdmin.dispose();
});
