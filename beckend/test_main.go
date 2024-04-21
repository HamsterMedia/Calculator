package main

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		{expr: "2 + 2", expected: "4"},
		{expr: "3 * 4", expected: "12"},
		{expr: "10 / 2", expected: "5"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Expr: %s", test.expr), func(t *testing.T) {
			result, err := Eval(test.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("got %s, want %s", result, test.expected)
			}
		})
	}
}

func TestEvalVars(t *testing.T) {
	tests := []struct {
		expr     string
		vars     map[string]interface{}
		expected string
	}{
		{
			expr:     "x + y",
			vars:     map[string]interface{}{"x": 3, "y": 5},
			expected: "8",
		},
		{
			expr:     "a * b",
			vars:     map[string]interface{}{"a": 2, "b": 4},
			expected: "8",
		},
		{
			expr:     "x / y",
			vars:     map[string]interface{}{"x": 10, "y": 2},
			expected: "5",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Expr: %s with vars: %v", test.expr, test.vars), func(t *testing.T) {
			result, err := EvalVars(test.expr, test.vars)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("got %s, want %s", result, test.expected)
			}
		})
	}
}
