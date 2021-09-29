package gopkg

import (
	"testing"
	"time"
)

func TestTimeCommon_Timestamp2String(t *testing.T) {
	t.Log(Timestamp2String(1568686616, 28889, "2006/01/02 15:04:05", time.Local))
	t.Log(Timestamp2String(1568686626, 98882, "2006/01/02 15:04:05", time.UTC))
}
