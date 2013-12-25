package inject

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuncCalled(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*SomeInterface)(nil), "abcd")

	var value SomeInterface
	PrepareFunc(func(arg SomeInterface) { value = arg })(registry)
	assert.Equal(t, "abcd", value)
}

func TestValueArg(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*string)(nil), "foo")

	PrepareFunc(func(arg string) {
		assert.Equal(t, "foo", arg)
	})(registry)
}

func TestLateBoundValue(t *testing.T) {
	prepared := PrepareFunc(func(arg string) {
		assert.Equal(t, "bar", arg)
	})

	registry := NewTypeRegistry()
	registry.Register((*string)(nil), "bar")
	prepared(registry)
}

func TestRedefinedValue(t *testing.T) {
	var value string
	registry := NewTypeRegistry()
	prepared := PrepareFunc(func(arg string) { value = arg })

	registry.Register((*string)(nil), "foo")
	prepared(registry)
	assert.Equal(t, "foo", value)

	registry.Register((*string)(nil), "bar")
	prepared(registry)
	assert.Equal(t, "bar", value)
}

func TestMultipleArgs(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*int)(nil), 1234)
	registry.Register((*string)(nil), "asdf")
	registry.Register((*SomeType)(nil), &SomeType{"foo"})

	PrepareFunc(func(str string, num int, some *SomeType) {
		assert.Equal(t, "asdf", str)
		assert.Equal(t, 1234, num)
		assert.Equal(t, &SomeType{"foo"}, some)
	})(registry)
}

func TestSingleReturnValue(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*int)(nil), 9000)

	result := PrepareFunc(func(val int) int { return val + 1 })(registry)
	assert.Equal(t, []interface{}{9001}, result)
}

func TestMultipleReturnValues(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*int)(nil), 9000)

	result := PrepareFunc(func(val int) (int, string) { return 1234, "hi" })(registry)
	assert.Equal(t, []interface{}{1234, "hi"}, result)
}

// Error Cases

func TestPrepareFuncRequireFunction(t *testing.T) {
	assert.Panics(t, func() {
		PrepareFunc(1234)
	})
}

func TestPrepareFuncNoVariadicSupport(t *testing.T) {
	assert.Panics(t, func() {
		PrepareFunc(func(arg ...string) {})
	})
}

func TestPrepareFuncNoArgs(t *testing.T) {
	result := PrepareFunc(func() int { return 123 })(NewTypeRegistry())
	assert.Equal(t, []interface{}{123}, result)
}

func TestUnboundValue(t *testing.T) {
	assert.Panics(t, func() {
		PrepareFunc(func(val string) {})(NewTypeRegistry())
	})
}
