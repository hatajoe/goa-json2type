package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

var funcFmt = `
func() {
	%s
	%s
}
`
var typeFmt = `var %s = Type("%s", func() {
	%s
})
`
var memberFmt = `
	Member("%s", %s, "", func() {})`
var attrFmt = `
	Attribute("%s", %s, "", func() {})`
var requiredFmt = `
	Required("%s")`
var arrayOfFmt = `ArrayOf(%s)`

func typeName(k string, e interface{}, types *[]string) string {
	switch e.(type) {
	case bool:
		return "Boolean"
	case json.Number:
		return "Integer"
	case string:
		return "String"
	case time.Time:
		return "DateTime"
	default:
		if reflect.ValueOf(e).Kind() == reflect.Slice {
			return fmt.Sprintf(arrayOfFmt, typeName(k, e.([]interface{})[0], types))
		} else if reflect.ValueOf(e).Kind() == reflect.Map {
			attrs := make([]string, len(e.(map[string]interface{})))
			i := 0
			for key, value := range e.(map[string]interface{}) {
				attrs[i] = fmt.Sprintf(memberFmt, key, typeName(key, value, types))
				i++
			}
			*types = append(*types, fmt.Sprintf(typeFmt, k, k, strings.Join(attrs, "")))
			return k
		}
	}
	panic(fmt.Sprintf("err: Unsupported type `%v' specified %v", reflect.ValueOf(e).Kind(), e))
}

func main() {
	t := flag.String("t", "media", "media or payload")
	flag.Parse()

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var v = map[string]interface{}{}
	decoder := json.NewDecoder(bytes.NewBuffer(buf))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		panic(err)
	}
	types := []string{}
	members := make([]string, len(v))
	requires := make([]string, len(v))

	i := 0
	paramFmt := memberFmt
	if *t == "media" {
		paramFmt = attrFmt
	}
	for k, e := range v {
		members[i] = fmt.Sprintf(paramFmt, k, typeName(k, e, &types))
		requires[i] = fmt.Sprintf(requiredFmt, k)
		i++
	}
	if len(types) > 0 {
		fmt.Print(strings.Join(types, ""))
	}
	fmt.Print(fmt.Sprintf(funcFmt, strings.Join(members, ""), strings.Join(requires, "")))
}
