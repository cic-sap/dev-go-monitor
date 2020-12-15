package monitor

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestWarpTransport(t *testing.T) {
	SetDebug(true)
	StartHttpClientTrace()
	req, err := http.NewRequest("GET", "https://httpbin.org/ip", nil)
	if err != nil {
		t.Fatal("err", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("err", err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("err", err)
	}
	log.Println("res", string(buf))

}
