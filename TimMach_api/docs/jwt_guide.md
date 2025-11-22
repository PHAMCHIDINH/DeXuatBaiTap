# JWT trong TimMach_api

## JWT là gì?
- JWT (JSON Web Token) là chuỗi text gồm 3 phần `header.payload.signature` được Base64URL và nối bằng dấu chấm.
- Token được ký (ở đây bằng HMAC-SHA256) để chứng minh tính toàn vẹn; nó **không** mã hóa nội dung, nên không nhét thông tin nhạy cảm vào payload.
- Máy chủ có thể xác thực người dùng mà không cần lưu session trên server (stateless).

## Token được tạo ở đâu?
- `modules/users/controllers.go` → struct `JWTMaker` tạo token khi đăng ký/đăng nhập.
- Thuật toán: `jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`, ký bằng secret.
- Secret lấy từ biến môi trường `JWT_SECRET` (xem `config/config.go`), default `"dev-secret"`; production phải thay bằng chuỗi ngẫu nhiên mạnh (>= 32 byte).
- TTL mặc định: 24h (set tại `main.go`: `JWTMaker{Secret: cfg.JWTSecret, TTL: 24 * time.Hour}`).
- Claims hiện có:
  - `user_id` (string UUID)
  - `email`
  - `exp` (expiry Unix)
  - `iat` (issued-at Unix)
  - `sub` (subject = user_id)

## Token được kiểm tra ở đâu?
- `middleware/auth.go` → `AuthMiddleware(secret)` gắn vào router.
- Luồng: đọc header `Authorization: Bearer <jwt>` → parse token với cùng secret → kiểm tra phương thức ký là HMAC → nếu hợp lệ, lấy `user_id` trong claims và set vào context (`c.Set("userID", uid)`).
- Các handler dùng `utils.UserIDFromContext(c)` để lấy userID đã set.

## Những API nào dùng JWT?
- Public: `POST /api/users/register` trả `token`, `POST /api/users/login` trả `access_token`.
- Sau đó client đính kèm header `Authorization: Bearer <token>` cho các route đã gắn middleware:
  - `GET /api/users/me`
  - Toàn bộ route của `patients` và `predictions` (đăng ký ở `main.go` với `authMiddleware`).

Ví dụ gọi:
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/users/me
```

## Cấu hình thực tế cần đặt
- `JWT_SECRET`: bắt buộc đổi ở môi trường thật; nên là key ngẫu nhiên 256 bit+.
- TTL: hiện cố định 24h; nếu muốn thay đổi, sửa `JWTMaker.TTL` khi khởi tạo (có thể đưa vào env nếu cần linh hoạt).
- Giờ hệ thống phải chuẩn để tránh lệch `exp/iat`.

## Kiến thức nên nắm khi tìm hiểu thêm
- Cấu trúc JWT, Base64URL, khác biệt “ký” vs “mã hóa”.
- Thuật toán ký: HS256 (đối xứng) vs RS256 (bất đối xứng) và cách quản lý key.
- Ý nghĩa các claim tiêu chuẩn `exp`, `iat`, `sub`, cùng claim tự định nghĩa (`user_id`, `email`).
- Cách truyền an toàn: luôn qua HTTPS, để trong header `Authorization`, tránh query string.
- Chiến lược thu hồi/rotation: đổi secret, duy trì bảng blacklist/allowlist token, hoặc rút ngắn TTL + refresh token.
- Xử lý lỗi thường gặp: token hết hạn → 401, sai chữ ký/phương thức → 401, thiếu header → 401.
