# Working Process — Smart Commerce Backend

> **Dokumen Pedoman Utama:**
> - [`docs/prd.md`](file:///c:/Users/traa/Desktop/Backend-AIC/docs/prd.md) — Product Requirements Document
> - [`docs/sdd.md`](file:///c:/Users/traa/Desktop/Backend-AIC/docs/sdd.md) — System Design Document
> - [`docs/SYSTEM_MAP.md`](file:///c:/Users/traa/Desktop/Backend-AIC/docs/SYSTEM_MAP.md) — System Map (kondisi aktual codebase)
>
> Ketiga dokumen di atas adalah **sumber kebenaran (source of truth)** proyek ini. Semua keputusan teknis harus bisa ditelusuri balik (traceable) ke salah satu dari tiga dokumen tersebut.
>
> **Workflow ini WAJIB diikuti oleh semua kontributor — termasuk agent AI mana pun yang mengerjakan repository ini.** Tidak ada pengecualian.

---

## Prinsip Inti

**PLAN FIRST, IMPLEMENT LATER.**

Tidak ada task yang boleh langsung masuk fase coding tanpa melalui fase perencanaan tertulis yang sudah disetujui oleh owner.

---

## Siklus Kerja

Setiap task mengikuti 8 tahapan berurutan. Tidak boleh melompati tahapan.

```
┌──────────────┐    ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│  1. Task      │───▶│  2. Reference    │───▶│  3. Planning     │───▶│  4. Approval     │
│     Intake    │    │     Check        │    │     Phase        │    │     Gate         │
└──────────────┘    └──────────────────┘    └──────────────────┘    └───────┬──────────┘
                                                                           │ owner says
                                                                           │ "approved"
                                                                           ▼
┌──────────────┐    ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│  8. Doc      │◀───│  7. Verification │◀───│  6. Code Review  │◀───│  5. Implement-   │
│     Sync     │    │     Phase        │    │     Gate         │    │     ation Phase  │
└──────────────┘    └──────────────────┘    └──────────────────┘    └──────────────────┘
```

Tambahan section yang berlaku **di sepanjang siklus**:
- **Section 9**: Handling Ambiguity / Conflict
- **Section 10**: Aturan Khusus Task ML/AI Integration

---

## 1. Task Intake

### 1.1 Sumber Task

Setiap task/fitur **WAJIB** berasal dari salah satu sumber berikut:

| Sumber | Contoh Rujukan | Kapan Dipakai |
| --- | --- | --- |
| FR-ID di `prd.md` | FR-03, FR-16 | Fitur baru sesuai requirement produk |
| Section di `sdd.md` | §6.3, §8.4, §9.2 | Implementasi arsitektur/desain teknis |
| Gap/bug/issue di `SYSTEM_MAP.md` | §13.1, §11.1, §14 | Perbaikan bug, alignment codebase ke desain |
| Permintaan eksplisit dari owner | Instruksi langsung | Ad-hoc task yang belum tercakup di dokumen |

### 1.2 Aturan

- Setiap task **WAJIB** punya rujukan (traceability) ke minimal satu dokumen di atas.
- Jika task tidak bisa ditelusuri ke dokumen mana pun → **STOP**, minta klarifikasi ke owner.
- Jika task berasal dari permintaan owner yang belum ada di PRD/SDD/system_map, catat bahwa ini adalah task ad-hoc dan tetap lanjutkan — tapi dokumentasikan di planning phase bahwa rujukan dokumennya belum ada.

---

## 2. Reference Check

### 2.1 Langkah Wajib

Sebelum menyusun rencana implementasi, **WAJIB** baca ulang bagian-bagian yang relevan dari ketiga dokumen:

1. **`prd.md`** — Baca FR-ID terkait beserta konteksnya (user flow, data model, prioritas Must/Should/Won't).
2. **`sdd.md`** — Baca section modul terkait: data schema (§6), API design (§8), business logic (§9), client interface (§7), dan aturan arsitektur (§2).
3. **`SYSTEM_MAP.md`** — Baca status implementasi terkini: modul sudah implemented/stub/belum ada (§5), endpoint aktual (§6), schema aktual (§7), known issues (§13), gap analysis (§11).

### 2.2 Hierarki Rujukan

Jika **tidak ada konflik** antar dokumen, hierarki rujukan adalah:

```
Codebase aktual  →  SYSTEM_MAP.md  →  sdd.md  →  prd.md
(paling dipercaya)                              (paling aspiratif)
```

Alasan: `SYSTEM_MAP.md` mencerminkan **kondisi aktual**, bukan aspirasi. Keputusan-keputusan yang sudah di-resolve (mis. `cmd/web/` vs `cmd/api/`, custom user ID vs UUID) tercatat di sana.

### 2.3 Jika Ditemukan Konflik

Jika ditemukan **konflik atau inkonsistensi** antar dokumen:

1. **STOP** — jangan lanjut ke planning.
2. Laporkan konflik secara eksplisit ke owner dengan format:

   > **Konflik ditemukan:**
   > - [Dokumen A, section X] menyatakan: "..."
   > - [Dokumen B, section Y] menyatakan: "..."
   > - Mohon keputusan sebelum saya lanjutkan.

3. Tunggu keputusan owner sebelum melanjutkan.

---

## 3. Planning Phase

### 3.1 Kewajiban

Sebelum menulis kode **APA PUN**, agent/kontributor **WAJIB** menyusun rencana implementasi tertulis yang mencakup:

| Item | Deskripsi | Wajib? |
| --- | --- | --- |
| **Rujukan** | FR-ID dan/atau section SDD yang di-address oleh task ini | ✅ |
| **File yang disentuh** | Daftar file yang akan dibuat / dimodifikasi (path lengkap) | ✅ |
| **Module & client interface** | Module mana yang terlibat, client interface mana yang digunakan (ref SDD §7) | ✅ |
| **Perubahan schema data** | Entity baru, kolom baru, migration file — jika ada | Jika relevan |
| **Endpoint baru/berubah** | Method, path, auth requirement, request/response body | Jika relevan |
| **Dampak ke modul lain** | Apakah perlu update client interface modul lain, atau perlu koordinasi dengan modul yang belum dibangun | ✅ |
| **Konvensi yang diikuti** | Konfirmasi bahwa implementasi akan mengikuti pola modular monolith (ref SDD §2.2) | ✅ |
| **Out of scope** | Hal-hal yang TIDAK termasuk dalam task ini | ✅ |
| **Open questions** | Hal-hal yang masih ambigu dan butuh keputusan owner | Jika ada |

### 3.2 Format

Rencana implementasi bisa ditulis dalam bentuk:
- Implementation plan artifact di chat
- Markdown di chat langsung

Yang penting: **tertulis, terstruktur, dan visible ke owner**.

### 3.3 Aturan

- Rencana ini **HARUS** ditampilkan ke owner.
- Rencana ini **HARUS** menunggu approval eksplisit sebelum lanjut ke tahap berikutnya.
- Jika owner memberikan revisi → update rencana → minta approval ulang.

---

## 4. Approval Gate

### 4.1 Aturan Utama

- Implementasi **HANYA** boleh dimulai setelah **owner** memberikan approval eksplisit.
- Bentuk approval yang valid: owner mengatakan "approved", "lanjut", "setuju", "ok lanjut", atau pernyataan serupa yang jelas menyetujui rencana.
- **Tanpa approval eksplisit, TIDAK ADA kode yang boleh ditulis.**

### 4.2 Siapa yang Approve

- **Selalu owner (Putra Rizky)** — tidak ada delegasi approval ke pihak lain.

### 4.3 Jika Ada Revisi

- Jika owner memberikan feedback/revisi terhadap rencana:
  1. Update rencana implementasi sesuai feedback.
  2. Tampilkan ulang rencana yang sudah direvisi.
  3. Tunggu approval ulang.
- Siklus ini diulang sampai owner approve.

---

## 5. Implementation Phase

### 5.1 Aturan Arsitektur Modular Monolith

Semua implementasi **WAJIB** mengikuti prinsip arsitektur di SDD §2.2:

| Prinsip | Aturan | Ref SDD |
| --- | --- | --- |
| **Data Isolation** | Setiap modul punya tabel sendiri. Tidak ada foreign key GORM lintas modul — referensi antar modul berupa plain reference ID (string), bukan FK database | §2.2 |
| **Inter-module Communication** | Akses data modul lain **hanya** lewat `<module>-client` interface. Dilarang query/JOIN langsung ke tabel modul lain | §2.2, §7 |
| **Independent Migration** | Setiap modul menjalankan `AutoMigrate` untuk tabelnya sendiri via `Migrate()` | §2.2 |
| **Module Contract** | Setiap modul mengimplementasikan interface `module.Module`: `Migrate() error` dan `RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler)` | §2.4, §9.2 (SYSTEM_MAP) |

### 5.2 Struktur Folder Modul Baru

Mengikuti pola codebase existing (ref SYSTEM_MAP §4):

```
internal/module/<nama>/
├── module.go                          # Wiring: New(), Migrate(), Client()
├── route.go                           # RegisterRoutes
├── client_impl.go                     # Implements <nama>-client.Client
└── src/
    ├── controller/
    │   └── <nama>_controller.go       # HTTP handler
    ├── entity/
    │   └── <nama>_entity.go           # GORM struct → tabel DB
    ├── model/
    │   ├── <nama>_request.go          # Request DTO
    │   ├── <nama>_response.go         # Response DTO
    │   └── converter/
    │       └── <nama>_converter.go    # Entity ↔ Response converter
    ├── repository/
    │   └── <nama>_repository.go       # Data access layer
    └── usecase/
        └── <nama>_usecase.go          # Business logic

internal/module/<nama>-client/
└── client.go                          # Public interface (kontrak untuk modul lain)
```

### 5.3 Konvensi Kode

| Aspek | Konvensi | Referensi |
| --- | --- | --- |
| **Response API** | Selalu gunakan `WebResponse[T]` dari `shared/response` | SDD §5, SYSTEM_MAP §9.4 |
| **Validasi input** | Gunakan `go-playground/validator` via struct tags | SDD §4 |
| **ID generation** | Custom format (bukan raw UUID) — ikuti pola `utils.GenerateUserId` | SYSTEM_MAP §4 |
| **Timestamp** | `int64` milli-epoch, GORM tag `autoCreateTime:milli` / `autoUpdateTime:milli` | SYSTEM_MAP §7.1 (OQ-3) |
| **Error response** | Format `WebResponse[T]` dengan `success: false`, HTTP status sesuai SDD §12 | SDD §12 |
| **Logging** | Gunakan Logrus (JSON structured), level sesuai SDD §15 | SDD §15 |
| **Registrasi modul** | Daftarkan modul baru di `cmd/web/main.go` | SDD §2.5 |

### 5.4 SQL Migration

- Setiap perubahan schema **WAJIB** disertai SQL migration file di `db/migrations/`.
- Format nama file: `<timestamp>_<deskripsi>.up.sql` dan `.down.sql`.
- Migration file adalah referensi skema resmi — jika ada perbedaan antara GORM entity dan migration, **migration yang benar** (ref SYSTEM_MAP §9.8).

---

## 6. Code Review Gate

### 6.1 Kapan Dilakukan

Setelah implementasi selesai dan **sebelum** masuk ke verification phase, hasil kode harus di-review oleh owner.

### 6.2 Yang Di-review

- Apakah implementasi sudah sesuai dengan rencana yang disetujui di Approval Gate.
- Apakah arsitektur modular monolith dipatuhi (data isolation, client interface, dll).
- Apakah ada perubahan yang menyimpang dari rencana — jika ya, harus dijelaskan alasannya.

### 6.3 Aturan

- Agent/kontributor **WAJIB** menampilkan ringkasan perubahan (file apa saja yang berubah, endpoint baru, schema baru) ke owner.
- Owner melakukan review dan memberikan salah satu:
  - **Approved** → lanjut ke Verification Phase.
  - **Revisi** → perbaiki sesuai feedback → tampilkan ulang → minta review lagi.

---

## 7. Verification Phase

### 7.1 Checklist Wajib

Sebelum task dianggap selesai, **SEMUA** item berikut harus terpenuhi:

#### Build & Test
- [ ] `go build ./...` — build berhasil tanpa error
- [ ] `go test ./...` — unit test lolos (termasuk test baru untuk kode yang ditulis)
- [ ] Unit test mencakup minimal: business logic di usecase dan skenario utama

#### Kesesuaian dengan Dokumen
- [ ] FR-ID terkait sudah terpenuhi sesuai definisi di `prd.md`
- [ ] Implementasi sesuai dengan desain di `sdd.md` (schema, API contract, business logic)
- [ ] Tidak ada penyimpangan dari rencana yang disetujui di Approval Gate — atau jika ada, sudah dijelaskan dan disetujui di Code Review Gate

#### Kepatuhan Arsitektur
- [ ] Data isolation: tidak ada FK GORM lintas modul
- [ ] Client interface: tidak ada query/JOIN langsung ke tabel modul lain
- [ ] Response format: semua endpoint menggunakan `WebResponse[T]`
- [ ] Validasi input: request body divalidasi via struct tags (go-playground/validator)
- [ ] Modul baru sudah terdaftar di `cmd/web/main.go`

#### Artefak
- [ ] SQL migration file tersedia di `db/migrations/` (jika ada perubahan schema)
- [ ] Client interface (`<nama>-client/client.go`) sudah didefinisikan (jika modul baru)

---

## 8. Documentation Sync

### 8.1 Aturan Utama

Setelah implementasi selesai dan terverifikasi, dokumen-dokumen referensi **harus diperbarui** agar tetap akurat.

### 8.2 `SYSTEM_MAP.md` — Update Wajib

`SYSTEM_MAP.md` mencerminkan **kondisi aktual codebase**, bukan aspirasi. Setiap kali codebase berubah, system map harus diupdate:

| Yang diupdate | Contoh |
| --- | --- |
| Status modul | ❌ Not Started → 🟡 Stub → ✅ Implemented |
| Endpoint aktual | Tambahkan endpoint baru ke §6.1 |
| Schema tabel | Tambahkan tabel baru ke §7.1 dengan detail field |
| Client interface | Update status di §3 Client Interface Inventory |
| Known issues | Tambahkan issue baru atau hapus yang sudah di-fix |
| Gap analysis | Update §11 jika gap sudah tertutup |
| Pending action items | Update §14 |

### 8.3 `prd.md` dan `sdd.md` — Update Kondisional

Jika implementasi mengungkap bahwa PRD atau SDD perlu diupdate (misalnya: keputusan teknis baru, konvensi yang berubah, gap yang ditemukan):

1. **JANGAN langsung edit.** 
2. Ajukan perubahan ke owner dengan format:

   > **Pengajuan update dokumen:**
   > - **Dokumen**: [prd.md / sdd.md]
   > - **Section**: [section mana]
   > - **Perubahan yang diusulkan**: [jelaskan apa yang ingin diubah dan alasannya]
   > - **Alasan**: [kenapa perubahan ini diperlukan]

3. Tunggu approval owner sebelum melakukan perubahan.

### 8.4 Timing

- Documentation sync dilakukan **segera setelah task selesai** (setelah Verification Phase lulus).
- Task **belum dianggap "done"** sampai documentation sync selesai.

---

## 9. Handling Ambiguity / Conflict

### 9.1 Aturan Mutlak

> **Jika ragu, TANYA. Jangan asumsikan.**

Aturan ini berlaku di **setiap tahapan** siklus kerja, tanpa pengecualian.

### 9.2 Prosedur

Ketika agent/kontributor menemukan ambiguitas, konflik, atau keraguan:

1. **BERHENTI** — jangan lanjutkan ke langkah berikutnya.
2. **Dokumentasikan** ambiguitas/pertanyaan secara eksplisit dan jelas.
3. **Tanyakan** ke owner.
4. **Tunggu** jawaban sebelum melanjutkan.

### 9.3 Larangan

Berikut hal-hal yang **DILARANG KERAS**:

- ❌ Mengasumsikan jawaban atas ambiguitas
- ❌ Memutuskan sendiri di luar scope yang sudah disetujui
- ❌ Mengubah dokumen referensi (PRD/SDD/SYSTEM_MAP) tanpa persetujuan owner
- ❌ Mengubah arsitektur atau pola desain tanpa diskusi dengan owner
- ❌ Menambahkan dependency baru tanpa menyebutkannya di planning phase
- ❌ Membuat tabel/kolom/endpoint yang tidak ada di rencana yang sudah disetujui

### 9.4 Contoh Situasi yang HARUS Ditanyakan

- Konflik antara PRD dan SDD (misalnya data model berbeda)
- Field/tabel yang tidak jelas tipenya atau kegunaannya
- Business logic yang ambigu atau punya lebih dari satu interpretasi
- Perubahan yang berdampak ke modul lain yang belum disetujui
- Temuan bug di modul existing yang tidak terkait task saat ini
- Keputusan yang akan membuat implementasi menyimpang dari rencana

---

## 10. Aturan Khusus: Task ML/AI Integration

### 10.1 Konteks

SDD §10 mendefinisikan kontrak integrasi backend dengan layanan ML eksternal melalui interface `mlclient`. Layanan ML (Whisper, NLP extraction, OCR, model restock) berada di luar scope backend — backend hanya bertanggung jawab atas **kontrak interface** dan **mock implementation**.

### 10.2 Aturan Tambahan untuk Task ML

Selain mengikuti seluruh siklus kerja di atas, task yang melibatkan integrasi ML/AI **WAJIB** memperhatikan:

| Aturan | Detail | Ref |
| --- | --- | --- |
| **Interface-first** | Definisikan `mlclient` interface terlebih dahulu di `internal/pkg/mlclient/`. Semua interaksi backend → ML hanya melalui interface ini | SDD §10.1 |
| **Mock wajib ada** | Sediakan mock implementation dari `mlclient` agar development modul `transaction` tidak terblokir oleh ketersediaan layanan ML | SDD §10, akhir paragraf |
| **Kontrak response** | Response dari layanan ML harus mengikuti format `ExtractedItem` yang didefinisikan di SDD §10.1–10.2 | SDD §10.2 |
| **Matching di backend** | Matching `raw_text` ke `product_id` dilakukan di backend (Golang), **bukan** oleh layanan ML | PRD §5 step 5, SDD §10 |
| **Tidak bergantung pada ML** | Backend harus bisa berjalan (start, serve endpoint lain) tanpa layanan ML aktif. Endpoint extract boleh return error jika ML tidak tersedia, tapi tidak boleh crash | — |

### 10.3 Planning Phase Tambahan

Untuk task ML integration, rencana implementasi juga harus mencakup:
- Interface method yang akan didefinisikan/digunakan
- Bagaimana mock implementation akan bekerja (data dummy apa yang dikembalikan)
- Bagaimana error handling jika layanan ML tidak tersedia (timeout, unreachable)

---

## Quick Reference Checklist

Gunakan checklist ini sebagai **reminder cepat** di setiap mulai task baru:

```
SEBELUM CODING:
  □ Task punya rujukan ke PRD / SDD / SYSTEM_MAP
  □ Sudah baca ulang bagian relevan dari ketiga dokumen
  □ Tidak ada konflik antar dokumen (atau sudah dilaporkan & di-resolve)
  □ Rencana implementasi sudah ditulis lengkap
  □ Owner sudah memberikan approval eksplisit

SAAT CODING:
  □ Mengikuti pola modular monolith (data isolation, client interface)
  □ Struktur folder sesuai pola existing
  □ WebResponse[T] untuk semua endpoint
  □ Validasi input via struct tags
  □ Tidak ada asumsi — semua keraguan sudah ditanyakan

SETELAH CODING:
  □ Code review oleh owner — approved
  □ go build ./... berhasil
  □ go test ./... lolos (termasuk unit test baru)
  □ Verification checklist lengkap terpenuhi
  □ SYSTEM_MAP.md sudah diupdate
  □ PRD/SDD update (jika perlu) sudah diajukan & disetujui
  □ Task boleh dianggap DONE
```

---

_Dokumen ini dibuat pada 13 Juli 2026. Berlaku untuk seluruh task berikutnya di repository Smart Commerce Backend._
