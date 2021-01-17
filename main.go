package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/runRandom"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	var userID string
	var userSolvedProblems []int
	problemsToSolve := runRandom.GetPage()

	router.POST("/", func(c *gin.Context) {
		userID = c.PostForm("account")
		var address string
		var message string
		var title string
		var account string
		if len(userID) == 0 {
			address = ""
			message = "Enter your ID!"
			title = ""
			account = "guest"
		} else {
			userSolvedProblems = runRandom.GetUserSolvedProblemInfo(userID)
			if userSolvedProblems[0] == -1 {
				address = ""
				message = "The user does not exist!"
				title = ""
				account = "guest"
			} else {
				problemTitle, url := runRandom.RunRandom(problemsToSolve, userSolvedProblems)
				address = url
				message = ""
				title = problemTitle + " : "
				account = userID
			}
		}
		runRandom.PrintUserID(userID)
		c.HTML(http.StatusOK, "main.tmpl.html", gin.H{
			"Address": address,
			"Message": message,
			"Title":   title,
			"Account": account,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.tmpl.html", gin.H{
			"Address": "",
			"Message": "Enter your ID!",
			"Title":   "",
			"Account": "guest",
		})
	})

	router.Run(":" + port)
}
