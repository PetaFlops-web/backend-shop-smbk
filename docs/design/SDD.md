# System Design Document — Backend

## Smart Commerce

_Desain Teknis Backend: Modular Monolith (Golang)_

| Field             | Detail                                                      |
| ----------------- | ----------------------------------------------------------- |
| Nama Proyek       | Smart Commerce                                              |
| Kompetisi         | AI Innovation Challenge (AIC)                               |
| Dokumen Referensi | Smart Commerce — PRD v0.1                                   |
| Penulis           | Putra Rizky                                                 |
| Scope Dokumen     | Backend only — arsitektur, data, API, business logic, infra |
| Status Dokumen    | Draft — modul FE & ML masih ditentukan tim masing-masing    |
| Versi             | 0.1                                                         |

---

## 1. Overview

### 1.1 Tujuan Dokumen

Dokumen ini menjelaskan desain teknis backend Smart Commerce: bagaimana sistem backend dibangun untuk memenuhi requirement pada PRD. Dokumen ini fokus pada bagian yang menjadi tanggung jawab penulis, yaitu Backend (BE). Bagian Frontend (FE) dan Machine Learning (ML) belum didetailkan karena tim untuk kedua area tersebut belum ditentukan — backend didesain agar dapat berjalan mandiri dan berintegrasi dengan FE/ML melalui kontrak API yang jelas begitu tim tersebut tersedia.

### 1.2 Scope

- **Termasuk:** arsitektur backend, desain data, desain API, business logic, integrasi dengan layanan ML (sebagai kontrak/interface, bukan implementasi model), cron job, keamanan, konfigurasi, logging, dan deployment backend.
- **Tidak termasuk:** implementasi UI/FE, implementasi/training model AI (Whisper, NLP extraction, model restock) — backend hanya mendefinisikan kontrak integrasi ke layanan tersebut.

### 1.3 Referensi

- Smart Commerce — Product Requirements Document (PRD), v0.1, 08 Juli 2026.
- Template arsitektur acuan: golang-modular-monolith (arttVinci) — https://github.com/arttVinci/golang-modular-monolith

---

## 2. Arsitektur Sistem

### 2.1 Pola Arsitektur: Modular Monolith

Backend dibangun dengan pola **Modular Monolith** — satu aplikasi/deployment unit, namun kode dipecah menjadi modul-modul yang terisolasi secara data dan hanya saling berkomunikasi lewat interface (client) yang eksplisit. Pola ini dipilih karena tim backend masih kecil (solo/awal), sehingga kompleksitas operasional microservices belum diperlukan, tapi struktur modular tetap memudahkan pemisahan tanggung jawab dan migrasi ke microservices di masa depan bila dibutuhkan.

### 2.2 Prinsip Utama

- **Data Isolation:** setiap modul punya tabel sendiri. Tidak ada foreign key GORM lintas modul — referensi antar modul hanya berupa ID (string/UUID) biasa, bukan FK database.
- **Inter-module Communication via Client Interface:** satu modul mengakses data modul lain hanya lewat interface publik (`*-client`) milik modul tersebut, tidak melakukan query/JOIN langsung ke tabel modul lain.
- **Independent Data Ownership:** setiap modul bertanggung jawab atas `AutoMigrate`/migration tabelnya sendiri.
- Setiap modul mengimplementasikan interface modul generik (`Migrate()` dan `RegisterRoutes()`) sehingga wiring di `main.go` seragam untuk semua modul.

### 2.3 Diagram Modul (High-Level)

