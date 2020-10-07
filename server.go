package main

import (
	"ekans/app"
	"ekans/app/api"
	"ekans/app/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	ec "github.com/oscrud/echo"
	"github.com/oscrud/oscrud"
	"github.com/oscrud/postgres"
)

func main() {
	db := getDatabaseClient()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	server := oscrud.NewOscrud()
	server.RegisterTransport(ec.New(e).UsePort(3001))
	// server.RegisterLogger(app.Logger{})

	server.RegisterBinder(
		string(""),
		time.Time{},
		func(raw interface{}) (interface{}, error) {
			text := raw.(string)
			// log.Println(text)
			return time.Parse("2006-01-02T15:04:05.999999999Z07:00", text)
		},
	)

	server.RegisterBinder(
		[]interface{}{},
		[]int{},
		func(raw interface{}) (interface{}, error) {
			text := raw.([]interface{})
			newData := make([]int, len(text))
			for idx, data := range text {
				value, _ := strconv.Atoi(fmt.Sprintf("%v", data))
				newData[idx] = value
			}
			return newData, nil
		},
	)

	server.SetState("db", db)
	pg := postgres.New()
	pg.UseEncode(postgres.EncodeTypeHex)
	pg.UsePagination(postgres.PaginationTypeDefault, true)
	pg.UseCustomResponse(customBody)

	endpoint := api.NewAPI(db)

	server.RegisterEndpoint("GET", "/Login", endpoint.Login)
	server.RegisterEndpoint("POST", "/Register", endpoint.Register)
	server.RegisterEndpoint("GET", "/Calander", endpoint.Calander)
	server.RegisterEndpoint("GET", "/CheckRosterExist", endpoint.CheckRosterExist)
	server.RegisterService("/User", pg.ToService(db, new(model.User)), new(oscrud.ServiceOptions))
	server.RegisterService("/Group", pg.ToService(db, new(model.Group)), new(oscrud.ServiceOptions))
	server.RegisterService("/Roster", pg.ToService(db, new(model.Roster)), new(oscrud.ServiceOptions))
	server.RegisterService("/Area", pg.ToService(db, new(model.Area)), new(oscrud.ServiceOptions))
	server.RegisterService("/ShiftRequest", pg.ToService(db, new(model.ShiftRequest)), new(oscrud.ServiceOptions))
	server.RegisterService("/Leave", pg.ToService(db, new(model.Leave)), new(oscrud.ServiceOptions))
	server.Start()
}

func customBody(ctx oscrud.Context, sr *postgres.ServiceResult, err error) oscrud.Context {
	if err != nil {
		return ctx.JSON(http.StatusOK, api.Rbody(err.Error(), false, api.Poem(true)))
	}
	return ctx.JSON(http.StatusOK, api.Rbody(sr, true, api.Poem(true)))
}

func getDatabaseClient() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "admin",
		Database: "postgres",
	})

	//db.AddQueryHook(queryDebugger{})

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	}

	err = app.RegisterModels(db)
	if err != nil {
		panic(err)
	}
	log.Println("Postgres connected succesfully.")
	return db
}

type queryDebugger struct{}

// BeforeQuery :
func (qd queryDebugger) BeforeQuery(evt *pg.QueryEvent) {
	query, err := evt.FormattedQuery()
	log.Println(query, err)
}

func (qd queryDebugger) AfterQuery(evt *pg.QueryEvent) {

}
