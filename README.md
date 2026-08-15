# 🥞 Martabak Tracker

Aplikasi web untuk mengelola dan melacak pesanan martabak secara real-time — dibangun dengan Go, Gin, dan GORM. Customer bisa memesan, melacak status pesanan secara live tanpa refresh manual, dan admin bisa mengelola semua pesanan dari satu dashboard.

## ✨ Fitur

- **Form pemesanan** — customer mengisi data diri + bisa menambahkan beberapa martabak sekaligus dalam satu pesanan (variant, ukuran, catatan tambahan)
- **Lacak pesanan real-time** — halaman tracking dengan status stepper (Order placed → Preparing → Baking → Quality Check → Ready), update otomatis lewat Server-Sent Events (SSE) tanpa perlu refresh
- **Cari pesanan** — customer bisa mencari pesanan lama lewat Order ID langsung dari halaman utama
- **Tandai selesai** — tombol "Done" otomatis aktif saat status sudah "Ready", menghapus pesanan setelah customer konfirmasi selesai
- **Dashboard admin** — login terproteksi, lihat semua pesanan, update status, hapus pesanan, notifikasi real-time saat ada pesanan baru masuk
- **Responsive** — tampilan menyesuaikan untuk perangkat mobile maupun desktop

## 🛠️ Tech Stack

| Kategori | Teknologi |
|---|---|
| Bahasa | Go 1.26.5 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) v1.12.0 |
| ORM | [GORM](https://gorm.io/) v1.31.2 |
| Database | [Turso](https://turso.tech/) (libSQL, kompatibel SQLite) — via [ekristen/gorm-libsql](https://github.com/ekristen/gorm-libsql) |
| Session | [gin-contrib/sessions](https://github.com/gin-contrib/sessions) v1.1.0 (disimpan di database) |
| Validasi input | [go-playground/validator](https://github.com/go-playground/validator) v10 |
| Autentikasi | bcrypt ([golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)) |
| ID pesanan | [teris-io/shortid](https://github.com/teris-io/shortid) |
| Environment variable | [joho/godotenv](https://github.com/joho/godotenv) (development lokal) |
| Template | `html/template` + `embed.FS` (ter-bundle langsung ke binary) |
| Styling | Tailwind CSS |
| Live update | Server-Sent Events (SSE) |
| Hot reload (dev) | [Air](https://github.com/air-verse/air) |

## 📁 Struktur Project

```
.
├── cmd/
│   └── server/          # Entry point aplikasi (package main)
│       ├── main.go
│       ├── admin.go       # Handler dashboard & login admin
│       ├── customer.go    # Handler order & tracking customer
│       ├── events.go      # SSE handler untuk live update
│       ├── handlers.go    # Struct Handler & dependency injection
│       ├── middleware.go  # Auth middleware
│       ├── notifications.go # Notification manager (pub-sub in-memory)
│       ├── routes.go      # Definisi semua route
│       ├── utils.go       # Config, session store, template loader
│       └── validators.go  # Custom validator (tipe & ukuran martabak)
├── internal/
│   └── models/
│       ├── models.go     # Koneksi database & migrasi
│       ├── order.go      # Model Order & OrderItem
│       └── user.go       # Model User & autentikasi
├── templates/
│   ├── *.tmpl            # Semua halaman HTML
│   ├── embed.go          # Embed template & static assets ke binary
│   └── static/images/    # Aset gambar
├── .air.toml             # Konfigurasi hot-reload development
├── go.mod
└── go.sum
```

## 🚀 Menjalankan di Lokal

### Prasyarat
- Go 1.26.5 atau lebih baru
- Akun [Turso](https://turso.tech/) (gratis) untuk database, atau bisa gunakan file SQLite lokal untuk development cepat

### Langkah instalasi

```bash
git clone https://github.com/karikuzoo/martabak-tracker.git
cd martabak-tracker
go mod tidy
```

Buat file `.env` di root project:

```env
DATABASE_URL=libsql://<nama-database-kamu>.turso.io?authToken=<token-kamu>
SESSION_SECRET_KEY=<string-random-yang-panjang>
PORT=8080
GIN_MODE=debug
```

Jalankan dengan hot-reload:

```bash
air
```

Atau tanpa hot-reload:

```bash
go run ./cmd/server
```

Buka `http://localhost:8080` di browser.

## 🌐 Deployment

Project ini di-deploy menggunakan:
- **[Vercel](https://vercel.com)** — hosting aplikasi (Go Framework Preset, mendukung server long-running)
- **[Turso](https://turso.tech/)** — database cloud (SQLite-compatible via libSQL)

Environment variable yang perlu di-set di dashboard Vercel:

| Variable | Keterangan |
|---|---|
| `DATABASE_URL` | Connection string Turso lengkap dengan auth token |
| `SESSION_SECRET_KEY` | String random untuk enkripsi session (jangan pakai default) |
| `GIN_MODE` | `release` untuk production |

## 📄 Lisensi

Belum ditentukan.
