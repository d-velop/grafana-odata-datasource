//go:build integration

package plugin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
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
	expectedResults  int // exact number of entries the server must return
}

type referenceSystem struct {
	name                string
	baseURL             string
	version             string // "V2", "V3", "V4"
	requiresLocalServer bool   // if true: skip when localhost:4004 is not reachable
	filterTests         []filterTestCase
}

var referenceSystems []referenceSystem

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
		httpClient:       &http.Client{Timeout: 30 * time.Second},
		baseUrl:          rs.baseURL,
		odataVersion:     rs.version,
		urlSpaceEncoding: "%20",
	}
}

func isPortOpen(host, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TestIntegration_ServiceRoot(t *testing.T) {
	for _, rs := range referenceSystems {
		rs := rs
		t.Run(rs.name, func(t *testing.T) {
			if rs.requiresLocalServer && !isPortOpen("localhost", "4004") {
				t.Skip("local test server not running — start with: cd test-server && pnpm start")
			}
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
			if rs.requiresLocalServer && !isPortOpen("localhost", "4004") {
				t.Skip("local test server not running — start with: cd test-server && pnpm start")
			}
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
			if rs.requiresLocalServer && !isPortOpen("localhost", "4004") {
				t.Skip("local test server not running — start with: cd test-server && pnpm start")
			}
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
