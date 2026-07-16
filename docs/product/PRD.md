# Product Requirements Document

## Smart Commerce

_Asisten AI untuk Pencatatan Transaksi, Untung-Rugi, Prediksi Restock & Promosi Otomatis bagi Pedagang UMKM_

| Field          | Detail                                                  |
| -------------- | ------------------------------------------------------- |
| Nama Proyek    | Smart Commerce                                          |
| Kompetisi      | AI Innovation Challenge (AIC)                           |
| Penulis        | Putra Rizky                                             |
| Tanggal Dibuat | 08 Juli 2026                                            |
| Status Dokumen | Draft — dalam pengembangan, sebagian scope masih ambigu |
| Versi          | 0.1                                                     |

---

## 1. Overview

### 1.1 Ringkasan

Smart Commerce adalah aplikasi web yang membantu pedagang UMKM (dengan studi kasus toko sembako) mengelola tokonya dengan lebih mudah. Alih-alih mencatat transaksi secara manual, pedagang cukup merekam suara atau memotret nota saat transaksi berlangsung. Sistem akan mengekstrak data transaksi (nama barang, jumlah, harga) secara otomatis, menghitung untung-rugi harian, memprediksi kapan produk perlu di-restock berdasarkan pola penjualan, dan memicu promosi otomatis berbasis aturan sederhana.

### 1.2 Problem Statement

- Pedagang UMKM (khususnya toko sembako) umumnya mencatat transaksi secara manual atau bahkan tidak mencatat sama sekali, sehingga sulit mengetahui untung-rugi harian secara akurat.
- Pedagang kesulitan menentukan kapan harus restock suatu produk karena tidak ada data pola penjualan yang terekam rapi.
- Belum ada mekanisme yang membantu pedagang memberi promosi ke pelanggan tetap secara otomatis dan konsisten.
- Proses input data transaksi manual (mengetik) memakan waktu dan mengganggu kecepatan layanan ke pembeli, terutama saat toko ramai.

### 1.3 Latar Belakang

Ide ini berangkat dari riset dan pengalaman pribadi penulis yang juga memiliki toko sembako, sehingga permasalahan yang diangkat merupakan permasalahan nyata di lapangan. Fokus AI pada produk ini bukan untuk menghitung untung-rugi (yang sifatnya kalkulasi sederhana: harga jual dikurangi harga modal), melainkan untuk (1) menangkap data transaksi dari suara/foto nota, dan (2) memprediksi kebutuhan restock berdasarkan histori penjualan.

---

## 2. Tujuan & Sasaran

### 2.1 Business Goals

- Menghadirkan solusi inovatif dan fungsional untuk kompetisi AI Innovation Challenge (AIC).
- Memvalidasi konsep asisten pencatatan transaksi berbasis suara/foto untuk segmen UMKM sembako.

### 2.2 User Goals

- Pedagang dapat mencatat transaksi dengan cepat tanpa mengetik manual, cukup lewat suara atau foto nota.
- Pedagang dapat melihat untung-rugi harian secara otomatis dan akurat.
- Pedagang mendapat rekomendasi kapan harus restock produk tertentu.
- Pedagang dapat memberikan promosi ke pelanggan tetap tanpa perlu menghitung manual.

### 2.3 Non-Goals (untuk MVP)

- Tidak membangun fitur "tutup hari / closing" — laporan bersifat read-only, dihitung on-the-fly.
- Tidak melakukan tracking identitas pembeli individual — logika promosi dibuat berbasis produk, bukan berbasis histori pelanggan.
- Tidak menggunakan AI untuk menghitung untung-rugi — cukup kalkulasi sederhana (harga jual − harga modal).
- Tidak menggunakan AI untuk promosi — cukup rule-based (if-else).

---

## 3. Scope MVP / Demo

### 3.1 In Scope

