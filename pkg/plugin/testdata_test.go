package plugin

import (
	"encoding/json"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func aQueryDataRequest(queryDataRequestBuilders ...func(*backend.QueryDataRequest)) backend.QueryDataRequest {
	// Default values
	request := &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries:       []backend.DataQuery{},
	}
	for _, build := range queryDataRequestBuilders {
		build(request)
	}
	return *request
}

func aDataQuery(refId string, builders ...func(*backend.DataQuery)) backend.DataQuery {
	query := backend.DataQuery{
		// Default values
		RefID:         refId,
		QueryType:     "",
		MaxDataPoints: 100,
		Interval:      10,
		TimeRange:     aOneDayTimeRange(),
	}
	for _, build := range builders {
		build(&query)
	}
	return query
}

func aQueryModel(builders ...func(*queryModel)) *queryModel {
	model := &queryModel{
		// Initialize with default values
		From: "2022-01-01T00:00:00Z",
		To:   "2022-01-02T00:00:00Z",
		TimeProperty: property{
			Name: "time",
			Type: odata.EdmDateTimeOffset,
		},
		EntitySet: entitySet{
			Name:       "Temperatures",
			EntityType: "TemperatureODataMock.Models.Temperature",
		},
		// No default properties
		// No default filter
	}
	for _, build := range builders {
		build(model)
	}
	return model
}

func aQueryDataResponse(builders ...func(*backend.QueryDataResponse)) backend.QueryDataResponse {
	// Default values
	response := &backend.QueryDataResponse{
		Responses: backend.Responses{},
	}
	for _, build := range builders {
		build(response)
	}
	return *response
}

func aDataResponse(builders ...func(*backend.DataResponse)) backend.DataResponse {
	// Default values
	response := &backend.DataResponse{}
	for _, build := range builders {
		build(response)
	}
	return *response
}

func aDataFrame(frameName string, builders ...func(*data.Frame)) *data.Frame {
	// Default values
	frame := &data.Frame{
		Name: frameName,
		Meta: &data.FrameMeta{},
	}
	frame.Meta.PreferredVisualization = data.VisTypeTable
	labels, _ := data.LabelsFromString("time=" + "time")
	frame.Fields = append(
		frame.Fields,
		data.NewField("time", labels, []*time.Time{}),
	)
	for _, build := range builders {
		build(frame)
	}
	return frame
}

func anOdataResponse(builders ...func(*odata.Response)) odata.Response {
	response := odata.Response{}
	for _, build := range builders {
		build(&response)
	}
	return response
}

func anEntity(builders ...func(map[string]interface{})) map[string]interface{} {
	var entity = map[string]interface{}{}
	for _, build := range builders {
		build(entity)
	}
	return entity
}

func withEntity(builders ...func(map[string]interface{})) func(n *odata.Response) {
	return func(response *odata.Response) {
		var entity = anEntity(builders...)
		response.Value = append(response.Value, entity)
	}
}

func withProp(name string, value interface{}) func(e map[string]interface{}) {
	return func(e map[string]interface{}) {
		e[name] = value
	}
}

func withDefaultEntity(builders ...func(map[string]interface{})) func(n *odata.Response) {
	return func(response *odata.Response) {
		var entity = anEntity(builders...)

		entity["time"] = "2022-01-02T00:00:00Z"
		entity["int32"] = 5.0
		entity["boolean"] = true
		entity["string"] = "Hello World!"

		response.Value = append(response.Value, entity)
	}
}

func withDataResponse(name string, builders ...func(*backend.DataResponse)) func(n *backend.QueryDataResponse) {
	return func(dataResponse *backend.QueryDataResponse) {
		var response = aDataResponse(builders...)
		dataResponse.Responses[name] = response
	}
}

func withErrorResponse(err error) func(n *backend.DataResponse) {
	return func(dataResponse *backend.DataResponse) {
		dataResponse.Error = err
	}
}

func withDefaultTestFrame(builders ...func(*data.Frame)) func(n *backend.DataResponse) {
	return func(dataResponse *backend.DataResponse) {
		var frame = aDataFrame("defaultTestFrame", builders...)
		frame.Fields = append(
			frame.Fields,
			data.NewField("int32", nil, []*int32{}),
			data.NewField("boolean", nil, []*bool{}),
			data.NewField("string", nil, []*string{}),
		)
		dataResponse.Frames = append(dataResponse.Frames, frame)
		values := make([]interface{}, 4)
		valueTime, _ := time.Parse(time.RFC3339Nano, "2022-01-02T00:00:00Z")
		valueInt := int32(5)
		valueBool := true
		valueString := "Hello World!"
		values[0] = &valueTime
		values[1] = &valueInt
		values[2] = &valueBool
		values[3] = &valueString
		frame.AppendRow(values...)
	}
}

