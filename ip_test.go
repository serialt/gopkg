package gopkg

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const checkMark = " OK! "
const ballotX = " ERROR! "

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ip address", IPGet(r))
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestIPCommon_Get(t *testing.T) {
	statusCode := http.StatusOK

	server := mockServer()
	defer server.Close()

	t.Log("Given the need to test downloading content.")
	{

		t.Logf("\tWhen checking \"%s\" for status code \"%d\"", server.URL, statusCode)
		{
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatal("\tShould be able to make the Get call.", ballotX, err)
			}
			t.Log("\t\tShould be able to make the Get call.", checkMark)
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode == statusCode {
				t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
			}
		}
	}
}

func TestIPCommon_Get_XRealIP(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Real-IP", "192.168.0.1")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}

func TestIPCommon_Get_XForwardedFor(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Forwarded-For", "192.168.0.2")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}

func TestIPCommon_Get_XForwardedForIP6(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Forwarded-For", "::1")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}

func TestGetPubIP(t *testing.T) {

	want := ""
	got, all := GetPubIP()
	assert.NotEqual(t, want, got)
	t.Logf("\t\tget public ip: %v, all: %v ", got, all)
}

func TestSubNetMaskToLen(t *testing.T) {
	want24 := 24
	got24, err := SubNetMaskToLen("255.255.255.0")
	assert.Empty(t, err)
	assert.Equal(t, want24, got24)

	want16 := 16
	got16, err := SubNetMaskToLen("255.255.0.0")
	assert.Empty(t, err)
	assert.Equal(t, want16, got16)
}

func TestLenToSubNetMask(t *testing.T) {
	want8 := "255.0.0.0"
	want32 := "255.255.255.255"
	got32 := LenToSubNetMask(32)
	assert.Equal(t, want32, got32)

	got8 := LenToSubNetMask(8)
	assert.Equal(t, want8, got8)
}

func TestIsPublicIPv4(t *testing.T) {
	ip1 := net.ParseIP("8.8.8.8")
	ip2 := net.ParseIP("119.29.29.29")
	ip3 := net.ParseIP("10.10.10.5")
	ip4 := net.ParseIP("127.0.0.1")
	assert.True(t, IsPublicIPv4(ip1))
	assert.True(t, IsPublicIPv4(ip2))
	assert.False(t, IsPublicIPv4(ip3))
	assert.False(t, IsPublicIPv4(ip4))
}
