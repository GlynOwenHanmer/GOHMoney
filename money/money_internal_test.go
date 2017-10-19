package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialiseIfRequired(t *testing.T) {
	m := Money{}
	assert.Nil(t, m.inner)
	initialiseIfRequired(&m)
	expected, _ := New(0, "")
	assert.Equal(t, *expected, m)
	assert.NotNil(t, m.inner)
	assert.Equal(t, int64(0), m.Amount())
	c, err := m.Currency()
	assert.NotNil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, c.Code, "")
}