func withBaseFrame(name string, builders ...func(*data.Frame)) func(n *backend.DataResponse) {
	return func(dataResponse *backend.DataResponse) {
		dataResponse.Frames = append(dataResponse.Frames, aDataFrame(name, builders...))
	}
}

func withField(name string, values interface{}) func(n *data.Frame) {
	return func(frame *data.Frame) {
		frame.Fields = append(frame.Fields, data.NewField(name, nil, values))
	}
}

func withRow(builders ...func(index int, f *data.Frame)) func(n *data.Frame) {
	return func(frame *data.Frame) {
		for i, build := range builders {
			if build != nil {
				build(i, frame)
			} else {
				frame.Fields[i].Append(nil)
			}
		}
	}
}

func withRowValue[T string | int32 | bool | time.Time](value T) func(index int, n *data.Frame) {
	return func(index int, frame *data.Frame) {
		frame.Fields[index].Append(&value)
	}
}

// --- Filter related ---
func someFilterConditions(builders ...func(*filterCondition)) []filterCondition {
	var conditions []filterCondition
	for _, build := range builders {
		condition := aFilterCondition()
		build(condition)
		conditions = append(conditions, *condition)
	}
	return conditions
}

func aFilterCondition(builders ...func(*filterCondition)) *filterCondition {
	condition := &filterCondition{}
	for _, build := range builders {
		build(condition)
	}
	return condition
}

func withFilterConditions(builders ...func(*filterCondition)) func(fs *queryModel) {
	return func(model *queryModel) {
		for _, build := range builders {
			condition := aFilterCondition()
			build(condition)
			model.FilterConditions = append(model.FilterConditions, *condition)
		}
	}
}

func withFilterCondition(property func(*property), operator string, val string) func(n *filterCondition) {
	return func(condition *filterCondition) {
		p := aProperty()
		property(&p)
		condition.Property = p
		condition.Operator = operator
		condition.Value = val
	}
}
func int32Eq5(condition *filterCondition) {
	condition.Property.Name = "int32"
	condition.Operator = "eq"
	condition.Value = "5"
}

func withDataQuery(refId string, builders ...func(*backend.DataQuery)) func(n *backend.QueryDataRequest) {
	return func(request *backend.QueryDataRequest) {
		var query = aDataQuery(refId, builders...)
		request.Queries = append(request.Queries, query)
	}
}

func withQueryModel(builders ...func(*queryModel)) func(n *backend.DataQuery) {
	return func(query *backend.DataQuery) {
		var qm = aQueryModel(builders...)
		model, _ := json.Marshal(qm)

		// Hack to make sure the value for query.JSON contains the queryModel and all base props from DataQuery
		// as this behavior is stated in the docs:
		// "JSON is the raw JSON query and includes the above properties as well as custom properties."
		base := make(map[string]interface{})
		json.Unmarshal(model, &base)

		base["refId"] = query.RefID
		base["queryType"] = query.QueryType
		base["maxDataPoints"] = query.MaxDataPoints
		base["interval"] = query.Interval
		base["timeRange"] = query.TimeRange

		baseRaw, _ := json.Marshal(base)

		var finalJson = baseRaw

		query.JSON = finalJson
	}
}

// --- Property related ---
func aProperty(builders ...func(*property)) property {
	p := property{}
	for _, build := range builders {
		build(&p)
	}
	return p
}

func withProperties(builders ...func(*property)) func(n *queryModel) {
	return func(model *queryModel) {
		for _, build := range builders {
			p := aProperty()
			build(&p)
			model.Properties = append(model.Properties, p)
		}
	}
}

// Properties
func int32Prop(p *property) {
	p.Name = "int32"
	p.Type = odata.EdmInt32
}
func booleanProp(p *property) {
	p.Name = "boolean"
	p.Type = odata.EdmBoolean
}
func stringProp(p *property) {
	p.Name = "string"
	p.Type = odata.EdmString
}

// Misc
func aOneDayTimeRange() backend.TimeRange {
	return backend.TimeRange{
		From: time.Date(2022, 4, 21, 12, 30, 50, 50, time.UTC),
		To:   time.Date(2022, 4, 21, 12, 30, 50, 50, time.UTC)}
}
