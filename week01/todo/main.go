package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

type Todo struct {
	ID int
	Title string
	Completed bool
	CreateAt time.Time
}

var db *gorm.DB

func main(){
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=todolist sslmode=disable")
	if err != nil{
		log.Fatal("Failed to connect to DB" + err.Error())
	}

	defer db.Close()

	db.LogMode(true)

	err = db.AutoMigrate(Todo{}).Error
	if err != nil{
		log.Fatal("Failed to migrate table todo")
	}

	router := gin.Default()

	router.GET("/todos", listTodos)
	router.POST("/todos", createTodo)
	router.GET("/todos/:id", getTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run(":8088")
}

func listTodos(c *gin.Context){
	var todos []Todo
	err := db.Find(&todos).Error

	if err != nil {
		c.String(500, "Failed to list todoList")
		return
	}
	c.JSON(200, todos)
}


func getTodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo Todo
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, todo)
	}
}

func createTodo(c *gin.Context) {

	var argument struct {
		Title string
	}
	//var todo Todo
	err := c.BindJSON(&argument)
	if err != nil {
		c.JSON(500, err.Error())
	}

	todo := Todo{
		Title: argument.Title,
	}
	err = db.Create(&todo).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, todo)
}

func deleteTodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo Todo
	d := db.Where("id = ?", id).Delete(&todo)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
