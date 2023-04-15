export class Key {
  propertyRef: PropertyRef;
  constructor(propertyRef: PropertyRef) {
    this.propertyRef = propertyRef;
  }
}

export class PropertyRef {
  name: string;
  constructor(name: string) {
    this.name = name;
  }
}

export class Property {
  name: string;
  type: string;
  nullable: string;
  constructor(name: string, type: string, nullable: string) {
    this.name = name;
    this.type = type;
    this.nullable = nullable;
  }
}

export class EntityType {
  name: string;
  key: Key;
  properties: Property[];
  constructor(name: string, key: Key, properties: Property[]) {
    this.name = name;
    this.key = key;
    this.properties = properties;
  }
}

export class EntitySet {
  name: string;
  entityType: string;
  constructor(name: string, entityType: string) {
    this.name = name;
    this.entityType = entityType;
  }
}

export class EntityContainer {
  name: string;
  entitySets: EntitySet[];
  constructor(name: string, entitySets: EntitySet[]) {
    this.name = name;
    this.entitySets = entitySets;
  }
}

export class Schema {
  namespace: string;
  entityTypes: EntityType[] | undefined;
  entityContainer: EntityContainer | undefined;
  constructor(namespace: string, entityTypes?: EntityType[], entityContainer?: EntityContainer) {
    this.namespace = namespace;
    this.entityTypes = entityTypes;
    this.entityContainer = entityContainer;
  }
}

export class DataServices {
  schemas: Schema[];
  constructor(schemas: Schema[]) {
    this.schemas = schemas;
  }
}
