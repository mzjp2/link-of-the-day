package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mzjp2/link-of-the-day/link"
	"github.com/mzjp2/link-of-the-day/storage"
)

func newLinkHandler(date time.Time, svc storage.Service) http.Handler {
	link := new(linkHandler)
	link.date = date
	link.svc = svc
	return link
}

type linkHandler struct {
	date time.Time
	svc  storage.Service
}

func (l *linkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, err := link.GetURL(l.svc, l.date)
		if err != nil {
			log.Print(err)
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
		}
		err = link.SaveURL(l.svc, r.Form.Get("url"), time.Now())
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	svc, err := storage.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	http.Handle("/link", newLinkHandler(time.Now(), svc))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
