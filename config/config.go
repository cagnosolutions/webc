package config

import "reflect"

type Config map[string]interface{}

func New(v interface{}, conf Config) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		panic("Construct requires a pointer!")
	}
	for name, value := range conf {
		fld := val.Elem().FieldByName(name)
		fld.Set(reflect.ValueOf(value))
	}
	return v
}
