package giu

import (
	"fmt"
	"testing"
)

type RequestPath struct {
	path       string
	exceptPath string
}

func fakeHandler() HandlerChain {
	return HandlerChain{
		func(c *Context) {
		}}
}

var requestSet = []RequestPath{
	{path: "/cmd/test", exceptPath: "/cmd/:tool/"},
	{path: "/cmd/test/", exceptPath: "/cmd/:tool/"},
	{path: "/cmd/test/3", exceptPath: "/cmd/:tool/:sub"},
	{path: "/cmd/who", exceptPath: "/cmd/:tool/"},
	{path: "/cmd/who/", exceptPath: "/cmd/:tool/"},
	{path: "/cmd/whoami", exceptPath: "/cmd/whoami"},
	{path: "/cmd/whoami/", exceptPath: "/cmd/whoami"},
	{path: "/cmd/whoami/r", exceptPath: "/cmd/:tool/:sub"},
	{path: "/cmd/whoami/r/", exceptPath: "/cmd/:tool/:sub"},
	{path: "/cmd/whoami/root", exceptPath: "/cmd/whoami/root"},
	{path: "/cmd/whoami/root/", exceptPath: "/cmd/whoami/root/"},
	{path: "/src/", exceptPath: "/src/*filepath"},
	{path: "/src/some/file.png", exceptPath: "/src/*filepath"},
	{path: "/search/", exceptPath: "/search/"},
	{path: "/search/someth!ng+in+ünìcodé", exceptPath: "/search/:query"},
	{path: "/search/gin", exceptPath: "/search/:query"},
	{path: "/search/gin-gonic", exceptPath: "/search/gin-gonic"},
	{path: "/search/google", exceptPath: "/search/google"},
	{path: "/files/js/inc/framework.js", exceptPath: "/files/:dir/*filepath"},
	{path: "/info/gordon/public", exceptPath: "/info/:user/public"},
	{path: "/info/gordon/project/go", exceptPath: "/info/:user/project/:project"},
	{path: "/info/gordon/project/golang", exceptPath: "/info/:user/project/golang"},
	{path: "/aa/aa", exceptPath: "/aa/*xx"},
	{path: "/ab/ab", exceptPath: "/ab/*xx"},
	{path: "/a", exceptPath: "/:cc"},
	//{path:,exceptPath: },
	//{path:,exceptPath: },

}

func TestGetValue(t *testing.T) {

	var pathSet = []string{
		"/",
		"/cmd/:tool/",
		"/cmd/:tool/:sub",
		"/cmd/whoami",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/src/*filepath",
		"/search/",
		"/search/:query",
		"/search/gin-gonic",
		"/search/google",
		"/files/:dir/*filepath",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
		"/info/:user/project/golang",
		"/aa/*xx",
		"/ab/*xx",
		"/:cc",
		"/c1/:dd/e",
		"/c1/:dd/e1",
		"/:cc/cc",
		"/:cc/:dd/ee",
		"/:cc/:dd/:ee/ff",
		"/:cc/:dd/:ee/:ff/gg",
		"/:cc/:dd/:ee/:ff/:gg/hh",
		"/get/test/abc/",
		"/get/:param/abc/",
		"/something/:paramname/thirdthing",
		"/something/secondthing/test",
		"/get/abc",
		"/get/:param",
		"/get/abc/123abc",
		"/get/abc/:param",
		"/get/abc/123abc/xxx8",
		"/get/abc/123abc/:param",
		"/get/abc/123abc/xxx8/1234",
		"/get/abc/123abc/xxx8/:param",
		"/get/abc/123abc/xxx8/1234/ffas",
		"/get/abc/123abc/xxx8/1234/:param",
		"/get/abc/123abc/xxx8/1234/kkdd/12c",
		"/get/abc/123abc/xxx8/1234/kkdd/:param",
		"/get/abc/:param/test",
		"/get/abc/123abd/:param",
		"/get/abc/123abddd/:param",
		"/get/abc/123/:param",
		"/get/abc/123abg/:param",
		"/get/abc/123abf/:param",
		"/get/abc/123abfff/:param",
	}
	var requestSet = []RequestPath{
		{path: "/cmd/test", exceptPath: "/cmd/:tool/"},
		{path: "/cmd/test/", exceptPath: "/cmd/:tool/"},
		{path: "/cmd/test/3", exceptPath: "/cmd/:tool/:sub"},
		{path: "/cmd/who", exceptPath: "/cmd/:tool/"},
		{path: "/cmd/who/", exceptPath: "/cmd/:tool/"},
		{path: "/cmd/whoami", exceptPath: "/cmd/whoami"},
		{path: "/cmd/whoami/", exceptPath: "/cmd/whoami"},
		{path: "/cmd/whoami/r", exceptPath: "/cmd/:tool/:sub"},
		{path: "/cmd/whoami/r/", exceptPath: "/cmd/:tool/:sub"},
		{path: "/cmd/whoami/root", exceptPath: "/cmd/whoami/root"},
		{path: "/cmd/whoami/root/", exceptPath: "/cmd/whoami/root/"},
		{path: "/src/", exceptPath: "/src/*filepath"},
		{path: "/src/some/file.png", exceptPath: "/src/*filepath"},
		{path: "/search/", exceptPath: "/search/"},
		{path: "/search/someth!ng+in+ünìcodé", exceptPath: "/search/:query"},
		{path: "/search/gin", exceptPath: "/search/:query"},
		{path: "/search/gin-gonic", exceptPath: "/search/gin-gonic"},
		{path: "/search/google", exceptPath: "/search/google"},
		{path: "/files/js/inc/framework.js", exceptPath: "/files/:dir/*filepath"},
		{path: "/info/gordon/public", exceptPath: "/info/:user/public"},
		{path: "/info/gordon/project/go", exceptPath: "/info/:user/project/:project"},
		{path: "/info/gordon/project/golang", exceptPath: "/info/:user/project/golang"},
		{path: "/aa/aa", exceptPath: "/aa/*xx"},
		{path: "/ab/ab", exceptPath: "/ab/*xx"},
		{path: "/a", exceptPath: "/:cc"},
		//{path:,exceptPath: },
		//{path:,exceptPath: },

	}
	root := newNode()
	for _, s := range pathSet {
		root.addRoute(s, fakeHandler())
	}
	for _, path := range requestSet {
		path_1 := path.exceptPath
		path_ := path.path
		fmt.Println(path_)
		fmt.Println(path_1)
		value := root.Search(path.path)

		if value.fullPath != path.exceptPath {
			t.Fail()
		}
	}
}

func TestCheckValid(t *testing.T) {
	for i := 0; i < len(requestSet); i++ {
		valid := checkSearchPath(requestSet[i].path)
		if !valid {
			t.Fail()
		}
	}
	return
}
