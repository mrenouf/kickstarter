package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"github.com/mrenouf/kickstarter/protos"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	client := http.Client{}
	list := protos.ProjectList{}
	projects := []*protos.Project{}
	page := 1
	for {
		values := url.Values{}
		values.Add("sort", "launch_date")
		values.Add("category_id", "16") // Technology
		values.Add("page", strconv.Itoa(page))
		url := "http://www.kickstarter.com/discover/advanced?" + values.Encode()
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		if err != nil {
			panic(err.Error())
		}

		if err != nil {
			panic(err.Error())
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		pl, err := parseBodyAsResult(resp)
		if err != nil {
			panic(err.Error())
		}
		projects = append(projects, pl.GetProjects()...)
		if len(projects) >= int(pl.GetTotalHits()) {
			break
		}
		fmt.Println(len(projects))
		page++
	}
	list.Projects = projects
	bytes, err := proto.Marshal(&list)
	if err != nil {
		panic(err.Error())
	}
	file, err := os.OpenFile("kickstarter.pb", os.O_CREATE|os.O_WRONLY, 0666)
	file.Write(bytes)
	file.Close()
}

func parseBodyAsResult(response *http.Response) (res *protos.ProjectList, err error) {
	var data []byte
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	res = &protos.ProjectList{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return
}
