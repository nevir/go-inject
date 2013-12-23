package inject

import (
	"fmt"
	"reflect"
)

// An Injectable wraps a function so that it can be called with type registries.
type Injectable func(registry *TypeRegistry) []reflect.Value

// Prepares a function for injection.
//
// Once prepared, the function can be Call()'d any number of times.
func PrepareFunc(function interface{}) Injectable {
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

	return func(registry *TypeRegistry) []reflect.Value {
		args := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			value := registry.get(argTypes[i])
			if !value.IsValid() {
				panic(fmt.Sprintf("A value for %v is not registered in the TypeRegistry.", argTypes[i]))
			}
			args[i] = value
		}

		return funcValue.Call(args)
	}
}

// Invokes the underlying function, injecting values present in registry that
// match the function's argument types.
//
// Any return values are returned.
func (i Injectable) Call(registry *TypeRegistry) []interface{} {
	reflectedResults := i(registry)

	results := make([]interface{}, len(reflectedResults))
	for i := range reflectedResults {
		results[i] = reflectedResults[i].Interface()
	}

	return results
}
