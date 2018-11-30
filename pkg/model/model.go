package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
)

// Manifest is the high level container for test suite definition
// It contains a list of all the rules required to be passed for conformance testing
// Each rule can have multiple testcases which contribute to testing that particular rule
// So essentially Manifest is a container
type Manifest struct {
	Context     string    `json:"@context"`         // JSONLD contest reference
	ID          string    `json:"@id"`              // JSONLD ID reference
	Type        string    `json:"@type"`            // JSONLD Type reference
	Name        string    `json:"name"`             // Name of the manifiest
	Description string    `json:"description"`      // Description of the Mainfest and what it contains
	BaseIri     string    `json:"baseIri"`          // Base Iri
	Sections    []Context `json:"section_contexts"` // Section specific contexts
	Rules       []Rule    `json:"rules"`            // All the rules in the Manifest
}

// Rule - Define a specific location within a specification that is being tested
// Rule also identifies all the tests that must be passed in order to show that the rule
// implementation in conformant with the specific section in the referenced specification
type Rule struct {
	ID           string           `json:"@id"`             // JSONLD ID reference
	Type         []string         `json:"@type,omitempty"` // JSONLD type reference
	Name         string           `json:"name"`            // A short meaningful name for this rule
	Purpose      string           `json:"purpose"`         // The purpose of this rule
	Specref      string           `json:"specref"`         // Description of area of spec/name/version/section under test
	Speclocation string           `json:"speclocation"`    // specific http reference to location in spec under test covered by this rule
	Tests        [][]TestCase     `json:"tests"`           // Tests - allows for many testcases - array of arrays - to be associated with this rule
	Executor     TestCaseExecutor // TestCaseExecutor interface allow different testcase execution strategies
}

// TestCaseExecutor defines an interface capable of executing a testcase
type TestCaseExecutor interface {
	//ExecuteTestCase(r *http.Request, t *TestCase, ctx *Context) (*http.Response, error)
	ExecuteTestCase(r *http.Request, t *TestCase, ctx *Context) (*http.Response, error)
}

// TestCase defines a test that will be run and needs to be passed as part of the conformance suite
// in order to determine implementation conformance to a specification.
// Testcase have three major sections
// Input:
//     Defines the inputs that are required by the testcase. This effectively involves preparing the http request object
// Context:
//     Provides a link between Discovery information and the testcase
// Expects:
//     Examines the http response to the testcase Input in order to determine if the expected conditions existing in the response
//     and therefore the testcase has passed
//
type TestCase struct {
	ID         string        `json:"@id,omitempty"`     // JSONLD ID Reference
	Type       []string      `json:"@type,omitempty"`   // JSONLD type array
	Name       string        `json:"name,omitempty"`    // Name
	Purpose    string        `json:"purpose,omitempty"` // Purpose of the testcase in simple words
	Input      Input         `json:"input,omitempty"`   // Input Object
	Context    Context       `json:"context,omitempty"` // Local Context Object
	Expect     Expect        `json:"expect,omitempty"`  // Expected object
	ParentRule *Rule         // Allows accessing parent Rule
	Request    *http.Request // The request that's been generated in order to call the endpoint
	Header     http.Header   // ResponseHeader
	Body       string        // ResponseBody
}

// Prepare a Testcase for execution at and endpoint,
// results in a standard http request that encapsulates the testcase request
// as defined in the test case object with any context inputs/replacements etc applied
func (t *TestCase) Prepare(ctx *Context) (*http.Request, error) {
	req, err := t.ApplyInput(ctx)
	if err != nil {
		return nil, err
	}
	req, err = t.ApplyContext() // Apply Context at end of creating request
	return req, err
}

// Validate takes the http response that results as a consequence of sending the testcase http
// request to the endpoint implementation. Validate is responsible for checking the http status
// code and running the set of 'Matches' within the 'Expect' object, to determine if all the
// match conditions are met - which would mean the validation passed.
// The context object is passed as part of the validation as its allows the match clauses to
// examine the request object and 'push' response variables into the context object for use
// in downstream test cases which are potentially part of this testcase sequence
// returns true - validation successful
//         false - validation unsuccessful
//         error - adds detail to validation failure
//         TODO - cater for returning multiple validation failures and explanations
//         NOTE: Vadiate will only return false if a check fails - no checks = true
func (t *TestCase) Validate(resp *http.Response, rulectx *Context) (bool, error) {
	if len(t.Body) == 0 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		t.Body = string(responseBody)
	}
	t.Header = resp.Header
	return t.ApplyExpects(resp, rulectx)
}

