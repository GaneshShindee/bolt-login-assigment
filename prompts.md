prompts.md — BoltShop OTP-Based User Login

This file documents the LLM prompts used during development of the BoltShop take-home assignment.

Phase 1 — Understand the Assignment

Prompt 1 — Analyze the take-home assignment

I need to build a take-home assignment for an OTP Based User Login system.

Requirements:

Registration flow:

Collect email, first name, and last name.

Register the user.

Generate a random 6-digit numeric code.

Display the code after registration so the user can use it later for login.

Recognition and checkout flow:

Checkout collects email, phone number, and shipping address.

Validate the email in real time while typing.

Once the email is complete and valid, run a recognition check in the background.

If the email belongs to a registered user, show a modal requesting the 6-digit code.

Allow the user to skip login and continue checkout.

Validate the code.

On success, log the user in and show their name.

Save checkout information to a database table.

No real payment processing is needed.

Deployment:

The application must be publicly accessible.

The repository must contain the complete source code.

Database schema must be committed as SQL files.

A prompts.md file must list the LLM prompts used during development.

Architecture:

Frontend, API, and database must be distinct layers.

Recommended stack:

TypeScript

React

Go

PostgreSQL

Help me create an implementation plan and identify the major technical risks.

Phase 2 — Architecture

Prompt 2 — Design the three-layer architecture

Design a clean architecture for this assignment using:

React + TypeScript + Vite for the frontend

Go for the API

PostgreSQL for persistence

The architecture should clearly separate:

Browser
→ Frontend
→ HTTP/JSON API
→ Business logic
→ PostgreSQL

Keep the implementation simple enough for a two-day take-home assignment but structured enough to demonstrate good engineering practices.

Prompt 3 — Define backend layering

Organize the Go backend using:

handler

service

biz

repository

entity

test fakes

Explain the responsibility of each layer and how dependencies should flow.

The business layer should define repository interfaces so business logic does not depend directly on PostgreSQL implementation details.

Prompt 4 — Define API endpoints

Design REST endpoints for:

user registration

email recognition

login-code verification

checkout submission

For every endpoint specify:

HTTP method

URL

request JSON

response JSON

validation

possible HTTP status codes

security considerations

Phase 3 — Database

Prompt 5 — Design the PostgreSQL schema

Create a PostgreSQL schema for:

users

Fields:

ID

email

first name

last name

login code hash

timestamps

Requirements:

email must be unique

email should be normalized to lowercase

login code should never be stored in plaintext

checkouts

Fields:

ID

user_id, nullable for guest checkout

email

mobile number

address fields

address label

created timestamp

Add appropriate:

primary keys

foreign keys

unique constraints

CHECK constraints

indexes

The schema must be checked into the repository as .sql files.

Prompt 6 — Database security review

Review the PostgreSQL schema for:

SQL injection concerns

invalid phone numbers

invalid email values

duplicate users

referential integrity

nullable guest user IDs

indexing for recent user addresses

row-level security

Use PostgreSQL constraints as the final data-integrity boundary.

Phase 4 — Registration

Prompt 7 — Implement secure 6-digit code generation

Implement registration in Go.

Requirements:

Validate email, first name, and last name.

Normalize email to lowercase.

Detect duplicate registrations.

Generate a random 6-digit numeric code using crypto/rand.

Hash the code with bcrypt.

Store only the bcrypt hash in PostgreSQL.

Return the code once in the registration response so it can be shown to the user for this assignment.

Never expose the stored hash to the frontend.

Use clear domain/business errors for validation and duplicates.

Prompt 8 — Implement the registration frontend

Build a React registration page.

Fields:

Email

First name

Last name

Requirements:

client-side validation

loading state

API error handling

successful registration state

clearly display the generated six-digit code

explain that the code is required for future login

Use React Hook Form and Zod where appropriate.

Phase 5 — Email Recognition

Prompt 9 — Real-time email validation

Implement a React email field that validates while the user types.

Behavior:

Incomplete email → show an appropriate state such as "Keep typing…"

Complete valid email → start recognition in the background

Do not block the checkout form

Debounce recognition by approximately 400 ms

Do not send an API request for every keystroke

The user should be able to continue entering phone and address information while recognition runs.

Prompt 10 — Prevent stale recognition responses

Implement robust cancellation for the email recognition request.

