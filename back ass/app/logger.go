package app

import (
	"fmt"

	"github.com/oscrud/oscrud"
)

// Logger :
type Logger struct {
}

// Log :
func (l Logger) Log(operation string, content string) {
	fmt.Println("Operation - ", operation)
	fmt.Println("Content - ", content)
}

// StartRequest :
func (l Logger) StartRequest(ctx oscrud.Context) {
	fmt.Println("**************************************")
	fmt.Println("RequestID - ", ctx.RequestID())
	fmt.Println("Method - ", ctx.Method())
	fmt.Println("Path - ", ctx.Path())
	fmt.Println("Params - ", ctx.Params())
	fmt.Println("State - ", ctx.State())
	fmt.Println("Header - ", ctx.Headers())
	fmt.Println("Query - ", ctx.Query())
	fmt.Println("Body - ", ctx.Body())
	fmt.Println("**************************************")
}

// EndRequest :
func (l Logger) EndRequest(ctx oscrud.Context) {
	fmt.Println("**************************************")
	fmt.Println("RequestID - ", ctx.RequestID())
	fmt.Println("Method - ", ctx.Method())
	fmt.Println("Path - ", ctx.Path())
	fmt.Println("State - ", ctx.State())

	res := ctx.Response()
	fmt.Println("Headers - ", res.ResponseHeaders())
	fmt.Println("Status - ", res.Status())
	fmt.Println("ContentType - ", res.ContentType())
	fmt.Println("Result - ", res.Exception())
	fmt.Println("Error - ", res.Exception())
	fmt.Println("**************************************")
}