// Input defines the content of the http request object used to execute the test case
// Input is built up typically from the openapi/swagger definition of the method/endpoint for a particualar
// specification. Additional properties/fields/headers can be added or change in order to setup the http
// request object of the specific test case. Once setup correctly,the testcase gives the http request object
// to the parent Rule which determine how to execute the requestion object. On execution an http response object
// is received and passed back to the testcase for validation using the Expects object.
type Input struct {
	Method     string          `json:"method,omitempty"`     // http Method that this test case uses
	Endpoint   string          `json:"endpoint,omitempty"`   // resource endpoint where the http object needs to be sent to get a response
	ContextGet ContextAccessor `json:"contextGet,omitempty"` // Allows retrieval of context variables an input parameters
}

// Context is intended to handle two types of object and make them available to various parts of the suite including
// testcases. The first set are objects created as a result of the discovery phase, which capture discovery model
// information like endpoints and conditional implementation indicators. The other set of data is information passed
// between a sequeuence of test cases, for example AccountId - extracted from the output of one testcase (/Accounts) and fed in
// as part of the input of another testcase for example (/Accounts/{AccountId}/transactions}
type Context map[string]interface{}

// Expect defines a structure for expressing testcase result expectations.
type Expect struct {
	StatusCode       int  `json:"status-code,omitempty"`       // Http response code
	SchemaValidation bool `json:"schema-validation,omitempty"` // Flag to indicate if we need schema validation -
	// provides the ability to switch off schema validation
	Matches    []Match         `json:"matches,omitempty"`    // An array of zero or more match items which must be 'passed' for the testcase to succeed
	ContextPut ContextAccessor `json:"contextPut,omitempty"` // allows storing of test response fragments in context variables
}

