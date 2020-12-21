package hjson

import "testing"

func TestAddValue(t *testing.T) {
	testCases := []struct {
		values map[string]interface{}
		value  Value
	}{
		{
			values: map[string]interface{}{
				"key":  nil,
				"key1": false,
				"key2": true,
				"key3": 2,
				"key4": "hello parser",
			},
			value: &JObject{
				values: map[string]Value{
					"key":  JNull{},
					"key1": JBool(false),
					"key2": JBool(true),
					"key3": JNumber(2),
					"key4": JString("hello parser"),
				},
			},
		},
		{
			values: map[string]interface{}{
				"key": map[string]string{
					"key7": "hello",
				},
				"key5": []interface{}{
					true,
					false,
					"key",
					"key",
					struct{}{},
				},
			},
			value: &JObject{
				values: map[string]Value{
					//"key":  JNull{},
					"key": &JObject{
						values: map[string]Value{
							"key7": JString("hello"),
						},
					},

					"key5": &JArray{
						elements: []Value{
							//JNull{},
							JBool(true),
							JBool(false),
							JString("key"),
							JString("key"),
							&JObject{},
						},
					},
				},
			},
		},
		{
			values: map[string]interface{}{
				"key": struct {
					Key string `json:"key"`
				}{
					Key: "hello",
				},
			},
			value: &JObject{
				values: map[string]Value{
					"key": &JObject{
						values: map[string]Value{
							"key": JString("hello"),
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		obj := NewObject()
		for k, v := range tc.values {
			AddValue(obj, k, v)
		}
		if err := check(tc.value, obj); err != nil {
			t.Fatalf("case:%d %s\n", i, err)
		}
	}
}
