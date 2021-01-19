package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {

	var buf []byte
	var err error

	// 用例 1
	buf, err = GetBuffer(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(buf))
	assert.Equal(t, 1, cap(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
	// 用例 2
	buf, err = GetBuffer(65536)
	assert.Nil(t, err)
	assert.Equal(t, 65536, cap(buf))
	assert.Equal(t, 65536, len(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
	// 用例 3
	buf, err = GetBuffer(65535)
	assert.Nil(t, err)
	assert.Equal(t, 65536, cap(buf))
	assert.Equal(t, 65535, len(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
	// 用例 4
	buf, err = GetBuffer(32768)
	assert.Nil(t, err)
	assert.Equal(t, 32768, cap(buf))
	assert.Equal(t, 32768, len(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
	// 用例 5
	buf, err = GetBuffer(32769)
	assert.Nil(t, err)
	assert.Equal(t, 65536, cap(buf))
	assert.Equal(t, 32769, len(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
	// 用例 6
	buf, err = GetBuffer(32767)
	assert.Nil(t, err)
	assert.Equal(t, 32768, cap(buf))
	assert.Equal(t, 32767, len(buf))
	err = PutBuffer(buf)
	assert.Nil(t, err)
	buf = nil
}