Requirements:

When the email changes before a request finishes, cancel the previous request using AbortController.

A response belonging to an old email must never cause the login modal to appear.

Recognition state should correspond only to the current email.

Implement this as a reusable hook named useEmailRecognition.

Prompt 11 — Recognition API privacy

Implement POST /api/recognize.

Input:

email

Output should contain only:

recognized: true|false

Do not return:

first name

last name

login code

other personal information

The endpoint should normalize the email and safely query the database.

Phase 6 — Login Modal

Prompt 12 — Build the code verification modal

When a recognized email is detected, show a login modal.

The modal must:

request the six-digit code

accept numeric input

show validation errors inside the modal

provide a login/verify button

provide a skip option

allow the user to continue checkout as a guest

Do not erase any checkout form data when the modal opens or closes.

Prompt 13 — Verify code securely

Implement the login/code verification business logic.

Input:

email

six-digit code

Requirements:

look up the user

compare the provided code with the bcrypt hash

authenticate only if the code matches

do not return the hash

use the same public error for unknown users and wrong codes where practical to reduce account enumeration

support rate limiting against brute-force attempts

Prompt 14 — Rate-limit login attempts

Implement brute-force protection for six-digit login codes.

Requirements:

limit failed attempts per email

allow at most 5 failed attempts within 15 minutes

return HTTP 429 when blocked

support concurrent requests safely

reset/expire entries appropriately

Explain why six-digit codes require rate limiting.

Phase 7 — Session Authentication

Prompt 15 — Design a session token

After successful code verification, issue a signed session token.

Requirements:

HMAC-SHA256 signing

expiration time

claims should identify the authenticated user

tampering must invalidate the token

backend must reject expired or malformed tokens

Use a token-based approach rather than a cookie because the frontend and backend are hosted on different origins.

Prompt 16 — Frontend session storage

Implement frontend authentication state.

Requirements:

store the session token in sessionStorage

restore authentication after a page refresh

authentication should end when the browser tab/session is closed

send the token to the API using an authorization header

do not store passwords or login codes in the frontend

Prompt 17 — Show authenticated user identity

After successful login:

close the code modal

show "Logged in as <First Last>" at the top of checkout

preserve all checkout fields already entered

allow the user to continue the checkout normally

The frontend should derive the visible name only after successful authentication.

Phase 8 — Checkout

Prompt 18 — Design the checkout form

Build a responsive checkout form containing:

email

10-digit mobile number

shipping address

address name/label

Address should be structured into useful fields such as:

house/street

area

city

state

PIN code

label: Home / Work / Other

The user should be able to submit as either:

authenticated customer

guest

Prompt 19 — Validate checkout data

Implement frontend and backend validation.

Rules should include:

valid email format

mobile must contain exactly 10 digits

mobile must begin with 6, 7, 8, or 9

required address fields must be present

PIN code must be valid

address label must be allowed

The frontend provides immediate feedback, but the API remains the trust boundary.

Prompt 20 — Persist checkout

Implement checkout submission.

Requirements:

validate the request

determine the authenticated user from the session token when available

save user_id for authenticated users

allow user_id = NULL for guest checkout

store all checkout fields in PostgreSQL

return a success response

do not implement payment processing

Phase 9 — Saved Addresses

Prompt 21 — Add saved addresses for logged-in users

Enhance checkout for returning users.

Requirements:

load recent addresses belonging to the logged-in user

show them as selectable address cards

select the newest address by default

allow the user to place another order quickly

add an index on (user_id, created_at) for efficient retrieval

Guest orders must not appear as a saved address for another user.

Prompt 22 — Preserve a newly typed address

Handle the case where a user enters a new address before login.

Requirements:

preserve the typed address when the login modal opens

after successful login, show the typed address as "Your new address"

switching between saved addresses and the new address must not lose unsaved form data

skipping login must also preserve the typed address

Phase 10 — Backend Implementation

Prompt 23 — Implement handler layer

Implement handler/server.go and the HTTP routes.

Responsibilities:

decode request JSON

validate basic transport input

call the appropriate service

translate domain errors into HTTP status codes

encode JSON responses

configure CORS

expose registration, recognition, login, and checkout routes

Expected status mappings should include:

400 Bad Request

401 Unauthorized

409 Conflict

429 Too Many Requests

