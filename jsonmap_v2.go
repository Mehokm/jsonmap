package jsonmap

import (
	"encoding/json"
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

type JSONMapV2 struct {
	Nodes []JSONNode
}

type JSONNode struct {
	Key       string
	ValueType Type
	Value     interface{}
}

func NewV2(b []byte) (JSONMapV2, error) {
	var m map[string]interface{}

	var jm JSONMapV2

	if err := json.Unmarshal(b, &m); err != nil {
		return jm, err
	}

	jm.Nodes = buildNodes(m)

	return jm, nil
}

func buildNodes(m map[string]interface{}) []JSONNode {
	var jns []JSONNode

	for key, value := range m {
		var j JSONNode

		if vv, ok := value.(map[string]interface{}); ok {
			j = JSONNode{key, Object, buildNodes(vv)}
		} else {
			j = buildNode(key, value)
		}

		jns = append(jns, j)
	}

	return jns
}

func buildNode(key string, value interface{}) JSONNode {
	j := JSONNode{Key: key}

	switch value.(type) {
	case float64:
		j.ValueType = Number
		j.Value = value.(float64)
	case string:
		j.ValueType = String
		j.Value = value.(string)
	case bool:
		j.ValueType = Bool
		j.Value = value.(bool)
	case nil:
		j.ValueType = Null
		j.Value = nil
	case map[string]interface{}:
		j.ValueType = Object
		j.Value = buildNodes(value.(map[string]interface{}))
	case []interface{}:
		values := value.([]interface{})

		var nodes []JSONNode

		for _, v := range values {

			switch v.(type) {
			case map[string]interface{}:
				nodes = append(nodes, buildNodes(v.(map[string]interface{}))...)
			default:
				nodes = append(nodes, buildNode(key, v))
			}
		}

		j.ValueType = Array
		j.Value = nodes
	}

	return j
}
