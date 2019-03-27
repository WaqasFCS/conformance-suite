// Code generated by mockery v1.0.0. DO NOT EDIT.

package generation

import discovery "bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
import logrus "github.com/sirupsen/logrus"
import mock "github.com/stretchr/testify/mock"
import model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"

// MockGenerator is an autogenerated mock type for the Generator type
type MockGenerator struct {
	mock.Mock
}

// GenerateManifestTests provides a mock function with given fields: log, config, _a2, ctx
func (_m *MockGenerator) GenerateManifestTests(log *logrus.Entry, config GeneratorConfig, _a2 discovery.ModelDiscovery, ctx *model.Context) TestCasesRun {
	ret := _m.Called(log, config, _a2, ctx)

	var r0 TestCasesRun
	if rf, ok := ret.Get(0).(func(*logrus.Entry, GeneratorConfig, discovery.ModelDiscovery, *model.Context) TestCasesRun); ok {
		r0 = rf(log, config, _a2, ctx)
	} else {
		r0 = ret.Get(0).(TestCasesRun)
	}

	return r0
}

// GenerateSpecificationTestCases provides a mock function with given fields: log, config, _a2, ctx
func (_m *MockGenerator) GenerateSpecificationTestCases(log *logrus.Entry, config GeneratorConfig, _a2 discovery.ModelDiscovery, ctx *model.Context) TestCasesRun {
	ret := _m.Called(log, config, _a2, ctx)

	var r0 TestCasesRun
	if rf, ok := ret.Get(0).(func(*logrus.Entry, GeneratorConfig, discovery.ModelDiscovery, *model.Context) TestCasesRun); ok {
		r0 = rf(log, config, _a2, ctx)
	} else {
		r0 = ret.Get(0).(TestCasesRun)
	}

	return r0
}