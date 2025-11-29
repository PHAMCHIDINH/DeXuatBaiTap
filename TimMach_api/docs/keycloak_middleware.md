# Giải thích `middleware/keycloak.go`

Middleware này bảo vệ các route bằng Keycloak OIDC access token. Nó nhận cấu hình runtime, xác thực token, ánh xạ user Keycloak vào DB, rồi gắn thông tin vào Gin context cho handler phía sau.

## Cấu hình đầu vào (`KeycloakConfig`)
- `BaseURL`: URL Keycloak (ví dụ `http://localhost:8081`).
- `Realm`: realm đang dùng (ví dụ `timmach`).
- `ClientID`: client cần kiểm tra (dùng khi verify audience; hiện bật `SkipClientIDCheck`).
- `SkipTLSVerify`: cho phép bỏ kiểm tra TLS (chỉ dev).
- `ExpectedAudiences`: chưa sử dụng; chỗ để ép audience nếu cần sau này.

## Các claim đọc từ token
- `sub` (bắt buộc): subject Keycloak, dùng để map vào `users.keycloak_id`.
- `email`, `preferred_username`: dùng làm email khi tạo user nếu có.
- `realm_access.roles`: lưu vào context để dùng phân quyền nếu cần.
- `resource_access`: đọc nhưng chưa dùng; có sẵn nếu muốn kiểm role theo client.

## Luồng xử lý từng request
1) Đọc header `Authorization`. Thiếu hoặc không phải Bearer → 401 `missing token`.
2) Verify token:
   - Tải OIDC metadata từ `BaseURL/realms/<Realm>` để lấy JWKS.
   - Dùng `oidc.Provider().Verifier` kiểm tra chữ ký + hạn và `aud` khớp `KEYCLOAK_CLIENT_ID` (có thể thêm `ExpectedAudiences` nếu cần nhiều aud).
   - Lỗi verify → 401 `invalid token`.
3) Parse claims vào struct `keycloakClaims`. Thiếu `sub` → 401.
4) Ánh xạ user vào DB (`MapKeycloakUser`):
   - Tìm `users.keycloak_id == sub`. Nếu có → dùng luôn.
   - Nếu chưa có, nhưng có email: tìm user theo email; nếu email chưa gắn Keycloak ID → gắn `keycloak_id=sub`; nếu đã gắn ID khác → báo lỗi.
   - Nếu vẫn chưa có user: tạo mới với ID format `USER_<YYYYMMDD>_<seq>`, lưu email (nếu có) và `keycloak_id`.
   - Lỗi DB khác → 500 `cannot map user`.
5) Đặt vào Gin context:
   - `userID`: ID nội bộ dùng cho toàn bộ handler để ràng buộc dữ liệu.
   - `userEmail`: nếu claim có email.
   - `keycloak_sub`: subject Keycloak.
   - `keycloak_roles`: danh sách role trong realm (nếu có).
6) `c.Next()` để handler phía sau xử lý.

## Nơi sử dụng
- Được khởi tạo trong `main.go` và truyền vào router của users/patients/predictions/exercises/reports/stats. Mọi endpoint này yêu cầu Bearer token Keycloak hợp lệ.

## Lưu ý/An toàn
- `SkipTLSVerify` chỉ dành cho dev khi Keycloak dùng chứng chỉ tự ký.
- Vì `SkipClientIDCheck=true`, token hợp lệ từ cùng realm đều được chấp nhận; nếu muốn khóa theo client/audience, cần đặt `SkipClientIDCheck=false` và thiết lập `ExpectedAudiences`.
- Roles mới chỉ được gắn vào context, chưa enforce; nếu cần RBAC, đọc `keycloak_roles` trong handler hoặc thêm middleware kiểm quyền.***
