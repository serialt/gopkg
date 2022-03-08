package gopkg

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	is := assert.New(t)

	tests := []interface{}{
		2,
		int8(2), int16(2), int32(2), int64(2),
		uint(2), uint8(2), uint16(2), uint32(2), uint64(2),
		float32(2.2), 2.3,
		"2",
		time.Duration(2),
	}
	errTests := []interface{}{
		nil,
		"2a",
		[]int{1},
	}

	// To int
	intVal, err := Int("2")
	is.Nil(err)
	is.Equal(2, intVal)

	intVal, err = ToInt("-2")
	is.Nil(err)
	is.Equal(-2, intVal)

	is.Equal(-2, MustInt("-2"))
	for _, in := range tests {
		is.Equal(2, MustInt(in))
	}
	for _, in := range errTests {
		is.Equal(0, MustInt(in))
	}

	// To uint
	uintVal, err := Uint("2")
	is.Nil(err)
	is.Equal(uint64(2), uintVal)

	_, err = ToUint("-2")
	is.Error(err)

	is.Equal(uint64(0), MustUint("-2"))
	for _, in := range tests {
		is.Equal(uint64(2), MustUint(in))
	}
	for _, in := range errTests {
		is.Equal(uint64(0), MustUint(in))
	}

	// To int64
	i64Val, err := ToInt64("2")
	is.Nil(err)
	is.Equal(int64(2), i64Val)

	i64Val, err = Int64("-2")
	is.Nil(err)
	is.Equal(int64(-2), i64Val)

	for _, in := range tests {
		is.Equal(int64(2), MustInt64(in))
	}
	for _, in := range errTests {
		is.Equal(int64(0), MustInt64(in))
	}
}

func TestToFloat(t *testing.T) {
	is := assert.New(t)

	tests := []interface{}{
		2,
		int8(2), int16(2), int32(2), int64(2),
		uint(2), uint8(2), uint16(2), uint32(2), uint64(2),
		float32(2), float64(2),
		"2",
		time.Duration(2),
	}
	for _, in := range tests {
		is.Equal(float64(2), MustFloat(in))
	}

	is.Equal(123.5, MustFloat("123.5"))
	is.Equal(float64(0), MustFloat("invalid"))

	fltVal, err := ToFloat("123.5")
	is.Nil(err)
	is.Equal(123.5, fltVal)

	fltVal, err = Float("-123.5")
	is.Nil(err)
	is.Equal(-123.5, fltVal)
}

func TestPercent(t *testing.T) {
	assert.Equal(t, float64(34), Percent(34, 100))
	assert.Equal(t, float64(0), Percent(34, 0))
	assert.Equal(t, float64(-100), Percent(34, -34))
}

func TestElapsedTime(t *testing.T) {
	nt := time.Now().Add(-time.Second * 3)
	num := ElapsedTime(nt)

	assert.Equal(t, 3000, int(MustFloat(num)))
}
func TestRandomInt(t *testing.T) {
	min, max := 100, 999

	for i := 0; i < 10; i++ {
		val := RandomInt(min, max)
		fmt.Println(val)

		assert.True(t, val >= min)
		assert.True(t, val <= max)
	}
}