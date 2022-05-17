package gofofa

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFofaURLFromEnv(t *testing.T) {
	os.Setenv("FOFA_SERVER", "https://1.1.1.1")
	os.Setenv("FOFA_EMAIL", "a@a.com")
	os.Setenv("FOFA_KEY", "123456")
	os.Unsetenv("FOFA_CLIENT_URL")
	assert.Equal(t, "https://1.1.1.1/?email=a@a.com&key=123456&version=v1", FofaURLFromEnv())

	// 异常
	os.Setenv("FOFA_CLIENT_URL", "\x7f")
	assert.Equal(t, "", FofaURLFromEnv())

	// 部分更新
	os.Setenv("FOFA_CLIENT_URL", "https://2.2.2.2/?email=b@b.com")
	assert.Equal(t, "https://2.2.2.2/?email=b@b.com&key=123456&version=v1", FofaURLFromEnv())

	// 全更新
	os.Setenv("FOFA_CLIENT_URL", "https://2.2.2.2/?email=b@b.com&key=000000&version=v2")
	assert.Equal(t, "https://2.2.2.2/?email=b@b.com&key=000000&version=v2", FofaURLFromEnv())
}
