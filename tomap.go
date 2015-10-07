package dbdb

import "reflect"

func ToMap(v interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	value := reflect.ValueOf(v)
	fields := structFields(value)
	for _, field := range fields {
		val := value.FieldByName(field.Name)
		var finalVal interface{}
		zero := reflect.Zero(val.Type()).Interface()
		current := val.Interface()
		if reflect.DeepEqual(current, zero) {
			continue
		}
		if v, ok := isStruct(val.Interface()); ok {
			finalVal = ToMap(v)
		} else {
			finalVal = val.Interface()
		}
		output[field.Name] = finalVal
	}
	return output
}

func structFields(value reflect.Value) []reflect.StructField {
	t := value.Type()
	var f []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		f = append(f, field)
	}
	return f
}

func isStruct(s interface{}) (interface{}, bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Invalid {
		return nil, false
	}
	return v.Interface(), v.Kind() == reflect.Struct
}
