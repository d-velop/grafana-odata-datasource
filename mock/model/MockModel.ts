export const MockModel = {
  schemas: [
    {
      namespace: "GrafanaMock",
      entityTypes: [
        {
          name: "Temperature",
          key: { propertyRef: { name: "Id" } },
          properties: [
            { name: "Id", type: "Edm.Guid", nullable: "false" },
            { name: "Time", type: "Edm.DateTimeOffset", nullable: "false" },
            { name: "Epoch", type: "Edm.Int64", nullable: "false" },
            { name: "Value1", type: "Edm.Double", nullable: "false" },
            { name: "Value2", type: "Edm.Double", nullable: "false" },
            { name: "Value3", type: "Edm.Double", nullable: "false" }
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
