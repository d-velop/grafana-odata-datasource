import React, { PureComponent } from 'react';
import { Button, InlineFormLabel, Input, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { ODataSource } from '../DataSource';
import { EntitySet, Metadata, ODataOptions, ODataQuery, Property, FilterOperators } from '../types';

type Props = QueryEditorProps<ODataSource, ODataQuery, ODataOptions>;

interface State {
  resetKey: number;
  metadata: Metadata | undefined;
  entitySets: Array<SelectableValue<EntitySet>>;
  timeProperties: Array<SelectableValue<Property>>;
  allProperties: Array<SelectableValue<Property>>;
  filterOperators: Array<SelectableValue<string>>;
}

enum PropertyKind {
  Time = 1,
  All = 2,
}

export class QueryEditor extends PureComponent<Props, State> {
  dataSource: ODataSource;

  constructor(props: Props) {
    super(props);
    this.dataSource = this.props.datasource;
    this.state = {
      resetKey: 0,
      metadata: undefined,
      entitySets: [],
      timeProperties: [],
      allProperties: [],
      filterOperators: [],
    };
  }

  componentDidMount() {
    this.dataSource.getResource('metadata').then((metadata: Metadata) => {
      const entityType = this.props.query.entitySet?.entityType;
      this.setState({
        metadata: metadata,
        entitySets: Object.values(metadata.entitySets).map((entitySet) => ({
          label: entitySet.name,
          value: entitySet,
        })),
        timeProperties: this.mapProperties(metadata, entityType, PropertyKind.Time),
        allProperties: this.mapProperties(metadata, entityType, PropertyKind.All),
      });
    });
    const filterOperators: Array<SelectableValue<string>> = FilterOperators.map((operator) => ({
      label: operator,
      value: operator,
    }));
    this.setState({
      filterOperators: filterOperators,
    });
  }

  mapProperties(metadata: Metadata | undefined, entityType: string | undefined, propertyKind: PropertyKind) {
    if (!metadata || !entityType || !metadata.entityTypes[entityType]) {
      return [];
    }
    return (propertyKind === PropertyKind.Time ? [{ label: '(None)', value: undefined as Property | undefined }] : [])
      .concat(metadata.entityTypes[entityType].properties
        .filter(property =>
          propertyKind === PropertyKind.All ||
          (propertyKind === PropertyKind.Time && property.type === 'Edm.DateTimeOffset')
        )
        .map(property => ({
          label: property.name,
          value: property,
        }))
      );
  }

  update = (query: ODataQuery) => {
    this.props.onChange(query);
    this.props.onRunQuery();
  };

  onEntitySetChange = (option: SelectableValue<EntitySet>) => {
    if (this.props.query.entitySet?.name === option?.value?.name) { return };
    const { metadata } = this.state;
    const entitySet = option?.value;
    const updatedQuery = { ...this.props.query, entitySet: entitySet };
    updatedQuery.timeProperty = null;
    updatedQuery.properties = [];
    updatedQuery.filterConditions = [];
    this.props.onChange(updatedQuery);
    this.setState((prevState) => (
      {
        timeProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.Time),
        allProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.All),
        resetKey: prevState.resetKey + 1
      }),
      this.props.onRunQuery
    );
  };

  onTimePropertyChange = (option: SelectableValue<Property>) => {
    if (this.props.query.timeProperty === option?.value) { return };
    this.update({ ...this.props.query, timeProperty: option?.value, });
  };

  onPropertyChange = (option: SelectableValue<Property>, index: number) => {
    if (this.props.query.properties![index] === option?.value) { return };
    const properties = [...this.props.query.properties!];
    properties[index] = option?.value ? { ...option.value } : { name: '', type: '' };
    this.update({ ...this.props.query, properties: properties });
  };

  onFilterConditionPropertyChange = (option: SelectableValue<Property>, index: number) => {
    if (this.props.query.filterConditions![index].property === option?.value) { return };
    const filterConditions = [...this.props.query.filterConditions!];
    filterConditions[index].property = option?.value ? { ...option.value } : { name: '', type: '' };
    this.update({ ...this.props.query, filterConditions: filterConditions });
  };

  onFilterConditionOperatorChange = (option: SelectableValue<string>, index: number) => {
    if (this.props.query.filterConditions![index].operator === option?.value) { return };
    const filterConditions = [...this.props.query.filterConditions!];
    filterConditions[index].operator = option?.value ? option.value : '';
    this.update({ ...this.props.query, filterConditions: filterConditions });
  };

  addProperty = () => {
    const properties = [...this.props.query.properties!];
    properties.push({ name: '', type: '' });
    this.update({ ...this.props.query, properties: properties });
  };

  removeProperty = (index: number) => {
    const properties = [...this.props.query.properties!];
    properties.splice(index, 1);
    this.update({ ...this.props.query, properties: properties });
  };

  addFilterCondition = () => {
    const filterConditions = [...this.props.query.filterConditions!];
    filterConditions.push({ property: { name: '', type: '' }, operator: '', value: '' });
    this.update({ ...this.props.query, filterConditions: filterConditions });
  };

  removeFilterCondition = (index: number) => {
    const filterConditions = [...this.props.query.filterConditions!];
    filterConditions.splice(index, 1);
    this.update({ ...this.props.query, filterConditions: filterConditions });
  };

  render() {
    const { entitySets, timeProperties, allProperties, filterOperators } = this.state;
    let property = null;
    const listProperties = this.props.query.properties?.map((selectedProperty, index) => {
      property = (
        <div className="gf-form-inline">
          <div className={'gf-form'}>
            <InlineFormLabel width={8} tooltip={'Add select'}>
              Select
            </InlineFormLabel>
            <Select
              value={allProperties.find((item) => item.value?.name === this.props.query.properties?.[index].name)}
              isClearable={true}
              placeholder="(Property)"
              onChange={(item) => this.onPropertyChange(item, index)}
              onBlur={this.props.onRunQuery}
              options={allProperties}
              isSearchable={false}
            />
            <Button variant={'secondary'} onClick={() => this.removeProperty(index)}>
              -
            </Button>
          </div>
        </div>
      );
      return property;
    });
    let filter = null;
    const listFilters = this.props.query.filterConditions?.map((filterCondition, index) => {
      filter = (
        <div className="gf-form-inline">
          <div className={'gf-form'}>
            <InlineFormLabel width={8} tooltip={'Add filter condition'}>
              {index === 0 ? 'Filter' : 'AND'}
            </InlineFormLabel>
            <Select
              value={allProperties.find(
                (item) => item.value?.name === this.props.query.filterConditions?.[index].property.name
              )}
              isClearable={true}
              placeholder="(Property)"
              onChange={(item) => this.onFilterConditionPropertyChange(item, index)}
              onBlur={this.props.onRunQuery}
              options={allProperties}
              isSearchable={false}
            />
            <Select
              value={
                filterCondition.operator
                  ? { label: filterCondition.operator, value: filterCondition.operator }
                  : undefined
              }
              isClearable={true}
              placeholder="(Operator)"
              onChange={(item) => this.onFilterConditionOperatorChange(item, index)}
              onBlur={this.props.onRunQuery}
              options={filterOperators}
              isSearchable={false}
            />
            <Input
              required={true}
              defaultValue={this.props.query.filterConditions?.[index].value}
              type="text"
              placeholder="(value)"
              onChange={(item) => {
                filterCondition = this.props.query.filterConditions![index];
                if (item.currentTarget.value === filterCondition.value) {
                  return;
                }
                if (item.currentTarget.value) {
                  filterCondition.value = item.currentTarget.value;
                }
              }}
              onBlur={this.props.onRunQuery}
            />
            <Button variant={'secondary'} onClick={() => this.removeFilterCondition(index)}>
              -
            </Button>
          </div>
        </div>
      );
      return filter;
    });
    return (
      <div>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel width={8} tooltip="Select an entity set for a list of available metrics.">
              Entity set
            </InlineFormLabel>
            <Select
              value={entitySets.find((o) => o.value?.name === this.props.query.entitySet?.name)}
              isClearable={true}
              placeholder="(Entity set)"
              onChange={this.onEntitySetChange}
              onBlur={this.props.onRunQuery}
              options={entitySets}
              isSearchable={false}
            />
            <InlineFormLabel width={8} tooltip="Time property">
              Time property
            </InlineFormLabel>
            <Select
              key={this.state.resetKey}
              value={timeProperties.find((o) => o.value?.name === this.props.query.timeProperty?.name)}
              isClearable={true}
              placeholder="(Property)"
              onChange={this.onTimePropertyChange}
              onBlur={this.props.onRunQuery}
              options={timeProperties}
              isSearchable={false}
            />
          </div>
        </div>
        {listProperties}
        <div className="gf-form-inline">
          <div className={'gf-form'}>
            <Button variant={'secondary'} onClick={this.addProperty}>
              + Select
            </Button>
          </div>
        </div>
        {listFilters}
        <div className={'gf-form'}>
          <Button variant={'secondary'} onClick={this.addFilterCondition}>
            + Filter condition
          </Button>
        </div>
      </div>
    );
  }
}
