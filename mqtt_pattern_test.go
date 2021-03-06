package mqttpattern_test

import "testing"

import mqttpattern "github.com/amir-yaghoubi/mqtt-pattern"

import "fmt"

func TestMattches(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		topic   string
		result  bool
	}{
		{name: "Supports empty patterns", pattern: "", topic: "foo/bar/baz", result: false},
		{name: "Supports empty patterns", pattern: "foo/bar/baz", topic: "", result: false},
		{name: "Supports patterns with no wildcards", pattern: "foo/bar/baz", topic: "foo/bar/baz", result: true},
		{name: "Doesn't match different pattern and topic", pattern: "foo/bar/baz", topic: "bar/foo/baz", result: false},
		{name: "Supports # at the beginning", pattern: "#", topic: "foo/bar/baz", result: true},
		{name: "Supports # at the end", pattern: "foo/#", topic: "foo/bar/baz", result: true},
		{name: "Supports # at the end and topic has no children", pattern: "foo/bar/#", topic: "foo/bar/baz", result: true},
		{name: "Doesn't support # wildcards with more after them", pattern: "#/foo/bar", topic: "foo/bar/baz", result: false},
		{name: "Supports patterns with + at the beginning", pattern: "+/bar/baz", topic: "foo/bar/baz", result: true},
		{name: "Supports patterns with + at the end", pattern: "foo/bar/+", topic: "foo/bar/baz", result: true},
		{name: "Supports patterns with + at the middle", pattern: "foo/+/baz", topic: "foo/bar/baz", result: true},
		{name: "Supports patterns with multiple wildcards", pattern: "foo/+/#", topic: "foo/bar/baz", result: true},
		{name: "Supports named wildcards", pattern: "foo/+something/#else", topic: "foo/bar/baz", result: true},
		{name: "Supports leading slashes", pattern: "/foo/bar/baz", topic: "/foo/bar/baz", result: true},
		{name: "Supports leading slashes with invalid topic", pattern: "/foo/bar", topic: "/bar/foo", result: false},
	}

	for _, tt := range testCases {
		if mqttpattern.Matches(tt.pattern, tt.topic) != tt.result {
			t.Error(tt.name)
		}
	}
}

func TestExtract(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		topic   string
		result  map[string]string
	}{
		{name: "Returns empty map if pattern is empty", pattern: "", topic: "foo/bar/baz", result: make(map[string]string)},
		{name: "Returns empty map if topic is empty", pattern: "foo/+bar/+baz", topic: "", result: make(map[string]string)},
		{name: "Returns empty map if pattern and topic are empty", pattern: "", topic: "", result: make(map[string]string)},
		{name: "Returns empty map if wildcards don't have label", pattern: "foo/+/#", topic: "foo/bar/baz", result: make(map[string]string)},
		{name: "Returns map with a rest of topic as string for # wildcard", pattern: "foo/#bar", topic: "foo/bar/baz", result: map[string]string{"bar": "bar/baz"}},
		{name: "Returns map with a string for + wildcard", pattern: "foo/+bar/+baz", topic: "foo/bar/baz", result: map[string]string{"bar": "bar", "baz": "baz"}},
		{name: "Parses params from all wildcards", pattern: "+foo/+bar/#baz", topic: "foo/bar/baz", result: map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}},
	}

	for _, tt := range testCases {
		result := mqttpattern.Extract(tt.pattern, tt.topic)

		rStr := fmt.Sprintf("%v", result)
		trStr := fmt.Sprintf("%v", tt.result)
		if rStr != trStr {
			t.Errorf("%s | expected %s but received %s", tt.name, trStr, rStr)
		}
	}
}

func TestExec(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		topic   string
		result  map[string]string
	}{
		{name: "Returns nil if doesn't match", pattern: "foo/bar", topic: "foo/bar/baz", result: nil},
		{name: "Returns params if they can be parsed", pattern: "foo/+bar/#baz", topic: "foo/bar/baz", result: map[string]string{"bar": "bar", "baz": "baz"}},
	}

	for _, tt := range testCases {
		result := mqttpattern.Exec(tt.pattern, tt.topic)

		rStr := fmt.Sprintf("%v", result)
		trStr := fmt.Sprintf("%v", tt.result)
		if rStr != trStr {
			t.Errorf("%s | expected %s but received %s", tt.name, trStr, rStr)
		}
	}
}

func TestFill(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		params  map[string]string
		result  string
	}{
		{name: "Returns empty string for empty patterns", pattern: "", params: map[string]string{"bar": "BAR"}, result: ""},
		{name: "Fills in pattern with both types of wildcards", pattern: "foo/+bar/#baz", params: map[string]string{"bar": "BAR", "baz": "BAZ"}, result: "foo/BAR/BAZ"},
		{name: "Fills missing + params with \"\"", pattern: "foo/+bar", params: make(map[string]string), result: "foo/"},
		{name: "Fills missing # params with \"\"", pattern: "foo/#bar", params: make(map[string]string), result: "foo/"},
	}

	for _, tt := range testCases {
		result := mqttpattern.Fill(tt.pattern, tt.params)
		if result != tt.result {
			t.Errorf("%s | expected %s but received %s", tt.name, tt.result, result)
		}
	}
}

func TestClean(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		result  string
	}{
		{name: "Returns empty string for empty patterns", pattern: "", result: ""},
		{name: "Works when there aren't any named parameter", pattern: "foo/+/bar/#", result: "foo/+/bar/#"},
		{name: "Removes named parameters", pattern: "foo/+something/bar/#otherthing", result: "foo/+/bar/#"},
	}

	for _, tt := range testCases {
		result := mqttpattern.Clean(tt.pattern)
		if result != tt.result {
			t.Errorf("%s | expected %s but received %s", tt.name, tt.result, result)
		}
	}
}
