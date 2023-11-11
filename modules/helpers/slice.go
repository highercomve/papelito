package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

// AllIn are all elements of and array inside another array
func AllIn(in, all []string) bool {
	allOfThen := false

	for _, x := range all {
		found := false
		for _, y := range in {
			if x == y {
				found = true
				break
			}
		}

		if !found {
			return found
		}

		allOfThen = found
	}

	return allOfThen
}

type updateFunc func(string, string)

// FlatMap flat a map structure
func FlatMap(i interface{}, prefix string) map[string]string {
	result := map[string]string{}

	update := func(k, value string) {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s/%s", prefix, key)
		}
		result[key] = value
	}

	flatten(i, "", update)

	return result
}

func flatten(i interface{}, prefix string, update updateFunc) {
	switch v := i.(type) {
	case map[string]interface{}:
		flattenMap(v, prefix, update)
	case map[string]string:
		flattenMapString(v, prefix, update)
	case []interface{}:
		flattenSlice(v, prefix, update)
	case string:
		update(prefix, v)
	default:
	}
}

func flattenMap(m map[string]interface{}, prefix string, update updateFunc) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, key)
		}
		flatten(v, key, update)
	}
}

func flattenMapString(m map[string]string, prefix string, update updateFunc) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, key)
		}
		flatten(v, key, update)
	}
}

func flattenSlice(vs []interface{}, prefix string, update updateFunc) {
	update(prefix+".len", fmt.Sprintf("%d", len(vs)))
	for i, v := range vs {
		key := fmt.Sprintf("%d", i)
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, key)
		}

		flatten(v, key, update)
	}
}

// IsNotEmptyString check if an string is not empty
func IsNotEmptyString(value string) bool {
	return value != ""
}

// FilterStrings filter and slice of strings
func FilterStrings(arr []string, cond func(string) bool) []string {
	result := []string{}

	for i := range arr {
		if cond(arr[i]) {
			result = append(result, strings.Trim(arr[i], " "))
		}
	}

	return result
}

// MergeMaps Given two maps, recursively merge right into left, NEVER replacing any key that already exists in left
func MergeMaps(left, right map[string]interface{}) map[string]interface{} {
	if left == nil {
		left = map[string]interface{}{}
	}
	if right == nil {
		right = map[string]interface{}{}
	}
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			// then we don't want to replace it - recurse
			if reflect.ValueOf(rightVal).Kind() == reflect.Map {
				left[key] = MergeMaps(leftVal.(map[string]interface{}), rightVal.(map[string]interface{}))
			} else {
				left[key] = rightVal
			}

		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}
