package people

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shunfenger-tech/tracing-test/chapter04/lib/model"
	"github.com/opentracing/opentracing-go"
)

const dburl = "root:mysqlpwd@tcp(127.0.0.1:3306)/chapter04"

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	db, err := sql.Open("mysql", dburl)
	if err!= nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err!= nil {
		log.Fatal("Cannot ping the db: %v/n", err)
	}

	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPerson(name string, span opentracing.Span) (model.Person, error){
	query := "SELECT title, description FROM people WHERE name=?"

	span = opentracing.GlobalTracer().StartSpan("get-person",
		opentracing.ChildOf(span.Context()),
		opentracing.Tag{Key: "db.statement", Value: query})
	defer span.Finish()

	rows, err := r.db.Query(query, name)
	if err!= nil {
		return model.Person{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, descr string
		err := rows.Scan(&title, &descr)
		if err != nil {
			return model.Person{}, err
		}
		return model.Person{
			Name: name,
			Title: title,
			Description: descr,
		}, nil
	}

	return model.Person{
		Name: name,
	}, nil
}

func (r *Repository) Close(){
	r.db.Close()
}
