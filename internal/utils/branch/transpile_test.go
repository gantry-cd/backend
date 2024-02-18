package branch

import (
	"regexp"
	"testing"
)

const (
	kubeLabelRegexp = `(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?`
)

func FuzzTranspile1123(f *testing.F) {
	testCases := []string{
		"test",
		"test-branch",
		"test.branch",
		"test/branch",
		"test_branch",
		"test:branch",
		"test;branch",
		"test,branch",
		"test branch",
		"test?branch",
		"test&branch",
		"test=branch",
		"test+branch",
		"test@branch",
		"test#branch",
		"test!branch",
		"test~branch",
		"test*branch",
		"test(branch",
		"test)branch",
		"test[branch",
		"test]branch",
		"test{branch",
		"test}branch",
		"test|branch",
		"test\\branch",
		"test^branch",
		"test%branch",
		"test`branch",
	}

	for _, s := range testCases {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		name := Transpile1123(s)

		ok, err := regexp.Match(kubeLabelRegexp, []byte(name))
		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("name %s is not valid", name)
		}

		decoded, err := TranspileBranchName(name)
		if err != nil {
			t.Fatal(err)
		}

		if s != decoded {
			t.Fatalf("decoded %s is not equal to original %s", decoded, s)
		}
	})
}
