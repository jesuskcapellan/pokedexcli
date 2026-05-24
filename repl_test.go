package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "OH NO THIS ISN'T THE CAR",
			expected: []string{"oh", "no", "this", "isn't", "the", "car"},
		},
		{
			input:    "  gibby      guppy    ",
			expected: []string{"gibby", "guppy"},
		},
		{
			input:    "foobar",
			expected: []string{"foobar"},
		},
		{
			input:    "foo bar",
			expected: []string{"foo", "bar"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Actual length: %d did not match expected length: %d.\nFAIL", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Actual word: %s did not match expected word: %s.\nFAIL", word, expectedWord)
			}
		}
	}
}
