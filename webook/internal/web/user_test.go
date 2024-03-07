package web

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	pwd := "hello#world123"

	// 对pwd进行加密
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	// 比较两个密码是否一致
	err = bcrypt.CompareHashAndPassword(password, []byte(pwd))
	assert.NoError(t, err)

}
