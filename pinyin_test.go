package sugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPinyin(t *testing.T) {
	data := "音乐"
	got := Pinyin(data)
	want := "yinle"
	assert.Equal(t, want, got)
}
