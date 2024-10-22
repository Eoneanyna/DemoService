package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
)

func StructByReflect(data map[string]interface{}, inStructPtr interface{}) error {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}
	// 遍历结构体
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		var column string
		// 得到tag中的字段名
		key := t.Tag.Get("gorm")
		reg := regexp.MustCompile(`column:(.*)`)
		result := reg.FindStringSubmatch(string(key))
		if result != nil {
			column = result[1]
		} else {
			column = t.Name
		}
		if v, ok := data[column]; ok && v != nil {
			// 检查是否需要类型转换
			dataType := reflect.TypeOf(v)
			structType := f.Type()
			if dataType.Kind() == structType.Kind() {
				f.Set(reflect.ValueOf(v))
			} else {
				if structType.String() == "int" {
					// 转换类型
					value, _ := strconv.Atoi(v.(string))
					f.Set(reflect.ValueOf((value)))
				} else if structType.String() == "float32" {
					value, _ := strconv.ParseInt(v.(string), 10, 32)
					f.Set(reflect.ValueOf((value)))
				} else {
					return errors.New("StructByReflect " + t.Name + "err")
				}
			}
		}
	}
	return nil
}
func Md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
