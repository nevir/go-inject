package inject

import (
	"fmt"
	"reflect"
)

// A function that is prepared for type injection via type registries.
//
// If the wrapped function returns any values, they will be returned in array
// form.
type PreparedFunc func(registry *TypeRegistry) []interface{}

// Prepares a function for injection.
//
// Once prepared, the function can be called any number of times with any type
// registries.
func PrepareFunc(function interface{}) PreparedFunc {
	signature := reflect.TypeOf(function)
	if signature.Kind() != reflect.Func {
		panic(fmt.Sprintf("PrepareFunc() requires a function. Got: %v <%T>", function, function))
	}
	if signature.IsVariadic() {
		panic(fmt.Sprintf("PrepareFunc() doesn't support variadic funcs. Got: %T", function))
	}

	numIn := signature.NumIn()
	if numIn == 0 {
		panic("PrepareFunc() requires a function with one or more argument.")
	}

	argTypes := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		argTypes[i] = signature.In(i)
	}

	funcValue := reflect.ValueOf(function)

	return func(registry *TypeRegistry) []interface{} {
		args := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			value := registry.get(argTypes[i])
			if !value.IsValid() {
				panic(fmt.Sprintf("A value for %v is not registered in the TypeRegistry.", argTypes[i]))
			}
			args[i] = value
		}

		reflectedResults := funcValue.Call(args)
		results := make([]interface{}, len(reflectedResults))
		for i := range reflectedResults {
			results[i] = reflectedResults[i].Interface()
		}

		return results
	}
}
