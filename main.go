package main

import (
	cs "go-money-tracker/controllers"
	"go-money-tracker/initializers"
	mw "go-money-tracker/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5173/", "https://vue-money-tracker-tt7q.vercel.app", "https://vue-money-tracker-tt7q.vercel.app/", "https://vue-money-tracker-tt7q-4iyn9q01h-wisnudwgs-projects.vercel.app", "https://vue-money-tracker-tt7q-4iyn9q01h-wisnudwgs-projects.vercel.app/"},
		// AllowOrigins:     []string{os.Getenv("FE_ORIGIN")},
		AllowMethods:     []string{"GET", "DELETE", "OPTIONS", "POST", "PATCH", "PUT"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		//   return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "money-tracker-backend",
		})
	})

	// USER-RELATED ROUTES
	r.POST("/register", cs.Register)
	r.POST("/login", cs.Login)
	r.GET("/validate-token", mw.RequireAuth, cs.ValidateToken)
	r.PUT("/update-user", mw.RequireAuth, cs.UpdateUser)
	r.DELETE("/delete-user/:uid", cs.DeleteUser)
	r.GET("/get-user/:uid", cs.ReadUser)

	// ENTRY-RELATED ROUTES
	r.POST("/create-entry", mw.RequireAuth, cs.CreateEntry)
	r.GET("/get-entry/:eid", mw.RequireAuth, cs.GetEntry)
	r.POST("/get-entries", mw.RequireAuth, cs.GetEntries)
	r.PUT("/update-entry/:eid", mw.RequireAuth, cs.UpdateEntry)
	r.DELETE("/delete-entry/:eid", mw.RequireAuth, cs.DeleteEntry)
	r.GET("/get-notes/:uid", mw.RequireAuth, cs.GetNotes)
	r.GET("/get-assets/:uid", mw.RequireAuth, cs.GetAssets)
	r.GET("/get-income-categories/:uid", mw.RequireAuth, cs.GetIncomeCategories)
	r.GET("/get-expense-categories/:uid", mw.RequireAuth, cs.GetExpenseCategories)

	r.Run()
}
