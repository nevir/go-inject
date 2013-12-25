// Package inject provides simple type/dependency injection.
//
// To inject values, you must first register them via a TypeRegistry:
//
//   registry := inject.NewTypeRegistry()
//   registry.Register(someValue, (*SomeType)(nil))
//
// Functions that you wish to inject into must also be pre-processed so that
// they can accept a TypeRegistry (and its values). This is done via
// PrepareFunc():
//
//   prepared := inject.PrepareFunc(func(value SomeType) {
//     // ...
//   })
//
// The result of PrepareFunc() is a function that accepts a TypeRegistry, and
// returns any value returned by the wrapped func:
//
//   results := prepared(registry)
//
package inject