- Registrasi/login pedagang via OAuth, lalu membuat toko (create store).
- Input & kelola daftar produk (nama, harga modal, harga jual, stok awal, satuan).
- Pencatatan transaksi via rekaman suara → Whisper → NLP extraction → matching produk → preview yang dapat diedit → konfirmasi.
- Pencatatan transaksi via foto nota → OCR → NLP extraction (model yang sama dengan jalur suara) → preview yang dapat diedit → konfirmasi.
- Decrement stok produk otomatis setiap kali transaksi dikonfirmasi.
- Laporan harian read-only: total omset, total untung, jumlah transaksi, produk terlaris — dihitung langsung dari query transaksi hari itu (tanpa tabel ringkasan terpisah).
- Prediksi restock otomatis via cron job setiap malam, tersimpan di `restock_predictions`.
- Promosi otomatis berbasis produk (rule if-else) via cron job setiap malam, tersimpan di `promotion_logs`.

### 3.2 Out of Scope (MVP)

- Fitur tutup hari / closing transaksi harian.
- Tracking & profil identitas pelanggan individual.
- Pengiriman notifikasi WhatsApp secara aktual ke pelanggan (kemungkinan fase berikutnya).
- Model AI untuk perhitungan untung-rugi.
- Model AI untuk logika promosi.

---

## 4. User Persona

| Persona       | Deskripsi                                                                                                                |
| ------------- | ------------------------------------------------------------------------------------------------------------------------ |
| Pedagang UMKM | Pemilik/penjaga toko sembako. Pengguna utama aplikasi, melakukan input produk, merekam transaksi, dan melihat laporan.   |
| Pembeli       | Tidak memiliki akun/login. Hanya sumber data transaksi (item, qty, harga) yang ditangkap lewat suara/foto nota pedagang. |

---

## 5. User Flow

### Tahap 1 — Onboarding & Setup Produk

1. Pedagang registrasi/login menggunakan OAuth, lalu membuat toko (create store).
2. Pedagang menginput daftar produk: nama produk, harga modal, harga jual, stok awal, satuan, dll. — `POST /products`.

### Tahap 2 — Pencatatan Transaksi (Jalur Suara)

3. Saat pembeli datang dan transaksi terjadi, pedagang menekan tombol rekam suara (bisa langsung per transaksi, atau direkam sekaligus di akhir hari saat tutup toko).
4. Pedagang mengucapkan detail transaksi, contoh: "Beras 1kg 17rb, telor 1kg 18rb".
5. Backend memproses: audio → Whisper (speech-to-text) → NLP/LLM extraction (model hasil tuning) → matching data produk menggunakan service Golang → hasil dikembalikan sebagai preview ke frontend.
6. Pedagang melihat preview hasil ekstraksi di layar.
7. Preview dapat diedit oleh pedagang — mengantisipasi kesalahan tangkap nama produk/qty/harga oleh Whisper atau NLP.
8. Pedagang mengonfirmasi data → `POST /transactions`.
9. Backend membuat baris baru di tabel `transactions` (`transaction_date`: hari ini).
10. Backend menyimpan detail per item di `transaction_items`, mencatat harga modal & harga jual saat transaksi terjadi (snapshot harga), beserta produk & qty-nya.
11. Backend melakukan decrement stok produk sesuai qty yang terjual pada transaksi tersebut.

> Catatan: jalur ini berlaku sama untuk input via foto nota (lihat Bagian 6.4) — hanya berbeda di titik masuk data (OCR menggantikan Whisper), lalu masuk ke NLP extraction yang sama.

### Tahap 3 — Laporan Harian (Read-only)

12. Pedagang membuka tab "Laporan Hari Ini".
13. Backend mengambil (GET) seluruh `transactions` dengan `transaction_date` = hari ini, milik store pedagang tersebut.
14. Backend menghitung total penjualan (omset), total untung, dan jumlah transaksi — dihitung langsung dari query, bukan dari tabel ringkasan terpisah.
15. Perhitungan untung menggunakan `harga_modal` & `harga_jual` yang tersimpan di `transaction_items` pada saat transaksi (bukan harga produk saat ini), agar laporan tetap akurat meski harga produk berubah di kemudian hari.
16. Backend juga menghitung produk apa saja yang laku dan terjual paling banyak hari itu.
17. Response ke frontend: total omset, total untung, produk terlaris, sisa stok, dst.

> MVP/demo laporan ini bersifat read-only saja, tanpa fitur tutup hari.

