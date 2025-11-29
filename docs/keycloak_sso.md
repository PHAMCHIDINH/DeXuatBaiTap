# Keycloak SSO cho TimMach

Tài liệu này mô tả cách dự án tích hợp Keycloak, cấu hình cần thiết và luồng đăng nhập/xác thực từ frontend → API → Keycloak → Postgres.

## Thành phần liên quan
- Dịch vụ Keycloak chạy trong `docker-compose.yml` (`quay.io/keycloak/keycloak:23.0.0`, cổng host `8081`, admin mặc định `admin/admin`, lệnh `start-dev`).
- API Go (Gin) đọc cấu hình Keycloak từ biến môi trường trong `TimMach_api/config/config.go`:
  - `KEYCLOAK_URL` (mặc định `http://localhost:8081`)
  - `KEYCLOAK_REALM` (mặc định `timmach`)
  - `KEYCLOAK_CLIENT_ID` (mặc định `timmach-webapp`)
  - `KEYCLOAK_CLIENT_SECRET` (không dùng cho flow login mới, giữ để tương thích nếu cần gọi client confidential)
  - `KEYCLOAK_SKIP_TLS_VERIFY` (chỉ bật cho môi trường dev)
- Cơ sở dữ liệu: migration `TimMach_api/db/migrations/20251126090000_add_keycloak_auth.sql` thêm cột `users.keycloak_id` (unique) và cho phép `password_hash` rỗng để chứa tài khoản SSO.

## Thiết lập Keycloak (dev)
1) Chạy `docker compose up -d keycloak` (hoặc toàn bộ stack). Mở http://localhost:8081, đăng nhập `admin/admin`.  
2) Tạo Realm mới tên `timmach`.  
3) Tạo Client (public hoặc confidential) tên `timmach-webapp`:
   - Bật **Direct Access Grants** để cho phép password grant.
   - Nếu chạy trên trình duyệt (public client), cấu hình `Web Origins` (`http://localhost:5173`) và `Valid Redirect URIs` (có thể `*` nếu không dùng redirect).
   - Nếu cần refresh token, bật “Offline Access” + “Allow Refresh Token”.
4) Tạo người dùng thử (username/email + password) và bật “Temporary” = off để không phải đổi mật khẩu lần đầu.  
5) Khởi động API với các biến env trên (Compose không tự seed realm/client, nên bước này bắt buộc).

## Luồng đăng nhập (FE gọi thẳng Keycloak - password grant)
Các file chính phía API: `TimMach_api/main.go`, `TimMach_api/middleware/keycloak.go`. Frontend gửi username/password trực tiếp tới token endpoint của Keycloak.

1) FE gọi POST `KEYCLOAK_URL/realms/<realm>/protocol/openid-connect/token` với body form: `grant_type=password`, `client_id=<timmach-webapp>`, `username`, `password` (thêm `client_secret` nếu là confidential).  
2) Keycloak trả `access_token` (+ `refresh_token` nếu bật). FE lưu `access_token` (localStorage/session) và gắn vào header `Authorization: Bearer <token>` cho mọi request API.  
3) API dùng `KeycloakMiddleware` để verify token, map user và xử lý.  
4) Khi token hết hạn, FE gọi lại token endpoint với `grant_type=refresh_token` để lấy token mới.

## Luồng bảo vệ API (Bearer token)
Các file chính: `TimMach_api/main.go`, `TimMach_api/middleware/keycloak.go`.

1) Mọi router chính (`/api/users/me`, `/patients`, `/predictions`, `/exercises`, `/reports`, `/stats`) được wrap bằng `KeycloakMiddleware`.  
2) Middleware:
   - Lấy JWKS từ `KEYCLOAK_URL/realms/<realm>` và verify chữ ký `access_token` (có thể bỏ qua TLS nếu `KEYCLOAK_SKIP_TLS_VERIFY=true` cho dev).  
   - Kiểm `aud` theo `KEYCLOAK_CLIENT_ID` (không bỏ qua ClientID check). Nếu cần nhiều aud, đặt `ExpectedAudiences`.  
   - Parse claims: `sub`, `email`, `preferred_username`, `realm_access.roles`.  
3) Ánh xạ user Keycloak vào DB (`MapKeycloakUser`):
   - Tìm `users.keycloak_id == sub`. Nếu có → dùng user đó.  
   - Nếu chưa có và claim email tồn tại → tìm `users.email`. Nếu trùng email nhưng chưa gán Keycloak ID → gán thêm `keycloak_id` để tái sử dụng tài khoản cũ; nếu email đã gắn ID khác → trả lỗi.  
   - Nếu không tìm thấy → tạo user mới với ID `USER_<YYYYMMDD>_<seq>`, lưu `email` (nếu có) và `keycloak_id`.  
4) Sau khi map, middleware set vào Gin context:
   - `userID`: ID nội bộ của user (dùng cho toàn bộ handler để ràng buộc dữ liệu).  
   - `userEmail`: email nếu claim có.  
   - `keycloak_sub`: subject Keycloak.  
   - `keycloak_roles`: danh sách role trong realm (hiện mới lưu vào context, chưa áp policy).  
5) Handler phía sau chỉ cần gọi `utils.UserIDFromContext` để lấy `userID` và tiếp tục xử lý nghiệp vụ.

## Mối quan hệ dữ liệu
- Bảng `users` vừa hỗ trợ tài khoản cũ (password) vừa hỗ trợ SSO: `password_hash` có thể rỗng, `keycloak_id` là unique.  
- Khi user đăng nhập bằng Keycloak lần đầu, bản ghi sẽ được tạo (hoặc gắn thêm `keycloak_id`) rồi tái sử dụng cho các request sau, đảm bảo mọi bảng khác vẫn khóa theo `user_id`.

## Lưu ý vận hành
- Backend không còn endpoint `/users/password-login`; FE phải tự đăng nhập/refresh token với Keycloak.
- Với HTTPS tự ký trong môi trường dev, bật `KEYCLOAK_SKIP_TLS_VERIFY=true` cho API; không bật trong production.  
- Nếu muốn dùng role để phân quyền, đọc `keycloak_roles` từ context và bổ sung kiểm tra tại handler hoặc middleware mới (hiện chưa enforce).  
- Legacy JWT (HMAC) vẫn còn mã nguồn (`middleware/auth.go`, `docs/token_flow.md`) nhưng server mặc định đang dùng Keycloak middleware cho toàn bộ route; không trộn hai chế độ cùng lúc.
