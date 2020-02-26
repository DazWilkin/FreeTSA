package main

import (
	"bytes"
	"crypto"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/digitorus/timestamp"
)

var (
	mesg = flag.String("mesg", "Frederik Jack is a bubbly Border Collie", "Message to be timestamped by FreeTSA")
	file = flag.String("file", "file", "Filename to use for the Timestamp Request (`/tmp/${FILE}.tsq`) and Response (`/tmp/${FILE}.tsr`)")
)

func main() {
	flag.Parse()

	// DER-encoded Request
	rqst, err := timestamp.CreateRequest(strings.NewReader(*mesg), &timestamp.RequestOptions{
		Hash: crypto.SHA256,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DER-encoded TSR: %x\n", rqst)

	// Persist the Request as a TSQ to verify using OpenSSL
	if err := save(fmt.Sprintf("%s.tsq", *file), rqst); err != nil {
		log.Println(err)
	}

	tsr, err := timestamp.ParseRequest(rqst)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Timestamp Request: \n%+v", tsr)
	log.Printf("Rqst HashedMessage: %x\n", tsr.HashedMessage)

	h := crypto.SHA256.New()
	h.Write([]byte(*mesg))
	s := h.Sum(nil)
	log.Printf("Calc HashedMessage: %x\n", s)

	if !bytes.Equal(tsr.HashedMessage, s) {
		log.Fatal(err)
	}

	// Submit the request to FreeTSR
	client := NewClient(&http.Client{})

	// DER-encoded Response
	resp, err := client.Request(rqst)
	if err != nil {
		log.Fatal(err)
	}

	// Persist the Response as a TSR to verify using OpenSSL
	if err := save(fmt.Sprintf("%s.tsr", *file), resp); err != nil {
		log.Println(err)
	}

	ts, err := timestamp.ParseResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Timestamp Response: \n%+v", ts)
	log.Printf("Resp HashedMessage: %x\n", tsr.HashedMessage)
}
func save(name string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("/tmp/%s", name), data, 0644)
}
