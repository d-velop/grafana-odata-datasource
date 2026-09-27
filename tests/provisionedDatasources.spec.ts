import { expect, test } from '@grafana/plugin-e2e';

for (const name of ['OData-Test', 'OData-Mock']) {
  test(`health check of provisioned data source ${name} should succeed`, async ({
    gotoDataSourceConfigPage,
    readProvisionedDataSource,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml', name });
    const configPage = await gotoDataSourceConfigPage(ds.uid);
    await expect(configPage.saveAndTest()).toBeOK();
  });
}
