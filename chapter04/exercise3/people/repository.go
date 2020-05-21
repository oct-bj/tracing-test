package people

import (
	"context"
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

func (r *Repository) GetPerson(ctx context.Context, name string) (model.Person, error){
	query := "SELECT title, description FROM people WHERE name=?"

	span, ctx := opentracing.StartSpanFromContext(ctx, "get-person",
		opentracing.Tag{Key: "db.statement", Value: query})
	defer span.Finish()

	rows, err := r.db.QueryContext(ctx, query, name)
	if err!= nil {
		return model.Person{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, descr string
		err := rows.Scan(ctx, &title, &descr)
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
