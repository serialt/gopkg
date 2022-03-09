package sugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {

	cnData := "红帽"
	enData := "red hat"
	t.Logf("测试中文翻译成英文")
	enTrans, _ := TranslateCh2En(cnData)
	t.Logf("中文原文： %v", cnData)
	t.Logf("翻译成英文结果： %v", enTrans)
	assert.Equal(t, enData, enTrans)

	t.Logf("测试英文翻译成中文")
	cnTrans, _ := TranslateEn2Ch(enData)
	t.Logf("英文原文： %v", enData)
	t.Logf("翻译成中文结果： %v", cnTrans)
	assert.Equal(t, cnData, cnTrans)

}
