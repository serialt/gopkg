package gopkg

import "testing"

type bTest struct {
	name string
	age  uint8
	male bool
}

type BTest struct {
	Name string
	Age  uint8
	Male bool
}

func TestByteCommon_GetBytes(t *testing.T) {
	if data, err := GetBytes(&bTest{name: "test", age: 18, male: true}); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := GetBytes(&BTest{Name: "test", Age: 18, Male: true}); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := GetBytes(100); nil != err { // [4 4 0 255 200]
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := GetBytes(true); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := GetBytes("100"); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
}

func TestByteCommon_IntToBytes(t *testing.T) {
	if data, err := IntToBytes(100); nil != err { // [0 0 0 100]
		t.Log(err)
	} else {
		t.Log(data)
	}
}

func TestByteCommon_BytesToInt(t *testing.T) {
	if data, err := IntToBytes(100); nil != err { // [0 0 0 100]
		t.Log(err)
	} else {
		t.Log(data)
		if dataInt, err := BytesToInt(data); nil != err { // [0 0 0 100]
			t.Log(err)
		} else {
			t.Log(dataInt)

		}
	}
}

func TestByteCommon_Uint16ToBytes(t *testing.T) {
	t.Log(Uint16ToBytes(100))
}

func TestByteCommon_BytesToUint16(t *testing.T) {
	data := Uint16ToBytes(100)
	t.Log(BytesToUint16(data))
}

func TestByteCommon_Uint32ToBytes(t *testing.T) {
	t.Log(Uint32ToBytes(100))
}

func TestByteCommon_BytesToUint32(t *testing.T) {
	data := Uint32ToBytes(100)
	t.Log(BytesToUint32(data))
}

func TestByteCommon_Uint64ToBytes(t *testing.T) {
	t.Log(Uint64ToBytes(100))
}

func TestByteCommon_BytesToUint64(t *testing.T) {
	data := Uint64ToBytes(100)
	t.Log(BytesToUint64(data))
}
