package main

import (
	_ "backend2/docs"
	"backend2/internal/handlers/image"
	httpSwagger "github.com/swaggo/http-swagger"

	//"backend2/internal/auth"
	"backend2/internal/db"
	_ "backend2/internal/dto"
	_ "backend2/internal/handlers/client"
	clienthandler "backend2/internal/handlers/client"
	_ "backend2/internal/handlers/image"
	_ "backend2/internal/handlers/product"
	producthandler "backend2/internal/handlers/product"
	_ "backend2/internal/handlers/supplier"
	suplierhandler "backend2/internal/handlers/supplier"
	"backend2/internal/repository"
	"backend2/internal/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

// @title        Shop API
// @version      1.0
// @description  Документация для API интернет-магазина
// @host         localhost:8080
// @BasePath     /api/v1
func main() {
	database, err := db.Connection()
	if err != nil {
		panic(err)
	}

	//tokenStore := auth.NewInMemoryTokenStore()
	//authUsecase := auth.NewAuthUsecase([]byte("my-secret-key"), tokenStore)
	//authHandler := a.NewAuthHandler(authUsecase)

	repoAdr := repository.NewAddressRepo(database)
	//
	clientRepo := repository.NewClientRepo(database)
	client := usecases.NewClient(clientRepo, repoAdr)
	clientHandler := clienthandler.NewClientHandler(client)
	//
	supplierRepo := repository.NewSupplier(database)
	supplier := usecases.NewSupplier(supplierRepo, repoAdr)
	supplierHandler := suplierhandler.NewSupplierHandler(supplier)
	//
	imgRepo := repository.NewImageRepo(database)
	img := usecases.NewImage(imgRepo)
	imgHandler := image.NewImageHandler(img)
	//
	productRepo := repository.NewProductRepo(database)
	product := usecases.NewProduct(productRepo, supplierRepo, imgRepo)
	productHandler := producthandler.NewProductHandler(product)
	//
	// основной роутер
	router := mux.NewRouter()
	router.StrictSlash(true)

	// открытый маршрут
	//router.HandleFunc("/token", authHandler.GetToken).Methods(http.MethodGet)
	//
	//// защищённые маршруты
	//protected := router.PathPrefix("/").Subrouter()
	//protected.Use(middleware.AuthMiddleware(authUsecase))
	//clients
	router.HandleFunc("/api/v1/clients", clientHandler.GetAllClients).Methods(http.MethodGet)          //+
	router.HandleFunc("/api/v1/client", clientHandler.CreateClient).Methods(http.MethodPost)           //+
	router.HandleFunc("/api/v1/client/{id}", clientHandler.UpdateClient).Methods(http.MethodPatch)     //+
	router.HandleFunc("/api/v1/client", clientHandler.GetClientsByNameSurname).Methods(http.MethodGet) //+
	router.HandleFunc("/api/v1/client/{id}", clientHandler.DeleteClient).Methods(http.MethodDelete)    //+
	//products
	router.HandleFunc("/api/v1/products", productHandler.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product/{id}", productHandler.GetProductById).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product", productHandler.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/product/{id}", productHandler.DeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/product/{id}", productHandler.ReduceProduct).Methods(http.MethodPatch)
	// supplier
	router.HandleFunc("/api/v1/supplier", supplierHandler.CreateSupplier).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/supplier/{id}", supplierHandler.GetSupplierById).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/suppliers", supplierHandler.GetAllSuppliers).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/supplier/{id}", supplierHandler.UpdateAddress).Methods(http.MethodPatch)
	router.HandleFunc("/supplier/{id}", supplierHandler.DeleteSupplierById).Methods(http.MethodDelete)

	//image
	router.HandleFunc("/api/v1/image/{id}", imgHandler.AddImage).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/image/{id}", imgHandler.GetImageById).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/image/{id}", imgHandler.UpdateImage).Methods(http.MethodPatch)
	router.HandleFunc("/api/v1/products/{id}/image", imgHandler.GetProductImageById).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/image/{id}", imgHandler.DeleteImage).Methods(http.MethodDelete)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
