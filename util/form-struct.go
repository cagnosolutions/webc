package util

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func FormToStruct(stuff interface{}, v url.Values) {
	s := reflect.ValueOf(stuff).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Kind() {
		case reflect.String:
			s.Field(i).SetString(v.Get(typeOfT.Field(i).Name))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			in, _ := strconv.ParseInt(v.Get(typeOfT.Field(i).Name), 10, 64)
			s.Field(i).SetInt(in)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, _ := strconv.ParseUint(v.Get(typeOfT.Field(i).Name), 10, 64)
			s.Field(i).SetUint(u)
		case reflect.Float32, reflect.Float64:
			f, _ := strconv.ParseFloat(v.Get(typeOfT.Field(i).Name), 64)
			s.Field(i).SetFloat(f)
		case reflect.Bool:
			b, _ := strconv.ParseBool(v.Get(typeOfT.Field(i).Name))
			s.Field(i).SetBool(b)
		case reflect.Map:
			s.Field(i).Set(reflect.MakeMap(s.Field(i).Type()))
		case reflect.Slice:
			ss := reflect.MakeSlice(s.Field(i).Type(), 0, 0)
			s.Field(i).Set(genSlice(ss, v.Get(typeOfT.Field(i).Name)))
		case reflect.Struct:
			s.Field(i).Set(reflect.Indirect(reflect.New(s.Field(i).Type())))
		}
	}
}

func genSlice(sl reflect.Value, val string) reflect.Value {
	vs := strings.Split(val, ",")
	for _, v := range vs {
		switch sl.Type().String() {
		case "[]string":
			sl = reflect.Append(sl, reflect.ValueOf(v))
		case "[]int":
			in, _ := strconv.ParseInt(v, 10, 0)
			sl = reflect.Append(sl, reflect.ValueOf(int(in)))
		case "[]int8":
			in, _ := strconv.ParseInt(v, 10, 8)
			sl = reflect.Append(sl, reflect.ValueOf(int8(in)))
		case "[]int16":
			in, _ := strconv.ParseInt(v, 10, 16)
			sl = reflect.Append(sl, reflect.ValueOf(int16(in)))
		case "[]int32":
			in, _ := strconv.ParseInt(v, 10, 32)
			sl = reflect.Append(sl, reflect.ValueOf(int32(in)))
		case "[]int64":
			in, _ := strconv.ParseInt(v, 10, 64)
			sl = reflect.Append(sl, reflect.ValueOf(int64(in)))
		case "[]uint":
			in, _ := strconv.ParseUint(v, 10, 0)
			sl = reflect.Append(sl, reflect.ValueOf(uint(in)))
		case "[]uint8":
			in, _ := strconv.ParseUint(v, 10, 8)
			sl = reflect.Append(sl, reflect.ValueOf(uint8(in)))
		case "[]uint16":
			in, _ := strconv.ParseUint(v, 10, 16)
			sl = reflect.Append(sl, reflect.ValueOf(uint16(in)))
		case "[]uint32":
			in, _ := strconv.ParseUint(v, 10, 32)
			sl = reflect.Append(sl, reflect.ValueOf(uint32(in)))
		case "[]uint64":
			in, _ := strconv.ParseUint(v, 10, 64)
			sl = reflect.Append(sl, reflect.ValueOf(uint64(in)))
		case "[]float32":
			in, _ := strconv.ParseFloat(v, 32)
			sl = reflect.Append(sl, reflect.ValueOf(float32(in)))
		case "[]float64":
			in, _ := strconv.ParseFloat(v, 64)
			sl = reflect.Append(sl, reflect.ValueOf(float64(in)))
		case "[]bool":
			b, _ := strconv.ParseBool(v)
			sl = reflect.Append(sl, reflect.ValueOf(b))
		}
	}
	return sl
}
