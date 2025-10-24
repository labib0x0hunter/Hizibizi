package main

import "errors"

// Add contact
// Delete contact
// Restore contact from trash
type Action struct {
	Undo func() error
}

type Stack struct {
	history []Action
}

func NewStack() *Stack {
	return &Stack{
		history: make([]Action, 0),
	}
}

func (s *Stack) Push(a Action) {
	s.history = append(s.history, a)
}

func (s *Stack) Undo() error {
	if len(s.history) == 0 {
		return errors.New("no last operations")
	}
	last := s.history[len(s.history) - 1]
	err := last.Undo()
	if err == nil {
		s.history = s.history[: len(s.history) - 1]
	}
	return err
}