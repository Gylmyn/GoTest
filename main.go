package main

import (
	"database/sql"
	"os"

	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int       `json:"price"`
}

func main() {
	connStr := "postgresql://Gylmyn:kLUA2E6lYWam@ep-delicate-pond-24139553.ap-southeast-1.aws.neon.tech/mobile_db?sslmode=require"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", welkam)
	e.GET("/data", getAllItems)
	e.POST("/mytable/add", addItem)
	e.DELETE("/mytable/delete/:id", deleteItem)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func welkam(c echo.Context) error {
	welkam := fmt.Sprintln("Welkam Walaw We \n 1. /data")
	return c.String(http.StatusOK, welkam)

}

func getAllItems(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM myitems")
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return err
		}
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, items)
}

func deleteItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID parameter is required")
	}

	_, err := db.Exec("DELETE FROM myitems WHERE id = $1", id)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("Item with ID %s deleted successfully\n", id))
}

func addItem(c echo.Context) error {
	var item Item
	if err := c.Bind(&item); err != nil {
		return err
	}

	id := uuid.New()

	_, err := db.Exec("INSERT INTO myitems(id, name, price) VALUES($1, $2, $3)", id, item.Name, item.Price)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, item)
}
