//go:build integration

package plugin

var northwindFilterTests = []filterTestCase{
	{
		name:      "date range on Orders.OrderDate",
		entitySet: "Orders",
		properties: []property{
			{Name: "OrderID"},
			{Name: "CustomerID"},
			{Name: "OrderDate"},
			{Name: "Freight"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "OrderDate"}, Operator: "ge", Value: "1996-07-01T00:00:00Z"},
			{Property: property{Name: "OrderDate"}, Operator: "le", Value: "1996-07-31T23:59:59Z"},
		},
		expectedResults: 22,
	},
	{
		name:      "decimal range on Products.UnitPrice",
		entitySet: "Products",
		properties: []property{
			{Name: "ProductID"},
			{Name: "ProductName"},
			{Name: "UnitPrice"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "UnitPrice"}, Operator: "ge", Value: "10"},
			{Property: property{Name: "UnitPrice"}, Operator: "le", Value: "50"},
		},
		expectedResults: 20,
	},
	{
		name:      "single equality on Order_Details.Discount",
		entitySet: "Order_Details",
		properties: []property{
			{Name: "OrderID"},
			{Name: "ProductID"},
			{Name: "Discount"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "Discount"}, Operator: "eq", Value: "0.03"},
		},
		expectedResults: 3,
	},
}

var parliFilterTests = []filterTestCase{
	{
		name:      "int64 range on Meeting.ID",
		entitySet: "Meeting",
		properties: []property{
			{Name: "ID"},
			{Name: "Language"},
			{Name: "MeetingNumber"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "ID"}, Operator: "ge", Value: "1000"},
			{Property: property{Name: "ID"}, Operator: "le", Value: "1010"},
		},
		expectedResults: 55,
	},
}

var usgsFilterTests = []filterTestCase{
	{
		name:      "int64 range on Sites.SiteID",
		entitySet: "Sites",
		properties: []property{
			{Name: "SiteID"},
			{Name: "SiteName"},
			{Name: "SiteTypeCode"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "SiteID"}, Operator: "ge", Value: "1"},
			{Property: property{Name: "SiteID"}, Operator: "le", Value: "100"},
		},
		expectedResults: 89,
	},
}

var odataSvcFilterTests = []filterTestCase{
	{
		name:      "date range on Products.ReleaseDate",
		entitySet: "Products",
		properties: []property{
			{Name: "ID"},
			{Name: "Name"},
			{Name: "ReleaseDate"},
			{Name: "Price"},
			{Name: "Rating"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "ReleaseDate"}, Operator: "ge", Value: "2000-01-01T00:00:00Z"},
			{Property: property{Name: "ReleaseDate"}, Operator: "le", Value: "2006-12-31T23:59:59Z"},
		},
		expectedResults: 6,
	},
	{
		name:      "double range on Products.Price",
		entitySet: "Products",
		properties: []property{
			{Name: "ID"},
			{Name: "Name"},
			{Name: "Price"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "Price"}, Operator: "ge", Value: "18"},
			{Property: property{Name: "Price"}, Operator: "le", Value: "25"},
		},
		expectedResults: 5,
	},
	{
		name:      "int16 filter on Products.Rating",
		entitySet: "Products",
		properties: []property{
			{Name: "ID"},
			{Name: "Name"},
			{Name: "Rating"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "Rating"}, Operator: "ge", Value: "4"},
		},
		expectedResults: 3,
	},
	{
		name:      "guid equality on Advertisements.ID",
		entitySet: "Advertisements",
		properties: []property{
			{Name: "ID"},
			{Name: "Name"},
		},
		filterConditions: []filterCondition{
			{Property: property{Name: "ID"}, Operator: "eq", Value: "f89dee73-af9f-4cd4-b330-db93c25ff3c7"},
		},
		expectedResults: 1,
	},
}

func init() {
	referenceSystems = append(referenceSystems,
		referenceSystem{
			name:        "Northwind V2 (services.odata.org)",
			baseURL:     "https://services.odata.org/V2/Northwind/Northwind.svc",
			version:     "V2",
			filterTests: northwindFilterTests,
		},
		referenceSystem{
			name:        "Northwind V3 (services.odata.org)",
			baseURL:     "https://services.odata.org/V3/Northwind/Northwind.svc",
			version:     "V3",
			filterTests: northwindFilterTests,
		},
		referenceSystem{
			name:        "Northwind V4 (services.odata.org)",
			baseURL:     "https://services.odata.org/V4/Northwind/Northwind.svc",
			version:     "V4",
			filterTests: northwindFilterTests,
		},
		referenceSystem{
			name:        "OData V3 demo (services.odata.org)",
			baseURL:     "https://services.odata.org/V3/OData/OData.svc",
			version:     "V3",
			filterTests: odataSvcFilterTests,
		},
		referenceSystem{
			name:        "OData V4 demo (services.odata.org)",
			baseURL:     "https://services.odata.org/V4/OData/OData.svc",
			version:     "V4",
			filterTests: odataSvcFilterTests,
		},
		referenceSystem{
			name:        "Swiss Parliament V2 (ws.parlament.ch)",
			baseURL:     "https://ws.parlament.ch/odata.svc",
			version:     "V2",
			filterTests: parliFilterTests,
		},
		referenceSystem{
			name:        "USGS Water Data V4 (waterdata.usgs.gov)",
			baseURL:     "https://dashboard.waterdata.usgs.gov/service/cwis/1.0/odata",
			version:     "V4",
			filterTests: usgsFilterTests,
		},
	)
}
