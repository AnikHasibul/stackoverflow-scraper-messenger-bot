package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Json structure for messenger incoming
type MessageModel struct {
	Object string
	Entry  []struct {
		Id        string
		Time      int64
		Messaging []struct {
			Sender struct {
				Id string
			}
			Reciepent struct {
				Id string
			}
			Message struct {
				Mid  string
				Text string
			}
		}
	}
}

// fb api v2.6 url
const (
	facebook_api = "https://graph.facebook.com/v2.6/me/messages?access_token="
)

func botHandler(w http.ResponseWriter, r *http.Request) {
	// Verify the webhook
	if r.Method == "GET" {
		challenge := r.URL.Query().Get(
			"hub.challenge",
		)
		token := r.URL.Query().Get(
			"hub.verify_token",
		)

		if token == VERIFY_TOKEN {
			// Success
			w.WriteHeader(200)
			w.Write([]byte(challenge))
		} else {
			// Oops
			w.WriteHeader(404)
			w.Write([]byte(
				"Not a verified request!",
			))
		}
		return
	}

	// Now have fun with json
	var M MessageModel
	if err := json.NewDecoder(r.Body).
		Decode(&M); err != nil {
		panic(err)
		return
	}

	if M.Object != "page" {
		return
	}

	// The data we need
	Sid := M.Entry[0].
		Messaging[0].
		Sender.Id

	Msg := M.Entry[0].
		Messaging[0].
		Message.Text

	Time := M.Entry[0].Time

	// Logging the message variables
	log.Printf(
		"MSG: %s || %s || %d \n",
		Sid,
		Msg,
		Time,
	)

	// Take a lower case input!
	input := strings.ToLower(Msg)

	// Don't response empty shits
	if Sid == "" &&
		Msg == "" &&
		Time == 0 {
		log.Println("Invalid Request!")
		return
	}

	// Now Ask it
	Ask(Sid, input)
}

// SendToFb sends a message to facebook
func SendToFb(id, text string) {
	// Declare a client
	client := &http.Client{}
	// Encode the text to json
	textb, err := json.Marshal(text)
	// Error? :/
	if err != nil {
		log.Fatal(err)
	}

	// Build it
	messageData := `
	{
    recipient: {
      id: "` + id + `"
    },
    message: {
      text: ` + string(textb) + `
    }
   }`
	// Prepare it
	data := []byte(messageData)
	url := facebook_api + ACCESS_TOKEN

	// Send It
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(data),
	)
	req.Header.Add(
		"Content-Type",
		"application/json",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Boom!
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RES:", string(body))
	defer resp.Body.Close()
}

func CheckErr(hints string, err error) error {
	if err != nil {
		log.Println(hints, err)
	}
	return err
}
