package user

import (
	"github.com/jinykim0x80/42api_gateway/internal"
	"testing"
)

func TestValid(t *testing.T) {
	var u Users
	u = Get()
	if err := util.ReadJSON("user_test.json", &u); err != nil {
		t.Fatal(err)
	}

	if len(u) == 0 {
		t.Fatal("Fail to loading user")
	}
	Set(u)

	tests := []struct {
		name     []string
		expected []string
	}{
		{[]string{"jinykim", "jinykim1", "jiny", "test1", "", "yepark"}, []string{"jinykim", "yepark"}},
		{[]string{"", "", "getnextline", "a"}, []string{}},
		{nil, []string{}},
	}

	f := 0
	for _, tt := range tests {
		var vu []string
		GetValid(tt.name, &vu)
		t.Logf("Expected:%v, Result:%v\n", tt.expected, vu)
		if len(vu) != len(tt.expected) {
			t.Fail()
			f++
		}
		vu = []string{}
	}
	t.Logf("\tCoverage:%d/%d", len(tests)-f, len(tests))
}
