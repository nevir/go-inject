`inject`
========

`inject` is a very simple [dependency injection](http://en.wikipedia.org/wiki/Dependency_injection)
library for [Go](http://golang.org/).

You can manage registries of values (keyed by type), which can be used to inject
those values into functions that request them.


Example Use
-----------

In order to inject values into functions, you must first prepare them via
`PrepareFunc`. This does as much up front work as possible so that you can
repeatedly call prepared functions with minimal overhead.

```go
indexRoute = inject.PrepareFunc(response http.ResponseWriter) {
  fmt.Fprint(response, "Homepage!")
}

articleRoute = inject.PrepareFunc(response http.ResponseWriter, path Path) {
  fmt.Fprintf(response, "You're at the %s article. Clearly.", path)
}
```

With `inject`, you manage a `TypeRegistry`, which acts as a lookup table for all
the values you wish to inject into functions:

```go
registry := inject.NewTypeRegistry()
registry.Register((*http.ResponseWriter)(nil), response)
registry.Register((*http.Request)(nil), request)
registry.Register((*Path)(nil), request.URL.Path)
```

To perform the injection, simply `Call` the prepared function:

```go
articleRoute.Call(registry)
```


Registry Hierarchies
--------------------

`TypeRegistry` instances can also be arranged into a hierarchy, which enables
values to be registered in different "scopes".

For example, you might want to make global values available:

```go
globalRegistry := inject.NewTypeRegistry()
globalRegistry.Register((*Env)(nil), currentEnv)
globalRegistry.Register((*User)(nil), guestUser)
```

While also providing values that might change more frequently:

```go
requestRegistry := globalRegistry.NewChild()
requestRegistry.Register((*RequestId)(nil), 12345)
requestRegistry.Register((*User)(nil), currentUser)
```

Note that child registries can override values specified by their parents.


License
-------

`inject` is MIT licensed. [See the accompanying file](MIT-LICENSE.md) for the
full text.
