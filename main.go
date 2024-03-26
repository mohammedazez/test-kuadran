package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

var db *sql.DB

func initDB() *sql.DB {
	connectionString := "username:password@tcp(127.0.0.1:3306)/database_name"

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")

	return db
}

type Data struct {
	Number []int `json:"number"`
}

func IsPalindrome(str string) string {
	ToLower := strings.ToLower(str)

	for i := 0; i < len(ToLower)/2; i++ {
		if string(ToLower[i]) != string(ToLower[len(ToLower)-1-i]) {
			return "Bukan Palindrome"
		}
	}

	return "Palindrome"
}

func PostBigToSmall(c echo.Context) error {
	u := new(Data)
	if err := c.Bind(u); err != nil {
		return err
	}

	sort.Slice(u.Number, func(i, j int) bool {
		return u.Number[j] < u.Number[i]
	})

	stmt, err := db.Prepare("INSERT INTO sorted_numbers (number) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, num := range u.Number {
		_, err := stmt.Exec(num)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusCreated, u.Number)
}

func PostAverageValue(c echo.Context) error {
	u := new(Data)
	if err := c.Bind(u); err != nil {
		return err
	}

	var total int
	for _, nums := range u.Number {
		total += nums
	}
	Average := total / len(u.Number)

	stmt, err := db.Prepare("INSERT INTO average_values (average) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Average)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, Average)
}

func PostPalindrome(c echo.Context) error {
	u := new(Data)
	if err := c.Bind(u); err != nil {
		return err
	}

	number := u.Number[0]
	numberStr := fmt.Sprintf("%d", number)

	palindromeStatus := IsPalindrome(numberStr)

	stmt, err := db.Prepare("INSERT INTO palindrome_status (number, status) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(number, palindromeStatus)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, palindromeStatus)
}

func GetInputNumbers(c echo.Context) error {
	rows, err := db.Query("SELECT number FROM input_numbers")
	if err != nil {
		return err
	}
	defer rows.Close()

	var numbers []int
	for rows.Next() {
		var number int
		if err := rows.Scan(&number); err != nil {
			return err
		}
		numbers = append(numbers, number)
	}

	return c.JSON(http.StatusOK, numbers)
}

func GetAverageNumber(c echo.Context) error {
	var average int
	err := db.QueryRow("SELECT AVG(number) FROM input_numbers").Scan(&average)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, average)
}

func main() {
	db = initDB()
	defer db.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// post palindrome
	e.POST("/palindrome", PostPalindrome)

	// post input bilangan
	e.POST("/bigtosmall", PostBigToSmall)

	// post rata-rata bilangan
	e.POST("/averagevalue", PostAverageValue)

	// get input bilangan
	e.GET("/inputnumbers", GetInputNumbers)

	// get bilangan rata-rata
	e.GET("/averagenumber", GetAverageNumber)

	e.Logger.Fatal(e.Start(":8080"))
}
