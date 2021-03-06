package jsonmap

import (
	"errors"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
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
			}
		}`,
	)

	j := New(json)

	a, err := j.Get("a")
	aT := "foo"

	if err != nil || a != aT {
		t.Errorf("var a (%v) does not equal var aT (%v)", a, aT)
	}

	c, err := j.Get("c")
	cT := map[string]string{"d": "baz"}

	if err != nil || reflect.DeepEqual(c, cT) {
		t.Errorf("var c (%v) does not equal var cT (%v)", c, cT)
	}

	c2, err := j.Get("c.d")
	c2T := "baz"

	if err != nil || c2 != c2T {
		t.Errorf("var c2 (%v) does not equal var c2T (%v)", c2, c2T)
	}

	e, err := j.Get("e.f.g.h")
	eT := "nested"

	if err != nil || e != eT {
		t.Errorf("var e (%v) does not equal var eT (%v)", e, eT)
	}

	notExist, err := j.Get("e.f.not.real.path")

	if notExist != nil && err != errors.New("jsonmap: key 'not' does not exist") {
		t.Error("expected 'does not exist' error")
	}
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
