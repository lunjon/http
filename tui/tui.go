package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Option[T any] struct {
	some  bool
	value T
}

func (o Option[T]) IsSome() bool {
	return o.some
}

func (o Option[T]) IsNone() bool {
	return !o.IsSome()
}

func (o Option[T]) Value() T {
	if !o.some {
		panic("No value")
	}
	return o.value
}

func (o Option[T]) Set(value T) Option[T] {
	return Option[T]{
		some:  true,
		value: value,
	}
}

func Start() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}
