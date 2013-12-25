package inject

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type SomeInterface interface{
	GetStuff() string
}

type SomeType struct {
	stuff string
}

func (t SomeType) GetStuff() string {
	return t.stuff
}

// The easiest way to express our types is just to expose them via a function
// argument. Takes a function with a single argument, and passes the matching
// value in the registry to the function.
func RegistryGet(t *testing.T, registry *TypeRegistry, getter interface{}) {
	getterType := reflect.TypeOf(getter)
	argType := getterType.In(0)
	reflectedValue := registry.get(argType)

	assert.True(t, reflectedValue.IsValid(), "Unregistered type %T", argType)

	reflect.ValueOf(getter).Call([]reflect.Value{reflectedValue})
}

// The same, but just asserts that the registry _doesn't_ contain the target
// type. The getter will _not_ be called.
func AssertRegistryMissing(t *testing.T, registry *TypeRegistry, getter interface{}) {
	getterType := reflect.TypeOf(getter)
	argType := getterType.In(0)
	reflectedValue := registry.get(argType)

	assert.False(t, reflectedValue.IsValid(), "Expected type %T to not be registered", argType)
}
