package hjson

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type JsonType int

const (
	invalid JsonType = iota
	typeObj
	typeArray
	typeString
	typeBool
	typeNumber
	typeNull
)

type Value interface {
	Type() JsonType
	accept(Walker)
	String() string
}

type JObject struct {
	values map[string]Value
}

type JArray struct {
	elements []Value
}

type JNumber int64

type JString string

type JBool bool

type JNull struct {
}

func (_ JObject) Type() JsonType {
	return typeObj
}
func (_ JArray) Type() JsonType {
	return typeArray
}
func (_ JString) Type() JsonType {
	return typeString
}
func (_ JBool) Type() JsonType {
	return typeBool
}
func (_ JNumber) Type() JsonType {
	return typeNumber
}

func (_ JNull) Type() JsonType {
	return typeNull
}

func NewObject() *JObject {
	return &JObject{
		values: make(map[string]Value),
	}
}

func NewArray() *JArray {
	return &JArray{
		elements: make([]Value, 0),
	}
}

func (a *JArray) addValue(v Value) {
	a.elements = append(a.elements, v)
}

func (o JObject) String() string {
	visitor := newNodeVisitor()
	o.accept(visitor)
	return visitor.buf.String()
}
func (a JArray) String() string {
	visitor := newNodeVisitor()
	a.accept(visitor)
	return visitor.buf.String()
}
func (s JString) String() string {
	return string(s)
}
func (n JNumber) String() string {
	return strconv.Itoa(int(n))
}
func (b JBool) String() string {
	if b {
		return "true"
	}
	return "false"
}
func (JNull) String() string {
	return "null"
}

func toJNumber(v int64) JNumber {
	return JNumber(v)
}

func converBaseType(value interface{}) (Value, bool) {
	if value == nil {
		return JNull{}, true
	}
	v := reflect.ValueOf(value)
	switch reflect.TypeOf(value).Kind() {
	//bool 必须改成reflect.TypeOf().Kind()
	case reflect.Bool:
		return JBool(v.Bool()), true
	case reflect.String:
		return JString(v.String()), true
	case reflect.Int8:
		return JNumber(v.Int()), true
	case reflect.Int16:
		return JNumber(v.Int()), true
	case reflect.Int32:
		return JNumber(v.Int()), true
	case reflect.Int64:
		return JNumber(v.Int()), true
	case reflect.Int:
		return JNumber(v.Int()), true
	case reflect.Uint8:
		return JNumber(v.Uint()), true
	case reflect.Uint16:
		return JNumber(v.Uint()), true
	case reflect.Uint32:
		return JNumber(v.Uint()), true
	case reflect.Uint64:
		return JNumber(v.Uint()), true
	case reflect.Uint:
		return JNumber(v.Uint()), true
	}
	return nil, false
}

//TODO: structToObj能处理匿名对象么
func structToObj(st interface{}) *JObject {
	obj := NewObject()
	t := reflect.TypeOf(st)
	value := reflect.ValueOf(st)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		fieldValue := value.Field(i)
		if name[0] >= 'A' && name[0] <= 'Z' {
			key := field.Tag.Get("json")
			if key != "" {
				name = key
			}
			//fieldValue.Interface的使用
			AddValue(obj, name, fieldValue.Interface())
		}
	}
	return obj
}

func sliceToArray(slice interface{}) *JArray {
	array := NewArray()
	value := reflect.ValueOf(slice)
	num := value.Len()
	for i := 0; i < num; i++ {
		tmp := value.Index(i)
		array.elements = append(array.elements, toJSONValue(tmp.Interface()))
	}
	return array
}

func mapToObj(m interface{}) *JObject {
	obj := NewObject()
	value := reflect.ValueOf(m)
	iter := value.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		obj.values[key] = toJSONValue(iter.Value().Interface())
	}
	return obj
}

func toJSONValue(value interface{}) Value {
	jv, ok := converBaseType(value)
	if ok {
		return jv
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Array:
		jv = sliceToArray(value)
	case reflect.Map:
		jv = mapToObj(value)
	case reflect.Slice:
		jv = sliceToArray(value)
	case reflect.Struct:
		jv = structToObj(value)
	case reflect.Ptr:
		v := reflect.ValueOf(value)
		elem := v.Elem()
		if v.IsNil() {
			jv = JNull{}
		} else {
			fmt.Printf("ptr %#v\n", value)
			jv = toJSONValue(elem.Interface())
		}
	default:
		panic(fmt.Errorf("invalid type:%T", value))
	}
	return jv
}

/*
TODO: 提供一些接口,那么应该提供哪些接口。
1 获取指定的元素 完成
2 添加元素 完成
3 解析字符串 完成
4 转换成字符串 完成
5 转换成TOML
*/
//===============================API========================

//AddArrayElement 添加元素到数组
func AddArrayElement(array *JArray, value interface{}) {
	jv := toJSONValue(value)
	array.elements = append(array.elements, jv)
}

//AddValue 添加新对象
func AddValue(obj *JObject, key string, value interface{}) {
	obj.values[key] = toJSONValue(value)
}

//GetObjField 从对象获取指定的键值对
func GetObjField(obj *JObject, key string) (Value, bool) {
	v, ok := obj.values[key]
	return v, ok
}

//Index 获取指定位置的值
func Index(value *JArray, index int) Value {
	if index < 0 || index > len(value.elements) {
		panic("")
	}
	return value.elements[index]
}

func ToValue(data []byte) (Value, error) {
	parser := newParser(bytes.NewBuffer(data))
	value, err := parser.parse()
	if err != nil {
		return nil, err
	}
	return value, nil
}
