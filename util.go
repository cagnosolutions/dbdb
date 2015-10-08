package dbdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func WriteStore(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("WriteStore() -> os.MkdirAll() -> %v\n", err)
	}
}

func WriteDoc(filepath string, doc interface{}) {
	data, err := json.Marshal(doc)
	if err != nil {
		log.Fatalf("WriteDoc() -> %v\n", err)
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		log.Fatalf("WriteDoc() -> ioutil.WriteFile() -> %v\n", err)
	}
}

func DeleteStore(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("DeleteStore() -> os.Remove() -> %v\n", err)
	}
}

func DeleteDoc(filepath string) {
	if err := os.Remove(filepath); err != nil {
		log.Fatalf("DeleteDoc() -> os.Remove() -> %v\n", err)
	}
}

func ToMap(v interface{}) map[string]interface{} {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
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
		output[field.Name] = finalVal
	}
	return output
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

func ToStruct(m map[string]interface{}, ptr interface{}) error {
	structVal := reflect.ValueOf(ptr)
	if structVal.Kind() != reflect.Ptr || structVal.IsNil() {
		return fmt.Errorf("Expected pointer didn't get one...")
	}
	for mFld, mVal := range m {
		setField(structVal.Elem(), mFld, mVal)
	}
	return nil
}

func setField(structVal reflect.Value, mFld string, mVal interface{}) {
	fld := structVal.FieldByName(mFld)
	if fld.IsValid() && fld.CanSet() {
		val := reflect.ValueOf(mVal)
		if fld.Type() == val.Type() {
			fld.Set(val)
		}
	}
}
