package inject

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Register

func TestRegisterStringValue(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*string)(nil), "abc")

	RegistryGet(t, registry, func(value string) {
		assert.Equal(t, "abc", value)
	})
}

func TestRegisterInterfaceValue(t *testing.T) {
	registry := NewTypeRegistry()
	registry.Register((*SomeInterface)(nil), 123)

	RegistryGet(t, registry, func(value SomeInterface) {
		assert.Equal(t, 123, value)
	})
}

func TestRegisterPointerValue(t *testing.T) {
	registry := NewTypeRegistry()
	foo := &SomeType{}
	registry.Register((*SomeType)(nil), foo)

	RegistryGet(t, registry, func(value *SomeType) {
		assert.Equal(t, foo, value)
	})
}

func TestRegisterPointerValuePersistence(t *testing.T) {
	registry := NewTypeRegistry()
	foo := &SomeType{"bar"}
	registry.Register((*SomeType)(nil), foo)

	RegistryGet(t, registry, func(value *SomeType) {
		assert.Equal(t, "bar", value.stuff)
		value.stuff = "thing"
	})
	RegistryGet(t, registry, func(value *SomeType) {
		assert.Equal(t, "thing", value.stuff)
	})
}

func TestRegisterRequireInterfacePointer(t *testing.T) {
	registry := NewTypeRegistry()
	assert.Panics(t, func() {
		registry.Register("asdf", "asdf")
	})
}

// NewChild

func TestNewChildInheritance(t *testing.T) {
	parent := NewTypeRegistry()
	child := parent.NewChild()
	grandchild := child.NewChild()
	parent.Register((*string)(nil), "parent")

	RegistryGet(t, child, func(value string) {
		assert.Equal(t, "parent", value)
	})
	RegistryGet(t, grandchild, func(value string) {
		assert.Equal(t, "parent", value)
	})
}

func TestNewChildInheritanceNoBackprop(t *testing.T) {
	parent := NewTypeRegistry()
	child := parent.NewChild()
	grandchild := child.NewChild()
	child.Register((*string)(nil), "child")

	AssertRegistryMissing(t, parent, func(value string) {})
	RegistryGet(t, grandchild, func(value string) {
		assert.Equal(t, "child", value)
	})
}

func TestNewChildOverrideParents(t *testing.T) {
	parent := NewTypeRegistry()
	child := parent.NewChild()
	grandchild := child.NewChild()
	parent.Register((*string)(nil), "parent")
	child.Register((*string)(nil), "child")

	RegistryGet(t, child, func(value string) {
		assert.Equal(t, "child", value)
	})
	RegistryGet(t, grandchild, func(value string) {
		assert.Equal(t, "child", value)
	})
}
