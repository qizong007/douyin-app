package util

import (
	"douyin-app/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWT(t *testing.T) {
	conf.InitConf("./conf/default_conf.yaml")
	InitJWTVal()
	userId := int64(1 << 50)
	for i := 0; i < 10000; i++ {

		token, err := GenerateToken(userId)
		assert.Equal(t, nil, err)

		c, err := ParseToken(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, userId, c.UserId)

		id, e := CheckToken(token)
		assert.Equal(t, userId, id)
		assert.Equal(t, nil, e)
	}
}
