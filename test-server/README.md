# Test server
The test server offers a set of basic functions that aid in developing the Grafana OData Data Source plugin.

The test server features three OData service roots:

1. `/mock/` provides a very simple mock service root
2. `/odata/v4/test` provides a full OData V4 service and is built on the
   [SAP Cloud Application Programming Model](https://cap.cloud.sap/)
3. `/odata/v2/test` based on
   [OData V2 Adapter](https://github.com/cap-js-community/odata-v2-adapter)

## Mock
This is a simple mock service root that returns test data. It provides endpoints for service document (`/mock/`)
and metadata (`/mock/$metadata`). In addition, it provides a simple entity set `/mock/temperatures` with support for
very basic predefined filter functions.

For an introduction to programming or extension see [mock/MockService.ts](mock/MockService.ts).

## OData V4 Test
Based on the Core Data Services (CDS) of the SAP Cloud Application Programming Model, this service offers extensive
OData query options.

To get an overview of the supported OData features in CDS/CAP, please visit
[Serving OData APIs](https://cap.cloud.sap/docs/advanced/odata).

The data model of the test server can be found in the file [db/data-model.cds](db/data-model.cds).

## OData V2 Test
Provides a OData V2 service that converts OData V2 requests into CDS OData V4 service calls and responses.
