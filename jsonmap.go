package jsonmap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JsonMap map[string]interface{}

func New(data []byte) JsonMap {
	var j JsonMap

	if err := json.Unmarshal(data, &j); err == nil {
		fmt.Println(err)

		return nil
	}

	return j
}

func (j JsonMap) Get(path string) (interface{}, error) {
	keys := strings.Split(path, ".")

	jj := j

	for i, key := range keys {
		if _, ok := jj[key]; !ok {
			return nil, fmt.Errorf("jsonmap: key '%v' does not exist", key)
		}

		if i == len(keys)-1 {
			return jj[key], nil
		}

		switch jj[key].(type) {
		case map[string]interface{}:
			jj = jj[key].(map[string]interface{})
		}
	}

	return jj, nil
}

func (j JsonMap) String(path string) (string, error) {
	val, err := j.Get(path)

	if err != nil || val == nil {
		return "", err
	}
	return val.(string), nil
}

func (j JsonMap) Int(path string) (int, error) {
	val, err := j.String(path)

	if err != nil {
		return 0, err
	}

	i, convErr := strconv.Atoi(val)
	if convErr != nil {
		return 0, convErr
	}
	return i, nil
}

func (j JsonMap) Float(path string) (float64, error) {
	val, err := j.String(path)

	if err != nil {
		return 0.0, err
	}

	f, convErr := strconv.ParseFloat(val, 64)
	if convErr != nil {
		return 0.0, convErr
	}
	return f, nil
}

func (j JsonMap) Bool(path string) (bool, error) {
	val, err := j.String(path)

	if err != nil {
		return false, err
	}

	b, convErr := strconv.ParseBool(val)
	if convErr != nil {
		return false, convErr
	}
	return b, nil
}

func (j JsonMap) Array(path string) ([]JsonMap, error) {
	val, err := j.Get(path)

	if err != nil || val == nil {
		return nil, err
	}
	maps := make([]JsonMap, 0)
	arr := val.([]interface{})
	for _, v := range arr {
		jMap := JsonMap(v.(map[string]interface{}))
		maps = append(maps, jMap)
	}
	return maps, nil
}

func (j JsonMap) Find(key string) interface{} {
	if strings.Contains(key, ".") {
		return findSubpath(key, j)
	} else {
		return find(key, j)
	}
}

func find(key string, tree map[string]interface{}) interface{} {
	treeCopy := tree

	if f, ok := treeCopy[key]; ok {
		return f
	}

	for _, val := range treeCopy {
		switch val.(type) {
		case map[string]interface{}:
			m := val.(map[string]interface{})

			found := find(key, m)

			if found != nil {
				return found
			}
		}
	}
	return nil
}

func findSubpath(sub string, j JsonMap) interface{} {
	segs := strings.Split(sub, ".")
	subpath := strings.Join(segs[1:], ".")

	f := find(segs[0], j)

	foundMap := JsonMap(f.(map[string]interface{}))

	found, err := foundMap.Get(subpath)

	if err != nil {
		return nil
	}

	return found
}
