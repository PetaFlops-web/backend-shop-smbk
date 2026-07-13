# AGENTS.md — Smart Commerce Backend

> **File ini adalah entry point bagi setiap agent AI yang mengerjakan repository ini.**
> Baca seluruh isi file ini sebelum melakukan apa pun.

---

## 1. Source of Truth

Tiga dokumen berikut adalah **pedoman utama** proyek ini. Semua keputusan teknis harus bisa ditelusuri balik (traceable) ke salah satu dari ketiganya:

| Dokumen | Fungsi |
| --- | --- |
| [`docs/prd.md`](docs/prd.md) | Product Requirements — FR-ID, user flow, data model, prioritas |
| [`docs/sdd.md`](docs/sdd.md) | System Design — arsitektur, schema, API, business logic, tech stack |
| [`docs/SYSTEM_MAP.md`](docs/SYSTEM_MAP.md) | System Map — kondisi **aktual** codebase, gap analysis, known issues |

Jika suatu task tidak bisa ditelusuri ke dokumen mana pun, **tanyakan ke owner** sebelum melanjutkan.

---

## 2. Proses Kerja Wajib

**Baca [`workflow.md`](.agents/workflow.md) SEBELUM memulai task apa pun.** Dokumen tersebut berisi SOP 8 tahap yang wajib diikuti tanpa pengecualian:

```
Task Intake → Reference Check → Planning → Approval Gate →
Implementation → Code Review → Verification → Documentation Sync
```

Dua prinsip inti yang tidak boleh dilanggar:

1. **PLAN FIRST, IMPLEMENT LATER** — tidak ada kode yang boleh ditulis tanpa rencana tertulis yang sudah disetujui owner.
2. **Jika ragu, TANYA. Jangan asumsikan** — setiap ambiguitas harus ditanyakan ke owner, bukan diputuskan sendiri (detail: workflow.md §9).

---

## 3. Task Tracking

Lihat [`docs/progress.md`](docs/progress.md) untuk mengetahui:
- Task mana yang sudah selesai (`[x]`) dan mana yang belum (`[ ]`)
- Task berikutnya yang harus dikerjakan (ikuti urutan dependency modul)
- Rujukan FR-ID / section SDD untuk setiap task

**`docs/progress.md` adalah satu-satunya tempat tracking task teknis.** Update file ini setelah menyelesaikan task (sesuai workflow.md §8 — Documentation Sync).

---

## 4. Konvensi Teknis (Ringkasan)

Backend menggunakan arsitektur **Modular Monolith** (Go/Golang). Aturan terpenting yang harus diketahui sejak awal:

- **Data Isolation**: setiap modul punya tabel sendiri, **tidak boleh** ada foreign key GORM lintas modul — gunakan plain reference ID.
- **Client Interface**: akses data modul lain **hanya** lewat `<module>-client` interface, **dilarang** query/JOIN langsung ke tabel modul lain.
- **Response Format**: semua endpoint menggunakan `WebResponse[T]` (di `shared/response`).
- **Validasi Input**: `go-playground/validator` via struct tags.
- **Timestamp**: `int64` milli-epoch (GORM `autoCreateTime:milli` / `autoUpdateTime:milli`).

Detail lengkap arsitektur, schema, dan API ada di [`docs/sdd.md`](docs/sdd.md). Konvensi coding detail ada di [`workflow.md §5`](.agents/workflow.md).
