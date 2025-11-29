# Module `users`

Quản lý thông tin user và trả về profile sau khi đã xác thực qua Keycloak.

## Routes
- `GET /api/users/me` (cần Keycloak middleware)

## Request/Response
- Me: header `Authorization: Bearer <Keycloak access token>` → 200 `UserResponse`

`UserResponse`:
```json
{"id": "...uuid", "email": "a@b.com", "created_at": "..."}
```

## Logic (controllers)
- `GetMe`: lấy `userID` từ context (Keycloak middleware đặt vào sau khi verify token và map user) → `GetUserByID`.

## Phụ thuộc
- DB queries: `GetUserByID` (từ `db/sqlc`).
- Middleware: Keycloak middleware bắt buộc cho `/me`.

## Lỗi chuẩn (JSON)
`{"error": "missing token"}`, `{"error": "invalid token"}`, `{"error": "cannot fetch user"}`...
