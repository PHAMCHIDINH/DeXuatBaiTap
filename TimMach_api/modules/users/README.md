# Module `users`

Quản lý tài khoản và xác thực (JWT).

## Routes
- `POST /api/users/register` (public)
- `POST /api/users/login` (public)
- `GET /api/users/me` (cần Auth middleware)

## Request/Response
- Register: body `{"email","password"}` → 201 `{"user": UserResponse,"token":string}`
- Login: body `{"email","password"}` → 200 `{"access_token":string,"user":UserResponse}`
- Me: header `Authorization: Bearer <jwt>` → 200 `UserResponse`

`UserResponse`:
```json
{"id": "...uuid", "email": "a@b.com", "created_at": "..."}
```

## Logic (controllers)
- `Register`: validate email/password → check trùng email (sqlc `GetUserByEmail`) → hash bcrypt → `CreateUser` → tạo JWT nếu có TokenService.
- `Login`: tìm user theo email → bcrypt check password → tạo JWT.
- `GetMe`: lấy `userID` từ context (middleware đặt vào) → `GetUserByID`.

## Phụ thuộc
- DB queries: `CreateUser`, `GetUserByEmail`, `GetUserByID` (từ `db/sqlc`).
- TokenService (mặc định `JWTMaker` HMAC).
- Middleware: yêu cầu Auth cho `/me`.

## Lỗi chuẩn (JSON)
`{"error": "invalid request body"}`, `{"error": "invalid credentials"}`, `{"error": "missing token"}`, ...