### Tahap 4 — Proses Otomatis Malam Hari (Cron Job)

18. **Prediksi Restock**: setiap malam, sistem melakukan loop ulang untuk memprediksi kebutuhan restock tiap produk berdasarkan histori penjualan, lalu menyimpan hasilnya ke `restock_predictions`.
19. **Promosi Otomatis**: setiap malam (proses terpisah dari restock), sistem mengecek produk mana yang memenuhi rule promo (logika if-else berbasis produk), lalu mencatat hasilnya ke `promotion_logs`. Trigger notifikasi (misalnya WA) ke pelanggan dapat dikembangkan pada fase berikutnya.

> Perbedaan penting: proses restock & promosi berjalan otomatis di background lewat cron job setiap malam. Sedangkan laporan harian (Tahap 3) hanya dihitung/muncul saat pedagang membuka tab laporan — tidak berjalan sebagai proses background.

---

## 6. Functional Requirements

### 6.1 Autentikasi & Toko

| ID    | Requirement                                              | Priority |
| ----- | -------------------------------------------------------- | -------- |
| FR-01 | Pedagang dapat registrasi/login menggunakan OAuth        | Must     |
| FR-02 | Pedagang dapat membuat toko (create store) setelah login | Must     |

### 6.2 Manajemen Produk

| ID    | Requirement                                                                                                        | Priority |
| ----- | ------------------------------------------------------------------------------------------------------------------ | -------- |
| FR-03 | Pedagang dapat menambah produk dengan atribut: nama, harga modal, harga jual, stok awal, satuan (`POST /products`) | Must     |
| FR-04 | Sistem menyimpan harga modal per produk sebagai referensi awal perhitungan untung-rugi                             | Must     |
| FR-05 | Sistem melakukan decrement stok otomatis setiap transaksi terkonfirmasi                                            | Must     |

### 6.3 Pencatatan Transaksi via Suara

| ID    | Requirement                                                                                                            | Priority |
| ----- | ---------------------------------------------------------------------------------------------------------------------- | -------- |
| FR-06 | Pedagang dapat merekam suara transaksi per transaksi atau sekaligus di akhir hari                                      | Must     |
| FR-07 | Sistem mengonversi audio menjadi teks menggunakan model Whisper                                                        | Must     |
| FR-08 | Sistem mengekstrak teks menjadi data terstruktur (item, qty, harga) menggunakan model NLP/LLM (fine-tuned dengan LoRA) | Must     |
| FR-09 | Sistem melakukan matching hasil ekstraksi dengan data produk yang sudah terdaftar (service Golang)                     | Must     |
| FR-10 | Sistem menampilkan hasil sebagai preview yang dapat diedit oleh pedagang sebelum dikonfirmasi                          | Must     |
| FR-11 | Pedagang dapat mengonfirmasi preview untuk menyimpan transaksi (`POST /transactions`)                                  | Must     |

### 6.4 Pencatatan Transaksi via Foto Nota

| ID    | Requirement                                                                                                            | Priority |
| ----- | ---------------------------------------------------------------------------------------------------------------------- | -------- |
| FR-12 | Pedagang dapat mengunggah/memotret foto nota transaksi                                                                 | Should   |
| FR-13 | Sistem melakukan OCR (model pre-trained, tanpa tuning) terhadap foto nota untuk mengekstrak teks                       | Should   |
| FR-14 | Hasil OCR diproses oleh model NLP extraction yang sama dengan jalur suara (reuse model, beda titik masuk)              | Should   |
| FR-15 | Hasil ekstraksi ditampilkan sebagai preview yang dapat diedit, mengikuti alur konfirmasi yang sama seperti jalur suara | Should   |

### 6.5 Laporan Harian

