import { expect, test, PanelEditPage } from '@grafana/plugin-e2e';
import { Page } from '@playwright/test';

async function openSelect(page: Page, panelEditPage: PanelEditPage, select: string) {
  await panelEditPage.getQueryEditorRow('A').getByRole('combobox', { name: select, exact: true }).focus();
  await page.keyboard.press('ArrowDown');
}

function queryResponseContaining(page: Page, text: string) {
  return page.waitForResponse(
    (res) => res.url().includes('/api/ds/query') && (res.request().postData() ?? '').includes(text)
  );
}

async function choose(page: Page, panelEditPage: PanelEditPage, select: string, option: string) {
  await openSelect(page, panelEditPage, select);
  await page.getByRole('option', { name: option, exact: true }).click();
}

test.beforeEach(async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml', name: 'OData-Test' });
  await panelEditPage.datasource.set(ds.name);
  await panelEditPage.setVisualization('Table');
});

test('should list entity sets from service metadata', async ({ panelEditPage, page }) => {
  await openSelect(page, panelEditPage, 'Entity set');
  await expect(page.getByRole('option', { name: 'Temperatures', exact: true })).toBeVisible();
  await expect(page.getByRole('option', { name: 'Rooms', exact: true })).toBeVisible();
});

test('should query entity set filtered by property', async ({ panelEditPage, page }) => {
  const row = panelEditPage.getQueryEditorRow('A');
  await choose(page, panelEditPage, 'Entity set', 'Rooms');
  await row.getByRole('button', { name: '+ Select' }).click();
  await choose(page, panelEditPage, 'Select property', 'name');
  await row.getByRole('button', { name: '+ Filter condition' }).click();
  await choose(page, panelEditPage, 'Filter property', 'name');
  await choose(page, panelEditPage, 'Filter operator', 'eq');
  const response = queryResponseContaining(page, 'Room 42');
  await row.getByRole('textbox', { name: 'Filter value' }).fill('Room 42');
  await row.getByRole('textbox', { name: 'Filter value' }).blur();

  await expect(response).toBeOK();
  await expect(panelEditPage.panel.data).toHaveText(['Room 42']);
});

test('should query entity set within dashboard time range', async ({ panelEditPage, page }) => {
  const row = panelEditPage.getQueryEditorRow('A');
  await choose(page, panelEditPage, 'Entity set', 'Temperatures');
  await choose(page, panelEditPage, 'Time property', 'time');
  await row.getByRole('button', { name: '+ Select' }).click();
  const response = queryResponseContaining(page, '"value1"');
  await choose(page, panelEditPage, 'Select property', 'value1');

  await expect(response).toBeOK();
  await expect(panelEditPage.panel.fieldNames).toContainText(['time', 'value1']);
  await expect(panelEditPage.panel.data.first()).toBeVisible();
});
