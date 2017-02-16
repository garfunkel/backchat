package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/ugjka/cleverbot-go"
)

func index(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("backchat.html")

	if err != nil {
		log.Fatal(err)
	}

	w.Write(data)
}

func chat(server *socketio.Server) {
	alice, err := cleverbot.New()

	if err != nil {
		log.Fatal(err)
	}

	bob, err := cleverbot.New()

	if err != nil {
		log.Fatal(err)
	}

	answer := "Think for me"

	for {
		for {
			answer, err = alice.Ask(answer)

			if err != nil {
				log.Print(err)

				time.Sleep(5 * time.Second)

				continue
			}

			break
		}

		log.Print("Alice: ", answer)

		server.BroadcastTo("chat", "Alice", answer)

		time.Sleep((time.Duration(len(answer)*100) * time.Millisecond) + time.Second)

		for {
			answer, err = bob.Ask(answer)

			if err != nil {
				log.Print(err)

				time.Sleep(5 * time.Second)

				continue
			}

			break
		}

		log.Print("Bob: ", answer)

		server.BroadcastTo("chat", "Bob", answer)

		time.Sleep((time.Duration(len(answer)*100) * time.Millisecond) + time.Second)
	}
}

func main() {
	chatStarted := false
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("Connection established")

		if !chatStarted {
			chatStarted = true

			go chat(server)
		}

		so.Join("chat")
	})

	http.HandleFunc("/", index)
	http.Handle("/socket.io/", server)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
