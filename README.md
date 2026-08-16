# Martabak Tracker

A web application for managing and tracking martabak orders in real time — built with Go, Gin, and GORM. Customers can place orders and track their status live without manual refreshing, while admins can manage all orders from a single dashboard.

## ✨ Features

- **Order form** — customers fill in their details and can add multiple martabak items in a single order (type, size, extra notes)
- **Real-time order tracking** — tracking page with a status stepper (Order placed → Preparing → Baking → Quality Check → Ready), updated automatically via Server-Sent Events (SSE), no manual refresh needed
- **Order search** — customers can look up a previous order by Order ID directly from the home page
- **Mark as done** — the "Done" button becomes active automatically once the status reaches "Ready", removing the order once the customer confirms completion
- **Admin dashboard** — protected login, view all orders, update order status, delete orders, real-time notification when a new order comes in
- **Responsive** — the layout adapts for both mobile and desktop devices

## 🛠️ Tech Stack

| Category | Technology |
|---|---|
| Language | Go 1.26.5 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) v1.12.0 |
| ORM | [GORM](https://gorm.io/) v1.31.2 |
| Database | [Turso](https://turso.tech/) (libSQL, SQLite-compatible) — via [ekristen/gorm-libsql](https://github.com/ekristen/gorm-libsql) |
| Sessions | [gin-contrib/sessions](https://github.com/gin-contrib/sessions) v1.1.0 (stored in the database) |
| Input validation | [go-playground/validator](https://github.com/go-playground/validator) v10 |
| Authentication | bcrypt ([golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)) |
| Order ID generation | [teris-io/shortid](https://github.com/teris-io/shortid) |
| Environment variables | [joho/godotenv](https://github.com/joho/godotenv) (local development) |
| Templates | `html/template` + `embed.FS` (bundled directly into the binary) |
| Styling | Tailwind CSS |
| Live updates | Server-Sent Events (SSE) |
| Hot reload (dev) | [Air](https://github.com/air-verse/air) |

## 🔌 How the API Works

All routes are defined in `routes.go` using the Gin router. Routes are split into two groups: **public** (accessible to anyone) and **admin** (protected by `AuthMiddleware`, login required — if the session is invalid, the request is automatically redirected to `/login`).

### Public Endpoints (Customer)

| Method | Endpoint | Purpose | Handler |
|---|---|---|---|
| `GET` | `/` | Displays the new-order form | `ServeNewOrderForm` |
| `POST` | `/new-order` | Creates a new order, generates an Order ID, redirects to the tracking page | `HandleNewOrderPost` |
| `GET` | `/customer/:id` | Order status tracking page based on Order ID | `serveCustomer` |
| `POST` | `/customer/:id/done` | Marks an order as done (only succeeds if the status is already "Ready"), then deletes it | `HandleOrderDone` |
| `GET` | `/notifications?orderId=xxx` | **SSE** (Server-Sent Events) connection — streams live status updates for a specific order to the customer's browser | `notificationHandler` |
| `GET` | `/static/*` | Serves static image/asset files (bundled directly into the binary via `embed.FS`) | — |

### Authentication Endpoints

| Method | Endpoint | Purpose | Handler |
|---|---|---|---|
| `GET` | `/login` | Displays the admin login form | `HandleLoginGet` |
| `POST` | `/login` | Verifies username & password (checked via bcrypt), creates a session | `HandleLoginPost` |
| `POST` | `/logout` | Clears the session, redirects to `/login` | `HandleLogout` |

### Admin Endpoints (login required)

| Method | Endpoint | Purpose | Handler |
|---|---|---|---|
| `GET` | `/admin` | Dashboard — displays all incoming orders | `ServeAdminDashboard` |
| `POST` | `/admin/order/:id/update` | Updates an order's status (Order placed → Preparing → etc.) | `HandleOrderPut` |
| `POST` | `/admin/order/:id/delete` | Manually deletes an order | `HandleOrderDelete` |
| `GET` | `/admin/notifications` | **SSE** connection — notifies the admin live whenever a new order comes in | `adminNotificationHandler` |

### Core Workflows

**1. Customer places an order**
```
Customer fills out the form on "/"
  → POST /new-order
  → Server validates input (name, phone, address, martabak type & size)
  → Saved to the database (Order ID generated via shortid)
  → Redirect to /customer/{orderID}
  → Admin notified via SSE ("new order")
```

**2. Live tracking without refresh (SSE)**

When a customer opens `/customer/:id`, the browser opens a persistent connection to `/notifications?orderId=xxx`. This connection is registered with the `NotificationManager` (an in-memory pub-sub built on Go channels). As soon as the admin updates the order's status from the dashboard, the server immediately pushes a message to every connection "subscribed" to that Order ID — the customer's page updates automatically without a manual reload.

**3. Admin authentication**
```
POST /login
  → Password checked with bcrypt.CompareHashAndPassword
  → On match, userID & username are stored in the session (persisted in the database via gorm-session-store)
  → Every request to /admin/* is checked by AuthMiddleware:
      - Verify the session has data
      - Verify the user still exists/is valid in the database
      - If invalid → redirect to /login
```

**4. Input validation**

Every order-creation request is validated on two layers:
- **Frontend** — JavaScript checks for empty fields/format before submission (instant feedback for better UX)
- **Backend** — `go-playground/validator` via struct tags (`binding:"required,min=..."`) plus custom validators (`valid_martabak_type`, `valid_martabak_size`) that ensure the martabak type & size match the allowed lists in `models.MartabakTypes`/`MartabakSizes`

## 📁 Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point (package main)
│       ├── main.go
│       ├── admin.go       # Admin dashboard & login handlers
│       ├── customer.go    # Customer order & tracking handlers
│       ├── events.go      # SSE handler for live updates
│       ├── handlers.go    # Handler struct & dependency injection
│       ├── middleware.go  # Auth middleware
│       ├── notifications.go # In-memory pub-sub notification manager
│       ├── routes.go      # All route definitions
│       ├── utils.go       # Config, session store, template loader
│       └── validators.go  # Custom validators (martabak type & size)
├── internal/
│   └── models/
│       ├── models.go     # Database connection & migrations
│       ├── order.go      # Order & OrderItem models
│       └── user.go       # User model & authentication
├── templates/
│   ├── *.tmpl            # All HTML pages
│   ├── embed.go          # Embeds templates & static assets into the binary
│   └── static/images/    # Image assets
├── .air.toml             # Development hot-reload configuration
├── go.mod
└── go.sum
```

## 🚀 Running Locally

### Prerequisites
- Go 1.26.5 or newer
- A [Turso](https://turso.tech/) account (free) for the database, or use a local SQLite file for quick development

### Installation

```bash
git clone https://github.com/karikuzoo/martabak-tracker.git
cd martabak-tracker
go mod tidy
```

Create a `.env` file in the project root:

```env
DATABASE_URL=libsql://<your-database-name>.turso.io?authToken=<your-token>
SESSION_SECRET_KEY=<a-long-random-string>
PORT=8080
GIN_MODE=debug
```

Run with hot-reload:

```bash
air
```

Or without hot-reload:

```bash
go run ./cmd/server
```

Open `http://localhost:8080` in your browser.

## 🌐 Deployment

This project is deployed using:
- **[Vercel](https://vercel.com)** — application hosting (Go Framework Preset, supports long-running servers)
- **[Turso](https://turso.tech/)** — cloud database (SQLite-compatible via libSQL)

Environment variables to set in the Vercel dashboard:

| Variable | Description |
|---|---|
| `DATABASE_URL` | Full Turso connection string including auth token |
| `SESSION_SECRET_KEY` | Random string for session encryption (do not use the default) |
| `GIN_MODE` | `release` for production |

## Login Access

For User

https://martabak-tracker.vercel.app/

For Admin

https://martabak-tracker.vercel.app/admin

Username: admin

Password: password1234
