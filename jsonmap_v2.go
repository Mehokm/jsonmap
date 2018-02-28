package jsonmap

import (
	"encoding/json"
	"strings"
)

type notFoundError struct {
	path string
}

func (nfe notFoundError) Error() string {
	return "JSONMapV2: path '" + nfe.path + "' does not exist"
}

type incorrectTypeError struct {
	t, path string
}

func (nte incorrectTypeError) Error() string {
	return "JSONMapV2: item for path '" + nte.path + "' is not of type " + nte.t
}

type JSONMapV2 map[string]interface{}

func NewV2(b []byte) (JSONMapV2, error) {
	var jm JSONMapV2

	if err := json.Unmarshal(b, &jm); err != nil {
		return jm, err
	}

	return jm, nil
}

func (jm JSONMapV2) Get(path string) (interface{}, bool) {
	keys := strings.Split(path, ".")

	return getFromMap(jm, keys)
}

func (jm JSONMapV2) Map(path string) (JSONMapV2, error) {
	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	if v, ok := it.(map[string]interface{}); ok {
		return JSONMapV2(v), nil
	}

	return nil, incorrectTypeError{"map", path}
}

func (jm JSONMapV2) String(path string) (string, error) {
	it, found := jm.Get(path)

	if !found {
		return "", notFoundError{path}
	}

	if v, ok := it.(string); ok {
		return v, nil
	}

	return "", incorrectTypeError{"string", path}
}

func (jm JSONMapV2) Number(path string) (float64, error) {
	it, found := jm.Get(path)

	if !found {
		return 0, notFoundError{path}
	}

	if v, ok := it.(float64); ok {
		return v, nil
	}

	return 0, incorrectTypeError{"number", path}
}

func (jm JSONMapV2) Bool(path string) (bool, error) {
	it, found := jm.Get(path)

	if !found {
		return false, notFoundError{path}
	}

	if v, ok := it.(bool); ok {
		return v, nil
	}

	return false, incorrectTypeError{"boolean", path}
}

// Array methods

func (jm JSONMapV2) Array(path string) ([]interface{}, error) {
	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	if v, ok := it.([]interface{}); ok {
		return v, nil
	}

	return nil, incorrectTypeError{"array", path}
}

func (jm JSONMapV2) MapArray(path string) ([]JSONMapV2, error) {
	var arrs []JSONMapV2

	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	for _, v := range it.([]interface{}) {
		if vv, ok := v.(map[string]interface{}); ok {
			arrs = append(arrs, vv)
		}
	}

	return arrs, nil
}

func (jm JSONMapV2) StringArray(path string) ([]string, error) {
	var arrs []string

	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	for _, v := range it.([]interface{}) {
		if vv, ok := v.(string); ok {
			arrs = append(arrs, vv)
		}
	}

	return arrs, nil
}

func (jm JSONMapV2) NumberArray(path string) ([]float64, error) {
	var arrs []float64

	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	for _, v := range it.([]interface{}) {
		if vv, ok := v.(float64); ok {
			arrs = append(arrs, vv)
		}
	}

	return arrs, nil
}

func (jm JSONMapV2) BoolArray(path string) ([]bool, error) {
	var arrs []bool

	it, found := jm.Get(path)

	if !found {
		return nil, notFoundError{path}
	}

	for _, v := range it.([]interface{}) {
		if vv, ok := v.(bool); ok {
			arrs = append(arrs, vv)
		}
	}

	return arrs, nil
}

func getFromMap(m map[string]interface{}, keys []string) (interface{}, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	value, ok := m[keys[0]]

	if ok && len(keys) > 1 {
		return get(value, keys[1:])
	} else if ok && len(keys) <= 1 {
		return value, ok
	}

	return nil, false
}

func get(it interface{}, keys []string) (interface{}, bool) {
	var vv interface{}

	switch it.(type) {
	case float64:
		vv = it.(float64)
	case string:
		vv = it.(string)
	case bool:
		vv = it.(bool)
	case nil:
		vv = nil
	case map[string]interface{}:
		m := it.(map[string]interface{})

		if _, ok := m[keys[0]]; !ok {
			return nil, false
		}

		if len(keys) == 1 {
			vv = m[keys[0]]
		} else {
			return getFromMap(m, keys)
		}
	case []interface{}:
		values := it.([]interface{})

		var vv []interface{}

		for _, v := range values {
			switch v.(type) {
			case map[string]interface{}:
				vv = append(vv, v.(map[string]interface{}))
			default:
				cv, found := get(v, keys)

				if found {
					vv = append(vv, cv)
				}
			}

		}
	default:
		return nil, false
	}

	return vv, true
}
