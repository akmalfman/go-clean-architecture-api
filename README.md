# Go API dengan Clean Architecture & JWT Auth

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT">
</p>

Repositori ini berisi kode backend API yang dibangun menggunakan Go (Golang), PostgreSQL, dan Docker. Proyek ini menerapkan prinsip **Clean Architecture** untuk memisahkan urusan (separation of concerns) dan memastikan kode yang *maintainable* dan *scalable*.

API ini mencakup fungsionalitas CRUD penuh untuk produk dan sistem autentikasi pengguna (Register & Login) yang aman menggunakan **JSON Web Tokens (JWT)**.

## 🚀 Fitur Utama

* **Clean Architecture:** Struktur proyek yang rapi dibagi menjadi lapisan `handler`, `service`, dan `repository`.
* **Full CRUD:** Operasi Create, Read, Update, & Delete untuk data `products`.
* **Autentikasi & Autorisasi:**
    * `POST /register` untuk pendaftaran user baru (password di-hash menggunakan `bcrypt`).
    * `POST /login` untuk autentikasi user dan mendapatkan JWT token.
* **JWT Middleware:** Rute-rute penting (Create, Update, Delete) diamankan dan memerlukan token `Bearer` yang valid.
* **Database:** Menggunakan PostgreSQL yang berjalan di dalam kontainer **Docker Compose** untuk *environment* pengembangan yang konsisten.
* **Modern Go:** Menggunakan `chi` untuk routing dan `pgx/v5` untuk driver database berperforma tinggi.

---

## 🛠️ Tumpukan Teknologi (Tech Stack)

* **Bahasa:** Go (Golang)
* **Database:** PostgreSQL
* **Kontainerisasi:** Docker (Docker Compose)
* **Go Libraries:**
    * `go-chi/chi/v5` (HTTP Router)
    * `jackc/pgx/v5` (PostgreSQL Driver & Toolkit)
    * `golang-jwt/jwt/v5` (Implementasi JWT)
    * `golang.org/x/crypto/bcrypt` (Password Hashing)

---

## 🏁 Cara Menjalankan (Getting Started)

Untuk menjalankan proyek ini di lokal, ikuti langkah-langkah berikut:

### Prasyarat

* [Go](https://go.dev/doc/install) (versi 1.18+)
* [Docker](https://www.docker.com/products/docker-desktop/)

### Instalasi

1.  **Clone repositori ini:**
    ```bash
    # Ganti [NAMA_USER_KAMU]/[NAMA_REPO_KAMU]
    git clone [https://github.com/](https://github.com/)[NAMA_USER_KAMU]/[NAMA_REPO_KAMU].git
    ```

2.  **Masuk ke direktori proyek:**
    ```bash
    cd [NAMA_REPO_KAMU]
    ```

3.  **Jalankan Database (Docker):**
    Pastikan Docker Desktop sudah berjalan. Lalu, jalankan perintah ini untuk membuat dan menjalankan kontainer PostgreSQL di *background*.
    ```bash
    docker-compose up -d
    ```

4.  **Install dependensi Go:**
    Perintah ini akan men-download semua *library* yang dibutuhkan (chi, pgx, jwt, dll.)
    ```bash
    go mod tidy
    ```

5.  **Jalankan Server Go:**
    ```bash
    go run main.go
    ```

✅ Server sekarang berjalan di `http://localhost:8080`.

---

## 🗺️ API Endpoints

Anda bisa menggunakan [Postman](https://www.postman.com/) atau *API Client* lain untuk menguji *endpoint* berikut.

### 👤 Auth (Autentikasi)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/register` | Mendaftarkan user baru. (Body: `email`, `password`) |
| `POST` | `/login` | Login user. (Body: `email`, `password`). Mengembalikan JWT. |

### 📦 Products (Produk)

| Method | Endpoint | Deskripsi | Auth? |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Menampilkan pesan selamat datang. | ❌ Tidak |
| `GET` | `/products` | Mendapatkan semua data produk. | ❌ Tidak |
| `POST` | `/products` | Membuat produk baru. (Body: `name`, `price`) | ✅ **Ya** (Bearer) |
| `PUT` | `/products/{id}` | Memperbarui produk berdasarkan ID. (Body: `name`, `price`) | ✅ **Ya** (Bearer) |
| `DELETE`| `/products/{id}` | Menghapus produk berdasarkan ID. | ✅ **Ya** (Bearer) |

---
