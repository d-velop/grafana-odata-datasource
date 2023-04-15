import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface ODataQuery extends DataQuery {
  entitySet?: EntitySet;
  timeProperty?: Property | null;
  properties?: Property[];
  filterConditions?: FilterCondition[];
  from: string;
  to: string;
}

export const defaultQuery: Partial<ODataQuery> = {};

export const FilterOperators: string[] = ['eq', 'ne', 'gt', 'ge', 'lt', 'le'];

export interface ODataOptions extends DataSourceJsonData {}

export interface Metadata {
  entityTypes: { [name: string]: EntityType };
  entitySets: { [name: string]: EntitySet };
}

export interface EntityType {
  name: string;
  qualifiedName: string;
  properties: Property[];
}

export interface EntitySet {
  name: string;
  entityType: string;
}

export interface Property {
  name: string;
  type: string;
}

export interface FilterCondition {
  property: Property;
  operator: string;
  value: string;
}