500 Internal Server Error

Do not expose database internals in HTTP responses.

Prompt 24 — Implement service layer

Create one service use case per major operation.

For example:

Register

Recognize

Login / VerifyCode

Checkout

The service layer should:

define request/response structures

coordinate business components

remain independent of PostgreSQL implementation details

CheckoutService.Checkout should resolve the authenticated user using the session token before saving the checkout.

Prompt 25 — Implement business layer

Implement the domain rules in biz.

Responsibilities include:

input validation

email normalization

code generation

password/code hashing and verification

rate limiting

token creation

token validation

authorization rules

The business layer should depend on repository interfaces rather than concrete PostgreSQL code.

Prompt 26 — Implement repository interfaces

Define repository interfaces such as:

UserRepo

CheckoutRepo

These interfaces should be declared in the business/domain side so the business logic can be tested independently from the real database.

Implement the concrete repositories using pgx.

Prompt 27 — Add in-memory repository fakes

Create in-memory repository implementations under a test package such as repotest/.

The goal is to run business logic tests without requiring a real PostgreSQL database.

Test behavior such as:

user registration

lookup

code verification

checkout persistence

Phase 11 — Security Review

Prompt 28 — Perform a security review

Review the complete application for:

SQL injection

plaintext code storage

code leakage

account enumeration

brute-force attacks

token forgery

expired token acceptance

CORS misconfiguration

insecure secrets

overexposed API responses

client-side-only validation

unsafe database access

Recommend practical fixes that are appropriate for a small take-home assignment.

Prompt 29 — Review authentication threat model

Analyze the authentication design.

Important facts:

code has 1,000,000 possibilities

code is stored as bcrypt hash

five failures per 15 minutes should block an email

successful login returns a signed token

token must expire

unknown email and wrong code should use the same public error where possible

Explain remaining limitations and practical improvements for production.

Phase 12 — CORS and Deployment Constraints

Prompt 30 — Design CORS configuration

The frontend is hosted on Vercel and the Go API is hosted separately.

Configure CORS so that:

only the known Vercel frontend origin is allowed

API methods and headers required by the app are permitted

credentials are not unnecessarily enabled

preflight requests work correctly

Explain why an authorization header with a signed token is useful in this cross-origin deployment.

Prompt 31 — Prepare Vercel frontend deployment

Prepare the frontend for Vercel.

Requirements:

Vite production build

VITE_API_URL environment variable

production API URL configured

no secrets included in frontend source

correct routing/build configuration

Prompt 32 — Prepare Render backend deployment

Prepare the Go API for Render.

Requirements:

production build/start command

DATABASE_URL

SESSION_SECRET

ALLOWED_ORIGIN

health check if appropriate

production CORS configuration

Keep deployment compatible with Render's free tier.

Document that free Render can sleep after inactivity and the first request after sleeping may take longer.

Phase 13 — Testing

Prompt 33 — Write Go unit tests

Create Go unit tests for:

six-digit code generation

token generation

token tampering

token expiration

email validation

phone validation

business rules

rate limiter

concurrent access to the rate limiter

Run:

go test -race ./...

Prompt 34 — Write frontend validation tests

Create Vitest tests for:

email validation

phone number validation

address formatting

checkout validation

other reusable validation rules

Run:

npm test

Prompt 35 — Add end-to-end HTTP login flow test

Create an HTTP-level end-to-end test covering:

register a user

obtain the generated login code

recognize the email

attempt login with an invalid code

verify that the error is returned correctly

login with the valid code

obtain a signed session token

submit a checkout using that token

verify that the checkout is associated with the user

Also test guest checkout where practical.

Prompt 36 — Test rate-limit concurrency

Stress the rate limiter using concurrent test requests.

Verify that:

failed attempts cannot race past the five-attempt threshold

concurrent requests are safe

blocked users receive 429

expired windows are eventually reusable

Run the race detector to catch synchronization bugs.

Phase 14 — UX and Responsive Design

Prompt 37 — Improve checkout UX

Review the checkout experience.

Requirements:

email recognition should happen quietly in the background

the user should not have to wait for recognition to fill the other fields

modal should be clear but not disruptive

skip should be obvious

existing form data must never be lost

authenticated user's name should be visible

saved addresses should be easy to select

Keep the product behavior unchanged while improving usability.