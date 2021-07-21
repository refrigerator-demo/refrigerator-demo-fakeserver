package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"fridge/src/middleware"
	"fridge/src/model"
)

type RestServer struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *RestServer) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if nil != err {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.Debug().DropTableIfExists(&model.User{}, &model.Inventory{}, &model.ItemType{}, &model.Item{})

	server.DB.Debug().AutoMigrate(&model.User{})
	server.DB.Debug().AutoMigrate(&model.ItemType{})
	server.DB.Debug().AutoMigrate(&model.Inventory{})
	server.DB.Debug().AutoMigrate(&model.Item{})

	server.DB.Debug().Model(&model.Inventory{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	server.DB.Debug().Model(&model.Item{}).AddForeignKey("type_id", "item_types(id)", "RESTRICT", "RESTRICT")
	server.DB.Debug().Model(&model.Item{}).AddForeignKey("inventory_id", "inventories(id)", "RESTRICT", "RESTRICT")

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *RestServer) initializeRoutes() {

	// Home Route
	server.Router.HandleFunc("/api", middleware.SetMiddlewareJSON(server.HealthCheck)).Methods("GET")

	// Login Route
	server.Router.HandleFunc("/api/login", middleware.SetMiddlewareJSON(server.Login)).Methods("POST")
	server.Router.HandleFunc("/api/login", server.AllowOrigin).Methods("OPTIONS")
	//Users routes
	server.Router.HandleFunc("/api/users", middleware.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/api/users", server.AllowOrigin).Methods("OPTIONS")
	//s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/api/users/{id}", middleware.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/api/users/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	//s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	server.Router.HandleFunc("/api/inventory", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(server.GetInventoryByID))).Methods("GET")
	server.Router.HandleFunc("/api/inventory", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(server.CreateInventory))).Methods("POST")
	server.Router.HandleFunc("/api/inventory", server.AllowOrigin).Methods("OPTIONS")
}

func (server *RestServer) Run(addr string) {

	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *RestServer) RunRESTServer() {

	var err error
	err = godotenv.Load()
	if nil != err {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":8000")
}
