package gopkg

import (
	"os"
	"testing"
)

func TestEnvCommon_Get(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	t.Log("HELLO =", EnvGet("HELLO"))
}

func TestEnvCommon_GetD(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	t.Log("HELLO =", EnvGetD("HELLO", "WORLD"))
	t.Log("WORLD =", EnvGetD("WORLD", "HELLO"))
}

func TestEnvCommon_GetInt(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := EnvGetInt("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := EnvGetInt("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetIntD(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", EnvGetIntD("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", EnvGetIntD("HELLO", 10))
}

func TestEnvCommon_GetInt64(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := EnvGetInt64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := EnvGetInt64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetInt64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", EnvGetInt64D("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", EnvGetInt64D("HELLO", 10))
}

func TestEnvCommon_GetUint64(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := EnvGetUint64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := EnvGetUint64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetUint64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", EnvGetUint64D("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", EnvGetUint64D("HELLO", 10))
}

func TestEnvCommon_GetFloat64(t *testing.T) {
	_ = os.Setenv("HELLO", "100.3254")
	i, _ := EnvGetFloat64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := EnvGetFloat64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetFloat64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100.3254")
	t.Log("HELLO =", EnvGetFloat64D("HELLO", 100.3254))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", EnvGetFloat64D("HELLO", 100.32541))
}

func TestEnvCommon_GetBool(t *testing.T) {
	_ = os.Setenv("HELLO", "true")
	t.Log("HELLO =", EnvGetBool("HELLO"))
	_ = os.Setenv("HELLO", "false")
	t.Log("HELLO =", EnvGetBool("HELLO"))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", EnvGetBool("HELLO"))
}
