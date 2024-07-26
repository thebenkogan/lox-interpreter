package evaluator

import "testing"

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	_, err := env.Get("a")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	env.Set("a", 1)
	assertEnv(t, env, "a", 1)

	env.Set("b", 2)
	assertEnv(t, env, "b", 2)

	env.Set("b", 3)
	assertEnv(t, env, "b", 3)

	env.CreateScope()

	assertEnv(t, env, "b", 3)

	env.Set("c", 4)
	assertEnv(t, env, "c", 4)

	env.Set("b", 5)
	assertEnv(t, env, "b", 5)

	env.ExitScope()

	_, err = env.Get("c")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	assertEnv(t, env, "b", 3)
}

func assertEnv(t *testing.T, env *Environment, key string, value any) {
	t.Helper()
	found, err := env.Get(key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found != value {
		t.Errorf("Expected %v for key %s, got %v", value, key, found)
	}
}
