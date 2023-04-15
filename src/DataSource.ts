import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { ODataOptions, ODataQuery } from './types';

export class ODataSource extends DataSourceWithBackend<ODataQuery, ODataOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<ODataOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: ODataQuery, scopedVars: ScopedVars) {
    const templateSrv = getTemplateSrv();

    query.filterConditions?.forEach((filterCondition) => {
      filterCondition.value = templateSrv.replace(filterCondition.value, scopedVars);
    });

    return {
      ...query,
    };
  }
}
