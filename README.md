# BoltShop: OTP-Based User Login

A checkout form that recognizes returning users by email and lets them log in with the 6-digit code they got at registration.

- **Live app:** https://bolt-login-assigment.vercel.app
- **API:** https://boltapp-api.onrender.com/api/health

> The API runs on Render's free tier, which sleeps after 15 minutes idle: the very first request can take ~30–50 s while it wakes up.

## Architecture

```
┌──────────────┐   HTTPS/JSON    ┌──────────────┐     SQL      ┌──────────────┐
│  frontend/   │ ──────────────▶ │  backend/    │ ───────────▶ │  Postgres    │
│  React + TS  │  Bearer token   │  Go net/http │    pgx       │  (Supabase)  │
│  (Vercel)    │                 │  (Render)    │              │  db/*.sql    │
└──────────────┘                 └──────────────┘              └──────────────┘
```

| Layer    | Tech                                  | Folder            |
|----------|---------------------------------------|-------------------|
| Frontend | React 19, TypeScript, Vite, React Router, React Hook Form + Zod | `frontend/` |
| API      | Go 1.26 stdlib `net/http`, pgx v5, bcrypt | `backend/`    |
| Database | PostgreSQL 16                         | `db/migrations/`  |

## Flows

**Registration** (`/register`): enter email, first name and last name. The API generates a random 6-digit code with `crypto/rand`, stores only its **bcrypt hash**, and returns the plain code **once** for the user to save.

**Recognition & login** (`/`, checkout):
1. The email field is checked on every keystroke with the same regex the API uses.
2. Once the email is complete and well-formed, the page calls `POST /api/recognize` in the background after a 400 ms debounce. In-flight requests are cancelled when the email changes, and the user can keep filling in phone and address.
3. If the email is recognized, a modal asks for the 6-digit code. **Skip** (or Esc, or a click outside the modal) goes back to the form as a guest. The user can reopen the modal from the "Log in with your code" link.
4. `POST /api/login` checks the code against the hash. On a mismatch, the modal shows an error. On a match, the modal closes and "Logged in as First Last" appears at the top of the form.
5. **Delivery details:** the address is collected as house/flat and street, area/landmark (optional), city, state, a 6-digit PIN code and a name for the address (**Home**, **Work** or a custom one via **Other**). Once logged in, the user's **saved phone and addresses from previous orders** are shown as cards (`GET /api/me/saved-details`). The most recent one prefills the form, but only if the user hasn't started typing. Every field stays editable, and "Use a new address" clears them. A first-time user simply fills the form.
6. **Place order** saves the email, phone, address and `user_id` (when logged in) to the `checkouts` table, and shows a receipt with exactly what was stored.

## API

| Method | Path             | Body                                  | Response |
|--------|------------------|---------------------------------------|----------|
| GET    | `/api/health`    | –                                     | `{status}` |
| POST   | `/api/register`  | `{email, first_name, last_name}`      | `201 {user, code}` · `409` if the email already exists |
| POST   | `/api/recognize` | `{email}`                             | `{recognized: bool}` |
| POST   | `/api/login`     | `{email, code}`                       | `{user, token}` · `401` on mismatch · `429` when rate-limited |
| GET    | `/api/me`        | `Authorization: Bearer <token>`       | `{user}` |
| POST   | `/api/checkout`  | `{email, phone, address: {line1, line2, city, state, pincode}}` (+ optional bearer) | `201 {id, user_id, email, phone, address, created_at}` |
| GET    | `/api/me/saved-details` | `Authorization: Bearer <token>` | `{saved: [{phone, address, last_used_at}]}`: distinct, newest first, max 3 |

## Backend architecture

```
backend/
  cmd/server/main.go         startup: env config, wires the layers, graceful shutdown
  internal/
    handler/                 HTTP transport: routes, CORS/logging, JSON in/out, error → status code
    service/                 use cases + API request/response DTOs; orchestrates biz calls
    biz/                     business rules: validation, register, recognize, login, checkout;
                             defines the UserRepo / CheckoutRepo interfaces and domain errors
    entity/                  domain models (User, Checkout): no JSON or DB tags
    repository/              Postgres implementation of the biz repo interfaces (pgx)
    repotest/                in-memory repos for tests
    auth/                    6-digit code + bcrypt, HMAC session tokens, login rate limiter
```

**Request flow:** `handler` → `service` → `biz` → `repository` → Postgres

Dependencies only point inward. `biz` never imports `service`, `handler` or `repository`: it owns the repository interfaces, and `repository` implements them. That means business rules can be tested without HTTP or a database. For example, the checkout use case in `service` combines two biz components: it resolves the optional session to a user (`UserBiz.Authenticate`), then records the order (`CheckoutBiz.Place`).

## Frontend structure

