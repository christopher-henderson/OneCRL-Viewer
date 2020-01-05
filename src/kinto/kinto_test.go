package kinto

import (
	"encoding/json"
	"net/http"
	"testing"
)

type Stuff struct {
	a string
}

func TestThin(t *testing.T) {
	//resp, err := client.Do(newRequest())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer resp.Body.Close()
	//var k KintoArray
	//err = json.NewDecoder(resp.Body).Decode(&k)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//q, _ := crtSh.BuildQuery(k.Data)
	//t.Log(q)
}

func TestLength(t *testing.T) {
	resp, err := http.DefaultClient.Do(newRequest())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var k KintoArray
	err = json.NewDecoder(resp.Body).Decode(&k)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(k.Data))
}
