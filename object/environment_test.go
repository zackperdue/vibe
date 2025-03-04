package object_test

import (
	"testing"

	"github.com/vibe-lang/vibe/object"
)

// TestNewEnvironment tests the NewEnvironment function
func TestNewEnvironment(t *testing.T) {
	env := object.NewEnvironment()

	if env == nil {
		t.Fatal("NewEnvironment() returned nil")
	}
}

// TestNewEnclosedEnvironment tests the NewEnclosedEnvironment function
func TestNewEnclosedEnvironment(t *testing.T) {
	outer := object.NewEnvironment()
	inner := object.NewEnclosedEnvironment(outer)

	if inner == nil {
		t.Fatal("NewEnclosedEnvironment() returned nil")
	}
}

// TestEnvironmentGet tests the Get method of the environment
func TestEnvironmentGet(t *testing.T) {
	env := object.NewEnvironment()

	// Set a variable in the environment
	intObj := &object.Integer{Value: 5}
	env.Set("x", intObj)

	// Test getting the variable
	retrievedObj, ok := env.Get("x")
	if !ok {
		t.Fatal("Failed to get variable from environment")
	}

	integer, ok := retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 5 {
		t.Fatalf("Expected integer value to be 5, got %d", integer.Value)
	}

	// Test getting a non-existent variable
	_, ok = env.Get("nonexistent")
	if ok {
		t.Fatal("Get() returned ok for non-existent variable")
	}
}

// TestEnvironmentSet tests the Set method of the environment
func TestEnvironmentSet(t *testing.T) {
	env := object.NewEnvironment()

	// Set a variable in the environment
	intObj := &object.Integer{Value: 5}
	env.Set("x", intObj)

	// Check if the variable is set correctly
	retrievedObj, ok := env.Get("x")
	if !ok {
		t.Fatal("Failed to get variable after setting it")
	}

	integer, ok := retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 5 {
		t.Fatalf("Expected integer value to be 5, got %d", integer.Value)
	}

	// Test updating an existing variable
	newIntObj := &object.Integer{Value: 10}
	env.Set("x", newIntObj)

	retrievedObj, ok = env.Get("x")
	if !ok {
		t.Fatal("Failed to get variable after updating it")
	}

	integer, ok = retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 10 {
		t.Fatalf("Expected updated integer value to be 10, got %d", integer.Value)
	}
}

// TestEnvironmentEnclosedGet tests variable retrieval from enclosed environments
func TestEnvironmentEnclosedGet(t *testing.T) {
	outer := object.NewEnvironment()
	inner := object.NewEnclosedEnvironment(outer)

	// Set a variable in the outer environment
	outerIntObj := &object.Integer{Value: 5}
	outer.Set("x", outerIntObj)

	// Test getting the variable from the inner environment
	retrievedObj, ok := inner.Get("x")
	if !ok {
		t.Fatal("Failed to get variable from outer environment through inner environment")
	}

	integer, ok := retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 5 {
		t.Fatalf("Expected integer value to be 5, got %d", integer.Value)
	}

	// Set a variable with the same name in the inner environment
	innerIntObj := &object.Integer{Value: 10}
	inner.Set("x", innerIntObj)

	// Test that the inner environment retrieves its own value, not the outer one
	retrievedObj, ok = inner.Get("x")
	if !ok {
		t.Fatal("Failed to get variable from inner environment")
	}

	integer, ok = retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 10 {
		t.Fatalf("Expected integer value to be 10, got %d", integer.Value)
	}

	// Test that the outer environment still has its original value
	retrievedObj, ok = outer.Get("x")
	if !ok {
		t.Fatal("Failed to get variable from outer environment")
	}

	integer, ok = retrievedObj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object, got %T", retrievedObj)
	}

	if integer.Value != 5 {
		t.Fatalf("Expected integer value to be 5, got %d", integer.Value)
	}
}

// TestEnvironmentNestedScopes tests multiple levels of nested environments
func TestEnvironmentNestedScopes(t *testing.T) {
	global := object.NewEnvironment()
	first := object.NewEnclosedEnvironment(global)
	second := object.NewEnclosedEnvironment(first)
	third := object.NewEnclosedEnvironment(second)

	// Set variables at different scope levels
	global.Set("a", &object.Integer{Value: 1})
	first.Set("b", &object.Integer{Value: 2})
	second.Set("c", &object.Integer{Value: 3})
	third.Set("d", &object.Integer{Value: 4})

	// Test that each level can access variables from outer scopes
	for _, tc := range []struct {
		env      *object.Environment
		name     string
		expected int64
		exists   bool
	}{
		{global, "a", 1, true},
		{global, "b", 0, false},
		{global, "c", 0, false},
		{global, "d", 0, false},

		{first, "a", 1, true},
		{first, "b", 2, true},
		{first, "c", 0, false},
		{first, "d", 0, false},

		{second, "a", 1, true},
		{second, "b", 2, true},
		{second, "c", 3, true},
		{second, "d", 0, false},

		{third, "a", 1, true},
		{third, "b", 2, true},
		{third, "c", 3, true},
		{third, "d", 4, true},
	} {
		retrievedObj, ok := tc.env.Get(tc.name)

		if tc.exists {
			if !ok {
				t.Fatalf("Environment %v failed to get variable %s that should exist",
					tc.env, tc.name)
			}

			integer, ok := retrievedObj.(*object.Integer)
			if !ok {
				t.Fatalf("Expected Integer object, got %T", retrievedObj)
			}

			if integer.Value != tc.expected {
				t.Fatalf("Expected integer value to be %d, got %d",
					tc.expected, integer.Value)
			}
		} else {
			if ok {
				t.Fatalf("Environment %v got variable %s that should not exist",
					tc.env, tc.name)
			}
		}
	}

	// Test variable shadowing by setting the same name at different levels
	global.Set("x", &object.Integer{Value: 100})
	first.Set("x", &object.Integer{Value: 200})
	second.Set("x", &object.Integer{Value: 300})

	for _, tc := range []struct {
		env      *object.Environment
		expected int64
	}{
		{global, 100},
		{first, 200},
		{second, 300},
		{third, 300}, // third inherits from second
	} {
		retrievedObj, ok := tc.env.Get("x")
		if !ok {
			t.Fatalf("Environment %v failed to get variable x", tc.env)
		}

		integer, ok := retrievedObj.(*object.Integer)
		if !ok {
			t.Fatalf("Expected Integer object, got %T", retrievedObj)
		}

		if integer.Value != tc.expected {
			t.Fatalf("Expected integer value to be %d, got %d",
				tc.expected, integer.Value)
		}
	}
}