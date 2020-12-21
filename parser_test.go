package hjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseObject(t *testing.T) {
	testCases := []struct {
		json  string
		value Value
	}{
		{
			json: `{"key":null,
				"key1":false,
				"key2":true,
				"key3":2,
				"key4":"hello parser",
				"key5":[null,true, false, "key", "key",[1, 2],{"key":"hello"}],
				"key6":{"key7":"hello"}}`,
			value: &JObject{
				values: map[string]Value{
					"key":  JNull{},
					"key1": JBool(false),
					"key2": JBool(true),
					"key3": JNumber(2),
					"key4": JString("hello parser"),
					"key6": &JObject{
						values: map[string]Value{
							"key7": JString("hello"),
						},
					},
					"key5": &JArray{
						elements: []Value{
							JNull{},
							JBool(true),
							JBool(false),
							JString("key"),
							JString("key"),
							&JArray{
								elements: []Value{
									JNumber(1),
									JNumber(2),
								},
							},
							&JObject{
								values: map[string]Value{
									"key": JString("hello"),
								},
							},
						},
					},
				},
			},
			/*
				{},
			*/
		},
	}
	for _, tc := range testCases {
		p := newParser(bytes.NewBufferString(tc.json))
		array, err := p.parse()
		if err != nil {
			t.Fatal(err)
		}
		if err := check(tc.value, array); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseArray(t *testing.T) {
	testCases := []struct {
		json  string
		value Value
	}{
		{
			json: `[null,true, false, "key", "key",[1, 2],{"key":"hello"}]`,
			value: &JArray{
				elements: []Value{
					JNull{},
					JBool(true),
					JBool(false),
					JString("key"),
					JString("key"),
					&JArray{
						elements: []Value{
							JNumber(1),
							JNumber(2),
						},
					},
					&JObject{
						values: map[string]Value{
							"key": JString("hello"),
						},
					},
				},
			},
		},
		/*
			{},
		*/
	}
	for _, tc := range testCases {
		p := newParser(bytes.NewBufferString(tc.json))
		array, err := p.parse()
		if err != nil {
			t.Fatal(err)
		}
		if err := check(tc.value, array); err != nil {
			t.Fatal(err)
		}
	}
}

func check(expect Value, got Value) error {
	if expect.Type() != got.Type() {
		return fmt.Errorf("expect:%v got:%v", expect, got)
	}
	switch got.Type() {
	case typeString:
		if expect.String() != got.String() {
			return fmt.Errorf("expect:%v got:%v", expect, got)
		}
	case typeBool:
		if expect.String() != got.String() {
			return fmt.Errorf("expect:%v got:%v", expect, got)
		}
	case typeArray:
		expectedArray := expect.(*JArray)
		gotArray := got.(*JArray)
		if len(expectedArray.elements) != len(gotArray.elements) {
			return fmt.Errorf("expect:%v got:%v", expect, got)
		}
		for i, e := range gotArray.elements {
			if err := check(expectedArray.elements[i], e); err != nil {
				return err
			}
		}

	case typeObj:
		expectedObj := expect.(*JObject)
		gotObj := got.(*JObject)
		values := expectedObj.values
		gotValues := gotObj.values
		if len(values) != len(gotValues) {
			return fmt.Errorf("expect:%v got:%v", expect, got)
		}
		for k, v := range values {
			if e, ok := gotValues[k]; !ok {
				return fmt.Errorf("expect:%v got:%v", expect, got)
			} else {
				if err := check(v, e); err != nil {
					return err
				}
			}
		}

	case typeNumber:
		if expect.String() != got.String() {
			return fmt.Errorf("expect:%v got:%v", expect, got)
		}
	}
	return nil
}

func TestInvalidJson(t *testing.T) {
	jsons := []string{
		``,
		`1`,
		`{"key"}`,
		`{key:123}`,
		`{"key":123df}`,
		`{"key":"kk\h"}`,
		`{"key":"kk`,
		`{`,
	}
	d := make(map[string]interface{})
	for i, j := range jsons {
		parser := NewParser(j)
		_, err := parser.parse()
		if err == nil {
			t.Fatalf("case:%d expect err, got nil", i)
		}
		t.Logf("case:%d %v\n", i, err)
		err = json.Unmarshal([]byte(j), &d)
		if err != nil {
			t.Logf("case:%d standard:%v\n", i, err)
		}
	}
}

/*
func TestJson(t *testing.T) {
	j := `{"key":12e+-1dddd`
	n := struct {
		Key string `json:"key"`
	}{}
	if err := json.Unmarshal([]byte(j), &n); err != nil {
		t.Logf("%v\n", err)
	}
	parser := newParser(bytes.NewBufferString(j))
	_, err := parser.parse()
	if err != nil {
		t.Fatal(err)
	}
}
*/
