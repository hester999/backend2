package main

import (
	"backend2/internal/handlers/image"
	//"backend2/internal/auth"
	"backend2/internal/db"
	clienthandler "backend2/internal/handlers/client"
	producthandler "backend2/internal/handlers/product"
	suplierhandler "backend2/internal/handlers/supplier"
	"backend2/internal/repository"
	"backend2/internal/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

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
	router.HandleFunc("/clients", clientHandler.GetAllClients).Methods(http.MethodGet)
	router.HandleFunc("/client", clientHandler.CreateClient).Methods(http.MethodPost)
	router.HandleFunc("/client/{id}", clientHandler.UpdateClient).Methods(http.MethodPatch)
	router.HandleFunc("/client", clientHandler.GetClientByNameSurname).Methods(http.MethodGet)
	router.HandleFunc("/client/{id}", clientHandler.DeleteClient).Methods(http.MethodDelete)
	//products
	router.HandleFunc("/products", productHandler.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/product/{id}", productHandler.GetProductById).Methods(http.MethodGet)
	router.HandleFunc("/product", productHandler.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/product/{id}", productHandler.DeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/product/{id}", productHandler.ReduceProduct).Methods(http.MethodPatch)
	// supplier
	router.HandleFunc("/supplier", supplierHandler.CreateSupplier).Methods(http.MethodPost)
	router.HandleFunc("/supplier/{id}", supplierHandler.GetSupplierById).Methods(http.MethodGet)
	router.HandleFunc("/suppliers", supplierHandler.GetAllSuppliers).Methods(http.MethodGet)
	router.HandleFunc("/supplier/{id}", supplierHandler.UpdateAddress).Methods(http.MethodPatch)
	router.HandleFunc("/supplier/{id}", supplierHandler.DeleteSupplierById).Methods(http.MethodDelete)

	//image
	router.HandleFunc("/image", imgHandler.AddImage).Methods(http.MethodPost)
	router.HandleFunc("/image/{id}", imgHandler.GetImageById).Methods(http.MethodGet)
	router.HandleFunc("/image/{id}", imgHandler.UpdateImage).Methods(http.MethodPatch)
	router.HandleFunc("/image/{id}", imgHandler.GetProductImageById).Methods(http.MethodGet)
	router.HandleFunc("/image/{id}", imgHandler.DeleteImage).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
