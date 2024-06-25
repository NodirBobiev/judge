package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/NodirBobiev/judge/internals/kafka"
)

var (
	kafkaClient *kafka.Client
	once        sync.Once
)

func main() {
	initializeKafkaClient()
	defer kafkaClient.Close()

	http.Handle("/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/upload", HandleUpload)

	fmt.Println("Starting server on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func initializeKafkaClient() {
	once.Do(func() {
		var err error
		kafkaClient, err = kafka.NewKafkaClient([]string{"localhost:29092"}) // Kafka broker address.
		if err != nil {
			panic(fmt.Sprintf("Failed to create Kafka client: %v", err))
		}
	})
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	topic := "file-uploads"
	message := kafka.Submission{
		Filename: handler.Filename,
		Content:  content,
	}

	err = kafkaClient.SendMessage(topic, message)
	if err != nil {
		http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully and sent to Kafka topic '%s'", topic)
}
