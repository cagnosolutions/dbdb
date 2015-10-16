package dbdb

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// used to transform a struct into a map
func ToMap(v interface{}) map[string]interface{} {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Type().String() == "map[string]interface {}" {
		return v.(map[string]interface{})
	}
	if value.Kind() != reflect.Struct {
		log.Fatalf("ToMap() -> value must be %q or %q", "struct", "&struct")
	}
	output := make(map[string]interface{})
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
		output[strings.ToLower(field.Name[0:1])+field.Name[1:]] = finalVal
	}
	return output
}

// helper used to check if current interface is a struct
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

// helper used to gather and return all exported struct fields
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

// used to transform a map into a struct
func ToStruct(m map[string]interface{}, ptr interface{}) error {
	structVal := reflect.ValueOf(ptr)
	if structVal.Kind() != reflect.Ptr || structVal.IsNil() {
		return fmt.Errorf("Expected pointer didn't get one...")
	}
	for mFld, mVal := range m {
		setField(structVal.Elem(), strings.Title(mFld), mVal)
	}
	return nil
}

// helper used to match map fields up with the supplied struct's fields
func setField(structVal reflect.Value, mFld string, mVal interface{}) {
	fld := structVal.FieldByName(mFld)
	if fld.IsValid() && fld.CanSet() {
		val := reflect.ValueOf(mVal)
		if fld.Type() == val.Type() {
			fld.Set(val)
		} else if fld.Type() != val.Type() && isFloat(val.Kind()) {
			setNumber(fld, val)
		}
	}
}

// helper used to marshal special cases of floats back into struct specific types
func setNumber(sv, mv reflect.Value) {
	switch sv.Kind() {
	case reflect.Float32:
		sv.SetFloat(mv.Interface().(float64))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sv.SetInt(int64(mv.Interface().(float64)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sv.SetUint(uint64(mv.Interface().(float64)))
	}
}

// helper used to check if a particular reflect kind is a float
func isFloat(k reflect.Kind) bool {
	if k == reflect.Float32 || k == reflect.Float64 {
		return true
	}
	return false
}
