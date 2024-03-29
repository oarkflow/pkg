package vue

import (
	"reflect"
)

type effect func()

var currentActiveEffect effect

type dep struct {
	subscribers map[uintptr]effect
}

func newDep() *dep {
	return &dep{subscribers: make(map[uintptr]effect)}
}

func (d *dep) track(update effect) {
	key := reflect.ValueOf(update).Pointer()
	d.subscribers[key] = currentActiveEffect
}

func (d *dep) trigger() {
	for _, effect := range d.subscribers {
		if effect != nil {
			effect()
		}
	}
}

/*
RefImpl (ref) is a reactive primitive that can be read and written onto
*/
type RefImpl[T any] struct {
	dep   *dep
	value T
}

func Ref[T any](initialValue T) *RefImpl[T] {
	return &RefImpl[T]{
		dep:   newDep(),
		value: initialValue,
	}
}

func (r *RefImpl[T]) GetValue() T {
	r.dep.track(currentActiveEffect)
	return r.value
}

func (r *RefImpl[T]) SetValue(newValue T) {
	r.value = newValue
	r.dep.trigger()
}

/*
computed is a ref that is computed by a getter
*/
type computed[T any] struct {
	dep     *dep
	compute func() T
}

func Computed[T any](computedValue func() T) *computed[T] {

	return &computed[T]{
		dep:     newDep(),
		compute: computedValue,
	}
}

func (c *computed[T]) GetValue() T {
	c.dep.track(currentActiveEffect)
	return c.compute()
}

// WatchEffect any
func WatchEffect(update effect) {
	var wrappedUpdate func()
	wrappedUpdate = func() {
		currentActiveEffect = wrappedUpdate
		update()
	}
	wrappedUpdate()
}
