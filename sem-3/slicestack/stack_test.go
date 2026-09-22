package slicestack

import (
	"errors"
	"testing"
)

func TestPushPeekAndPop(t *testing.T) {
	var stack []int

	stack = Push(stack, 10)
	stack = Push(stack, 20)

	value, err := Peek(stack)
	if err != nil {
		t.Fatalf("Peek returned error: %v", err)
	}

	if value != 20 {
		t.Fatalf("Peek returned %d, want 20", value)
	}

	if len(stack) != 2 {
		t.Fatalf("Peek changed stack length to %d, want 2", len(stack))
	}

	value, stack, err = Pop(stack)
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if value != 20 {
		t.Fatalf("Pop returned %d, want 20", value)
	}

	if len(stack) != 1 || stack[0] != 10 {
		t.Fatalf("remaining stack = %v, want [10]", stack)
	}
}

func TestEmptyStack(t *testing.T) {
	if _, err := Peek(nil); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Peek(nil) error = %v, want ErrEmpty", err)
	}

	value, rest, err := Pop(nil)
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("Pop(nil) error = %v, want ErrEmpty", err)
	}

	if value != 0 || rest != nil {
		t.Fatalf("Pop(nil) = (%d, %v), want (0, nil)", value, rest)
	}

	empty := []int{}

	value, rest, err = Pop(empty)
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("Pop([]) error = %v, want ErrEmpty", err)
	}

	if value != 0 || rest == nil || len(rest) != 0 {
		t.Fatalf("Pop([]) = (%d, %#v), want (0, non-nil empty slice)", value, rest)
	}
}

func TestZeroValueIsAValidValue(t *testing.T) {
	stack := Push(nil, 0)

	value, rest, err := Pop(stack)
	if err != nil || value != 0 || len(rest) != 0 {
		t.Fatalf("Pop(Push(nil, 0)) = (%d, %v, %v), want (0, [], nil)", value, rest, err)
	}
}
