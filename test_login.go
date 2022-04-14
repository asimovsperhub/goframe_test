package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	data := make(url.Values)
	data["name"] = []string{"asimov12345678"}
	data["password"] = []string{"asimov123"}
	data["Nickename"] = []string{"王1234"}
	resp, err := http.PostForm("http://127.0.0.1:8000/register", data)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

//curl -H "Content-Type: application/json" -X POST -d '{"name": "asimov", "password":"asimov123"}' "http://127.0.0.1:8000/login"
