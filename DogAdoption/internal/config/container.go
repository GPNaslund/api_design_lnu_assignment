package config

type provider struct {
	value any
}

type Lifecycle int

const (
	Singleton Lifecycle = iota
	Transient
)

// Simple IoC container
type Container struct {
	singletons map[string]any
	transients map[string]provider
}

// Constructor function.
func NewContainer() *Container {
	return &Container{
		singletons: make(map[string]any),
		transients: make(map[string]provider),
	}
}

// Method for providing singleton objects.
func (c *Container) ProvideSingleton(name string, constructorFunc func() any) {
	c.singletons[name] = constructorFunc()
}

// Method for providing transient object constructor.
func (c *Container) ProvideTransient(name string, constructorFunc func() any) {
	c.transients[name] = provider{value: constructorFunc}
}

// Method for resolving the specified object.
func (c *Container) Resolve(name string, lifecycle Lifecycle) any {
	switch lifecycle {
	case Singleton:
		instance, exists := c.singletons[name]
		if !exists {
			panic("no singleton found with name " + name)
		}
		return instance
	case Transient:
		provider, exists := c.transients[name]
		if !exists {
			panic("no transient provider found with name " + name)
		}
		return provider.value.(func() any)()
	default:
		panic("Unsupported lifecycle")
	}
}
