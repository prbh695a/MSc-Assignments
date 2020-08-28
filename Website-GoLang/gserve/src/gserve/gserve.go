package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var tpl *template.Template
var tplFuncMap template.FuncMap
var serverName = ""

//This Function is used to extract celltype from document:Kapitel1 which is document
func cellType(s string) string {
	arr := strings.Split(s, ":")
	s1 := arr[0]
	//s2 := arr[1]
	return s1
}

//This Function is used to extract celltype from document:Kapitel1 which is kapitel1
func cellValue(s string) string {
	arr := strings.Split(s, ":")
	//s1 := arr[0]
	s2 := arr[1]
	return s2
}

//This function initialize the template to be used for GET request and bind user defined functions
//to extract celltype and cellValue with these
func init() {
	fmap := template.FuncMap{
		"cellType":  cellType,
		"cellValue": cellValue,
	}
	tpl = template.Must(template.New("main").Funcs(fmap).ParseGlob("templates/*"))
}

//This function encodes the data sent by POST request into Base-64
func encoder(unencodedJSON []byte) []byte {
	var unencodedRows RowsType
	json.Unmarshal(unencodedJSON, &unencodedRows)
	// encode fields in Go objects
	encodedRows := unencodedRows.encode()
	// convert encoded Go objects to JSON
	encodedJSON, _ := json.Marshal(encodedRows)
	return encodedJSON
}

//This function decodes the data received by GET request from Base-64 to JSON to structure
func decoder(bodyBytes []byte) interface{} {
	var encodedRows64 EncRowsType
	json.Unmarshal(bodyBytes, &encodedRows64)

	decodedRows, _ := encodedRows64.decode()
	fmt.Println("decodedRows,", decodedRows)

	return decodedRows
}

//This Function help to POST Request to Hbase
func postToHbase(r *http.Request) {
	unencodedJSON, _ := ioutil.ReadAll(r.Body)
	payload := bytes.NewReader(encoder(unencodedJSON))
	req, _ := http.NewRequest(http.MethodPut, "http://"+"hbase"+":8080/se2:library/fakerow", payload)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
}

//This Function help to get the scanner details to scan the Hbase
func getScanner() string {
	payload1 := strings.NewReader("<Scanner batch=\"10\"/>")
	req1, _ := http.NewRequest(http.MethodPut, "http://"+"hbase"+":8080/se2:library/scanner", payload1)
	req1.Header.Set("Content-Type", "text/xml")
	req1.Header.Set("Accept", "text/plain")
	req1.Header.Set("Accept-Encoding", "identity")
	resp1, _ := http.DefaultClient.Do(req1)
	defer resp1.Body.Close()
	locationURL, _ := resp1.Location()
	return locationURL.String()

}

//This functions GET data from Hbase
func getFromHbase(r *http.Request) []byte {
	scanner := getScanner()
	req2, _ := http.NewRequest(http.MethodGet, scanner, nil)
	req2.Header.Set("Accept", "application/json")
	resp2, _ := http.DefaultClient.Do(req2)
	defer resp2.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp2.Body)
	return bodyBytes
}

//This function handles dynamic requests (non nginx)
func dynamic(w http.ResponseWriter, r *http.Request) {
	postToHbase(r)
	bodyBytes := getFromHbase(r)
	decodedRows := decoder(bodyBytes)
	tpl.ExecuteTemplate(w, "dataTemplate.gohtml", decodedRows)
	fmt.Fprintf(w, "Proudly served by *****%s*****", serverName)
}

func must(err error) {
	if err != nil {
		fmt.Printf("\n Error during zookeeper connect %+v", err)
	}
}

func connect() *zk.Conn {
	zksStr := "zookeeper:2181"
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second)
	must(err)
	//println(err)

	return conn
}

func main() {
	serverName = os.Getenv("server_name")
	conn := connect()
	defer conn.Close()

	//Delay logic so that Hbase connections are up
	for conn.State() != zk.StateHasSession {
		fmt.Printf("\n Wating for zookeeper to be up %s :", serverName)
		time.Sleep(30 * time.Second)
	}

	flags := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)
	path := "/mirror" + serverName
	existance, _, _ := conn.Exists(path)
	if !existance {
		_, err := conn.Create("/mirror/"+serverName, []byte(serverName+":9090"), flags, acl)
		must(err)
	} else {
		println("Node already exists hence not creating**************************")
	}
	http.HandleFunc("/", dynamic)
	http.ListenAndServe(":9090", nil)

}
