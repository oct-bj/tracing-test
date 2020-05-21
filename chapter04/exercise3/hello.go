package main

import (
	"log"
	"strings"
	"net/http"
	"context"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/shunfenger-tech/tracing-test/chapter04/exercise3/people"
	"github.com/opentracing/opentracing-go"
	"github.com/shunfenger-tech/tracing-test/chapter04/lib/tracing"
)

var (
	repo *people.Repository
)

func main(){
	repo = people.NewRepository()
	defer repo.Close()

	tr, closer := tracing.Init("go-2-hello")
	defer closer.Close()
	opentracing.SetGlobalTracer(tr)

	http.HandleFunc("/sayHello/", handleSayHello)

	log.Print("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSayHello(w http.ResponseWriter, r *http.Request){
	span := opentracing.GlobalTracer().StartSpan("say-hello")
	defer span.Finish()
	context := opentracing.ContextWithSpan(r.Context(), span)

	name := strings.TrimPrefix(r.URL.Path, "/sayHello/")
	greeting, err := SayHello(context, name)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(otlog.Error(err))

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetTag("response", greeting)
	w.Write([]byte(greeting))
}

func SayHello(context context.Context, name string) (string, error){
	person, err := repo.GetPerson(context, name)
	if err!= nil {
		return "", err
	}

	opentracing.SpanFromContext(context).LogKV("name", person.Name,
		"title", person.Title,
		"description", person.Description)

	return FormatGreeting(
		context,
		person.Name,
		person.Title,
		person.Description,
	), nil
}

func FormatGreeting(context context.Context, name, title, description string) string {
	span, context := opentracing.StartSpanFromContext(context,"format-greeting")
	defer span.Finish()

	response := "Hello, "
	if title != "" {
		response += title + " "
	}

	response += name + "!"

	if description != "" {
		response += " " + description
	}
	return response
}
