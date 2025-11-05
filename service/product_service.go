package service

import (
	"errors"
	"first-project/models"
	"first-project/repository"
)

// ProductService adalah interface untuk logika bisnis produk
type ProductService interface {
	GetAllProducts() ([]models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(id int, product models.Product) (models.Product, error)
	DeleteProduct(id int) error
}

// productService adalah implementasi dari ProductService
type productService struct {
	repo repository.ProductRepository
}

// NewProductService membuat instance baru dari productService
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	// Di sini bisa ada logika lain, misal: caching
	return s.repo.FindAll()
}

func (s *productService) CreateProduct(product models.Product) (models.Product, error) {
	// Ini adalah Logika Bisnis (Validasi)
	if product.Name == "" || product.Price <= 0 {
		return models.Product{}, errors.New("nama dan harga tidak boleh kosong/nol")
	}

	newID, err := s.repo.Save(product)
	if err != nil {
		return models.Product{}, err
	}
	product.ID = newID
	return product, nil
}

func (s *productService) UpdateProduct(id int, product models.Product) (models.Product, error) {
	if product.Name == "" || product.Price <= 0 {
		return models.Product{}, errors.New("nama dan harga tidak boleh kosong/nol")
	}

	updatedID, err := s.repo.Update(id, product)
	if err != nil {
		return models.Product{}, err
	}
	product.ID = updatedID
	return product, nil
}

func (s *productService) DeleteProduct(id int) error {
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}
