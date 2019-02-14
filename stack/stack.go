// stack.go holds a simple stack which can hold strings.
//

package stack

import (
	"errors"
	"sync"
)

// Stack holds the stack-data, protected by a mutex
type Stack struct {
	lock sync.Mutex
	s    []string
}

// New returns a new stack (for holding strings)
func New() *Stack {
	return &Stack{sync.Mutex{}, make([]string, 0)}
}

// Push adds a new item to our stack.
func (s *Stack) Push(v string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

// Pop returns an item from our stack.
func (s *Stack) Pop() (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return "", errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

// Empty returns `true` if our stack is empty.
func (s *Stack) Empty() bool {

	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	return (l == 0)
}
