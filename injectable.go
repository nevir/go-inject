package inject

import (
	"fmt"
	"reflect"
)

// An Injectable wraps a function so that it can be called with type registries.
type Injectable func(registry *TypeRegistry)

// Prepares a function for injection.
func PrepareFunc(function interface{}) Injectable {
	signature := reflect.TypeOf(function)
	if signature.Kind() != reflect.Func {
		panic(fmt.Sprintf("PrepareFunc() requires a function. Got: %v <%T>", function, function))
	}
	if signature.IsVariadic() {
		panic(fmt.Sprintf("PrepareFunc() doesn't support variadic funcs. Got: %T", function))
	}

	numIn := signature.NumIn()
	argTypes := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		argTypes[i] = signature.In(i)
	}

	funcValue := reflect.ValueOf(function)

	return func(registry *TypeRegistry) {
		args := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			value := registry.get(argTypes[i])
			if !value.IsValid() {
				panic(fmt.Sprintf("A value for %v is not registered in the TypeRegistry.", argTypes[i]))
			}
			args[i] = value
		}

		funcValue.Call(args)
	}
}

// Invokes the underlying function, passing the values in registry that match
// the type signature of the function.
//
// Any return values are returned as an array of interfaces.
func (i Injectable) Call(registry *TypeRegistry) ([]interface{}, error) {
	i(registry)

	return nil, nil
}
