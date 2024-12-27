# Sistem Manajemen Acara REST API

Hallo 👋 berikut Dokumentasi **Event Management System REST API**. API ini memfasilitasi pengelolaan acara, pengguna, tiket, pembayaran, feedback, notifikasi, dashboard, dan peserta. API ini menyediakan endpoint untuk berbagai operasi guna mendukung alur kerja manajemen acara yang mulus.

---

## Daftar Isi

1. [Fitur](#fitur)
2. [Memulai](#memulai)
3. [Endpoint API](#endpoint-api)
   - [Manajemen Pengguna](#manajemen-pengguna)
   - [Manajemen Acara](#manajemen-acara)
   - [Manajemen Tiket](#manajemen-tiket)
   - [Pembayaran](#pembayaran)
   - [Umpan Balik](#umpan-balik)
   - [Notifikasi](#notifikasi)
   - [Dasbor](#dasbor)
   - [Manajemen Peserta](#manajemen-peserta)
4. [Teknologi yang Digunakan](#teknologi-yang-digunakan)
5. [Instruksi Pengaturan](#instruksi-pengaturan)
6. [Lisensi](#lisensi)

---

## Fitur

- **Manajemen Pengguna**: Mengelola pendaftaran, autentikasi, dan profil pengguna.
- **Manajemen Acara**: Membuat, memperbarui, menghapus, dan melihat acara.
- **Manajemen Tiket**: Menerbitkan, melihat, dan mengelola tiket acara.
- **Integrasi Pembayaran**: Pemrosesan dan pelacakan pembayaran yang aman.
- **Sistem Umpan Balik**: Mengumpulkan dan mengelola umpan balik pengguna.
- **Sistem Notifikasi**: Mengirim email atau notifikasi push kepada pengguna.
- **Dasbor**: Melihat analitik dan wawasan acara.
- **Manajemen Peserta**: Mendaftarkan dan mengelola peserta untuk acara.

---

## Memulai

Untuk memulai menggunakan API ini, Anda perlu mengkloning repositori dan mengatur dependensi yang diperlukan. Ikuti [Instruksi Pengaturan](#instruksi-pengaturan) untuk memulai.

---

## Endpoint API

### Manajemen Pengguna

- **POST /users/register**: Mendaftarkan pengguna baru.
- **POST /users/login**: Mengautentikasi pengguna.
- **GET /users/:id**: Mendapatkan detail pengguna.
- **PUT /users/:id**: Memperbarui informasi pengguna.
- **DELETE /users/:id**: Menghapus akun pengguna.

### Manajemen Acara

- **POST /events**: Membuat acara baru.
- **GET /events**: Mendapatkan daftar acara.
- **GET /events/:id**: Mendapatkan detail acara.
- **PUT /events/:id**: Memperbarui acara.
- **DELETE /events/:id**: Menghapus acara.

### Manajemen Tiket

- **POST /tickets**: Menerbitkan tiket untuk acara.
- **GET /tickets**: Melihat semua tiket.
- **GET /tickets/:id**: Melihat detail tiket tertentu.
- **PUT /tickets/:id**: Memperbarui informasi tiket.
- **DELETE /tickets/:id**: Membatalkan tiket.

### Pembayaran

- **POST /payments**: Memulai pembayaran.
- **GET /payments/:id**: Melihat detail pembayaran.
- **GET /payments**: Mendapatkan daftar semua pembayaran.
- **PUT /payments/:id**: Memperbarui status pembayaran.

### Umpan Balik

- **POST /feedback**: Mengirimkan umpan balik untuk acara.
- **GET /feedback**: Mendapatkan umpan balik untuk acara.

### Notifikasi

- **POST /notifications**: Mengirim notifikasi kepada pengguna.
- **GET /notifications**: Mendapatkan daftar notifikasi.

### Dasbor

- **GET /dashboard**: Mendapatkan analitik dan wawasan acara.

### Manajemen Peserta

- **POST /participants**: Mendaftarkan peserta untuk acara.
- **GET /participants/:eventId**: Mendapatkan daftar peserta untuk acara tertentu.
- **DELETE /participants/:id**: Menghapus peserta dari acara.

---

## Teknologi yang Digunakan

- **Framework Backend**: [Golang](https://golang.org/) dengan [Gin](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) atau [MongoDB](https://www.mongodb.com/)
- **Autentikasi**: JSON Web Tokens (JWT)
- **Integrasi Pembayaran**: Stripe atau PayPal API
- **Layanan Notifikasi**: Firebase Cloud Messaging atau penyedia email (misalnya, SendGrid)

---

## Instruksi Pengaturan

1. Clone repositori:

   ```bash
   git clone https://github.com/yourusername/event-management-system.git
   ```

2. Masuk ke direktori proyek:

   ```bash
   cd event-management-system
   ```

3. Instal dependensi:

   ```bash
   go mod tidy
   ```

4. Konfigurasi variabel lingkungan dengan membuat file `.env`:

   ```
   PORT=8888
   DATABASE_URL=your_database_url
   JWT_SECRET=your_secret_key
   PAYMENT_API_KEY=your_payment_api_key
   NOTIFICATION_API_KEY=your_notification_api_key
   ```

5. Jalankan server:

   ```bash
   go run main.go
   ```

6. Akses API di `http://localhost:8888`.
