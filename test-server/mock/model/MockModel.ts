export const MockModel = {
  schemas: [
    {
      namespace: "GrafanaMock",
      entityTypes: [
        {
          name: "Temperature",
          key: { propertyRef: { name: "id" } },
          properties: [
            { name: "id", type: "Edm.Guid", nullable: "false" },
            { name: "time", type: "Edm.DateTimeOffset", nullable: "false" },
            { name: "epoch", type: "Edm.Int64", nullable: "false" },
            { name: "value1", type: "Edm.Double", nullable: "false" },
            { name: "value2", type: "Edm.Double", nullable: "false" },
            { name: "value3", type: "Edm.Double", nullable: "false" }
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
