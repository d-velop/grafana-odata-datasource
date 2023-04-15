import { Metadata } from './types';

export class XmlMetadataReader {
  body: string;
  parser: DOMParser;
  result: Metadata;

  constructor(body: string) {
    this.body = body;
    this.parser = new DOMParser();
    this.result = {
      entityTypes: {},
      entitySets: {},
    };
  }

  read() {
    const root = this.parser.parseFromString(this.body, 'application/xml');
    const edmx = root.firstElementChild;
    this.readDataServices(edmx!.firstElementChild);
    return this.result;
  }

  readDataServices(dataServices: Element | null) {
    if (!dataServices) {
      return;
    }
    for (let i = 0; i < dataServices.children.length; i++) {
      const schema = dataServices.children.item(i);
      this.readSchema(schema);
    }
  }

  readSchema(schema: Element | null) {
    if (!schema) {
      return;
    }
    const namespace = schema.getAttribute('Namespace');
    for (let i = 0; i < schema.children.length; i++) {
      const item = schema.children.item(i);
      if (item!.tagName === 'EntityType') {
        this.readEntityType(item, namespace);
      } else if (item!.tagName === 'EntityContainer') {
        this.readEntityContainer(item, namespace);
      }
    }
  }

  readEntityType(entityType: Element | null, namespace: string | null) {
    if (!entityType) {
      return;
    }
    const entityTypeName = entityType.getAttribute('Name')!;
    const qualifiedName = `${namespace}.${entityTypeName}`;
    const properties = [];
    for (let i = 0; i < entityType.children.length; i++) {
      const item = entityType.children.item(i);
      if (item!.tagName === 'Property') {
        const propertyName = item!.getAttribute('Name')!;
        const propertyType = item!.getAttribute('Type')!;
        properties.push({ name: propertyName, type: propertyType });
      }
    }
    this.result.entityTypes[qualifiedName] = {
      name: entityTypeName,
      qualifiedName: qualifiedName,
      properties: properties,
    };
  }

  readEntityContainer(entityContainer: Element | null, namespace: string | null) {
    if (!entityContainer) {
      return;
    }
    for (let i = 0; i < entityContainer.children.length; i++) {
      const item = entityContainer.children.item(i);
      if (item!.tagName === 'EntitySet') {
        const entitySetName = item!.getAttribute('Name')!;
        const entitySetType = item!.getAttribute('EntityType')!;
        this.result.entitySets[entitySetName] = { name: entitySetName, entityType: entitySetType };
      }
    }
  }
}
