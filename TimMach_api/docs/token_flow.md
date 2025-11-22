# Luồng hoạt động của Token (JWT) trong TimMach

## Tóm tắt
- Hệ thống dùng JWT (HMAC HS256) để xác thực stateless.
- Token chứa `user_id` (chuỗi, ví dụ `USER_20251121_001`) và `email`.
- Token được tạo khi đăng ký/đăng nhập, gửi về client, client đính kèm vào header `Authorization: Bearer <token>` cho các API cần bảo vệ.

## Nơi tạo token
- Module: `modules/auth/jwt.go`
  - Struct `JWTMaker{Secret, TTL}` tạo token với claims:
    - `user_id`, `email`, `sub` (=user_id)
    - `exp` (TTL mặc định 24h), `iat`
  - Secret lấy từ env `JWT_SECRET`.
- Được gọi khi:
  - Register: `users.Register` → tạo user → `Tokens.GenerateToken`
  - Login: `users.Login` → xác thực mật khẩu → `Tokens.GenerateToken`

## Nơi kiểm tra token
- Middleware: `middleware/auth.go`
  - Đọc header `Authorization` (Bearer)
  - Parse JWT với cùng secret, kiểm tra HMAC
  - Nếu hợp lệ: set `c.Set("userID", <user_id>)` vào Gin context
  - Nếu lỗi: trả 401 với thông báo (`missing token` / `invalid token`)
- Các handler lấy userID từ context qua `utils.UserIDFromContext`.

## Phạm vi áp dụng
- Public:
  - `POST /api/users/register`
  - `POST /api/users/login`
- Có JWT:
  - `GET /api/users/me`
  - Tất cả route `/api/patients/...`
  - Tất cả route `/api/patients/:id/predict` và history

## Cấu hình/Env
- `JWT_SECRET`: bắt buộc set an toàn ở môi trường thật.
- TTL mặc định: 24h (đặt trong `main.go` khi tạo `auth.JWTMaker`).

## Luồng mẫu
1) Đăng ký/Đăng nhập → nhận `token`/`access_token` trong response.
2) Client lưu token (localStorage) và gửi kèm header `Authorization: Bearer <token>` cho các API cần JWT.
3) Middleware auth kiểm tra token, gắn `userID` vào context.
4) Handler dùng `utils.UserIDFromContext` để ràng buộc dữ liệu theo user.

## Lưu ý an toàn
- Luôn truyền token qua HTTPS.
- Không gửi token trong query string.
- Đổi `JWT_SECRET` khi triển khai production, có thể rotate định kỳ.
