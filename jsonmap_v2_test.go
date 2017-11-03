package jsonmap

import (
	"fmt"
	"testing"
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
				}
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
		"l": ["why", "is", "this", "hard"]
	}`,
)

func TestV2(t *testing.T) {
	jm, err := NewV2(testy)

	if err != nil {
		fmt.Println(err)
	}

	for _, node := range jm.Nodes {
		if node.Key == "l" {
			fmt.Println(node.Value.([]JSONNode)[0].Value)
		} else if node.Key == "k" {
			fmt.Println(node.Value.([]JSONNode)[0].Key)
		} else if node.Key == "e" {
			fmt.Println(node.ValueType)
		}
	}

}
