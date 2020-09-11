package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"unicode"

	"github.com/go-pg/pg/v10"
)

// API :
type API struct {
	db *pg.DB
}

// NewAPI :
func NewAPI(db *pg.DB) *API {
	return &API{db}
}

// rBodyStruct :
type rBodyStruct struct {
	Result  interface{} `json:"Result"`
	Ok      bool        `json:"Ok"`
	Message string      `json:"Message"`
}

// Rbody :
func Rbody(data interface{}, Ok bool, Message string) interface{} {
	temp := new(rBodyStruct)
	temp.Result = data
	temp.Ok = Ok
	temp.Message = Message
	return temp
}

// poemStruct :
type poemStruct struct {
	Content string `json:"Content"`
}

// Poem :
func Poem(Ok bool) string {
	// return "aaa"
	ps := new(poemStruct)
	if Ok {
		poem, err := http.Get("https://v1.jinrishici.com/all.json")
		if err != nil {
			body, _ := ioutil.ReadAll(poem.Body)
			json.Unmarshal(body, &ps)
		} else {
			ps.Content = "NO WIFI !!!!"
		}
	} else {
		poem, err := http.Get("https://api.lovelive.tools/api/SweetNothings")
		if err != nil {
			body, _ := ioutil.ReadAll(poem.Body)
			json.Unmarshal(body, &ps)
		} else {
			ps.Content = "NO WIFI !!!!"
		}
	}

	return ps.Content
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
