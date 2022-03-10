package gopkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSize(t *testing.T) {
	tests := []struct {
		args uint64
		want string
	}{
		{346, "346B"},
		{3467, "3.39K"},
		{346778, "338.65K"},
		{1200346778, "1.12G"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, DataSize(tt.args))
	}
}

func TestPrettyJSON(t *testing.T) {
	tests := []interface{}{
		map[string]int{"a": 1},
		struct {
			A int `json:"a"`
		}{1},
	}
	want := `{
    "a": 1
}`
	for _, sample := range tests {
		got, err := PrettyJSON(sample)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}
}

func TestStringsToInts(t *testing.T) {
	is := assert.New(t)

	ints, err := StringsToInts([]string{"1", "2"})
	is.Nil(err)
	is.Equal("[]int{1, 2}", fmt.Sprintf("%#v", ints))

	_, err = StringsToInts([]string{"a", "b"})
	is.Error(err)
}

func TestHowLongAgo(t *testing.T) {
	tests := []struct {
		args int64
		want string
	}{
		{-36, "unknown"},
		{36, "36 secs"},
		{346, "5 mins"},
		{3467, "57 mins"},
		{346778, "4 days"},
		{1200346778, "13892 days"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, HowLongAgo(tt.args))
	}
}
