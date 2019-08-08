package main

import (
	"fmt"
	"net/url"
	"os"
)

//
//type URL struct {
//	Scheme     string
//	Opaque     string    // encoded opaque data
//	User       *Userinfo // username and password information
//	Host       string    // host or host:port
//	Path       string    // path (relative paths may omit leading slash)
//	RawPath    string    // encoded path hint (see EscapedPath method)
//	ForceQuery bool      // append a query ('?') even if RawQuery is empty
//	RawQuery   string    // encoded query values, without '?'
//	Fragment   string    // fragment for references, without '#'
//}
func main(){
	fmt.Printf("PATH=%s\n",  os.Getenv("PATH"))

	u := &url.URL{Path:"www.baidu.com/code/a=3&b=4"}
	fmt.Println(u.RawPath)

	fmt.Println(u.EscapedPath())
}
