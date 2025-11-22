# JWT Flow (Frontend)

## Tóm tắt
- Backend trả JWT khi đăng ký/đăng nhập.
- Frontend lưu token trong `localStorage` và set vào Axios default header `Authorization: Bearer <token>`.
- Khi logout, xoá token + clear header.
- Request 401 hiện tại không auto-logout; nếu cần, bắt lỗi trong hooks/API để logout thủ công.

## Các bước luồng
1) Đăng nhập (`POST /Users/Login`): response `{access_token, user}`.
2) Trong `AuthContext.login`:
   - Lưu token vào `localStorage`.
   - Gọi `setAuthToken(token)` để gắn header mặc định.
   - Lưu user vào state.
3) Đăng ký tương tự (resp `{token, user}`).
4) Mỗi request qua Axios client sẽ tự đính `Authorization` nếu token đang set.
5) Khi 401 (nếu backend trả), component/hook gọi `logout()` để xoá token.
6) Logout: xoá token khỏi storage, clear header, reset state.

## File chính
- `src/api/client.ts`: axios instance + `setAuthToken(token|null)`.
- `src/context/AuthContext.tsx`: quản lý token/user, login/register/logout, set header.
- `src/utils/storage.ts`: lưu/đọc/xoá token trong localStorage.

## Lưu ý
- Token không tự refresh; cần login lại khi hết hạn.
- Luôn dùng HTTPS ở môi trường thật; không pass token qua query string.
