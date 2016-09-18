package recurse

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// old version
func recurseStruct(obj interface{}, depth int) {

	v := reflect.ValueOf(obj)
	// if v.Kind() == reflect.Ptr {
	// 	v = v.Elem()
	// }

	indent := strings.Repeat("\t", depth)
	fmt.Printf("%s   %-16v  %+v  -  \n", indent, v.Type().Name(), v)
	for i := 0; i < v.NumField(); i++ {

		fieldName := v.Type().Field(i).Name
		subType := v.Field(i).Type().Name()
		realVal := v.Field(i).Interface()

		log.Printf("%s\t%d: %-12s %-12s = %+v  \n", indent, i, fieldName, subType, realVal)

		anon := v.Type().Field(i).Anonymous
		kind := v.Type().Kind()

		if anon && kind == reflect.Struct {
			recurseStruct(realVal, depth+1)
		}
	}
}
