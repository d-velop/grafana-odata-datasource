import { DataSourceConfigPage, test as base } from "@grafana/plugin-e2e";
import { createDataSourceViaAPI } from "@grafana/plugin-e2e/dist/fixtures/commands/createDataSource";
import { v4 as uuidv4 } from "uuid";

/*
 A patch version of createDataSourceConfigPage that uses waitUntil: 'domcontentloaded' instead of 'networkidle'.
 This avoids issues with endless or ongoing requests (e.g., WebSockets, polling, etc.) that prevent the 'networkidle'
 event from firing.

 Perhaps it is a bug that only occurs with Grafana version >= 11.5.2?
*/

type CreateDsArgs = {
  type: string;
};

export const test = base.extend<{
  createDataSourceConfigPageNew: (args: CreateDsArgs) => Promise<DataSourceConfigPage>;
}>({
  createDataSourceConfigPageNew: async (
    {page, grafanaAPIClient, selectors, grafanaVersion, request},
    use,
    testInfo
  ) => {
    let datasourceConfigPage: DataSourceConfigPage | undefined;
    const createDataSourceConfigPageFn = async (args: CreateDsArgs) => {
      const datasource = await createDataSourceViaAPI(grafanaAPIClient, {
        type: args.type,
        name: `${args.type}-${uuidv4()}`,
      });
      datasourceConfigPage = new DataSourceConfigPage(
        {page, selectors, grafanaVersion, request, testInfo},
        datasource
      );
      await datasourceConfigPage.goto({waitUntil: 'domcontentloaded'});
      return datasourceConfigPage;
    };
    await use(createDataSourceConfigPageFn);
    await datasourceConfigPage?.deleteDataSource();
  },
});
