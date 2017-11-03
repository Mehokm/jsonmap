package jsonmap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSONMap is a wrapper type for generic JSON decoding
type JSONMap map[string]interface{}

func New(data []byte) JSONMap {
	var j JSONMap

	if err := json.Unmarshal(data, &j); err != nil {
		fmt.Println(err)

		return nil
	}

	return j
}

func (j JSONMap) Get(path string) (interface{}, error) {
	keys := strings.Split(path, ".")

	jj := j

	for i, key := range keys {
		value, exists := jj[key]

		if !exists {
			return nil, fmt.Errorf("jsonmap: key '%v' does not exist", key)
		}

		if i == len(keys)-1 {
			return value, nil
		}

		switch value.(type) {
		case map[string]interface{}:
			jj = value.(map[string]interface{})
		}
	}

	return jj, nil
}

func (j JSONMap) String(path string) (string, error) {
	val, err := j.Get(path)

	if err != nil || val == nil {
		return "", err
	}
	return val.(string), nil
}

func (j JSONMap) Int(path string) (int, error) {
	val, err := j.String(path)

	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(val)

	if err != nil {
		return 0, err
	}

	return i, nil
}

func (j JSONMap) Float(path string) (float64, error) {
	val, err := j.String(path)

	if err != nil {
		return 0.0, err
	}

	f, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return 0.0, err
	}

	return f, nil
}

func (j JSONMap) Bool(path string) (bool, error) {
	val, err := j.String(path)

	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(val)

	if err != nil {
		return false, err
	}

	return b, nil
}

func (j JSONMap) Array(path string) ([]JSONMap, error) {
	val, err := j.Get(path)

	if err != nil || val == nil {
		return nil, err
	}

	jArrays, ok := val.([]interface{})

	if !ok {
		return nil, fmt.Errorf("jsonmap: key '%v' is not of type array", path)
	}

	jMaps := make([]JSONMap, len(jArrays))

	for i, arr := range jArrays {
		jMaps[i] = JSONMap(arr.(map[string]interface{}))
	}

	return jMaps, nil
}

func (j JSONMap) Find(key string) interface{} {
	if strings.Contains(key, ".") {
		return findSubpath(key, j)
	}

	return find(key, j)
}

func find(key string, m map[string]interface{}) interface{} {
	mm := m

	if found, ok := mm[key]; ok {
		return found
	}

	for _, val := range mm {
		switch val.(type) {
		case map[string]interface{}:
			sm := val.(map[string]interface{})

			found := find(key, sm)

			if found != nil {
				return found
			}
		}
	}
	return nil
}

func findSubpath(path string, j JSONMap) interface{} {
	keys := strings.Split(path, ".")

	subpath := strings.Join(keys[1:], ".")

	jj := JSONMap(find(keys[0], j).(map[string]interface{}))

	if found, err := jj.Get(subpath); err == nil {
		return found
	}

	return nil
}
