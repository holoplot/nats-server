// Copyright 2016-2023 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package subject

import (
	"testing"
)

func checkBool(b, expected bool, t *testing.T) {
	t.Helper()
	if b != expected {
		t.Fatalf("Expected %v, but got %v\n", expected, b)
	}
}

func TestIsValidSubject(t *testing.T) {
	checkBool(isValidSubject("."), false, t)
	checkBool(isValidSubject(".foo"), false, t)
	checkBool(isValidSubject("foo."), false, t)
	checkBool(isValidSubject("foo..bar"), false, t)
	checkBool(isValidSubject(">.bar"), false, t)
	checkBool(isValidSubject("foo.>.bar"), false, t)
	checkBool(isValidSubject("foo"), true, t)
	checkBool(isValidSubject("foo.bar.>"), true, t)
	checkBool(isValidSubject("*"), true, t)
	checkBool(isValidSubject(">"), true, t)
	checkBool(isValidSubject("foo*"), true, t)
	checkBool(isValidSubject("foo**"), true, t)
	checkBool(isValidSubject("foo.**"), true, t)
	checkBool(isValidSubject("foo*bar"), true, t)
	checkBool(isValidSubject("foo.*bar"), true, t)
	checkBool(isValidSubject("foo*.bar"), true, t)
	checkBool(isValidSubject("*bar"), true, t)
	checkBool(isValidSubject("foo>"), true, t)
	checkBool(isValidSubject("foo.>>"), true, t)
	checkBool(isValidSubject("foo>bar"), true, t)
	checkBool(isValidSubject("foo.>bar"), true, t)
	checkBool(isValidSubject("foo>.bar"), true, t)
	checkBool(isValidSubject(">bar"), true, t)
}

func TestIsLiteralSubject(t *testing.T) {
	checkBool(isLiteralSubject("foo"), true, t)
	checkBool(isLiteralSubject("foo.bar"), true, t)
	checkBool(isLiteralSubject("foo*.bar"), true, t)
	checkBool(isLiteralSubject("*"), false, t)
	checkBool(isLiteralSubject(">"), false, t)
	checkBool(isLiteralSubject("foo.*"), false, t)
	checkBool(isLiteralSubject("foo.>"), false, t)
	checkBool(isLiteralSubject("foo.*.>"), false, t)
	checkBool(isLiteralSubject("foo.*.bar"), false, t)
	checkBool(isLiteralSubject("foo.bar.>"), false, t)
}

func TestSubjectNew(t *testing.T) {
	for _, test := range []struct {
		subject        string
		expectsSuccess bool
	}{
		{"foo.bar", true},
		{"foo..bar", false},
	} {
		t.Run("", func(t *testing.T) {
			s, err := New(test.subject)
			if test.expectsSuccess {
				if err != nil {
					t.Fatalf("subject.New(\"%s\" failed: %v", test.subject, err)
				}
				if s == nil {
					t.Fatalf("subject.New(\"%s\" returned nil", test.subject)
				}
			} else {
				if s != nil || err == nil {
					t.Fatalf("subject.New(\"%s\" should have failed", test.subject)
				}
			}
		})
	}
}

func TestSubjectIsSubsetMatch(t *testing.T) {
	for _, test := range []struct {
		subject string
		test    string
		result  bool
	}{
		{"foo.bar", "foo.bar", true},
		{"foo.*", ">", true},
		{"foo.*", "*.*", true},
		{"foo.*", "foo.*", true},
		{"foo.*", "foo.bar", false},
		{"foo.>", ">", true},
		{"foo.>", "*.>", true},
		{"foo.>", "foo.>", true},
		{"foo.>", "foo.bar", false},
	} {
		t.Run("", func(t *testing.T) {
			s, err := New(test.subject)
			if err != nil {
				t.Fatalf("subject.New(\"%s\" failed: %v", test.subject, err)
			}
			o, err := New(test.test)
			if err != nil {
				t.Fatalf("subject.New(\"%s\" failed: %v", test.test, err)
			}
			if res := s.IsSubsetMatch(o); res != test.result {
				t.Fatalf("Subject %q subset match of %q, should be %v, got %v",
					test.test, test.subject, test.result, res)
			}
		})
	}
}

func TestSubjectIsLiteral(t *testing.T) {
	for _, test := range []struct {
		subject string
		result  bool
	}{
		{"foo", true},
		{"foo.bar.*", false},
		{"foo.bar.>", false},
		{"*", false},
		{">", false},
		// The followings have widlcards characters but are not
		// considered as such because they are not individual tokens.
		{"foo*", true},
		{"foo**", true},
		{"foo.**", true},
		{"foo*bar", true},
		{"foo.*bar", true},
		{"foo*.bar", true},
		{"*bar", true},
		{"foo>", true},
		{"foo>>", true},
		{"foo.>>", true},
		{"foo>bar", true},
		{"foo.>bar", true},
		{"foo>.bar", true},
		{">bar", true},
	} {
		t.Run(test.subject, func(t *testing.T) {
			s, err := New(test.subject)
			if err != nil {
				t.Fatalf("subject.New(\"%s\" failed: %v", test.subject, err)
			}
			if res := s.IsLiteral(); res != test.result {
				t.Fatalf("IsLiteral() for subject \"%s\", should be %v, got %v", test.subject, test.result, res)
			}
		})
	}
}
