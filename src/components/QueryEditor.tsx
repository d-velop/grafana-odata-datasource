import React, { PureComponent } from 'react';
import { Button, InlineFormLabel, LegacyForms, Input } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { ODataSource } from '../DataSource';
import { EntitySet, FilterCondition, Metadata, ODataOptions, ODataQuery, Property, FilterOperators } from '../types';

const { Select } = LegacyForms;

type Props = QueryEditorProps<ODataSource, ODataQuery, ODataOptions>;

interface State {
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

  update = (updatedQuery: ODataQuery) => {
    this.props.onChange(updatedQuery);
    this.props.onRunQuery();
  };

  onEntitySetChange = (option: SelectableValue<EntitySet>) => {
    if (this.props.query.entitySet?.name === option.value?.name) {
      return;
    }
    const { query } = this.props;
    const { metadata } = this.state;
    const entitySet = option.value;
    const updatedQuery: ODataQuery = { ...query, entitySet, timeProperty: null, properties: [] };
    this.setState(
      {
        timeProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.Time),
        allProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.All),
      },
      () => this.update(updatedQuery)
    );
  };

  onTimePropertyChange = (option: SelectableValue<Property>) => {
    if (this.props.query.timeProperty === option.value) {
      return;
    }
    this.update({ ...this.props.query, timeProperty: option.value });
  };

  onPropertyChange = (option: SelectableValue<Property>, index: number) => {
    if (this.props.query.properties![index] === option.value) {
      return;
    }
    if (option.value != null) {
      const properties = [...this.props.query.properties!];
      properties[index] = option.value;
      this.update({ ...this.props.query, properties });
    }
  };

  addProperty = () => {
    const properties = [...(this.props.query.properties ?? []), { name: '', type: '' }];
    this.update({ ...this.props.query, properties });
  };

  removeProperty = (index: number) => {
    const properties = [...this.props.query.properties!];
    properties.splice(index, 1);
    this.update({ ...this.props.query, properties });
  };

  addFilterCondition = () => {
    const filterConditions = [
      ...(this.props.query.filterConditions ?? []),
      { property: { name: '', type: '' }, operator: '', value: '' },
    ];
    this.update({ ...this.props.query, filterConditions });
  };

  removeFilterCondition = (index: number) => {
    const filterConditions = [...this.props.query.filterConditions!];
    filterConditions.splice(index, 1);
    this.update({ ...this.props.query, filterConditions });
  };

  onFilterPropertyChange = (option: SelectableValue<Property>, index: number) => {
    const filterCondition = this.props.query.filterConditions![index];
    if (option.value?.name === filterCondition.property.name) {
      return;
    }
    if (option.value) {
      const type = this.state.allProperties.find((item) => item.value?.name === option.value?.name)?.value?.type ?? '';
      const updatedCondition: FilterCondition = {
        ...filterCondition,
        property: { name: option.value.name, type },
      };
      const filterConditions = [...this.props.query.filterConditions!];
      filterConditions[index] = updatedCondition;
      this.update({ ...this.props.query, filterConditions });
    }
  };

  onFilterOperatorChange = (option: SelectableValue<string>, index: number) => {
    const filterCondition = this.props.query.filterConditions![index];
    if (option.value === filterCondition.operator) {
      return;
    }
    if (option.value) {
      const filterConditions = [...this.props.query.filterConditions!];
      filterConditions[index] = { ...filterCondition, operator: option.value };
      this.update({ ...this.props.query, filterConditions });
    }
  };

  onFilterValueChange = (value: string, index: number) => {
    const filterCondition = this.props.query.filterConditions![index];
    if (value === filterCondition.value) {
      return;
    }
    if (value) {
      const filterConditions = [...this.props.query.filterConditions!];
      filterConditions[index] = { ...filterCondition, value };
      this.update({ ...this.props.query, filterConditions });
    }
  };

  render() {
    const { entitySets, timeProperties, allProperties, filterOperators } = this.state;
    let property = null;
    const listProperties = this.props.query.properties?.map((selectedProperty, index) => {
      property = (
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
              onChange={(item) => this.onFilterPropertyChange(item, index)}
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
              onChange={(item) => this.onFilterOperatorChange(item, index)}
              onBlur={this.props.onRunQuery}
              options={filterOperators}
              isSearchable={false}
            />
            <Input
              required={true}
              defaultValue={this.props.query.filterConditions?.[index].value}
              type="text"
              placeholder="(value)"
              onChange={(item) => this.onFilterValueChange(item.currentTarget.value, index)}
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
