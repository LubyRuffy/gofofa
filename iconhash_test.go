package gofofa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	faviconOkHandler = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon.ico", "/favicon1.ico", "/favicon.jpg", "/favicon.gif", "/favicon.png", "/favicon.bmp":
			http.ServeFile(w, r, "./data"+r.URL.Path)
			return
		case "/favicon_noheader.ico":
			d, _ := os.ReadFile("./data/favicon.ico")
			w.Write(d)
			return
		case "/":
			// 主要用于测试favicon
			w.Write([]byte(`
<!doctype html>
<html>
  <head >
    <title>网络空间测绘，网络空间安全搜索引擎，网络空间搜索引擎，安全态势感知 - FOFA网络空间测绘系统</title>
	<link rel="icon" as="image" type="image/x-icon" href="/favicon1.ico">
  </head>
  <body >
    hello world
  </body>
</html>

`))
			return
		}
	}

	faviconRelHandler = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon_rel.ico":
			http.ServeFile(w, r, "./data/favicon.ico")
			return
		case "/":
			// 主要用于测试favicon
			w.Write([]byte(`
<!doctype html>
<html>
  <head >
    <title>网络空间测绘，网络空间安全搜索引擎，网络空间搜索引擎，安全态势感知 - FOFA网络空间测绘系统</title>
	<link rel="icon" as="image" type="image/x-icon" href="/favicon_rel.ico">
  </head>
  <body >
    hello world
  </body>
</html>

`))
			return
		}
	}

	faviconAbsHandler = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon_abs.ico":
			http.ServeFile(w, r, "./data/favicon.ico")
			return
		case "/":
			// 主要用于测试favicon
			w.Write([]byte(`
<!doctype html>
<html>
  <head >
    <title>网络空间测绘，网络空间安全搜索引擎，网络空间搜索引擎，安全态势感知 - FOFA网络空间测绘系统</title>
	<link rel="icon" as="image" type="image/x-icon" href="http://` + r.Host + `/favicon_abs.ico">
  </head>
  <body >
    hello world
  </body>
</html>

`))
			return
		}
	}

	faviconNoHTMLHandler = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon.ico":
			http.ServeFile(w, r, "./data"+r.URL.Path)
			return
		case "/":
			w.Write([]byte(`hello world`))
			return
		}
	}

	favicon404Handler = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			// 主要用于测试favicon
			w.Write([]byte(`
<!doctype html>
<html>
  <head >
    <title>网络空间测绘，网络空间安全搜索引擎，网络空间搜索引擎，安全态势感知 - FOFA网络空间测绘系统</title>
	<link rel="icon" as="image" type="image/x-icon" href="http://` + r.Host + `/favicon_404.ico">
  </head>
  <body >
    hello world
  </body>
</html>

`))
			return
		}
	}
)

func TestIconHash(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(faviconOkHandler))
	defer ts.Close()

	var hash string
	var err error

	// 本地文件不存在
	hash, err = IconHash("./data/aaa.ico")
	assert.Contains(t, err.Error(), "icon url is not valid url")

	// 本地文件存在
	hash, err = IconHash("./data/favicon.ico")
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)

	// URL存在，有header
	hash, err = IconHash(ts.URL + "/favicon.ico")
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)
	// URL存在，png格式
	hash, err = IconHash(ts.URL + "/favicon.png")
	assert.Nil(t, err)
	assert.Equal(t, "-343282923", hash)
	// URL存在，gif格式
	hash, err = IconHash(ts.URL + "/favicon.gif")
	assert.Nil(t, err)
	assert.Equal(t, "-466535725", hash)
	// URL存在，jpg格式
	hash, err = IconHash(ts.URL + "/favicon.jpg")
	assert.Nil(t, err)
	assert.Equal(t, "-366292100", hash)
	// URL存在，bmp格式
	hash, err = IconHash(ts.URL + "/favicon.bmp")
	assert.Nil(t, err)
	assert.Equal(t, "-1520915571", hash)

	//URL存在，没有header
	hash, err = IconHash(ts.URL + "/favicon_noheader.ico")
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)

	// URL不是图片文件，而是一个页面
	ts1 := httptest.NewServer(http.HandlerFunc(faviconNoHTMLHandler))
	defer ts1.Close()
	hash, err = IconHash(ts1.URL)
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)

	// 相对地址
	ts2 := httptest.NewServer(http.HandlerFunc(faviconRelHandler))
	defer ts2.Close()
	hash, err = IconHash(ts2.URL)
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)
	// 绝对地址
	ts3 := httptest.NewServer(http.HandlerFunc(faviconAbsHandler))
	defer ts3.Close()
	hash, err = IconHash(ts3.URL)
	assert.Nil(t, err)
	assert.Equal(t, "-247388890", hash)

	// 404地址
	ts4 := httptest.NewServer(http.HandlerFunc(favicon404Handler))
	defer ts4.Close()
	hash, err = IconHash(ts4.URL)
	assert.Contains(t, err.Error(), "can not find any icon")

}
