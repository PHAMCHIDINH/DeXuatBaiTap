
# Kiến Trúc Backend Tim Mạch – Mô Tả Chi Tiết Các Module

## 1. Module `users`
### Chức năng
- Đăng ký, đăng nhập, xác thực người dùng.
- Xử lý thông tin người dùng.

### Luồng xử lý API
1. Register → validate → hash password → save DB → return UserResponse.
2. Login → validate → check password → generate JWT → return LoginResponse.
3. Me → get userID from JWT → fetch user → return profile.

---

## 2. Module `patients`
### Chức năng
- Quản lý hồ sơ bệnh nhân.
- CRUD bệnh nhân.

### Luồng xử lý API
1. Create → extract userID → validate → save via sqlc → return PatientResponse.
2. List → get userID → find all patients of user → return list.
3. Detail → check patient belongs to user → return patient.
4. Update → validate → update DB → return updated patient.
5. Delete → check permission → delete.

---

## 3. Module `predictions`
### Chức năng
- Gọi ML FastAPI để dự đoán tim mạch.
- Lưu lịch sử dự đoán.

### Luồng xử lý API
1. Predict:
   - Get userID + patientID
   - Validate patient ownership
   - Build MLRequest
   - POST → ML FastAPI
   - Receive probability + risk_level
   - Save prediction to DB
   - Return PredictionResponse
2. History:
   - Validate ownership
   - List predictions by patient
   - Return array of PredictionResponse

---

## 4. Module `stats`
### Chức năng
- Thống kê nhanh cho dashboard: tổng số bệnh nhân của bác sĩ hiện tại.
- Đếm số bệnh nhân theo risk_level mới nhất (high/medium/low/none) dựa trên bảng `predictions`.

### Luồng xử lý API
- `GET /stats` (cần JWT): lấy userID từ middleware, đếm tổng bệnh nhân (`patients`), đếm nhóm nguy cơ với latest prediction của từng patient, trả về `total_patients` + `risk_counts`.

---

## 5. Module `exercises`
### Chức năng
- Quản lý template bài tập: tạo/list `exercise_templates` theo risk_level.
- Lưu & truy xuất kế hoạch gợi ý sau khi predict: `exercise_recommendations`.

### Luồng xử lý API
- `POST /exercise-templates` (JWT): tạo template mới (name/intensity/duration/freq/target_risk_level/tags).
- `GET /exercise-templates` (JWT): xem danh sách template (phục vụ gợi ý).
- `GET /patients/:id/recommendations` (JWT): kiểm tra sở hữu patient → trả các kế hoạch đã lưu (mỗi lần predict sẽ lưu một bản ghi).

---

## 6. Luồng hoạt động tổng thể

```
Client → /patients/:id/predict
   → Auth middleware (decode JWT)
   → patients module (verify ownership)
   → predictions module (call ML)
   → save prediction
   → return JSON response
```

---

## 7. Kết luận

Các module được phân chia rõ ràng theo domain:
- users → auth
- patients → thông tin bệnh nhân
- predictions → gọi ML + lưu lịch sử
- stats → tổng hợp số liệu cho dashboard
- exercises → template bài tập + khuyến nghị

Dễ mở rộng, bảo trì và tích hợp microservices.
