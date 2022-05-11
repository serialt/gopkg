package gopkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	data := []string{"hello", "world", "github"}
	res := Contains(data, "github")
	assert.True(t, res)
}

func TestIndexSlice(t *testing.T) {
	data := []string{"hello", "world", "github"}
	res, ok := IndexSlice(data, "github")
	assert.Equal(t, res, 2)
	assert.True(t, ok)

}

func TestSubslice(t *testing.T) {
	data := []string{"hello", "world", "github"}
	subslice := []string{"world", "github"}
	ok := IsSubSlice(subslice, data)
	assert.True(t, ok)
}

func TestDiffSlice(t *testing.T) {
	data := []string{"hello", "world", "github"}
	subslice := []string{"world", "github"}
	diffslice := []string{"worlds", "github"}

	res, ok := DiffSlice(subslice, data)
	assert.True(t, ok)
	assert.Empty(t, res)

	res, ok = DiffSlice(diffslice, data)
	assert.False(t, ok)
	assert.NotEmpty(t, res)

}