// ApplyInput - creates an HTTP request for this test case
// The reason why we're doing this is that a testcase behaves like an http object
// It produces an http.Request - which can be sent to a server
// It consumes and http.Response - which it uses to validate the response against "Expects"
// TestCase lifecycle:
//     Create a Testcase Object
//     Create / retrieve the http request object
//     Apply context information to the request object
//     Rule - manages passing the request object from the testcase to an appropriate endpoint handler (like the proxy)
//     Rule - receives http response from endpoint and provides it back to testcase
//     Testcase evaluates the http response object using its 'Expects' clause
//     Testcase passes or fails depending on the 'Expects' outcome
func (t *TestCase) ApplyInput(rulectx *Context) (*http.Request, error) {
	// NOTE: This is an initial implementation to get things moving - expect a lot of change in this function
	var err error

	err = t.Input.ContextGet.GetValues(t, rulectx)

	if &t.Input.Endpoint == nil || &t.Input.Method == nil { // we don't have a value input object
		return nil, errors.New("Testcase Input empty")
	}
	if t.Input.Method != "GET" { // Only get Supported Initially
		return nil, errors.New("Testcase Method Only support GET currently")
	}
	req, err := http.NewRequest(t.Input.Method, t.Input.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	t.Request = req // store the request in the testcase

	return req, err
}

// ApplyContext - at the end of ApplyInputs on the testcase - we have an initial http request object
// ApplyContext, applys context parameters to the http object.
// Context parameter typically involve variables that originaled in discovery
// The functionality of ApplyContext will grow significantly over time.
func (t *TestCase) ApplyContext() (*http.Request, error) {
	base := t.Context.Get("baseurl")
	if base != nil {
		urlWithBase, err := url.Parse(base.(string) + t.Input.Endpoint) // expand url in request to be full pathname including Discovery endpoint info from context
		if err != nil {
			return nil, errors.New("Error parsing context baseURL: (" + base.(string) + ")")
		}
		t.Request.URL = urlWithBase
	}
	return t.Request, nil
}

// ApplyExpects runs the Expects section of the testcase to evaluate if the response from the system under test passes or fails
// The Expects section of a testcase can contain multiple conditions that need to be met to pass a testcase
// When a test fails, ApplyExpects is responsible for reporting back information about the failure, why it occured, where it occured etc.
//
// The ApplyExpect section is also responsible for running and contextPut clauses.
// contextPuts are responsible for updated context variables with values selected from the test case response
// contextPuts will only be executed if the ApplyExpects standards match tests pass
// if any of the ApplyExpects match tests fail - ApplyExpects returns false and contextPuts aren't executed
func (t *TestCase) ApplyExpects(res *http.Response, rulectx *Context) (bool, error) {
	if res == nil { // if we've not got a response object to check, always return false
		return false, errors.New("nil http.Response - cannot process ApplyExpects")
	}

	if t.Expect.StatusCode != res.StatusCode { // Status codes don't match
		return false, fmt.Errorf("(%s):%s: HTTP Status code does not match: expected %d got %d", t.ID, t.Name, t.Expect.StatusCode, res.StatusCode)
	}

	for _, match := range t.Expect.Matches {
		checkResult, got := match.Check(t)
		if checkResult == false {
			return false, got
		}
	}

	_, err := t.Expect.ContextPut.PutValues(t, rulectx)

	return true, err
}

// Get the key form the Context map - currently assumes value converts easily to a string!
func (c Context) Get(key string) interface{} {
	return c[key]
}

// Put a value indexed by 'key' into the context. The value can be any type
func (c Context) Put(key string, value interface{}) {
	c[key] = value
}

// NewTestCase creates an Context thats initialised correctly with a map structure
// which holds the context parameters
func NewTestCase() *TestCase {
	var t TestCase
	t.Context = make(map[string]interface{})
	return &t
}

// GetIncludedPermission returns the list of permission names that need to be included
// in the access token for this testcase. See permission model docs for more information
//
func (t *TestCase) GetIncludedPermission() []string {
	var result []string
	if t.Context["permissions"] != nil {
		permissionArray := t.Context["permissions"].([]interface{})
		for _, permissionName := range permissionArray {
			result = append(result, permissionName.(string))
		}
		return result
	}

	// for defaults to apply there should be no permissions no permissions_excluded specified
	if t.Context["permissions"] == nil && t.Context["permissions_excluded"] == nil {
		// Attempt to get default permissions
		perms := GetPermissionsForEndpoint(t.Input.Endpoint)
		if len(perms) > 1 { // need to figure out default
			for _, p := range perms { // find default permission
				if p.Default == true {
					return []string{p.Permission}
				}
			}
		} else {
			if len(perms) > 0 { // only one permission so return that
				return []string{perms[0].Permission}
			}
		}
		return []string{} // no defaults - no permisions
	}

	if t.Context["permissions"] == nil {
		return []string{}
	}
	return result
}

// GetExcludedPermissions return a list of excluded permissions
func (t *TestCase) GetExcludedPermissions() []string {
	var permissionArray []interface{}
	var result []string
	if t.Context["permissions_excluded"] == nil {
		return []string{}
	}
	permissionArray = t.Context["permissions_excluded"].([]interface{})
	if permissionArray == nil {
		return []string{}
	}
	for _, permissionName := range permissionArray {
		result = append(result, permissionName.(string))
	}
	return result
}

// GetPermissions returns a list of Permission objects associated with a testcase
func (t *TestCase) GetPermissions() (included, excluded []string) {
	included = t.GetIncludedPermission()
	excluded = t.GetExcludedPermissions()
	return
}

// Various helpers - main to dump struct contents to console

func (m *Manifest) String() string {
	return fmt.Sprintf("MANIFEST\nName: %s\nDescription: %s\nRules: %d\n", m.Name, m.Description, len(m.Rules))
}

func (r *Rule) String() string {
	return fmt.Sprintf("RULE\nName: %s\nPurpose: %s\nSpecRef: %s\nSpec Location: %s\nTests: %d\n",
		r.Name, r.Purpose, r.Specref, r.Speclocation, len(r.Tests))
}

// Dump - TestCase helper
func (t *TestCase) Dump(print bool) {
	if print {
		fmt.Printf("TESTCASE\nID: %s\nName: %s\nPurpose: %s\n", t.ID, t.Name, t.Purpose)
		fmt.Printf("Input: ")
		pkgutils.DumpJSON(t.Input)
		fmt.Printf("Context: ")
		pkgutils.DumpJSON(t.Context)
		fmt.Printf("Expect: ")
		pkgutils.DumpJSON(t.Expect)
	}
}

// RunTests - runs all the tests for aTestRule
func (r *Rule) RunTests() {
	for _, testSequence := range r.Tests {
		for _, tester := range testSequence {
			// testcase.ApplyInput
			// testcase.ApplyContext
			// testcase.ApplyExpects
			_ = tester // placeholder
			fmt.Println("Test Result: ", true)
		}
	}
}

// GetPermissionSets returns the inclusive and exclusive permission sets required
// to run the tests under this rule.
// Initially the granulatiy of permissionSets will be set at rule level, meaning that one
// included set and one excluded set will cover all the testcases with a rule.
// In future iterations it may be desirable to have per testSequence permissionSets as this
// would allow a finer grained mix of negative permission testing
func (r *Rule) GetPermissionSets() (included, excluded []string) {
	includedSet := NewPermissionSet("included", []string{})
	excludedSet := NewPermissionSet("excluded", []string{})
	for _, testSequence := range r.Tests {
		for _, test := range testSequence {
			i, x := test.GetPermissions()
			includedSet.AddPermissions(i)
			excludedSet.AddPermissions(x)
		}
	}

	return includedSet.GetPermissions(), excludedSet.GetPermissions()
}

// Execute the testcase
// For the rule this effectively equates to sending the assembled http request from
// the testcase to an endpoint (typically ASPSP implemetation) and getting an http.Response
// The http.Request at this point will contain the fully assembled request from a testcase point of view
// - testcase will have likely pulled out appropriate access_tokens/permissions
// - rule will have the opportunited to further decorate this request before passing on
func (r *Rule) Execute(req *http.Request, tc *TestCase) (*http.Response, error) {
	return r.Executor.ExecuteTestCase(req, tc, &Context{})
}