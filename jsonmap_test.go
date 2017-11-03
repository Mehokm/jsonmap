package jsonmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJson = []byte(
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
				}
			}
		},
		"arr": [
			{
				"name": "blah"
			},
			{
				"name": "bleh"
			}
		]
	}`,
)

func TestGet(t *testing.T) {
	j := New(testJson)

	a, err := j.Get("a")

	assert.NoError(t, err)
	assert.Equal(t, "foo", a)

	c, err := j.Get("c")

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"d": "baz"}, c)

	c2, err := j.Get("c.d")

	assert.NoError(t, err)
	assert.Equal(t, "baz", c2)

	e, err := j.Get("e.f.g.h")

	assert.NoError(t, err)
	assert.Equal(t, "nested", e)

	shouldntExist, err := j.Get("e.f.not.real.path")

	assert.Nil(t, shouldntExist)
	assert.EqualError(t, err, "jsonmap: key 'not' does not exist")

	// arr, err := j.Array("arr")

	// if err != nil {
	// 	t.Error("expected array")
	// }

	// fmt.Println(arr[0].String("name"))
}

func TestFind(t *testing.T) {
	json := []byte(
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
					}
				}
			},
			"j":{
				"d":"baraz"
			}
		}`,
	)

	j := New(json)

	d := j.Find("h")
	dT := "nested"

	if d != dT {
		t.Errorf("expected: %v, actual: %v", dT, d)
	}

	gh := j.Find("g.h")
	ghT := "nested"

	if gh != ghT {
		t.Errorf("expected: %v, actual: %v", ghT, gh)
	}

	cd := j.Find("c.d")
	cdT := "baz"

	if cd != cdT {
		t.Errorf("expected: %v, actual: %v", cdT, cd)
	}

	jd := j.Find("j.d")
	jdT := "baraz"

	if jd != jdT {
		t.Errorf("expected: %v, actual: %v", jdT, jd)
	}
}
