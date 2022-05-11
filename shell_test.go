package gopkg

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCommandPath(t *testing.T) {
	got := "/usr/bin/cat"
	want, err := FindCommandPath("cat")
	if err != nil {
		t.Errorf("cat find cmd path %v", err)
	}
	if IsLinux() {
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected:%v, got:%v", want, got)
		}
	}
	if IsMac() {
		got = "/bin/cat"
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected:%v, got:%v", want, got)
		}
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

func TestUserDir(t *testing.T) {
	dir := UserHomeDir()
	assert.NotEmpty(t, dir)

	dir = UserDir("sub-path")
	assert.Contains(t, dir, "/sub-path")

	dir = UserCacheDir("my-logs")
	assert.Contains(t, dir, ".cache/my-logs")

	dir = UserConfigDir("my-conf")
	assert.Contains(t, dir, ".config/my-conf")

}

func TestWorkdir(t *testing.T) {
	assert.NotEmpty(t, Workdir())
}

func TestLoginUser(t *testing.T) {
	cu := LoginUser()
	assert.NotEmpty(t, cu)

	fu := MustFindUser(cu.Username)
	assert.NotEmpty(t, fu)
	assert.Equal(t, cu.Uid, fu.Uid)
}
