package results

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTestResult(t *testing.T) {
	result := NewTestCaseResult("123", true)
	assert.Equal(t, "123", result.Id)
	assert.True(t, result.Pass)

	result = NewTestCaseResult("321", false)
	assert.Equal(t, "321", result.Id)
	assert.False(t, result.Pass)
}

func TestNewTestFailResult(t *testing.T) {
	result := NewTestCaseFail("id")
	assert.Equal(t, "id", result.Id)
	assert.False(t, result.Pass)
}