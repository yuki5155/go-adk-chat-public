# Security Audit: 3 Critical Areas

## 1. ID Token Validation (iss/aud/exp/email_verified/hd)

**`backend/internal/infrastructure/auth/google/validator.go`**

Google's `idtoken.Validate()` library handles **iss**, **aud**, **exp**, and signature verification automatically.

**`backend/internal/application/auth/google_login.go:40-52`** — `email_verified` is properly enforced:

```go
if !oauthUser.EmailVerified {
    return nil, shared.ErrUnverifiedEmail
}
```

| Claim | Status | Notes |
|-------|--------|-------|
| `iss` | Validated | via `idtoken.Validate()` |
| `aud` | Validated | matched against `uc.clientID` |
| `exp` | Validated | via `idtoken.Validate()` |
| `email_verified` | Validated | explicit check in login use case |
| **`hd`** | **Not validated** | Any Google account can login |

**Finding:** `hd` (hosted domain) claim is not checked. If this is intended for enterprise/organization use, any `@gmail.com` account can currently authenticate. Consider adding optional `hd` validation if domain restriction is needed.

---

## 2. Cookie/JWT (SameSite/CSRF/Refresh)

**`backend/internal/presentation/http/handlers/auth_handler.go:156-180`** — Cookie settings:

```go
c.SetCookie("access_token", accessToken, 900, "/", domain, secure, true /*HttpOnly*/)
c.SetCookie("refresh_token", refreshToken, 604800, "/", domain, secure, true /*HttpOnly*/)
```

| Attribute | Status | Notes |
|-----------|--------|-------|
| `HttpOnly` | Set | XSS cannot steal tokens |
| `Secure` | Conditional | `true` in production |
| **`SameSite`** | **Not set** | Gin's `SetCookie()` doesn't set it — browser defaults apply |

**This is the biggest gap.** Gin's `c.SetCookie()` has no `SameSite` parameter. You need to use `http.SetCookie()` directly:

```go
http.SetCookie(c.Writer, &http.Cookie{
    Name:     "access_token",
    Value:    accessToken,
    Path:     "/",
    Domain:   domain,
    MaxAge:   900,
    Secure:   secure,
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode, // <- this is missing
})
```

**CSRF:** The CORS middleware in `middleware/cors.go` accepts `X-CSRF-Token` header but **never validates it**. Combined with the missing `SameSite`, CSRF protection is weak.

**Token Refresh** (`jwt/service.go`): Properly implemented — 15min access / 7-day refresh, signing method validation (prevents algorithm confusion), proper expiration handling. However, **no server-side revocation** exists — logout only clears client cookies, so stolen tokens remain valid until expiry.

**CORS** (`middleware/cors.go`): Allows wildcard `"*"` in `AllowedOrigins` while also setting `Access-Control-Allow-Credentials: true`. Browsers block this combination, but the implementation falls through to using the first allowed origin, which could mask issues.

---

## 3. SSE Endpoint Authorization

**Two SSE paths exist:**

| Path | Auth Method | Role Check |
|------|-------------|------------|
| Gin HTTP (`handlers/chat_handler.go:233`) | `Auth` + `RequireSubscriber` middleware | subscriber/admin/root |
| Lambda (`cmd/lambda/chat-stream/main.go`) | Manual cookie parsing + `validateAuth()` | `isAllowedRole()` check |

Both paths:

- Extract JWT from `access_token` cookie
- Validate signature + expiration
- Enforce role-based access (subscriber/admin/root)
- Bind `userID` from claims to the command (prevents cross-user data access)

**Lambda-specific concern** (`chat-stream/main.go`):

```go
origin := req.Headers["origin"]
if origin == "" {
    origin = "*"  // <- falls back to wildcard
}
```

If `Origin` header is missing, CORS opens to `*`. In practice API Gateway may always forward it, but the fallback is unsafe.

**No SSE-specific leak vectors found** — the stream is POST-based (not GET with `EventSource`), so it can't be opened via `<img>` or `<script>` tags. The `credentials: 'include'` + cookie auth pattern is correct for this design.

---

## Summary of Recommended Fixes (Priority Order)

| Priority | Issue | Fix |
|----------|-------|-----|
| **High** | Missing `SameSite` on auth cookies | Use `http.SetCookie()` with `SameSiteLaxMode` |
| **High** | CORS wildcard fallback in Lambda | Reject requests with no/unknown origin |
| **Medium** | No CSRF token validation | Either enforce `SameSite=Strict` or implement double-submit cookie |
| **Medium** | No token revocation | Add DynamoDB-based blacklist on logout |
| **Low** | Missing `hd` claim validation | Add optional hosted domain check if org-only access needed |
| **Low** | CORS allows `"*"` in config | Remove wildcard support when `Allow-Credentials: true` |
