package dbdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
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
		output[strings.ToLower(field.Name[0:1])+field.Name[1:]] = finalVal
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
		setField(structVal.Elem(), strings.Title(mFld), mVal)
	}
	return nil
}

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

func isFloat(k reflect.Kind) bool {
	if k == reflect.Float32 || k == reflect.Float64 {
		return true
	}
	return false
}
