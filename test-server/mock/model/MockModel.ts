export const MockModel = {
  schemas: [
    {
      namespace: "GrafanaMock",
      entityTypes: [
        {
          name: "Temperature",
          key: { propertyRef: { name: "id" } },
          properties: [
            { name: "id",              type: "Edm.Guid",            nullable: "false" },
            { name: "time",            type: "Edm.DateTimeOffset",  nullable: "false" },
            { name: "sampledAt",       type: "Edm.DateTimeOffset",  nullable: "true"  },
            { name: "measurementDate", type: "Edm.Date",            nullable: "true"  },
            { name: "epoch",           type: "Edm.Int64",           nullable: "false" },
            { name: "sensorId",        type: "Edm.Int32",           nullable: "true"  },
            { name: "qualityCode",     type: "Edm.Int16",           nullable: "true"  },
            { name: "value1",          type: "Edm.Double",          nullable: "false" },
            { name: "value2",          type: "Edm.Double",          nullable: "false" },
            { name: "value3",          type: "Edm.Double",          nullable: "false" },
            { name: "pressure",        type: "Edm.Decimal",         nullable: "true"  },
            { name: "isOutdoor",       type: "Edm.Boolean",         nullable: "true"  },
            { name: "unit",            type: "Edm.String",          nullable: "true"  }
          ]
        }
      ],
      entityContainer: undefined,
    },
    {
      namespace: "Default",
      entityTypes: undefined,
      entityContainer: {
        name: "Container",
        entitySets: [
          {
            name: "Temperatures",
            entityType: "GrafanaMock.Temperature"
          }
        ]
      }
    }
  ]
}
