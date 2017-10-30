package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/autlamps/delay-frontend-confirm/confirm"
)

var dburl string

func init() {
	flag.StringVar(&dburl, "DATABASE_URL", "", "database url")
	flag.Parse()

	if dburl == "" {
		dburl = os.Getenv("DATABASE_URL")
	}
}

func main() {
	c := confirm.Conf{dburl}

	r, err := confirm.Create(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":5000", r))
}
