package relkw

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GetRelKw to get a array of relevent keywords
func GetRelKw(keyword string) (sliceKw []string, err error) {
	keyword = strings.ReplaceAll(keyword, " ", "+")
	mapKeywords := make(map[string]bool)
	chkeywords := make(chan []string)
	chFinished := make(chan bool)
	// ----------------------------- google ----------------------------
	for x := 0; x <= 25; x++ {
		go func(x int, chkeywords chan []string, chFinshed chan bool) {
			resSlice, err := getGoogleJSON("http://suggestqueries.google.com/complete/search?client=chrome&hl=kr&q=" + keyword + "+" + string(rune('a'+x)))
			if err != nil {
				fmt.Printf("error: %v\n", err)
				chkeywords <- nil
				chFinshed <- true
			}
			chkeywords <- resSlice
			chFinshed <- true
		}(x, chkeywords, chFinished)
	}
	for x := 0; x <= 25; {
		select {
		case keywords := <-chkeywords:
			for _, keyword := range keywords {
				if !strings.Contains(keyword, "xaml") {
					mapKeywords[keyword] = true
				}
			}
		case <-chFinished:
			x++
		}
	}
	// ----------------------------- bing ----------------------------
	for x := 0; x <= 25; x++ {
		go func(x int, chkeywords chan []string, chFinshed chan bool) {
			resSlice, err := getBingJSON("https://www.bing.com/AS/Suggestions?pt=page.home&cp=1&cvid=" +
				randomStrFromCharset(22, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") +
				"&qry=" + keyword + "+" + string(rune('a'+x)))
			if err != nil {
				fmt.Printf("error: %v\n", err)
				chkeywords <- nil
				chFinshed <- true
			}
			chkeywords <- resSlice
			chFinshed <- true
		}(x, chkeywords, chFinished)
	}
	for x := 0; x <= 25; {
		select {
		case keywords := <-chkeywords:
			for _, keyword := range keywords {
				mapKeywords[keyword] = true
			}
		case <-chFinished:
			x++
		}
	}
	for keyword := range mapKeywords {
		sliceKw = append(sliceKw, keyword)
	}
	return sliceKw, nil
}

// Contact of a youtuber
type Contact struct {
	chan_url      string
	title         string
	facebook      string
	facebookGroup string
	twitter       string
	instagram     string
	email         string
}

// GetRelYt to get all info from youtube
func GetRelYt(sliceKw []string) (contacts []Contact, err error) {
	//
	for _, keyword := range sliceKw {
		r, err := http.Get("http://http://ec2-54-161-234-228.compute-1.amazonaws.com:3000/search?search=" + strings.ReplaceAll(keyword, " ", "+") + "&getallpage=true")
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		_, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// ----------------------------- scrape ----------------------------
func getGoogleJSON(url string) (resSlice []string, err error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	resSlice = strings.Split(strings.ReplaceAll(between(string(contents), ",[", "],["), "\"", ""), ",")
	return resSlice, nil
}
func getBingJSON(url string) (resSlice []string, err error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("%v\n")
	beforeSlice := strings.Split(string(contents), "query=\"")
	for i, str := range beforeSlice {
		if i == 0 {
			continue
		}
		resSlice = append(resSlice, before(str, "\""))
	}
	return resSlice, nil
}

// ----------------------------- additional ----------------------------
func randomStrFromCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
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
func before(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}