```
                         +------------------------+
                         |      Fiber HTTP        |
                         |   (cmd/api/main.go)    |
                         +-----------+------------+
                                     |
        +---------------+------------+------------+---------------+
        |               |            |             |               |
   +---------+   +-----------+  +-----------+ +-----------+  +-----------+
   |  auth   |   |   store   |  |  product  | |transaction| |  report   |
   +---------+   +-----------+  +-----------+ +-----------+  +-----------+
        |               |            |             |               |
   +-----------+  +-----------+ +-----------+
   | restock   |  | promotion |  |ml-client  |  (interface ke
   +-----------+  +-----------+  +-----------+   layanan ML eksternal:
                                                  Whisper + NLP extraction)

  Semua modul mengakses modul lain HANYA lewat <module>-client interface,
  tidak lewat query langsung ke tabel modul lain.
```

_Diagram di atas menunjukkan modul beserta arah dependency utama, bukan diagram jaringan/infrastruktur._

### 2.4 Struktur Proyek (Project Structure)

```
smart-commerce-backend/
├── cmd/
│   └── api/
│       └── main.go              # entrypoint, wiring semua modul
├── config.json                  # konfigurasi (dibaca oleh Viper)
├── internal/
│   ├── modules/
│   │   ├── auth/                # Autentikasi & data user
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   └── module.go
│   │   ├── auth-client/         # kontrak publik modul auth
│   │   ├── store/               # Manajemen toko
│   │   ├── store-client/
│   │   ├── product/             # Katalog produk & stok
│   │   ├── product-client/
│   │   ├── transaction/         # Transaksi & item transaksi
│   │   ├── transaction-client/
│   │   ├── report/              # Laporan harian (read-only, on-the-fly)
│   │   ├── restock/             # Prediksi restock (cron)
│   │   └── promotion/           # Promosi otomatis (cron)
│   ├── jobs/
│   │   ├── restock_job.go       # cron: prediksi restock tiap malam
│   │   └── promotion_job.go     # cron: cek rule promo tiap malam
│   ├── pkg/
│   │   └── mlclient/            # HTTP client ke layanan ML (Whisper/NLP/OCR)
│   └── shared/
│       ├── config/              # setup GORM, Viper, Fiber, Logrus
│       ├── middleware/          # JWT auth middleware, request logging
│       ├── module/              # interface Module generik
│       ├── repository/          # generic Repository[T] base
│       └── response/            # WebResponse[T] standar
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

### 2.5 Menambah Modul Baru

1. Buat folder `internal/modules/<nama>/` dan `internal/modules/<nama>-client/`.
2. Definisikan interface publik di `<nama>-client/client.go` — ini yang dipakai modul lain untuk berkomunikasi.
3. Implementasikan `module.go` yang memenuhi interface `module.Module`: `Migrate() error` dan `RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler)`.
4. Daftarkan modul tersebut di `cmd/api/main.go`.

---

## 3. Pemecahan Modul (Module Breakdown)

| Modul       | Tanggung Jawab                                                                                     | Tabel yang Dimiliki                                   |
| ----------- | -------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| auth        | Registrasi/login pedagang via OAuth, penerbitan & validasi JWT, profil user                        | users                                                 |
| store       | Pembuatan & kepemilikan toko                                                                       | stores                                                |
| product     | Master data produk, harga modal/jual, stok, decrement stok                                         | products                                              |
| transaction | Ekstraksi preview transaksi (proxy ke ML), penyimpanan transaksi & item transaksi (snapshot harga) | transactions, transaction_items                       |
| report      | Agregasi laporan harian (omset, untung, produk terlaris) — read-only, dihitung on-the-fly          | — (tanpa tabel sendiri, query via transaction-client) |
| restock     | Cron job prediksi kebutuhan restock berdasarkan histori penjualan                                  | restock_predictions                                   |
| promotion   | Cron job cek rule promo berbasis produk (if-else)                                                  | promotion_logs                                        |

---

## 4. Tech Stack

| Komponen         | Pilihan                               | Fungsi                                                                                    |
| ---------------- | ------------------------------------- | ----------------------------------------------------------------------------------------- |
| Bahasa           | Go (Golang)                           | Bahasa utama backend                                                                      |
| HTTP Framework   | Fiber                                 | Routing & HTTP server                                                                     |
| ORM              | GORM                                  | Mapping struct ke database, migration                                                     |
| Validasi Input   | go-playground/validator               | Validasi request DTO (struct tag based)                                                   |
| Konfigurasi      | Viper                                 | Load config.json / environment variables                                                  |
| Logging          | Logrus                                | Structured logging                                                                        |
| Containerization | Docker                                | Packaging & deployment                                                                    |
| Database         | MySQL/PostgreSQL (via GORM dialector) | Penyimpanan data relasional — dialector dapat disesuaikan; contoh referensi memakai MySQL |
| Scheduler        | robfig/cron (atau setara)             | Menjalankan cron job restock & promosi tiap malam                                         |

---

## 5. Format Response API (Standar)

Seluruh endpoint backend mengembalikan response dengan format seragam melalui generic `WebResponse[T]`:

```go
type WebResponse[T any] struct {
	Data    T             `json:"data"`
	Message string        `json:"message,omitempty"`
	Success bool          `json:"success,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
}
```

Contoh response sukses (single object):

```json
{
  "data": { "id": "prod_123", "nama_produk": "Beras 5kg" },
  "message": "Produk berhasil dibuat",
  "success": true
}
```

Contoh response list dengan paging:

```json
{
  "data": [{ "id": "prod_123" }, { "id": "prod_124" }],
  "success": true,
  "paging": { "page": 1, "page_size": 20, "total_data": 42, "total_page": 3 }
}
```

Contoh response error:

```json
{
  "data": null,
  "message": "Stok produk tidak mencukupi",
  "success": false
}
```

---

## 6. Desain Data

> Mengikuti prinsip data isolation: kolom yang berakhiran `_id` lintas modul (mis. `store_id` di tabel `products`) adalah **plain reference ID**, BUKAN foreign key database. Validasi keberadaan ID tersebut dilakukan lewat pemanggilan client interface modul pemilik, bukan lewat constraint FK.

### 6.1 Modul `auth` — tabel `users`

| Field                   | Tipe          | Keterangan  |
| ----------------------- | ------------- | ----------- |
| id                      | string (UUID) | Primary key |
| username                | string        | Username    |
| email                   | string        | Unique      |
| created_at / updated_at | timestamp     | Audit field |

### 6.2 Modul `store` — tabel `stores`

| Field                   | Tipe          | Keterangan                                 |
| ----------------------- | ------------- | ------------------------------------------ |
| id                      | string (UUID) | Primary key                                |
| owner_id                | string        | Reference ke users.id (plain ID, bukan FK) |
| store_name              | string        | Nama toko                                  |
| created_at / updated_at | timestamp     | Audit field                                |

### 6.3 Modul `product` — tabel `products`

| Field                   | Tipe          | Keterangan                                   |
| ----------------------- | ------------- | -------------------------------------------- |
| id                      | string (UUID) | Primary key                                  |
| store_id                | string        | Reference ke stores.id (plain ID)            |
| product_name            | string        | Nama produk                                  |
| cost_price              | bigint        | Harga modal saat ini (rupiah, tanpa desimal) |
| selling_price           | bigint        | Harga jual saat ini                          |
| stock                   | int           | Stok saat ini, di-decrement tiap transaksi   |
| unit                    | string        | mis. kg, pcs, liter                          |
| created_at / updated_at | timestamp     | Audit field                                  |

### 6.4 Modul `transaction` — tabel `transactions`

| Field            | Tipe          | Keterangan                                     |
| ---------------- | ------------- | ---------------------------------------------- |
| id               | string (UUID) | Primary key                                    |
| store_id         | string        | Reference ke stores.id (plain ID)              |
| transaction_date | date          | Tanggal transaksi (untuk query laporan harian) |
| source           | string (enum) | voice \| photo                                 |
| created_at       | timestamp     | Audit field                                    |

### 6.5 Modul `transaction` — tabel `transaction_items`

| Field                  | Tipe          | Keterangan                                                             |
| ---------------------- | ------------- | ---------------------------------------------------------------------- |
| id                     | string (UUID) | Primary key                                                            |
| transaction_id         | string        | FK ke transactions.id — valid karena satu modul yang sama              |
| product_id             | string        | Reference ke products.id (plain ID, lintas modul)                      |
| product_name_snapshot  | string        | Nama produk saat transaksi, untuk jaga-jaga bila produk diubah/dihapus |
| qty                    | int           | Jumlah unit terjual                                                    |
| cost_price_snapshot    | bigint        | Snapshot harga modal saat transaksi — dasar hitung untung              |
| selling_price_snapshot | bigint        | Snapshot harga jual saat transaksi                                     |

### 6.6 Modul `restock` — tabel `restock_predictions`

| Field                  | Tipe          | Keterangan                                                             |
| ---------------------- | ------------- | ---------------------------------------------------------------------- |
| id                     | string (UUID) | Primary key                                                            |
| store_id               | string        | Reference ke stores.id                                                 |
| product_id             | string        | Reference ke products.id                                               |
| predicted_restock_date | date          | Estimasi tanggal produk perlu di-restock                               |
| predicted_qty          | int           | Estimasi kuantitas restock                                             |
| avg_daily_sold         | decimal       | Rata-rata penjualan harian yang dipakai sebagai basis prediksi (audit) |
| created_at             | timestamp     | Kapan prediksi ini digenerate (per-run cron)                           |

### 6.7 Modul `promotion` — tabel `promotion_logs`

| Field          | Tipe          | Keterangan                                                                  |
| -------------- | ------------- | --------------------------------------------------------------------------- |
| id             | string (UUID) | Primary key                                                                 |
| store_id       | string        | Reference ke stores.id                                                      |
| product_id     | string        | Reference ke products.id                                                    |
| rule_triggered | string        | Nama/kode rule yang terpenuhi, mis. QTY_MONTHLY_THRESHOLD                   |
| period         | string        | Periode evaluasi, mis. "2026-07" — mencegah duplikat log dalam periode sama |
| created_at     | timestamp     | Kapan log ini dibuat                                                        |

---

## 7. Komunikasi Antar Modul (Client Interfaces)

Setiap modul mengekspos interface publik yang dipakai modul lain, tanpa membuka akses langsung ke repository/tabel internalnya.

| Client Interface   | Method Utama                                                                                          | Dipakai Oleh                                     |
| ------------------ | ----------------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| auth-client        | GetUserByID(id), ValidateToken(token)                                                                 | middleware JWT, store                            |
| store-client       | GetStoreByOwnerID(ownerID), GetStoreByID(id)                                                          | product, transaction, report, restock, promotion |
| product-client     | GetByID(id), ListByStoreID(storeID), DecrementStock(id, qty), Search(storeID, keyword) untuk matching | transaction, restock, promotion, report          |
| transaction-client | CreateTransaction(input), ListByStoreAndDate(storeID, date), ListItemsByStoreAndDateRange(...)        | report, restock, promotion                       |

> Contoh: modul `report` **tidak** melakukan JOIN SQL ke tabel `products`, melainkan memanggil `transaction-client` untuk data transaksi & item, lalu `product-client` bila perlu detail produk tambahan (mis. sisa stok terkini).

---

## 8. Desain API

### 8.1 Modul `auth`

| Method | Endpoint             | Deskripsi                           | Auth |
| ------ | -------------------- | ----------------------------------- | ---- |
| POST   | /api/users           | Register                            | -    |
| POST   | /api/users/\_login   | Login                               | -    |
| GET    | /api/users/\_current | Ambil profil user yang sedang login | JWT  |

### 8.2 Modul `store`

| Method | Endpoint    | Deskripsi                               | Auth |
| ------ | ----------- | --------------------------------------- | ---- |
| POST   | /api/stores | Membuat toko baru untuk user yang login | JWT  |
| GET    | /api/stores | Ambil data toko milik user yang login   | JWT  |

### 8.3 Modul `product`

| Method | Endpoint          | Deskripsi                          | Auth |
| ------ | ----------------- | ---------------------------------- | ---- |
| POST   | /api/products     | Tambah produk baru (FR-03)         | JWT  |
| GET    | /api/products     | List produk milik toko (paginated) | JWT  |
| GET    | /api/products/:id | Detail satu produk                 | JWT  |
| PUT    | /api/products/:id | Update produk (harga, stok, dll)   | JWT  |
| DELETE | /api/products/:id | Hapus produk                       | JWT  |

### 8.4 Modul `transaction`

| Method | Endpoint                        | Deskripsi                                                                                                           | Auth |
| ------ | ------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ---- |
| POST   | /api/transactions/extract/voice | Upload audio → proxy ke layanan ML (Whisper+NLP) → matching produk → return preview (belum tersimpan) (FR-06–FR-10) | JWT  |
| POST   | /api/transactions/extract/photo | Upload foto nota → proxy ke layanan OCR+NLP → matching produk → return preview (FR-12–FR-15)                        | JWT  |
| POST   | /api/transactions               | Konfirmasi & simpan transaksi + item + decrement stok (FR-11, FR-05)                                                | JWT  |
| GET    | /api/transactions               | List riwayat transaksi toko (paginated, untuk kebutuhan debugging/QA)                                               | JWT  |

### 8.5 Modul `report`

| Method | Endpoint                           | Deskripsi                                                                                                                       | Auth |
| ------ | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ---- |
| GET    | /api/reports/daily?date=YYYY-MM-DD | Laporan harian: total omset, total untung, jumlah transaksi, produk terlaris, sisa stok (FR-16–FR-20). Default date = hari ini. | JWT  |

### 8.6 Modul `restock` & `promotion` (read-only untuk FE, proses utama via cron)

| Method | Endpoint                 | Deskripsi                                              | Auth |
| ------ | ------------------------ | ------------------------------------------------------ | ---- |
| GET    | /api/restock-predictions | List prediksi restock terbaru untuk toko (FR-21–FR-22) | JWT  |
| GET    | /api/promotions          | List promotion_logs terbaru untuk toko (FR-23–FR-24)   | JWT  |

---

## 9. Business Logic & Algoritma

### 9.1 Flow Ekstraksi Transaksi (Suara/Foto) — Preview

```
POST /transactions/extract/voice atau /transactions/extract/photo
1. Terima file (audio/foto) dari request
2. Kirim file ke ml-client (interface ke layanan ML eksternal)
     -> voice: ml-client.TranscribeAndExtract(audio)
     -> photo: ml-client.OcrAndExtract(photo)
3. Terima structured output dari ML: [{item_text, qty, harga}, ...]
4. Untuk tiap item_text, panggil product-client.Search(storeID, item_text)
     -> dapatkan kandidat product_id + skor kemiripan (fuzzy match)
5. Susun preview response:
   [{ raw_text, matched_product_id, matched_product_name,
      qty, harga_terdeteksi, match_confidence, is_editable: true }, ...]
6. Return preview via WebResponse[T] -- BELUM ditulis ke database
```

> Preview tidak dipersist. Backend hanya menyimpan hasil setelah pedagang mengonfirmasi lewat `POST /transactions` (FR-10, FR-15).

### 9.2 Flow Konfirmasi Transaksi (Persist)

```
POST /transactions
Input: { store_id, items: [{ product_id, qty, harga_jual_final }] }
1. Validasi input (validator): store_id valid milik user login, items tidak kosong
2. BEGIN TRANSACTION (DB transaction, bukan business transaction)
3. Untuk setiap item:
     a. product := product-client.GetByID(item.product_id)
     b. Jika product.stok < item.qty -> ROLLBACK, return error "stok tidak cukup"
     c. product-client.DecrementStock(item.product_id, item.qty)
4. Buat 1 row transactions { store_id, transaction_date: today, source }
5. Untuk setiap item, buat row transaction_items {
     transaction_id, product_id, product_name_snapshot: product.nama_produk,
     qty: item.qty,
     harga_modal_saat_transaksi: product.harga_modal,   // snapshot, FR-17
     harga_jual_saat_transaksi: item.harga_jual_final,  // snapshot, FR-17
   }
6. COMMIT TRANSACTION
7. Return WebResponse sukses berisi transaction_id
```

### 9.3 Perhitungan Laporan Harian (On-the-fly)

```
GET /reports/daily?date=YYYY-MM-DD
1. items := transaction-client.ListItemsByStoreAndDate(storeID, date)
2. total_omset   := SUM(items.qty * items.harga_jual_saat_transaksi)
3. total_untung  := SUM(items.qty * (items.harga_jual_saat_transaksi
                                    - items.harga_modal_saat_transaksi))
4. jumlah_transaksi := COUNT(DISTINCT items.transaction_id)
5. produk_terlaris  := GROUP BY product_id, SUM(qty) DESC LIMIT 5
6. sisa_stok        := product-client.ListByStoreID(storeID)  // stok terkini
7. Return WebResponse { total_omset, total_untung, jumlah_transaksi,
                        produk_terlaris, sisa_stok }
```

> Poin penting (FR-17): perhitungan untung **selalu** memakai `harga_modal_saat_transaksi` & `harga_jual_saat_transaksi` dari `transaction_items`, bukan `products.harga_modal`/`harga_jual` saat ini — supaya laporan historis tidak berubah saat harga produk di-update di kemudian hari.

### 9.4 Cron Job — Prediksi Restock (Malam Hari)

```
Jadwal: setiap malam (mis. 23:00), berjalan per store
Untuk tiap store:
  Untuk tiap product milik store:
    1. history := transaction-client.ListItemsByProduct(
                     productID, lookback: 30 hari terakhir)
    2. avg_daily_sold := SUM(history.qty) / 30
    3. Jika avg_daily_sold == 0 -> skip (belum cukup data, cold-start)
    4. days_until_stockout := product.stok / avg_daily_sold
    5. Jika days_until_stockout <= THRESHOLD (mis. 3 hari):
         predicted_restock_date := today + days_until_stockout
         predicted_qty := avg_daily_sold * RESTOCK_WINDOW (mis. 7 hari)
         simpan row baru ke restock_predictions
```

> `THRESHOLD` dan `RESTOCK_WINDOW` sebaiknya dibuat configurable (Viper) agar mudah di-tuning tanpa redeploy kode.

### 9.5 Cron Job — Promosi Otomatis (Malam Hari)

```
Jadwal: setiap malam, terpisah dari job restock (FR-23)
Untuk tiap store:
  Untuk tiap product milik store:
    1. qty_bulan_ini := transaction-client.SumQtyByProductInMonth(
                            productID, currentMonth)
    2. Evaluasi rule (if-else, contoh):
         if qty_bulan_ini >= RULE_THRESHOLD_QTY:
             rule_triggered = "QTY_MONTHLY_THRESHOLD"
    3. Jika rule terpenuhi DAN belum ada log untuk (product_id, period):
         simpan row baru ke promotion_logs { product_id, rule_triggered, period }
    // Trigger notifikasi WA ke pelanggan: di luar scope MVP (FR-25)
```

---

## 10. Integrasi dengan Layanan ML (Kontrak, Bukan Implementasi)

Karena tim ML belum ditentukan, backend mendefinisikan kontrak integrasi berupa interface `mlclient` yang dapat diimplementasikan menghadap layanan ML apa pun (in-house service, endpoint model hasil fine-tuning, atau third-party), selama mengikuti kontrak request/response berikut.

### 10.1 Interface `mlclient` (usulan)

```go
type MLClient interface {
	TranscribeAndExtract(ctx context.Context, audio []byte) ([]ExtractedItem, error)
	OcrAndExtract(ctx context.Context, image []byte) ([]ExtractedItem, error)
}

type ExtractedItem struct {
	RawText string  `json:"raw_text"`   // teks item mentah hasil ekstraksi
	Qty     float64 `json:"qty"`
	Harga   int64   `json:"harga"`
}
```

### 10.2 Kontrak Response (JSON) yang Diharapkan dari Layanan ML

```json
{
  "items": [
    { "raw_text": "beras", "qty": 1, "harga": 17000 },
    { "raw_text": "telor", "qty": 1, "harga": 18000 }
  ]
}
```

- Backend tidak bergantung pada detail internal model (Whisper/LoRA/OCR) — hanya pada kontrak request (file) & response (list item mentah) di atas.
- Matching `raw_text` ke `product_id` dilakukan di backend (Golang), bukan oleh layanan ML — sesuai PRD ("matching data produk menggunakan service Golang").
- Bila layanan ML belum tersedia saat development, `mlclient` dapat di-mock (mock implementation) agar development transaction module tidak terblokir.

---

## 11. Security Design

- Autentikasi: OAuth login menghasilkan JWT yang dipakai untuk semua request terproteksi via header `Authorization: Bearer <token>`.
- Middleware JWT (`internal/shared/middleware`) memvalidasi token dan menyuntikkan `user_id` ke context request sebelum masuk ke handler.
- Otorisasi tingkat data: setiap query di modul product/transaction/report/restock/promotion selalu difilter berdasarkan `store_id` milik user yang login (dicek lewat store-client), untuk mencegah akses lintas toko.
- Validasi input: seluruh request body divalidasi memakai go-playground/validator berbasis struct tag sebelum diproses handler/service.
- Data sensitif (kredensial, token) tidak di-log secara penuh oleh Logrus.

---

## 12. Error Handling & Konvensi Response

| HTTP Status | Kondisi                                 | Contoh Kasus                                   |
| ----------- | --------------------------------------- | ---------------------------------------------- |
| 400         | Bad Request — validasi input gagal      | Field wajib kosong, format salah               |
| 401         | Unauthorized — token tidak ada/invalid  | Header Authorization kosong                    |
| 403         | Forbidden — akses ke resource toko lain | store_id bukan milik user login                |
| 404         | Not Found                               | product_id tidak ditemukan                     |
| 409         | Conflict                                | Stok tidak mencukupi saat konfirmasi transaksi |
| 500         | Internal Server Error                   | Kegagalan koneksi database / layanan ML        |

> Semua response error tetap memakai format `WebResponse[T]` dengan `success: false` dan `message` yang informatif (lihat Bagian 5).

---

## 13. Desain Cron Job

- Scheduler dijalankan dalam proses yang sama dengan aplikasi utama (in-process scheduler, mis. robfig/cron) untuk MVP — dapat dipisah jadi worker terpisah bila beban bertambah.
- Dua job berjalan independen dan terjadwal terpisah: `restock_job` dan `promotion_job` (sesuai catatan PRD bahwa keduanya adalah proses berbeda).
- Setiap job idempotent per hari — `promotion_job` mengecek kombinasi (product_id, period) sebelum insert untuk menghindari duplikasi log bila job dijalankan ulang.
- Logging hasil tiap run cron (jumlah produk diproses, jumlah prediksi/promo baru, durasi) dicatat via Logrus untuk observabilitas.

---

## 15. Logging Strategy

- Logrus dipakai dengan format JSON terstruktur agar mudah diparse (siap untuk log aggregator bila dibutuhkan).
- Level log: **Info** untuk request masuk/keluar & hasil cron job, **Warn** untuk kondisi tidak ideal (mis. match_confidence rendah), **Error** untuk kegagalan (DB, layanan ML).
- Setiap log request menyertakan `request_id`, `user_id` (bila ada), dan `store_id` untuk memudahkan tracing per toko.

---

## 16. Deployment

### 16.1 Dockerfile (multi-stage, ringkasan)

```dockerfile
# Stage 1: build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o smart-commerce-api ./cmd/api

# Stage 2: runtime
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/smart-commerce-api .
COPY config.json .
EXPOSE 8080
CMD ["./smart-commerce-api"]
```

### 16.2 docker-compose.yml (ringkasan)

```yaml
services:
  api:
    build: .
    ports: ["8080:8080"]
    environment:
      - DATABASE_DSN=user:pass@tcp(db:3306)/smart_commerce
    depends_on: [db]
  db:
    image: mysql:8
    environment:
      MYSQL_DATABASE: smart_commerce
      MYSQL_ROOT_PASSWORD: rootpass
    ports: ["3306:3306"]
```

> Pilihan MySQL di atas mengikuti template acuan; dapat diganti PostgreSQL lewat dialector GORM tanpa mengubah struktur modul.

---

## 17. Testing Strategy

- Unit test per modul: service & business logic (mis. perhitungan untung, matching produk) diuji terisolasi dengan mock repository/client.
- Integration test: endpoint utama (create product, confirm transaction, daily report) diuji terhadap database test (mis. SQLite in-memory atau container MySQL sementara).
- Test khusus untuk skenario decrement stok & race condition (dua transaksi bersamaan pada produk yang sama) untuk memastikan stok tidak minus.
- Test idempotency cron job promosi (dijalankan dua kali dalam periode sama tidak menghasilkan duplikat).

---

## 18. Assumptions & Dependencies ke Tim Lain

- Tim ML akan menyediakan layanan yang mengikuti kontrak request/response pada Bagian 10 — detail arsitektur model (Whisper, NLP extraction, OCR, model restock) berada di luar tanggung jawab dokumen ini.
- Tim FE akan mengonsumsi API sesuai kontrak pada Bagian 8 & format `WebResponse[T]` pada Bagian 5.
- Provider OAuth spesifik (Google, dll.) masih perlu dikonfirmasi — desain auth module dibuat cukup generik untuk mendukung provider apa pun via `oauth_provider` + `oauth_id`.

---

## 19. Open Questions (Backend)

- Nilai default THRESHOLD restock (hari) dan ambang batas rule promo — masih perlu didiskusikan dengan pemilik produk/bisnis.
- Apakah perlu endpoint terpisah untuk pedagang mengoreksi/override hasil matching produk secara massal (bulk correction), di luar edit per-item pada preview?
- Strategi penyimpanan file audio/foto nota mentah (disimpan permanen untuk audit/re-training model, atau dibuang setelah diproses) — berdampak pada kebutuhan storage (mis. object storage) di layer backend.
- Apakah dibutuhkan rate limiting pada endpoint extract (voice/photo) mengingat kemungkinan biaya pemrosesan ML per-request cukup mahal.

---

## 20. Appendix — Glossary

| Istilah           | Keterangan                                                                                                                                      |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Modular Monolith  | Pola arsitektur satu deployment unit dengan modul-modul terisolasi secara data & komunikasi lewat interface                                     |
| Client Interface  | Interface publik yang diekspos tiap modul untuk diakses modul lain, tanpa membuka akses langsung ke tabel                                       |
| Snapshot Harga    | Penyimpanan harga_modal & harga_jual pada saat transaksi terjadi di transaction_items, agar tidak berubah walau harga produk di-update kemudian |
| On-the-fly Report | Laporan yang dihitung langsung dari query saat diminta, bukan dari tabel ringkasan yang di-precompute                                           |
| mlclient          | Interface backend untuk berkomunikasi dengan layanan ML eksternal (Whisper, NLP extraction, OCR)                                                |
