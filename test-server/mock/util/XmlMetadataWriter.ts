import {create} from "xmlbuilder2";
import {DataServices} from "../Metamodel";

export function ToXml(dataServices: DataServices)
{
  let ds = create({ version: '1.0', encoding: 'utf-8' })
    .ele('edmx:Edmx', { Version: '4.0', 'xmlns:edmx': 'https://docs.oasis-open.org/odata/ns/edmx' })
    .ele('edmx:DataServices');

  for (let schema of dataServices.schemas) {
    let s = ds.ele("Schema", { Namespace: schema.namespace, xmlns: 'https://docs.oasis-open.org/odata/ns/edm' });
    if (schema.entityTypes != null)
    {
      for (let entityType of schema.entityTypes) {
        let et = s.ele('EntityType', { Name: entityType.name });
        et.ele('Key')
          .ele('PropertyRef', { Name: entityType.key.propertyRef.name });
        for (let property of entityType.properties) {
          et.ele('Property', { Name: property.name, Type: property.type, Nullable: property.nullable });
        }
      }
    }
    if (schema.entityContainer != null)
    {
      let entityContainer = schema.entityContainer;
      let ec = s.ele('EntityContainer', { Name: entityContainer.name });
      for (let entitySet of entityContainer.entitySets) {
        ec.ele('EntitySet', { Name: entitySet.name, EntityType: entitySet.entityType });
      }
    }
  }
  return ds.end({ prettyPrint: true });
}
