import React, { PureComponent } from 'react';
import { Button, InlineFormLabel, Input, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { ODataSource } from '../DataSource';
import { EntitySet, FilterCondition, Metadata, ODataOptions, ODataQuery, Property, FilterOperators } from '../types';

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
    if (!metadata || !entityType) {
      return [];
    }
    return metadata.entityTypes[entityType].properties
      .filter(
        (property) =>
          propertyKind === PropertyKind.All ||
          (propertyKind === PropertyKind.Time && property.type === 'Edm.DateTimeOffset')
      )
      .map((property) => ({
        label: property.name,
        value: property,
      }));
  }

  update = () => {
    this.props.onChange(this.props.query);
    this.props.onRunQuery();
  };

  resetSelection() {
    const { query } = this.props;
    query.timeProperty = null;
    query.properties = [];
  }

  onEntitySetChange = (option: SelectableValue<EntitySet>) => {
    if (this.props.query.entitySet?.name === option?.value?.name) {
      return;
    }
    const { query } = this.props;
    const { metadata } = this.state;
    const entitySet = option?.value;
    query.entitySet = entitySet;
    this.resetSelection();
    this.setState(
      {
        timeProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.Time),
        allProperties: this.mapProperties(metadata, entitySet?.entityType, PropertyKind.All),
      },
      this.update
    );
  };

  onTimePropertyChange = (option: SelectableValue<Property>) => {
    if (this.props.query.timeProperty === option?.value) {
      return;
    }
    const { query } = this.props;
    query.timeProperty = option?.value;
    this.update();
  };

  onPropertyChange = (option: SelectableValue<Property>, index: number) => {
    if (this.props.query.properties![index] === option?.value) {
      return;
    }
    const properties: Property[] = this.props.query.properties!;
    if (option?.value != null) {
      properties[index] = option.value;
    } else {
      properties[index] = { name: '', type: '' };
    }
    this.update();
  };

  onFilterConditionPropertyChange = (option: SelectableValue<Property>, index: number) => {
    const filterCondition = this.props.query.filterConditions![index];
    if (option?.value === filterCondition.property) {
      return;
    }
    if (option?.value) {
      filterCondition.property = option.value;
    } else {
      filterCondition.property = { name: '', type: '' };
    }
    this.update();
  };

  onFilterConditionOperatorChange = (option: SelectableValue<string>, index: number) => {
    const filterCondition = this.props.query.filterConditions![index];
    if (option?.value === filterCondition.operator) {
      return;
    }
    if (option?.value) {
      filterCondition.operator = option.value! as string;
    } else {
      filterCondition.operator = '';
    }
    this.update();
  };

  addProperty = () => {
    const properties: Property[] = this.props.query.properties!;
    properties.push({ name: '', type: '' });
    this.update();
  };

  removeProperty = (index: number) => {
    const properties: Property[] = this.props.query.properties!;
    properties.splice(index, 1);
    this.update();
  };

  addFilterCondition = () => {
    this.props.query.filterConditions = this.props.query.filterConditions ?? [];
    const filterConditions: FilterCondition[] = this.props.query.filterConditions!;
    filterConditions.push({ property: { name: '', type: '' }, operator: '', value: '' });
    this.update();
  };

  removeFilterCondition = (index: number) => {
    const filterConditions: FilterCondition[] = this.props.query.filterConditions!;
    filterConditions.splice(index, 1);
    this.update();
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
