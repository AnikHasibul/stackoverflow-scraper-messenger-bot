package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type m struct {
	Object string
	Entry  []entry
}
type entry struct {
	Id        string
	Time      int64
	Messaging []messaging
}

type messaging struct {
	sender struct {
		Id string
	}
	reciepent struct {
		Id string
	}
	message struct {
		Mid  string
		Text string
	}
}

const (
	facebook_api = "https://graph.facebook.com/v2.6/me/messages?access_token="
)

var (
	Answer string
)

func botHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		challenge := r.URL.Query().Get("hub.challenge")
		token := r.URL.Query().Get("hub.verify_token")

		if token == VERIFY_TOKEN {
			w.WriteHeader(200)
			w.Write([]byte(challenge))
			return
		} else {
			w.WriteHeader(404)
			w.Write([]byte("Not a verified request!"))
			return
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body:", http.StatusInternalServerError)
		CheckErr("Error on post data:", err)
		return
	}

	var M m

	if err := json.Unmarshal(body, &M); err != nil {
		log.Println("Error on decoding json from post request", err)
		return
	}

	if M.Object != "page" {
		return
	}

	Sid := M.Entry[0].Messaging[0].sender.Id

	Msg := M.Entry[0].Messaging[0].message.Text
	Time := M.Entry[0].Time
	log.Println("MSG: %s || %s || %d", Sid, Msg, Time)
	input := strings.ToLower(Msg)
	if Sid == "" && Msg == "" && Time == 0 {
		log.Println("Invalid Request!")
		return
	}
	Ask(Sid, input)
}

func SendToFb(id, text string) {
	client := &http.Client{}
	textb, jerr := json.Marshal(text)
	if CheckErr("JsonEncodeOnFbOutput", jerr) != nil {
		return
	}
	messageData := `
	{
    recipient: {
      id: "` + id + `"
    },
    message: {
      text: ` + string(textb) + `
    }
   }`
	data := []byte(messageData)
	url := facebook_api + ACCESS_TOKEN
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	CheckErr("FbRequestErr", err)
	resp, err := client.Do(req)
	CheckErr("FbSendErr", err)
	body, BodyErr := ioutil.ReadAll(resp.Body)
	CheckErr("ErrorOnReadingBodyForFbSend: ", BodyErr)
	log.Println("RES:", string(body))
	defer resp.Body.Close()
}

func CheckErr(hints string, err error) error {
	if err != nil {
		log.Println(hints, err)
	}
	return err
}
