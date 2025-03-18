import { expect } from '@grafana/plugin-e2e';
import { ODataOptions } from "../src/types";
import { test } from "./patch";

test('smoke: should render config editor', async ({ createDataSourceConfigPageNew, readProvisionedDataSource, page }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await createDataSourceConfigPageNew({ type: ds.type });
  await expect(page.getByLabel('Name')).toBeVisible()
});

test('"Save & test" should be successful when configuration is valid', async ({
  createDataSourceConfigPageNew,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<ODataOptions>({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPageNew({ type: ds.type });
  await page.getByPlaceholder('http://localhost:5000/odata').fill('http://test-server:4004/odata/v4/test');
  await expect(configPage.saveAndTest()).toBeOK();
});

test('"Save & test" should fail when configuration is invalid', async ({
  createDataSourceConfigPageNew,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<ODataOptions>({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPageNew({ type: ds.type });
  await page.getByPlaceholder('http://localhost:5000/odata').fill('');
  await expect(configPage.saveAndTest()).not.toBeOK();
  await expect(configPage).toHaveAlert('error', { hasText: 'Health check failed' });
});
