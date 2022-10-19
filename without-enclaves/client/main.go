package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// upload a file on the PSI server thanks to a HTTP Post request
//
// receive in return either a message inviting to wait for the
// uploading of the second fellow or the private set intersection
// result
func main() {

	log.Println("main() function started")

	// parse the command line options
	remoteURL := flag.String("remoteUrl", "http://localhost:8080/upload",
		"The targeted url for uploading the file")
	fileName := flag.String("file", "data.csv",
		"Filename to upload")
	flag.Parse()

	// prepare the data to send
	values := map[string]io.Reader{
		"myFile": mustOpen(*fileName),
	}

	// upload the file on the server
	err := upload(*remoteURL, values)
	if err != nil {
		log.Panic(err)
	}
}

// Upload a file on a http server
func upload(url string, values map[string]io.Reader) (err error) {
	log.Println("upload() function started")

	// catch the time in order to compute the duration at the end of that function
	start := time.Now()

	// prepare the form to submit to the URL
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// add the file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			// add other potential fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// close the multipart writer, otherwise the request will
	// not terminate the boundary properly
	w.Close()

	// preparation of the Post request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	// set the content type, this will contain the boundary
	req.Header.Set("Content-Type", w.FormDataContentType())

	// submit the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}

	// display the duration of the upload (and potentially the PSI computation)
	duration := time.Since(start)
	log.Println("[duration (in ms)] : upload + potentially the PSI computation : ", duration.Milliseconds())

	// save the result of the POST request in a file
	out, err := os.Create("result.txt")
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, res.Body)

	return
}

// open a file, and generate a Panic error in case there is an issue
func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		log.Panic(err)
	}
	return r
}
