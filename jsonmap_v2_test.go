package jsonmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testy = []byte(
	`{
		"a":"foo",
		"b":"bar",
		"c":{
			"d":"baz"
		},
		"e":{
			"f":{
				"g":{
					"h":"nested"
				},
				"i": 10.01
			}
		},
		"j":{
			"d":"baraz"
		},
		"k": [
			{
				"name":"joe"
			},
			{
				"name":"test"
			}
		],
		"l": ["why", "is", "this", "hard"],
		"m": ["foo", 1, "bar", "2", 3],
		"n": false,
		"o":{
			"p":{
				"q":true
			}
		},
		"s": 123.456,
		"t":{
			"q":{
				"r": ["first", "second", "third", "fourth"],
				"w": [1, 2, 3.4, 4]
			}
		},
		"v": [1, 2, 4.5, 6],
		"x": [true, true, false],
		"y":{
			"z": [false, false, true]
		}
	}`,
)

func TestJSONMapV2Get(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, found := jm.Get("a")

	assert.True(t, found)
	assert.Equal(t, "foo", node)

	// top level array

	node, found = jm.Get("l")

	assert.True(t, found)
	assert.Equal(t, []interface{}{"why", "is", "this", "hard"}, node)

	// nested

	node, found = jm.Get("e.f.i")

	assert.True(t, found)
	assert.Equal(t, 10.01, node)

	// nested deep

	node, found = jm.Get("e.f.g.h")

	assert.True(t, found)
	assert.Equal(t, "nested", node)
}

func TestJSONMapV2Get_returnsFalseWhenNotFound(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, found := jm.Get("z")

	assert.Nil(t, node)
	assert.False(t, found)

	// nested

	node, found = jm.Get("e.f.g.z")

	assert.Nil(t, node)
	assert.False(t, found)
}

func TestJSONMapV2Map(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.Map("c")

	assert.Equal(t, JSONMapV2(map[string]interface{}{"d": "baz"}), node)
	assert.NoError(t, err)

	// nested

	node, err = jm.Map("e.f.g")

	assert.Equal(t, JSONMapV2(map[string]interface{}{"h": "nested"}), node)
	assert.NoError(t, err)
}

func TestJSONMapV2String(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.String("a")

	assert.Equal(t, string("foo"), node)
	assert.NoError(t, err)

	// nested

	node, err = jm.String("e.f.g.h")

	assert.Equal(t, string("nested"), node)
	assert.NoError(t, err)
}

func TestJSONMapV2Bool(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.Bool("n")

	assert.Equal(t, false, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.Bool("o.p.q")

	assert.Equal(t, true, node)
	assert.NoError(t, err)
}

func TestJSONMapV2Number(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.Number("s")

	assert.Equal(t, 123.456, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.Number("e.f.i")

	assert.Equal(t, 10.01, node)
	assert.NoError(t, err)
}

func TestJSONMapV2Array(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.Array("l")

	assert.Equal(t, []interface{}{"why", "is", "this", "hard"}, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.Array("t.q.r")

	assert.Equal(t, []interface{}{"first", "second", "third", "fourth"}, node)
	assert.NoError(t, err)
}

func TestJSONMapV2StringArray(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.StringArray("l")

	assert.Equal(t, []string{"why", "is", "this", "hard"}, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.StringArray("t.q.r")

	assert.Equal(t, []string{"first", "second", "third", "fourth"}, node)
	assert.NoError(t, err)
}

func TestJSONMapV2BoolArray(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.BoolArray("x")

	assert.Equal(t, []bool{true, true, false}, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.BoolArray("y.z")

	assert.Equal(t, []bool{false, false, true}, node)
	assert.NoError(t, err)
}

func TestJSONMapV2NumberArray(t *testing.T) {
	jm, err := NewV2(testy)

	assert.NoError(t, err)

	// top level

	node, err := jm.NumberArray("v")

	assert.Equal(t, []float64{1, 2, 4.5, 6}, node)
	assert.NoError(t, err)

	// nested

	node, err = jm.NumberArray("t.q.w")

	assert.Equal(t, []float64{1, 2, 3.4, 4}, node)
	assert.NoError(t, err)
}
