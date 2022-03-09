package sugar

import (
	"testing"
)

func TestMd5Sum(t *testing.T) {
	md5String := "i love imau, i am serialt"

	got := Md5Sum([]byte(md5String))
	t.Log("Got md5 sum:", got)
	want := "410bc65560b789803b8844efe03ecd3a"
	t.Log("Want md5 sum:", want)
	if got != want {
		t.Fatalf("Md5Sum faild:  got->%s   want->%s", got, want)
	}

}

func TestMd5String(t *testing.T) {
	md5String := "i love imau, i am serialt"

	got := Md5String(md5String)
	t.Log("Got md5 sum:", got)
	want := "410bc65560b789803b8844efe03ecd3a"
	t.Log("Want md5 sum:", want)
	if got != want {
		t.Fatalf("Md5Sum faild:  got->%s   want->%s", got, want)
	}

}

func TestMd5SumUpper(t *testing.T) {
	md5String := "i love imau, i am serialt"

	got := Md5SumUpper([]byte(md5String))
	t.Log("Got md5 sum:", got)
	want := "410BC65560B789803B8844EFE03ECD3A"
	t.Log("Want md5 sum:", want)
	if got != want {
		t.Fatalf("Md5SumUpper faild:  got->%s   want->%s", got, want)
	}

}

func TestMdStringUpper(t *testing.T) {
	md5String := "i love imau, i am serialt"

	got := Md5StringUpper(md5String)
	t.Log("Got md5 sum:", got)
	want := "410BC65560B789803B8844EFE03ECD3A"
	t.Log("Want md5 sum:", want)
	if got != want {
		t.Fatalf("Md5SumUpper faild:  got->%s   want->%s", got, want)
	}

}
