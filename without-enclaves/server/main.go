package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// data related to one client
type Client struct {
	data          []string
	isFree        bool
	isReadyForPSI bool
}

var client1 Client
var client2 Client
var psi []string

// Listen on an end-point to get the data
// When the end-point has been called twice, the server
// compute the private set intersection, and return the result
// to the second caller
func main() {
	log.Println("PSI server started")

	// parse the command line options
	port := flag.String("port", "8080", "The port for the server")
	flag.Parse()

	// initialize the context for both clients
	initClientData()

	http.HandleFunc("/upload", uploadFileClient)
	http.ListenAndServe(":"+*port, nil)
}

// initialize the context for both clients
func initClientData() {
	client1.isFree = true
	client2.isFree = true
	client1.isReadyForPSI = false
	client2.isReadyForPSI = false
}

// handler for the /upload-client endpoint
func uploadFileClient(w http.ResponseWriter, r *http.Request) {
	if client1.isFree {
		client1.isFree = false
		processClientRequest(w, r, &client1)
	} else if client2.isFree {
		client2.isFree = false
		processClientRequest(w, r, &client2)
		// PSI is done, we restart the context for a next PSI computation
		initClientData()
	} else {
		w.Header().Set("Content-Type", "plain/text")
		fmt.Fprintf(w, "Sorry, I am busy, next time may be?\n")
	}
}

// process the data uploaded by one client
func processClientRequest(w http.ResponseWriter, r *http.Request, client *Client) {
	processPostRequest(w, r, client)
}

// process the Post Request sent by one client
func processPostRequest(w http.ResponseWriter, r *http.Request, client *Client) {
	log.Println("uploadFile() started")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	// read the file
	// 64 is the size of for one SHA256 record
	count := (handler.Size / 64)
	log.Println("count : ", count)
	client.data = make([]string, count)

	// read the file, line by line
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	i := 0
	for fileScanner.Scan() {
		client.data[i] = fileScanner.Text()
		i++
	}

	client.isReadyForPSI = true

	w.Header().Set(
		"Content-Type",
		"plain/text",
	)
	if client1.isReadyForPSI && client2.isReadyForPSI {
		computePSI()
		for _, row := range psi {
			io.WriteString(w, row+"\n")
		}

		initClientData()
	} else {
		fmt.Fprintf(w, "You are the first one, I am waiting for your partner\n")
	}
	log.Println("uploadFile() finished")
}

// compute the PSI on the data provided by both clients
func computePSI() {
	log.Println("computePSI() started")

	// catch the time in order to compute the duration at the end of that function
	start := time.Now()

	i, j := 0, 0
	for (i < len(client1.data)) && (j < len(client2.data)) {
		if client1.data[i] == client2.data[j] {
			psi = append(psi, client1.data[i])
			i++
			j++
		} else {
			if client1.data[i] < client2.data[j] {
				i++
			} else {
				j++
			}
		}
	}

	// display the duration of the PSI computation
	duration := time.Since(start)
	log.Println("[duration (in ms)] : PSI computation : ", duration.Milliseconds())

	log.Println("computePSI() finished")
}