| ID    | Requirement                                                                                                                                                                       | Priority |
| ----- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FR-16 | Sistem menghitung total omset, total untung, dan jumlah transaksi hari berjalan secara on-the-fly dari data `transactions` & `transaction_items` (bukan tabel ringkasan terpisah) | Must     |
| FR-17 | Perhitungan untung menggunakan `harga_modal` & `harga_jual` yang tersimpan di `transaction_items` saat transaksi terjadi, bukan harga produk saat ini                             | Must     |
| FR-18 | Sistem menampilkan produk terlaris (paling banyak terjual)                                                                                                                        | Must     |
| FR-19 | Sistem menampilkan sisa stok produk                                                                                                                                               | Should   |
| FR-20 | Laporan bersifat read-only untuk MVP — tidak ada fitur tutup hari/closing                                                                                                         | Must     |

### 6.6 Prediksi Restock (Otomatis, Cron Job)

| ID    | Requirement                                                                                                         | Priority |
| ----- | ------------------------------------------------------------------------------------------------------------------- | -------- |
| FR-21 | Sistem menjalankan proses prediksi restock untuk tiap produk setiap malam (cron job), berdasarkan histori penjualan | Must     |
| FR-22 | Hasil prediksi disimpan ke tabel `restock_predictions`                                                              | Must     |

### 6.7 Promosi Otomatis (Otomatis, Cron Job)

| ID    | Requirement                                                                                                                                                         | Priority    |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| FR-23 | Sistem mengecek setiap malam (cron job) produk mana saja yang memenuhi rule promo berbasis produk (logika if-else, contoh: produk X terjual > N unit dalam sebulan) | Must        |
| FR-24 | Hasil pengecekan yang memenuhi rule dicatat ke tabel `promotion_logs`                                                                                               | Must        |
| FR-25 | Trigger notifikasi/WA ke pelanggan atas promosi — di luar scope MVP, dapat dikembangkan di fase berikutnya                                                          | Won't (MVP) |

---

## 7. Data Model (Ringkasan)

| Tabel               | Field Utama                                                                                | Keterangan                                                                                               |
| ------------------- | ------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------- |
| stores              | id, owner_id, nama_toko                                                                    | Toko milik pedagang, dibuat saat onboarding                                                              |
| products            | id, store_id, nama_produk, harga_modal, harga_jual, stok, satuan                           | Master data produk per toko                                                                              |
| transactions        | id, store_id, transaction_date, created_at                                                 | Header transaksi, satu baris per transaksi                                                               |
| transaction_items   | id, transaction_id, product_id, qty, harga_modal_saat_transaksi, harga_jual_saat_transaksi | Detail item per transaksi; menyimpan snapshot harga agar laporan tetap akurat meski harga produk berubah |
| restock_predictions | id, product_id, predicted_restock_date, predicted_qty, created_at                          | Hasil prediksi restock, digenerate cron job tiap malam                                                   |
| promotion_logs      | id, product_id, rule_triggered, created_at                                                 | Log produk yang memenuhi rule promo, digenerate cron job tiap malam                                      |

> Catatan: identitas pembeli tidak disimpan/di-track pada MVP ini — data yang ditangkap dari suara/foto nota hanya item, qty, dan harga.

---

## 8. Komponen AI / Model

### 8.1 Ringkasan Model

| Komponen                | Pendekatan                      | Fungsi                                                                             |
| ----------------------- | ------------------------------- | ---------------------------------------------------------------------------------- |
| Whisper (fine-tuned)    | Fine-tuning model Whisper       | Mengubah rekaman suara pedagang menjadi teks                                       |
| Model NLP Extraction    | Fine-tuning LLM dengan LoRA     | Mengubah teks (dari Whisper maupun OCR) menjadi data terstruktur: item, qty, harga |
| OCR (foto nota)         | Model pre-trained, tanpa tuning | Mengekstrak teks dari foto nota untuk diteruskan ke model NLP extraction yang sama |
| Model Prediksi Restock  | Dibangun/dilatih from scratch   | Memprediksi kapan produk perlu di-restock berdasarkan pola penjualan historis      |
| Logika Promosi          | Rule-based (if-else), bukan AI  | Menentukan produk yang memenuhi syarat promo                                       |
| Perhitungan Untung-Rugi | Kalkulasi sederhana (bukan AI)  | harga_jual − harga_modal, dihitung langsung dari data transaksi                    |

### 8.2 Pipeline — Jalur Suara

