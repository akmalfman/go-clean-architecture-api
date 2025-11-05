package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"first-project/handler"
	"first-project/middleware"
	"first-project/repository"
	"first-project/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connString := "postgresql://root:root@localhost:5432/db_latihan"
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Gagal konek ke database (pool): %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal ping database: %v\n", err)
	}
	fmt.Println("✅ Sukses terhubung ke PostgreSQL (via Pool)!")

	productRepo := repository.NewProductRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)

	productRepo.SetupTable()
	userRepo.SetupUserTable()
	fmt.Println("✅ Tabel 'products' siap.")

	productService := service.NewProductService(productRepo)
	authService := service.NewAuthService(userRepo)

	productHandler := handler.NewProductHandler(productService)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	authHandler.RegisterRoutes(r)

	r.Get("/", productHandler.HandleHome)
	r.Get("/products", productHandler.HandleGetProducts)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)

		r.Post("/products", productHandler.HandleCreateProduct)
		r.Put("/products/{id}", productHandler.HandleUpdateProduct)
		r.Delete("/products/{id}", productHandler.HandleDeleteProduct)
	})

	port := ":8080"
	log.Printf("🚀 Server berjalan di http://localhost%s\n", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Server gagal berjalan: %v\n", err)
	}
}
