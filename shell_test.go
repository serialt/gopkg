package gopkg

import (
	"reflect"
	"testing"
)

func TestFindCommandPath(t *testing.T) {
	got := "/usr/bin/cat"
	want, err := FindCommandPath("cat")
	if err != nil {
		t.Errorf("cat find cmd path %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected:%v, got:%v", want, got)
	}
}

func TestRunCmd(t *testing.T) {
	want := "serialt"
	got, err := RunCmd("echo -en serialt")
	if err != nil {
		t.Errorf("cat find cmd path %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected: %v -- %T, got:%v -- %T", want, want, got, got)
	}
}

func TestRunCMD(t *testing.T) {
	want := "serialt"
	got, err := RunCMD("echo -en serialt")
	if err != nil {
		t.Errorf("cat find cmd path %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected: %v -- %T, got:%v -- %T", want, want, got, got)
	}
}
