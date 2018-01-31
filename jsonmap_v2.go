package jsonmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	String Type = iota
	Number
	Object
	Array
	Bool
	Null
)

type Type int

type JSONMapV2 map[string]JSONNode

type JSONMapV3 map[string]interface{}

type JSONNode struct {
	Key       string
	ValueType Type
	value     interface{}
}

// func NewV2(b []byte) (JSONMapV2, error) {
// 	var m map[string]interface{}

// 	var jm JSONMapV2

// 	if err := json.Unmarshal(b, &m); err != nil {
// 		return jm, err
// 	}

// 	jm = buildMap(m)

// 	return jm, nil
// }

// func buildMap(m map[string]interface{}) JSONMapV2 {
// 	var jns = make(JSONMapV2)

// 	for key, value := range m {
// 		var j JSONNode

// 		if vv, ok := value.(map[string]interface{}); ok {
// 			j = JSONNode{key, Object, buildMap(vv)}
// 		} else {
// 			j = buildNode(key, value)
// 		}

// 		jns[key] = j
// 	}

// 	return jns
// }

// func buildNode(key string, value interface{}) JSONNode {
// 	j := JSONNode{}

// 	switch value.(type) {
// 	case float64:
// 		j.ValueType = Number
// 		j.value = value.(float64)
// 	case string:
// 		j.ValueType = String
// 		j.value = value.(string)
// 	case bool:
// 		j.ValueType = Bool
// 		j.value = value.(bool)
// 	case nil:
// 		j.ValueType = Null
// 		j.value = nil
// 	case map[string]interface{}:
// 		j.ValueType = Object
// 		j.value = buildMap(value.(map[string]interface{}))
// 	case []interface{}:
// 		values := value.([]interface{})

// 		var nodes []interface{}

// 		for _, v := range values {

// 			switch v.(type) {
// 			case map[string]interface{}:
// 				nodes = append(nodes, buildMap(v.(map[string]interface{})))
// 			default:
// 				nodes = append(nodes, buildNode(key, v))
// 			}
// 		}

// 		j.ValueType = Array
// 		j.value = nodes
// 	}

// 	return j
// }

func NewV3(b []byte) (JSONMapV3, error) {
	var jm JSONMapV3

	if err := json.Unmarshal(b, &jm); err != nil {
		return jm, err
	}

	return jm, nil
}

func (jm JSONMapV3) Get(path string) (interface{}, bool) {
	keys := strings.Split(path, ".")

	value, found := getFromMap(jm, keys)

	return value, found
}

func (jm JSONMapV3) Map(path string) (JSONMapV3, error) {
	var m JSONMapV3
	var err error

	it, found := jm.Get(path)

	fmt.Println(reflect.TypeOf(it))
	if found {
		if v, ok := it.(map[string]interface{}); ok {
			m = JSONMapV3(v)
		} else {
			err = fmt.Errorf("jsonmap: item for path '%v' is not of type map", path)
		}
	} else {
		err = fmt.Errorf("jsonmap: path '%v' does not exist", path)
	}

	return m, err
}

func (jm JSONMapV3) String(path string) (string, error) {
	var s string
	var err error

	it, found := jm.Get(path)

	if found {
		if v, ok := it.(string); ok {
			s = v
		} else {
			err = fmt.Errorf("jsonmap: item for path '%v' is not of type string", path)
		}
	} else {
		err = fmt.Errorf("jsonmap: path '%v' does not exist", path)
	}

	return s, err
}

func (jm JSONMapV3) Number(path string) (float64, error) {
	var f float64
	var err error

	it, found := jm.Get(path)

	if found {
		if v, ok := it.(float64); ok {
			f = v
		} else {
			err = fmt.Errorf("jsonmap: item for path '%v' is not of type number", path)
		}
	} else {
		err = fmt.Errorf("jsonmap: path '%v' does not exist", path)
	}

	return f, err
}

func (jm JSONMapV3) Bool(path string) (bool, error) {
	var b bool
	var err error

	it, found := jm.Get(path)

	if found {
		if v, ok := it.(bool); ok {
			b = v
		} else {
			err = fmt.Errorf("jsonmap: item for path '%v' is not of type bool", path)
		}
	} else {
		err = fmt.Errorf("jsonmap: path '%v' does not exist", path)
	}

	return b, err
}

func (jm JSONMapV3) Array(path string) ([]interface{}, error) {
	var arrs []interface{}
	var err error

	it, found := jm.Get(path)

	if found {
		if v, ok := it.([]interface{}); ok {
			arrs = v
		} else {
			err = fmt.Errorf("jsonmap: item for path '%v' is not of type array", path)
		}
	} else {
		err = fmt.Errorf("jsonmap: path '%v' does not exist", path)
	}

	return arrs, err
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
		if len(keys) == 1 {
			vv = it.(map[string]interface{})[keys[0]]
		} else {
			return getFromMap(it.(map[string]interface{}), keys)
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
