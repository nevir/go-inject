package inject

import (
	"fmt"
	"reflect"
)

// A TypeRegistry manages a mapping between types and an instance of each
// registered type.
type TypeRegistry struct {
	mapping map[reflect.Type]reflect.Value
	parent  *TypeRegistry
}

// Returns a new TypeRegistry with no registered types.
func NewTypeRegistry() *TypeRegistry {
	registry := &TypeRegistry{
		mapping: make(map[reflect.Type]reflect.Value),
	}
	registry.Register(registry, (*TypeRegistry)(nil))

	return registry
}

// Registers a type and its instance.
//
// typePtr must be a pointer to an object (typically nil) that is cast to the
// desired interface, while value can be any object that satisfies that
// interface. Typical call pattern:
//
//   registry.Register(value, (*SomeType)(nil))
//
// For simplicity, and because we are forced to handle type pointers, Register
// infers the type that value will be exposed as based on the types of typePtr
// value:
//
// * If typePtr points to an interface, value will be directly exposed as that
//   interface.
//
// * If value is a pointer, it will be exposed as typePtr's given type.
//
// * If value is a value, it will be exposed as a value of the type that typePtr
//   points to.
//
// If a value is already registered for the type, it will be overridden.
func (r *TypeRegistry) Register(value, typePtr interface{}) {
	reflectedValue := reflect.ValueOf(value)
	reflectedType := reflect.TypeOf(typePtr)
	if reflectedType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Register() requires a type pointer. Did you mean to call Register(..., (*%T)(nil))?", typePtr))
	}

	if reflectedType.Elem().Kind() == reflect.Interface || reflectedValue.Kind() != reflect.Ptr {
		reflectedType = reflectedType.Elem()
	}

	r.mapping[reflectedType] = reflectedValue.Convert(reflectedType)
}

// Creates a child registry.
//
// Type registries can be organized into a hierarchy to enable modeling of
// various scopes. For example, a web server might have a global registry that
// manages all of its singletons. It would create child registries per request
// to manage request-specific values.
func (r *TypeRegistry) NewChild() *TypeRegistry {
	registry := NewTypeRegistry()
	registry.parent = r

	return registry
}

// Private Implementation

// Returns the value associated with the given type.
//
// If this registry does not contain a value of that type, its ancestors will
// be checked.
func (r *TypeRegistry) get(reflectedType reflect.Type) reflect.Value {
	value := r.mapping[reflectedType]
	if !value.IsValid() && r.parent != nil {
		return r.parent.get(reflectedType)
	}

	return value
}
