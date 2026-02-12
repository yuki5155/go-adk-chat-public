# Security Checklist

Tracking issue: [#12](https://github.com/yuki5155/go-adk-chat/issues/12)

---

## 1. Cookie Configuration

- [x] `HttpOnly` set on `access_token` cookie
- [x] `HttpOnly` set on `refresh_token` cookie
- [x] `Secure` flag enabled in production
- [x] `Path` set to `/`
- [x] `Domain` configurable via `COOKIE_DOMAIN`
- [ ] **`SameSite=Lax` (or `Strict`) explicitly set on `access_token`**
- [ ] **`SameSite=Lax` (or `Strict`) explicitly set on `refresh_token`**
- [ ] Migrate from `gin.SetCookie()` to `http.SetCookie()` for `SameSite` support
- [ ] Consistency between `Max-Age` and `Expires` (countermeasure for older browser differences)
- [ ] Consider `__Host-` / `__Secure-` cookie prefix
  - `__Host-` requires: `Secure`, no `Domain`, `Path=/` — incompatible with current `COOKIE_DOMAIN` cross-subdomain policy
  - If "same host only" is acceptable, `__Host-` is a strong hardening measure
  - Otherwise `__Secure-` (requires `Secure` only) is a viable middle ground

**Files:**
- `backend/internal/presentation/http/handlers/auth_handler.go`

---

## 2. CSRF Protection

- [x] CORS middleware present
- [x] `X-CSRF-Token` header listed in `Access-Control-Allow-Headers`
- [ ] **CSRF token generation (server-side)**
- [ ] **CSRF token validation on state-changing requests (POST/PUT/DELETE)**
- [ ] Choose approach: Double Submit Cookie / `SameSite=Strict` + Origin check
- [ ] Minimum `Origin` / `Referer` header check (effective as insurance even with Double Submit)
  - Reject if `Origin` is absent on state-changing operations (POST/PUT/DELETE)

**Files:**
- `backend/internal/presentation/http/middleware/cors.go`

---

## 3. Origin / CORS Policy

- [x] Origin-based CORS filtering implemented
- [x] `Access-Control-Allow-Credentials: true` set
- [ ] **Remove wildcard (`*`) from `AllowedOrigins` in production config**
- [ ] **Lambda: reject requests with missing/unknown `Origin` (currently falls back to `*`)**
- [ ] Validate `Origin` header against explicit allowlist in all environments
- [ ] Fix Preflight (OPTIONS) behavior — confirm disallowed origins do not receive `200` with permissive headers
- [ ] Add `Vary: Origin` response header (prevents cache mixing when CloudFront/proxy is involved)

**Files:**
- `backend/internal/presentation/http/middleware/cors.go`
- `backend/cmd/lambda/chat-stream/main.go`

---

## 4. ID Token Validation (Google OAuth)

- [x] `iss` (issuer) validated via `idtoken.Validate()`
- [x] `aud` (audience) validated against Google Client ID
- [x] `exp` (expiration) validated via `idtoken.Validate()`
- [x] `email_verified` enforced in login use case
- [x] Signature verified using Google's public keys
- [ ] **`hd` (hosted domain) claim check — if corporate-only access is required**
- [ ] **Server-side email domain allowlist (primary check, `hd` is supplementary)**
- [ ] Clarify `nonce` handling policy (use it / don't use it)
  - Depends on "Sign-in with Google" (GIS) flow choice
  - Not mandatory, but documenting the design decision explicitly strengthens the security posture

**Files:**
- `backend/internal/infrastructure/auth/google/validator.go`
- `backend/internal/application/auth/google_login.go`

---

## 5. JWT Token Security

- [x] HS256 (HMAC-SHA256) signing
- [x] Signing method validation (prevents algorithm confusion)
- [x] `exp` claim set and validated
- [x] `iat` and `nbf` claims set
- [x] `iss` claim set (`go-google-auth`)
- [x] `sub` claim set to user ID
- [x] Token type differentiation (`access` vs `refresh`)
- [ ] **JWT secret sourced from AWS Secrets Manager in production**
- [ ] **Minimum secret length enforced (>=256 bits)**
- [ ] Key rotation strategy (`kid` header can conceptually exist even with HMAC)
  - Even if not implemented immediately, documenting "rotate keys via X process in the future" makes operations easier

**Files:**
- `backend/internal/infrastructure/auth/jwt/service.go`

---

## 6. Token Refresh & Revocation

- [x] Access token: 15-minute expiry
- [x] Refresh token: 7-day expiry
- [x] Refresh endpoint validates token type
- [x] Expired refresh token clears cookies
- [ ] **Refresh token rotation (issue new refresh token on each use)**
- [ ] **Reuse detection (revoke token family if rotated token is replayed)**
- [ ] **Server-side token blacklist (DynamoDB with TTL)**
- [ ] **Logout invalidates tokens server-side (not just cookie clear)**
- [ ] Per-device/session refresh token management (token family design)
  - Foundation for "logout all devices at once" capability
  - Each device gets its own token family; revoking the family revokes the device

**Files:**
- `backend/internal/infrastructure/auth/jwt/service.go`
- `backend/internal/application/auth/refresh_token.go`
- `backend/internal/application/auth/logout.go`

---

## 7. SSE Endpoint Authorization

- [x] Gin: `Auth` middleware validates JWT from cookie
- [x] Gin: `RequireSubscriber` middleware enforces role
- [x] Lambda: manual cookie parsing + JWT validation
- [x] Lambda: role check (`subscriber`/`admin`/`root`)
- [x] User ID bound from claims (prevents cross-user access)
- [x] POST-based SSE (not vulnerable to `<img>`/`<script>` tag injection)
- [ ] **Lambda: fix CORS origin fallback (see #3 above)**

**Files:**
- `backend/internal/presentation/http/handlers/chat_handler.go`
- `backend/internal/presentation/http/middleware/auth.go`
- `backend/cmd/lambda/chat-stream/main.go`

---

## 8. General

- [x] Secrets loaded from environment variables
- [x] Production config uses HTTPS
- [x] Error messages do not leak internal details to client
- [ ] **Rate limiting on auth endpoints (login, refresh)**
  - At minimum: API Gateway usage plan / WAF / per-route throttling
- [ ] **Audit logging for authentication events**
- [ ] **Dependency vulnerability scanning (e.g., `govulncheck`)**
- [ ] PII masking in logs (email, subject, token values)
  - Define what "don't leak" means concretely — reduces accidental exposure
  - Ensure `log.Printf` calls in auth flows do not output raw tokens or email addresses
