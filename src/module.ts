import { DataSourcePlugin } from '@grafana/data';
import { ODataSource } from './DataSource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { ODataQuery, ODataOptions } from './types';

export const plugin = new DataSourcePlugin<ODataSource, ODataQuery, ODataOptions>(ODataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
