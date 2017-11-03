//
package main
 
import (
	"fmt"
	"time"
	"encoding/xml"
	"os"
	"log"
	"io/ioutil"
	"net/http"
)

type WebConf struct {
	ID string `xml:"id"`
	Addr string `xml:"addr"`
	Log LogConf `xml:"log"`
}

type LogConf struct {
	Path string `xml:"path"`
}

type Web struct {
	Http HttpServer
}

type HttpServer struct {
	Addr string
}

type DlServer struct {
	
}
type StatusServer struct {
	
}

var logptr *log.Logger

func main(){
	var web Web 
	cfgfile, err := os.Open("./conf.xml")
	if err != nil {
		fmt.Println("open conf file failed, err : ", err)
		return
	}
	defer func(){

		cfgfile.Close()
	}()

	data, err := ioutil.ReadAll(cfgfile)
	if err != nil {
		fmt.Println("read cfg file failed, err : ", err)
		return
	}

	v := WebConf{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Println("Unmarshal failed, err : ", err)
		return
	}
	fmt.Println(v, string(data), v.Addr)	
	
	logfile, err := os.OpenFile(v.Log.Path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err : ", err)
		return
	}
	defer func(){
		logfile.Close()
	}()
	

	logptr = log.New(logfile, "[web]", log.Ldate|log.Ltime|log.Lshortfile)
	if logptr == nil {
		return 
	}

	web.Http.Addr = v.Addr
	
	logptr.Println("test")
	logptr.Println("test2")

	go httpServer(&web.Http)
	time.Sleep(10000 * time.Second)
}

func httpServer(hs *HttpServer){
	var dls DlServer
	var ss StatusServer
	http.Handle("/download", dls)
	http.Handle("/status", ss)
	log.Fatal(http.ListenAndServe(hs.Addr, nil))
}


func (server DlServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logptr.Println("download")
	for i, v := range r.Header {
		fmt.Println(i, v)
	}
	fmt.Println(r.Method)
	fmt.Println(r.Body)

	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	w.Write([]byte(r.URL.Path + string(body)))
}

func (server StatusServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new req")
	logptr.Println("status")
	w.Write([]byte(r.URL.Path))
}
