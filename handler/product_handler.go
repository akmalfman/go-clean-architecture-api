package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"first-project/models"
	"first-project/service"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Selamat datang di API Clean Architecture!"})
}

func (h *ProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	h.respondWithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	newProduct, err := h.service.CreateProduct(p)
	if err != nil {
		if err.Error() == "nama dan harga tidak boleh kosong/nol" {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Gagal membuat produk")
		}
		return
	}
	h.respondWithJSON(w, http.StatusCreated, newProduct)
}

func (h *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "ID produk tidak valid")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	updatedProduct, err := h.service.UpdateProduct(id, p)
	if err != nil {
		if err.Error() == "no rows in result set" {
			h.respondWithError(w, http.StatusNotFound, "Produk tidak ditemukan")
		} else if err.Error() == "nama dan harga tidak boleh kosong/nol" {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Gagal update produk")
		}
		return
	}
	h.respondWithJSON(w, http.StatusOK, updatedProduct)
}

func (h *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "ID produk tidak valid")
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		if err.Error() == "produk tidak ditemukan" {
			h.respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Gagal delete produk")
		}
		return
	}
	h.respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandler) parseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