1. Input audio dari pedagang (rekaman suara transaksi).
2. Preprocessing: konversi audio ke format yang sesuai untuk model (misalnya spectrogram).
3. Encoder/Decoder Whisper: audio diubah menjadi teks (speech-to-text).
4. Detokenization: normalisasi angka/istilah lisan menjadi teks yang natural (mendekati cara manusia mengetik).
5. Teks hasil transkripsi diteruskan ke model NLP extraction untuk diubah menjadi data terstruktur (item, qty, harga).

### 8.3 Pipeline — Jalur Foto Nota

1. Input foto nota dari pedagang.
2. OCR (model pre-trained) mengekstrak teks mentah dari foto.
3. Teks hasil OCR diteruskan ke model NLP extraction yang sama dengan jalur suara — reuse model, hanya berbeda titik masuk data.
4. Lanjut ke tahap preview & konfirmasi yang sama seperti jalur suara.

### 8.4 Pendekatan Fine-tuning Model NLP Extraction

1. Load dataset (contoh transkrip transaksi & anotasi item/qty/harga).
2. Konfigurasi target output: item, qty, harga.
3. Fine-tuning menggunakan metode LoRA terhadap model LLM dasar.
4. Training — model disesuaikan dengan kebutuhan fitur ekstraksi transaksi.
5. Simpan model hasil tuning (format seperti `.pkl`, `.h5`, `.keras`, atau TFJS sesuai kebutuhan deployment).

> Catatan status: pipeline ini masih berupa rancangan awal (draft) dan dapat berubah menyesuaikan hasil eksplorasi teknis selama pengembangan.

---

## 9. Technical Considerations

- Backend menggunakan Golang, khususnya untuk proses matching hasil ekstraksi NLP terhadap data produk yang terdaftar.
- Model Whisper dan model NLP extraction di-hosting/serving terpisah dari backend utama, diakses melalui pemanggilan service/API.
- Cron job berjalan terpisah setiap malam untuk dua proses independen: prediksi restock dan pengecekan promosi — keduanya tidak terkait dengan laporan harian yang bersifat on-demand.
- Decrement stok wajib terjadi tepat setelah transaksi dikonfirmasi (`POST /transactions`) agar data stok akurat secara real-time.
- Snapshot harga modal & harga jual wajib disimpan di `transaction_items` pada saat transaksi terjadi, agar laporan historis tidak terpengaruh perubahan harga produk di kemudian hari.
- Perlu ditentukan strategi penyimpanan file audio/foto nota (sementara vs permanen) serta kebijakan retensinya.

---

## 10. Open Questions / Hal yang Masih Ambigu

- Definisi rule promo per produk secara detail (contoh ambang batas qty/frekuensi) belum final.
- Mekanisme dan channel notifikasi promosi ke pelanggan (WA API atau lainnya) belum ditentukan — di luar scope MVP.
- Sumber dan strategi data training untuk model prediksi restock belum ditentukan (potensi cold-start problem karena data histori penjualan awal masih minim).
- Format penyimpanan model NLP extraction & Whisper (pkl/h5/keras/tfjs) serta strategi deployment/serving-nya belum final.

---

## 11. Risks & Mitigasi

| Risiko                                                                                         | Mitigasi                                                                                     |
| ---------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| Whisper salah menangkap ucapan, terutama dialek/aksen pasar lokal                              | Preview hasil ekstraksi dibuat dapat diedit sebelum dikonfirmasi (FR-10)                     |
| Model NLP extraction salah mengenali nama produk/qty/harga                                     | Sama seperti di atas — preview & konfirmasi manual sebagai lapisan verifikasi                |
| Data histori penjualan awal masih sedikit sehingga prediksi restock kurang akurat (cold start) | Perlu strategi fallback (misal aturan sederhana) sebelum data historis mencukupi untuk model |
| Kualitas foto nota buram/tidak jelas menyebabkan hasil OCR buruk                               | Preview & edit manual sebelum konfirmasi transaksi                                           |

---

## 12. Timeline

Timeline detail belum ditentukan pada versi draft ini — menyesuaikan jadwal kompetisi AI Innovation Challenge (AIC). Akan diperbarui pada revisi berikutnya.
