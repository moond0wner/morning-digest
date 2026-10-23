package core_testutils

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("получили %v, ожидали %v", got, want)
	}
}
