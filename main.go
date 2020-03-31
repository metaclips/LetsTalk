package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/FinalYearProject/controller"
)

func main() {
	router := httprouter.New()

	router.GET("/", controller.HomePage)
	router.GET("/login", controller.LoginPage)

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	os.Setenv("DBHost", "mongodb+srv://admin:admin@cluster0-kcclx.gcp.mongodb.net/test?retryWrites=true&w=majority")
	// fmt.Println(os.Getenv("DBHost"), 1)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}
