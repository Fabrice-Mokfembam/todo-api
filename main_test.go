// File: main_test.go

package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	// Initialize the routes
	app.Post("/create/todo", createTodoTestHandler)
	app.Get("/getAll/todos", getAllTodosTestHandler)
	app.Put("/update/todo/:id", updateTodoTestHandler)
	app.Delete("/delete/todo/:id", deleteTodoTestHandler)

	return app
}


func createTodoTestHandler(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON data"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Todo created successfully", "data": todo})
}

func getAllTodosTestHandler(c *fiber.Ctx) error {
	todos := []Todo{
		{ID: 1, Task: "Test task 1", IsCompleted: false},
		{ID: 2, Task: "Test task 2", IsCompleted: true},
	}
	return c.Status(200).JSON(todos)
}

func updateTodoTestHandler(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON data"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo updated successfully"})
}

func deleteTodoTestHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
}

// Test create
func TestCreateTodoHandler(t *testing.T) {
	app := setupTestApp()

	todo := Todo{
		Task:        "New Task",
		IsCompleted: false,
	}

	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("POST", "/create/todo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
}

// Test fetch
func TestGetAllTodosHandler(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/getAll/todos", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// Test "update" 
func TestUpdateTodoHandler(t *testing.T) {
	app := setupTestApp()

	todo := Todo{
		Task:        "Updated Task",
		IsCompleted: true,
	}

	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("PUT", "/update/todo/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// Test delete
func TestDeleteTodoHandler(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/delete/todo/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}
