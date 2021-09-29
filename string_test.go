package gopkg

import (
	"reflect"
	"testing"
)

func TestStringIsEmpty(t *testing.T) {
	want := "serialt"
	got := StringIsEmpty("")
	if reflect.DeepEqual(want, got) {
		t.Errorf("expected:%v, got:%v", want, got)
	}

}
