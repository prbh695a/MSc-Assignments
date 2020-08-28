package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var final []string
var count = 1

func must(err error) {
	if err != nil {
		//panic(err)
		fmt.Printf("%+v From must \n", err)
	}
}

//TO Connect to zookeeper
func connect() *zk.Conn {
	zksStr := "zookeeper:2181"
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second)
	must(err)
	//println(err)

	return conn
}

//This Function is handling request if /library then proxy to gserve instance
//Else serve from nginx
func handler() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if req.URL.Path == "/library" {
			count++ //This variable helps to implement round robin in incremental manner
			req.URL.Scheme = "http"
			req.URL.Host = final[count%(len(final))]
		} else {
			req.URL.Scheme = "http"
			req.URL.Host = "nginx"
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

//This Function is monitoring gserve1 and gserve1 and any state changes done by them are notified to grproxy

func mirror(conn *zk.Conn, path string) (chan []string, chan error) {
	snapshots := make(chan []string)
	errors := make(chan error)
	go func() {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()
	return snapshots, errors
}

func main() {
	conn := connect()

	defer conn.Close()
	for conn.State() != zk.StateHasSession {
		fmt.Printf("grproxy is loading the Zookeeper \n...")
		time.Sleep(10 * time.Second)
	}
	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)
	path := "/mirror"

	existance, _, _ := conn.Exists(path)

	if !existance {
		//conn.Delete("/mirror", -1)
		Cpath, err := conn.Create(path, []byte(""), flags, acl)
		must(err)
		fmt.Printf("create: %+v\n", Cpath)
	}

	snapshots, errors := mirror(conn, "/mirror")
	go func() {
		for {
			select {
			case snapshot := <-snapshots:
				fmt.Printf("%+v\n", snapshot)
				var temp []string
				for _, element := range snapshot {
					data, _, err := conn.Get("/mirror/" + element)
					must(err)
					println("current data=", string(data))
					temp = append(temp, string(data))
					count = 1
					//println("current temp", temp)
				}
				final = temp
				for j, i := range final {
					fmt.Print("index", j, "i=", i)
				}
			case err := <-errors:
				panic(err)
			}
		}
	}()

	proxy := handler()
	http.ListenAndServe(":80", proxy)
}