```
frontend/src/
  pages/        Register, Checkout: compose the pieces below
  components/   LoginModal, EmailStatus, SavedDetailsPicker, AddressFields, AddressLabelPicker,
                OrderSummary, OrderReceipt, LoggedInBadge, Field (label + input + error)
  hooks/        useEmailRecognition (debounced background lookup)
                useDeliveryChoice (saved vs. new address, keeping a typed draft)
  auth/         session context: token in sessionStorage, restored on refresh
  api.ts        typed API client        validation.ts   Zod schemas + LIMITS shared by the forms
```

## Tests

**Frontend:** `npm test` in `frontend/` (Vitest) covers the email check used for recognition, the checkout rules (PIN code, phone, required fields, trimming) and address formatting.

**Backend:** `go test -race ./...` in `backend/`:
- **`auth`:** codes are always 6 digits and random; the bcrypt hash matches only the right code; tokens round-trip and reject tampering, a wrong secret and expiry; the limiter blocks after N failures, its window slides, reset works, expired entries are swept so memory stays bounded, and it's safe under concurrent use.
- **`biz`:** validation tables; register normalizes input, stores only the hash and rejects duplicates case-insensitively; recognize; login with a wrong code, an unknown email or the right code; authenticate; rate limiting; checkout normalization and validation (incl. PIN code); saved details are distinct, newest first and per user.
- **`service`:** checkout links the order to the user for a valid token, and saves it as a guest for no token or a bad one; saved details require a login and are empty before the first order.
- **`handler`:** a full HTTP **login-flow test** (register → 409 duplicate → recognize → 401 wrong code → login → `/me` → checkout linked vs. guest → saved details); 429 rate limiting; 400s for bad input; CORS allow/deny.

## Design decisions

- **The code is hashed (bcrypt), never stored in plain text.** It works like a password, so a database leak doesn't expose login codes. As a result, the code can only be shown once.
- **Brute-force protection.** A 6-digit code has only 1M possibilities, so login allows 5 failed attempts per email per 15 minutes (in-memory limiter).
- **Bearer token instead of a cookie.** The frontend (vercel.app) and API (onrender.com) are on different sites. Browsers increasingly block third-party cookies, so the API returns an HMAC-SHA256 signed token (`userId.expiry.sig`, 24h TTL). The frontend keeps it in `sessionStorage`, so the session survives a refresh but ends when the tab closes.
- **Recognition returns no personal data.** `/recognize` only returns a boolean. The user's name is revealed only after the correct code is entered.
- **Row Level Security is enabled on both tables, with no policies.** This blocks Supabase's auto-generated public REST API (the anon/authenticated keys) from reading user data. The Go API connects as the table owner, which RLS doesn't restrict, so it remains the only way in.
- **Emails are normalized** (trimmed and lowercased) before they're stored or looked up.
- **Validation runs on both sides.** The frontend uses Zod schemas (`frontend/src/validation.ts`) with React Hook Form for UX, and the Go API validates again with the same rules because it's the trust boundary.

## Running locally

Requires Go 1.26+, Node 20+ and a Postgres database: a local install, or simply your Supabase project.

```bash
# 1. Database: run db/migrations/001_users.sql then 002_checkouts.sql against your Postgres
psql "$DATABASE_URL" -f db/migrations/001_users.sql -f db/migrations/002_checkouts.sql

# 2. API → http://localhost:8080
cd backend
cp .env.example .env            # set DATABASE_URL
export $(cat .env | xargs) && go run ./cmd/server
go test -race ./...            # API tests (no database needed)

# 3. Frontend → http://localhost:5173
cd frontend
cp .env.example .env.local
npm install && npm run dev
npm test                        # frontend unit tests
```

## Deployment

1. **Supabase (DB):** create a project, open **SQL Editor**, and run `db/migrations/001_users.sql` then `002_checkouts.sql`. Copy the connection string from **Connect → Session pooler** (IPv4-compatible).
2. **Render (API):** click **New → Blueprint** and choose this repo (it uses `render.yaml`), or create a **Web Service** with root dir `backend`, runtime **Go**, build command `go build -o app ./cmd/server` and start command `./app`. Set these env vars:
   - `DATABASE_URL`: the Supabase string (append `?sslmode=require` if it's missing)
   - `SESSION_SECRET`: a long random string
   - `ALLOWED_ORIGIN`: your Vercel URL, e.g. `https://boltapp.vercel.app` (comma-separate multiple values)
3. **Vercel (frontend):** import the repo, set **Root Directory = `frontend`** (framework: Vite), and add the env var `VITE_API_URL=https://YOUR-API.onrender.com`. Deploy.

> Render's free tier sleeps after 15 minutes of inactivity, so the first request can take about 30–50 s.

## Database schema

See [`db/migrations/`](db/migrations/): `users` (with `login_code_hash`) and `checkouts` (a nullable `user_id` FK for guest checkouts, and a structured address with a PIN-code check; indexed on `(user_id, created_at)` for the saved-details lookup).
