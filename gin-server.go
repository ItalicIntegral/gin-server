package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type todo struct {
	Id        string `json:"id"`
	Created   string `json:"created"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{"1", "2022-01-01", "Clean Room", false},
	{"2", "2022-03-01", "Vaccuum", true},
	{"3", "2022-01-01", "Clean Room", false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func getTodoById(id string) (*todo, error) {
	for i, t := range todos {
		if t.Id == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func addTodo(context *gin.Context) {
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	todos = append(todos, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	todo.Completed = !todo.Completed
	context.IndentedJSON(http.StatusOK, todo)
}

type Document struct {
	Id      int    `form:"Document.Id" binding:"required"`
	Type    string `form:"Document.Type" binding:"required"`
	Company Company
}

type Company struct {
	Name     string `form:"Document.Company.Name" binding:"required"`
	Address1 string `form:"Document.Company.Address1" binding:"required"`
	Address2 string `form:"Document.Company.Address2" binding:"required"`
	Phone    string `form:"Document.Company.Phone" binding:"required"`
	Email    string `form:"Document.Company.Email" binding:"required"`
	Website  string `form:"Document.Company.Website" binding:"required"`
}

func main() {
	// var customer Customer
	var document Document
	r := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/static", "./static")

	// Todos
	r.GET("/todos", getTodos)
	r.GET("/todos/:id", getTodo)
	r.PATCH("/todos/:id", toggleTodoStatus)
	r.POST("/todos", addTodo)

	r.GET("/", func(c *gin.Context) {
		//c.String(http.StatusOK, "pong")
		c.JSON(http.StatusOK, "This is the index")
	})
	r.POST("/upload", func(c *gin.Context) {
		// Bind file
		if err := c.ShouldBind(&document); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%+v", document))
		examplePdf(&document)
	})
	r.Run(":8080")
}

func examplePdf(d *Document) {
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Helvetica", "B", 24)
	pdf.Cell(1.25, .3, d.Type)
	pdf.SetFont("Helvetica", "", 20)
	pdf.Cell(1.25, .3, fmt.Sprint(" ", d.Id))

	// Company Info
	// pdf.SetXY(2, 2)
	pdf.Ln(0.5)
	pdf.SetFont("Helvetica", "", 14)
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Name))
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Phone))
	pdf.Ln(0.18)
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Address1))
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Email))
	pdf.Ln(0.18)
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Address2))
	pdf.Cell(2.5, .3, fmt.Sprint(" ", d.Company.Website))

	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("error")
	}
}
