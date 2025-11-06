package repository

import (
	"context"
	"log"

	"github.com/akmalfsalman/go-clean-architecture-api/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	SetupTable()
	FindAll() ([]models.Product, error)
	FindByID(id int) (models.Product, error)
	Save(product models.Product) (int, error)
	Update(id int, product models.Product) (int, error)
	Delete(id int) (int64, error)
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) SetupTable() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price INT
	);`
	_, err := r.db.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Gagal membuat tabel: %v\n", err)
	}
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	query := "SELECT id, name, price FROM products ORDER BY id ASC"
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *productRepository) Save(product models.Product) (int, error) {
	query := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`
	var newID int
	err := r.db.QueryRow(
		context.Background(),
		query,
		product.Name,
		product.Price,
	).Scan(&newID)
	return newID, err
}

func (r *productRepository) Update(id int, product models.Product) (int, error) {
	query := `UPDATE products SET name = $1, price = $2 WHERE id = $3 RETURNING id`
	var updatedID int
	err := r.db.QueryRow(
		context.Background(),
		query,
		product.Name,
		product.Price,
		id,
	).Scan(&updatedID)
	return updatedID, err
}

func (r *productRepository) Delete(id int) (int64, error) {
	query := `DELETE FROM products WHERE id = $1`
	ct, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

func (r *productRepository) FindByID(id int) (models.Product, error) {
	query := "SELECT id, name, price FROM products WHERE id = $1"
	var p models.Product
	err := r.db.QueryRow(context.Background(), query, id).Scan(&p.ID, &p.Name, &p.Price)
	return p, err // Akan error jika 'no rows'
}
