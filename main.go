package main

import (
    "database/sql"
    "fmt"
    "log"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2"
    _ "github.com/go-sql-driver/mysql" 
)

type Todo struct {
    ID          int    `json:"id"`
    Task        string `json:"task"`
    IsCompleted bool   `json:"is_completed"`
}

func main() {
    fmt.Println("Todo API is starting...")

    app := fiber.New()
    app.Use(cors.New())

    db, err := connectToMySQL()
    if err != nil {
        log.Fatal("Could not connect to MySQL:", err)
    }
    defer db.Close()

    
    app.Get("/", func(c *fiber.Ctx)error{
         return c.SendString("working great")
    })

    app.Post("/create/todo", func(c *fiber.Ctx) error {
        return createTodoHandler(c, db)
    })


    app.Get("/getAll/todos", func(c *fiber.Ctx) error {
        return getAllTodosHandler(c, db)
    })

    app.Put("/update/todo/:id", func(c *fiber.Ctx) error {
        return updateTodoHandler(c, db)
    })

    app.Delete("/delete/todo/:id", func(c *fiber.Ctx) error {
        return deleteTodoHandler(c, db)
    })

    if err := app.Listen(":5000"); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

func connectToMySQL() (*sql.DB, error) {
    dsn := "Thiago:Thiago+123.@tcp(localhost:3306)/todo_go_app"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}


//create

func createTodoHandler(c *fiber.Ctx, db *sql.DB) error {
    todo := new(Todo)

    if err := c.BodyParser(todo); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON data"})
    }

    query := "INSERT INTO todos (task, is_completed) VALUES (?, ?)"
    _, err := db.Exec(query, todo.Task, todo.IsCompleted)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create todo"})
    }

    return c.Status(201).JSON(fiber.Map{"message": "Todo created successfully"})
}


//fetch
func getAllTodosHandler(c *fiber.Ctx, db *sql.DB) error {
    rows, err := db.Query("SELECT id, task, is_completed FROM todos")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todos"})
    }
    defer rows.Close()

    var todos []Todo

    for rows.Next() {
        var todo Todo
        var isCompleted int 

        if err := rows.Scan(&todo.ID, &todo.Task, &isCompleted); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error scanning todos"})
        }

        todo.IsCompleted = isCompleted == 1
        todos = append(todos, todo)
    }

    return c.JSON(todos)
}


//update
func updateTodoHandler(c *fiber.Ctx, db *sql.DB) error {
    id := c.Params("id")
    todo := new(Todo)

    if err := c.BodyParser(todo); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON data"})
    }

    query := "UPDATE todos SET task = ?, is_completed = ? WHERE id = ?"
    _, err := db.Exec(query, todo.Task, todo.IsCompleted, id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
    }

    return c.JSON(fiber.Map{"message": "Todo updated successfully"})
}


//delete
func deleteTodoHandler(c *fiber.Ctx, db *sql.DB) error {
    id := c.Params("id")

    query := "DELETE FROM todos WHERE id = ?"
    _, err := db.Exec(query, id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
    }

    return c.JSON(fiber.Map{"message": "Todo deleted successfully"})
}
