package schema

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/pkg/errors"
)

// OpenAPI3Validator - type
type OpenAPI3Validator struct {
	router routers.Router
	doc    *openapi3.T
}

// RequestWrapper -
type RequestWrapper struct {
	Method      string
	URL         string
	ContentType string
	Body        string
}

// ResponseWrapper -
type ResponseWrapper struct {
	Status      int
	ContentType string
	Body        string
}

// internal validation parameters
type validateParams struct {
	httpReq    *http.Request
	route      *routers.Route
	pathParams map[string]string
	statusCode int
	header     http.Header
	body       []byte
}

var headerCT = http.CanonicalHeaderKey("Content-Type")

// NewOpenAPI3Validator - Create a router current just for v3.1.8 of the specifications
// preferring yaml for the spec file
func NewOpenAPI3Validator(specName, version string) (Validator, error) {
	// TODO: @Release - update to use version compare instead of magic number
	if version != "v3.1.8" {
		return nil, fmt.Errorf("NewOpenAPI3Validator - unsupported version: %s", version)
	}
	return buildValidator(specName)
}

// NewRawOpenAPI3Validator -
func NewRawOpenAPI3Validator(specName, version string) (OpenAPI3Validator, error) {
	// TODO: @Release - update to use version compare instead of magic number
	if version != "v3.1.8" {
		return OpenAPI3Validator{}, fmt.Errorf("NewOpenAPI3Validator - unsupported version: %s", version)
	}
	return buildValidator(specName)
}

func buildValidator(specName string) (OpenAPI3Validator, error) {
	router, doc, err := getRouterForSpec(specName)
	return OpenAPI3Validator{router: router, doc: doc}, err
}

// IsRequestProperty -
func (v OpenAPI3Validator) IsRequestProperty(checkmethod, checkpath, propertyPath string) (bool, string, error) {
	spec := v.doc
	for path, props := range spec.Paths {
		for method, op := range getOas3Operations(props) {
			if path == checkpath && method == checkmethod && op.RequestBody != nil {
				for _, param := range op.RequestBody.Value.Content {
					schema := param.Schema.Value
					found, objType := findPropertyInOas3Schema(schema, propertyPath, "")
					if found {
						return true, objType, nil
					}
				}
			}
		}
	}

	return false, "", nil
}

func getRouterForSpec(spec string) (routers.Router, *openapi3.T, error) {
	var filename string
	switch spec {
	// TODO: @Release - update to use version compare instead of magic number
	case "Account and Transaction API Specification":
		filename = "spec/v3.1.8/account-info-openapi.json"
	case "Payment Initiation API":
		filename = "spec/v3.1.8/payment-initiation-openapi.json"
	case "Confirmation of Funds API Specification":
		filename = "spec/v3.1.8/confirmation-funds-openapi.json"
	case "OBIE VRP Profile":
		filename = "spec/v3.1.8/variable-recurring-payments-openapi.json"
	default:
		return nil, nil, errors.New("Cannot get router for spec: " + spec)
	}

	// TODO: @Release - update to a pattern similar to Swagger handling
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(filename)
	if err != nil {
		doc, err = loader.LoadFromFile("pkg/schema/" + filename)
		if err != nil {
			doc, err = loader.LoadFromFile("../../pkg/schema/" + filename)
			if err != nil {
				return nil, nil, fmt.Errorf("Cannot Load OpenApi Spec from file %s, %s", filename, err)
			}
		}
	}
	err = doc.Validate(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot Load OpenApi Spec from file %s, %s", filename, err)
	}

	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot Load OpenApi Router for %s file %s", spec, filename)
	}

	return router, doc, nil
}

// Validate - validates the response
func (v OpenAPI3Validator) Validate(r HTTPResponse) ([]Failure, error) {
	failures := []Failure{}

	serverPath := v.doc.Servers[0].URL
	var path string
	serverIndex := strings.Index(r.Path, serverPath)
	if serverIndex != -1 {
		path = r.Path[serverIndex:]
	} else {
		path = serverPath + r.Path
	}

	httpReq, err := createHTTPReq(r.Method, path)
	if err != nil {
		return nil, err
	}

	route, pathParams, err := v.router.FindRoute(httpReq)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("OpenApi3Validator: error converting body %s", err)
	}

	// check body
	params := validateParams{
		httpReq:    httpReq,
		route:      route,
		pathParams: pathParams,
		statusCode: r.StatusCode,
		header:     r.Header,
		body:       body,
	}

	// accumulate failures
	err = v.validateResponse(params)
	if err != nil {
		return nil, fmt.Errorf("Validate error response:  %s", err.Error())
	}

	return failures, nil
}

func (v OpenAPI3Validator) validateResponse(params validateParams) error {
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    params.httpReq,
		PathParams: params.pathParams,
		Route:      params.route,
	}

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 params.statusCode,
		Header:                 params.header,
		Options: &openapi3filter.Options{
			ExcludeRequestBody:    true,
			IncludeResponseStatus: true,
			MultiError:            false,
		},
	}

	if len(params.body) > 0 {
		responseValidationInput.SetBodyBytes(params.body)
	}

	return openapi3filter.ValidateResponse(context.Background(), responseValidationInput)
}

func (v OpenAPI3Validator) findTestRoute(req *http.Request) (*routers.Route, map[string]string, error) {
	route, pathParams, err := v.router.FindRoute(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %s - findTestRoute:  %s", req.Method, req.URL, err)
	}
	return route, pathParams, err
}

func createHTTPReq(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, strings.NewReader(""))
	req.Header = http.Header{"Content-type": []string{"application/json; charset=utf-8"}}
	return req, err
}

// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOas3Operations(props *openapi3.PathItem) map[string]*openapi3.Operation {
	ops := map[string]*openapi3.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}

	// Keep those != nil
	for key, op := range ops {
		if op == nil {
			delete(ops, key)
		}
	}
	return ops
}

func findPropertyInOas3Schema(sc *openapi3.Schema, propertyPath, previousPath string) (bool, string) {
	for k, j := range sc.Properties {
		var element string
		if len(previousPath) == 0 {
			element = k
		} else {
			element = previousPath + "." + k
		}

		if element == propertyPath {
			return true, fmt.Sprintf("%s", j.Value.Type)
		}

		ret, propType := findPropertyInOas3Schema(j.Value, propertyPath, element)
		if ret {
			return true, propType
		}
	}
	return false, ""
}