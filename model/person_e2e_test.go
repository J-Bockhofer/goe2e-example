package model_test

import (
	"goe2e-example/model"
	"net/http"
	"testing"

	goe2e "github.com/J-Bockhofer/goe2e/pkg"

	"github.com/stretchr/testify/assert"
)

var (
	responseDict = map[string]string{
		"personName": "name",
	}

	env = map[string]interface{}{
		"personName": "",
	}
)

// https://mickey.dev/posts/go-build-tags-testing/
// https://martinyonathann.medium.com/integration-unit-and-e2e-testing-in-golang-3e957f9920dd

func TestPersonPostSimple(t *testing.T) {
	p := model.Person{
		Name: "john",
		Age:  32,
	}
	rc := &goe2e.TestConfig{
		Name: "POST /persons",
		SpecOpts: []goe2e.SpecOption{
			goe2e.WithMethod(http.MethodPost),
			goe2e.WithUrl("http://localhost:8080/persons/"),
			goe2e.WithJSON(&p),
		},
		RequestMods: []goe2e.RequestModifier{
			goe2e.WithContentType(goe2e.ContentHeaderJSON),
		},
		PreTestStatements: []goe2e.TestStatement{
			{Description: "request not nil", Statement: func(t *testing.T, rh *goe2e.RequestHandler) {
				assert.NotNil(t, rh.GetRequest())
			}},
		},
		PostTestStatements: []goe2e.TestStatement{
			{Description: "body not nil", Statement: func(t *testing.T, rh *goe2e.RequestHandler) {
				assert.NotNil(t, rh.ResponseBody)
			}},
			{Description: "status 202", Statement: func(t *testing.T, rh *goe2e.RequestHandler) {
				assert.Equal(t, http.StatusAccepted, rh.Response.StatusCode)
			}},
		},
	}
	goe2e.TestRequest(t, rc)
}

func TestSecurePersonPost(t *testing.T) {
	// if os.Getenv("E2E") == "" {
	// 	t.Skip()
	// }
	p := model.Person{
		Name: "john",
		Age:  32,
	}
	tc := &goe2e.TestConfig{
		Name: "POST /persons",
		SpecOpts: []goe2e.SpecOption{
			goe2e.WithMethod(http.MethodPost),
			goe2e.WithUrl("http://localhost:8080/secure/persons/"),
			goe2e.WithJSON(&p),
		},
		RequestMods: []goe2e.RequestModifier{
			goe2e.WithContentType(goe2e.ContentHeaderJSON),
			goe2e.WithHeaders(goe2e.D{
				"Authorization": "john",
			}),
		},
		ResponseBodyMods: []goe2e.ResponseBodyModifier{
			goe2e.ResponseJSONToEnv(env, responseDict),
		},
		PostTestStatements: []goe2e.TestStatement{
			{Description: "body not nil", Statement: func(t *testing.T, rh *goe2e.RequestHandler) {
				assert.NotNil(t, rh.ResponseBody)
			}},
			{Description: "status 202", Statement: goe2e.TestStatusCode(http.StatusAccepted)},
			{Description: "name written to env", Statement: func(t *testing.T, rh *goe2e.RequestHandler) {
				assert.Equal(t, env["personName"], "john")
			}},
		},
	}
	goe2e.TestRequest(t, tc)
}
