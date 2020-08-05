package jsonhelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"utils/common"
)

func GetInt(tree map[string]interface{}, name string) (int, error) {
	item := tree[name]
	if item == nil {
		return 0, fmt.Errorf("%s 不存在", name)
	}
	str := fmt.Sprintf("%.0f", item)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func GetUInt8(tree map[string]interface{}, name string) (uint8, error) {
	r, err := GetUInt64(tree, name)
	return uint8(r), err
}

func GetUInt16(tree map[string]interface{}, name string) (uint16, error) {
	r, err := GetUInt64(tree, name)
	return uint16(r), err
}
func GetUInt32(tree map[string]interface{}, name string) (uint32, error) {
	r, err := GetUInt64(tree, name)
	return uint32(r), err
}

func GetUInt64(tree map[string]interface{}, name string) (uint64, error) {
	item := tree[name]
	if item == nil {
		return 0, fmt.Errorf("%s 不存在", name)
	}
	//str := fmt.Sprintf("%.0f", item)
	//result, err := strconv.ParseUint(str, 0, 64)
	//if err != nil {
	//		return 0, err
	//	}
	return commontools.InterfaceToUInt64(item)
}

func GetString(tree map[string]interface{}, name string) (string, error) {
	item := tree[name]
	if item == nil {
		return "", fmt.Errorf("%s 不存在", name)
	}
	result := fmt.Sprintf("%s", item)
	return result, nil
}
func ToStringInterface(item interface{}) (map[string]string, error) {
	var result map[string]string
	result, ok := item.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("数据转换出现异常 to map[string]string")
	}
	return result, nil
}
func ToMapInterface(item interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	result, ok := item.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("数据转换出现异常 to map[string]interface{}")
	}
	return result, nil
}
func ToInterfaces(item interface{}) ([]interface{}, error) {
	var result []interface{}
	result, ok := item.([]interface{})
	if !ok {
		return nil, fmt.Errorf("数据转换出现异常 to []interface{}")
	}
	return result, nil
}
func IteratorJson() {
	var v map[string]interface{}
	var list interface{}
	list = v["list"]

	b, _ := json.Marshal(list)

	fmt.Println(string(b))
	for k, v := range v {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
