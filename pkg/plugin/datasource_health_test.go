package plugin

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckHealth(t *testing.T) {
	im := managerMock{}
	ds := ODataSource{&im}

	is := ODataSourceInstance{GetOC("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})}

	im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	assert.NoError(t, err)
	assert.Equal(t, result.Status, backend.HealthStatusOk)
}

func TestCheckHealthWithError(t *testing.T) {
	im := managerMock{}
	ds := ODataSource{&im}

	is := ODataSourceInstance{GetOC("/not/found", nil)}

	im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	assert.NoError(t, err)
	assert.Equal(t, result.Status, backend.HealthStatusError)
}

func TestCheckHealthTimeout(t *testing.T) {
	im := managerMock{}
	ds := ODataSource{&im}

	is := ODataSourceInstance{GetOC("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})}
	im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	assert.NoError(t, err)
	assert.Equal(t, result.Status, backend.HealthStatusError)
}
