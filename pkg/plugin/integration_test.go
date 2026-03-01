//go:build integration

package plugin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type filterTestCase struct {
	name             string
	entitySet        string
	properties       []property
	filterConditions []filterCondition
	expectedResults  int // expectation: exact number of entries the server must return
}

type referenceSystem struct {
	name        string
	baseURL     string
	version     string // "V2", "V3", "V4"
	filterTests []filterTestCase
}

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

var referenceSystems = []referenceSystem{
	{
		name:        "Northwind V2 (services.odata.org)",
		baseURL:     "https://services.odata.org/V2/Northwind/Northwind.svc",
		version:     "V2",
		filterTests: northwindFilterTests,
	},
	{
		name:        "Northwind V3 (services.odata.org)",
		baseURL:     "https://services.odata.org/V3/Northwind/Northwind.svc",
		version:     "V3",
		filterTests: northwindFilterTests,
	},
	{
		name:        "Northwind V4 (services.odata.org)",
		baseURL:     "https://services.odata.org/V4/Northwind/Northwind.svc",
		version:     "V4",
		filterTests: northwindFilterTests,
	},
	{
		name:        "OData V3 demo (services.odata.org)",
		baseURL:     "https://services.odata.org/V3/OData/OData.svc",
		version:     "V3",
		filterTests: odataSvcFilterTests,
	},
	{
		name:        "OData V4 demo (services.odata.org)",
		baseURL:     "https://services.odata.org/V4/OData/OData.svc",
		version:     "V4",
		filterTests: odataSvcFilterTests,
	},
	{
		name:        "Swiss Parliament V2 (ws.parlament.ch)",
		baseURL:     "https://ws.parlament.ch/odata.svc",
		version:     "V2",
		filterTests: parliFilterTests,
	},
	{
		name:        "USGS Water Data V4 (waterdata.usgs.gov)",
		baseURL:     "https://dashboard.waterdata.usgs.gov/service/cwis/1.0/odata",
		version:     "V4",
		filterTests: usgsFilterTests,
	},
}

type propertyTypeMap map[string]map[string]string

func fetchPropertyTypes(t *testing.T, c *ODataClientImpl) propertyTypeMap {
	t.Helper()
	resp, err := c.GetMetadata(context.Background())
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var edmx odata.Edmx
	require.NoError(t, xml.Unmarshal(body, &edmx))

	etProps := map[string]map[string]string{}
	for _, ds := range edmx.DataServices {
		for _, schema := range ds.Schemas {
			for _, et := range schema.EntityTypes {
				props := map[string]string{}
				for _, p := range et.Properties {
					props[p.Name] = p.Type
				}
				etProps[schema.Namespace+"."+et.Name] = props
			}
		}
	}

	result := propertyTypeMap{}
	for _, ds := range edmx.DataServices {
		for _, schema := range ds.Schemas {
			for _, ec := range schema.EntityContainers {
				for _, es := range ec.EntitySet {
					if props, ok := etProps[es.EntityType]; ok {
						result[es.Name] = props
					}
				}
			}
		}
	}
	return result
}

func withResolvedTypes(types propertyTypeMap, entitySet string, properties []property, filterConditions []filterCondition) ([]property, []filterCondition) {
	propTypes := types[entitySet]

	resolvedProps := make([]property, len(properties))
	for i, p := range properties {
		p.Type = propTypes[p.Name]
		resolvedProps[i] = p
	}

	resolvedConds := make([]filterCondition, len(filterConditions))
	for i, c := range filterConditions {
		c.Property.Type = propTypes[c.Property.Name]
		resolvedConds[i] = c
	}

	return resolvedProps, resolvedConds
}

func fetchUnfilteredCount(t *testing.T, rs referenceSystem, entitySet string) int {
	t.Helper()
	c := newIntegrationClient(rs)

	if rs.version == "V4" {
		rawURL := fmt.Sprintf("%s/%s/$count", c.baseUrl, entitySet)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
		require.NoError(t, err)
		req.Header.Set("Accept", "text/plain")
		c.addVersionHeaders(req)
		resp, err := c.httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		count, err := strconv.Atoi(strings.Trim(string(body), "\ufeff\r\n\t "))
		require.NoError(t, err)
		return count
	}

	rawURL := fmt.Sprintf("%s/%s?$inlinecount=allpages&$top=1", c.baseUrl, entitySet)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	c.addVersionHeaders(req)
	resp, err := c.httpClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var v2Result struct {
		D struct {
			Count string `json:"__count"`
		} `json:"d"`
	}
	if json.Unmarshal(body, &v2Result) == nil && v2Result.D.Count != "" {
		count, err := strconv.Atoi(v2Result.D.Count)
		require.NoError(t, err)
		return count
	}

	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &raw))
	var countStr string
	require.NoError(t, json.Unmarshal(raw["odata.count"], &countStr))
	count, err := strconv.Atoi(countStr)
	require.NoError(t, err)
	return count
}

func newIntegrationClient(rs referenceSystem) *ODataClientImpl {
	return &ODataClientImpl{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		baseUrl:      rs.baseURL,
		odataVersion: rs.version,
	}
}

func TestIntegration_ServiceRoot(t *testing.T) {
	for _, rs := range referenceSystems {
		rs := rs
		t.Run(rs.name, func(t *testing.T) {
			client := newIntegrationClient(rs)
			resp, err := client.GetServiceRoot(context.Background())
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestIntegration_Metadata(t *testing.T) {
	for _, rs := range referenceSystems {
		rs := rs
		t.Run(rs.name, func(t *testing.T) {
			client := newIntegrationClient(rs)
			resp, err := client.GetMetadata(context.Background())
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			var edmx odata.Edmx
			err = xml.Unmarshal(body, &edmx)
			require.NoError(t, err, "EDMX must be parseable")

			var totalEntityTypes int
			for _, ds := range edmx.DataServices {
				for _, schema := range ds.Schemas {
					totalEntityTypes += len(schema.EntityTypes)
				}
			}
			assert.Greater(t, totalEntityTypes, 0, "expected at least one entity type in metadata")
		})
	}
}

func TestIntegration_QueryWithFilter(t *testing.T) {
	for _, rs := range referenceSystems {
		rs := rs
		t.Run(rs.name, func(t *testing.T) {
			client := newIntegrationClient(rs)
			types := fetchPropertyTypes(t, client)

			for _, tc := range rs.filterTests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					props, conds := withResolvedTypes(types, tc.entitySet, tc.properties, tc.filterConditions)

					resp, err := client.Get(context.Background(), tc.entitySet, props, conds)
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, resp.StatusCode,
						"server rejected filter — check filter syntax for version %s", rs.version)

					body, err := io.ReadAll(resp.Body)
					require.NoError(t, err)
					defer resp.Body.Close()

					entries, err := odata.MapToResponse(body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedResults, len(entries))

					unfilteredTotal := fetchUnfilteredCount(t, rs, tc.entitySet)
					assert.Greater(t, unfilteredTotal, len(entries),
						"unfiltered total (%d) must exceed filtered result count (%d)",
						unfilteredTotal, len(entries))
				})
			}
		})
	}
}
