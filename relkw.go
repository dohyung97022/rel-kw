package relkw

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetRelKw to get a array of relevent keywords
func GetRelKw(keyword string) (relKw []string, err error) {
	chkeywords := make(chan interface{})
	chFinished := make(chan bool)
	for x := 0; x <= 25; x++ {
		go func(x int, chkeywords chan interface{}, chFinshed chan bool) {
			var keywords interface{}
			err := getJSON("http://suggestqueries.google.com/complete/search?client=chrome&hl=kr&q="+keyword+"+"+string(rune('a'+x)), keywords)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				chkeywords <- nil
				chFinshed <- true
			}
			chkeywords <- keywords
			chFinshed <- true
		}(x, chkeywords, chFinished)
	}
	for x := 0; x <= 25; {
		select {
		case keywords := <-chkeywords:
			fmt.Printf(fmt.Sprintf("%v", keywords))
			fmt.Printf("\n")
		case <-chFinished:
			x++
		}
	}
	return nil, nil
}

// Use like getJSON("http://example.com", &struct)
func getJSON(url string, jsnInter interface{}) error {
	fmt.Printf("%v\n", url)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", strings.ReplaceAll(between(string(contents), ",[", "],["), "\"", ""))

	return nil
}

func between(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}
