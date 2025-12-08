package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/kissanjamgit/pornbox"
	"resty.dev/v3"
)

type Strmup struct {
	url string
}

func New(url string) Strmup {
	return Strmup{url: url}
}

func (s *Strmup) Video() (cr pornbox.ContentResource, err error) {
	id := regexp.MustCompile("[^/]+$").FindString(s.url)
	return VideoByID(id)
}

func VideoByID(id string) (cr pornbox.ContentResource, err error) {
	client := resty.New()
	defer client.Close()
	url_ := fmt.Sprintf("https://strmup.to/ajax/stream?filecode=%s", id)
	resp, err := client.R().Get(url_)
	if err != nil {
		fmt.Println(err)
		return
	}
	var value any
	err = json.Unmarshal(resp.Bytes(), &value)
	if err != nil {
		return
	}
	map_, ok := value.(map[string]any)
	if !ok {
		err = fmt.Errorf("value is't map, type: %s, value: %v", reflect.TypeOf(value).String(), value)
		return
	}
	url, ok := map_["streaming_url"].(string)
	if !ok {
		err = fmt.Errorf("url.TypeOf != string value: %v", value)
		return
	}
	title, ok := map_["title"].(string)
	if !ok {
		err = fmt.Errorf("title.TypeOf != string value: %v", value)
		return
	}

	cr = pornbox.ContentResource{Name: func() string {
		if title != "" {
			return title
		}
		return fmt.Sprintf("strmup_video_%s", id)
	}(), Url: url}
	return
}

func main() {
	id := flag.String("id", "", "")
	url := flag.String("i", "", "")
	name := flag.Bool("s", false, "")

	flag.Parse()
	if *id == "" && *url == "" {
		fmt.Println(`*id == "" && *url == ""`)
		os.Exit(1)
	}
	strmup := New(*url)
	cr, err := strmup.Video()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if *name {
		fmt.Printf("%s %s", cr.Name, cr.Url)
		return
	}
	fmt.Println(cr.Url)
}
