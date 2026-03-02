//go:build integration

package plugin

var testPrimitivesFilterTests = []filterTestCase{
	{
		name:      "boolean equality",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "boolean"},
			{Name: "int32"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "boolean"}, Operator: "eq", Value: "true"},
		},
		expectedResults: 50,
	},
	{
		name:      "int32 range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "int32"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "int32"}, Operator: "ge", Value: "50"},
		},
		expectedResults: 50,
	},
	{
		name:      "int64 range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "int64"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "int64"}, Operator: "ge", Value: "50000000000"},
		},
		expectedResults: 50,
	},
	{
		name:      "int16 range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "int16"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "int16"}, Operator: "ge", Value: "50"},
		},
		expectedResults: 50,
	},
	{
		name:      "decimal range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "decimal"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "decimal"}, Operator: "ge", Value: "25"},
		},
		expectedResults: 50,
	},
	{
		name:      "double range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "double"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "double"}, Operator: "gt", Value: "0"},
		},
		expectedResults: 50,
	},
	{
		name:      "string equality",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "string"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "string"}, Operator: "eq", Value: "item-050"},
		},
		expectedResults: 1,
	},
	{
		name:      "dateTimeOffset range",
		entitySet: "TestPrimitives",
		properties: []property{
			{Name: "guid"},
			{Name: "dateTimeOffset"},
			{Name: "date"},
			{Name: "int16"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "dateTimeOffset"}, Operator: "ge", Value: "2024-01-01T00:30:00Z"},
		},
		expectedResults: 70,
	},
}

func init() {
	referenceSystems = append(referenceSystems,
		referenceSystem{
			name:                "Local test server V4 (localhost:4004)",
			baseURL:             "http://localhost:4004/odata/v4/test",
			version:             "V4",
			requiresLocalServer: true,
			filterTests:         testPrimitivesFilterTests,
		},
		referenceSystem{
			name:                "Local test server V2 (localhost:4004)",
			baseURL:             "http://localhost:4004/odata/v2/test",
			version:             "V2",
			requiresLocalServer: true,
			filterTests:         testPrimitivesFilterTests,
		},
	)
}
