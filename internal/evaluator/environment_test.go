package evaluator

import "testing"

func value(v any) Value {
	return &ValueLiteral{Literal: v}
}

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	_, err := env.Get("a")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if env.Set("a", value(1)) == nil {
		t.Errorf("Expected error, got nil")
	}

	env.Declare("a", value(1))
	assertEnv(t, env, "a", 1)

	if env.Set("a", value(2)) != nil {
		t.Errorf("Expected no error for assigning to a")
	}

	env.Declare("b", value(2))
	assertEnv(t, env, "b", 2)

	env.Declare("b", value(3))
	assertEnv(t, env, "b", 3)

	env.CreateScope()

	assertEnv(t, env, "b", 3)

	env.Declare("c", value(4))
	assertEnv(t, env, "c", 4)

	env.Declare("b", value(5))
	assertEnv(t, env, "b", 5)

	if env.Set("a", value(5)) != nil {
		t.Errorf("Expected no error for assigning to a")
	}

	env.ExitScope()

	_, err = env.Get("c")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	assertEnv(t, env, "b", 3)

	assertEnv(t, env, "a", 5)
}

func assertEnv(t *testing.T, env *Environment, key string, value any) {
	t.Helper()
	found, err := env.Get(key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	val := found.(*ValueLiteral)
	if val.Literal != value {
		t.Errorf("Expected %v for key %s, got %v", value, key, found)
	}
}
