# Tổng quan dự án TimMach

Dự án gồm 3 thành phần chính: API (Go), frontend (React/Vite), và service ML (FastAPI). Dưới đây là bức tranh tổng thể, luồng nghiệp vụ và chi tiết triển khai.

## Bức tranh tổng thể (mục tiêu sản phẩm)
- Quản lý bệnh nhân tim mạch: mỗi bác sĩ có danh sách bệnh nhân riêng (gắn `user_id`), lưu tên/giới tính/ngày sinh.
- Dự đoán nguy cơ tim mạch bằng ML: nhập 11 thông số (tuổi, chiều cao, cân nặng, huyết áp, cholesterol, glucose, smoke, alco, active…) → gửi FastAPI → lấy probability + risk_level.
- Lưu lịch sử & raw_features: mỗi lần dự đoán lưu vào bảng `predictions` kèm dữ liệu đầu vào, tạo thành history theo thời gian.
- Gợi ý tập luyện theo risk_level: bảng `exercise_templates` chứa template cho low/medium/high; backend chọn template và lưu `exercise_recommendations`.
- Quản lý tài khoản: đăng ký/đăng nhập, JWT auth, `GET /users/me` lấy profile.
- Định hướng mở rộng: xuất PDF hồ sơ & giáo án, lịch tập chi tiết, dashboard thống kê risk group, audit log/roles.

## Màn hình & tính năng chính (hiện có vs. dự kiến)
- Dashboard: hiện số bệnh nhân (top 5); dự kiến thêm thống kê risk group, quick action tới nhóm “high risk”.
- Patients list: bảng bệnh nhân, phân trang (limit/offset), tạo/sửa/xóa; dự kiến thêm tìm kiếm/lọc theo risk.
- Patient detail: thông tin cơ bản + prediction gần nhất; nút Predict, Edit, History; dự kiến thêm “Xuất PDF hồ sơ/giáo án”.
- Predict: form 11 trường → hiển thị card kết quả (probability, risk_label) + card gợi ý tập luyện.
- History: biểu đồ Recharts diễn biến probability% theo thời gian + bảng các lần dự đoán.
- Profile: xem email, ngày tạo; dự kiến đổi mật khẩu/thông tin.
- Báo cáo (dự kiến): export CSV/Excel/PDF cho bệnh nhân hoặc danh sách.

## Luồng người dùng chính
1) Bác sĩ lần đầu vào hệ thống  
   - Đăng ký hoặc đăng nhập (JWT lưu localStorage, Axios tự gắn header).  
   - Thấy Dashboard với số bệnh nhân = 0.

2) Tạo bệnh nhân & dự đoán lần đầu  
   - Vào Patients → Add new (POST `/patients`, gắn `user_id`).  
   - Mở bệnh nhân → Predict → nhập 11 input → POST `/patients/:id/predict`.  
   - Backend gọi ML, lưu `predictions` + `exercise_recommendations` (chọn template theo risk).  
   - Frontend hiển thị kết quả + giáo án gợi ý.

3) Theo dõi tái khám  
   - Predict lại với dữ liệu mới → thêm dòng `predictions`.  
   - Tab History: biểu đồ probability theo thời gian + bảng chi tiết.  
   - (Tương lai) Export PDF hồ sơ/giáo án cho bệnh nhân.

## Luồng hệ thống (ví dụ Predict)
- Browser (React): POST `/api/patients/:id/predict` kèm JWT.
- API Go (Gin): middleware auth lấy user; kiểm tra quyền sở hữu patient; gọi ML `POST http://ml:8000/predict`; lưu `predictions` + `exercise_recommendations`; trả JSON.
- ML FastAPI: nhận 11 feature, tính thêm BMI/bp_ratio, dùng RandomForest → trả probability + risk_level.
- Database Postgres: lưu users/patients/predictions/exercise_templates/exercise_recommendations.

## Đích v1 (hiện trạng)
- Bác sĩ có thể đăng ký/đăng nhập, quản lý bệnh nhân riêng, chạy dự đoán nguy cơ tim mạch, nhận kế hoạch tập luyện gợi ý, xem lịch sử theo thời gian.  
- Kiến trúc tách lớp rõ: React (UI) → Go API (business + DB) → FastAPI ML (model), đóng gói bằng Docker Compose.

## Kiến trúc & Docker Compose
- File `docker-compose.yml` dựng 3 service:
  - `db`: Postgres 16, user/pass `postgres`, DB `heartdb`, volume `db_data`, publish cổng host `5432`.
  - `ml`: build từ `ml-python/dockerfile`, FastAPI cổng `8000`, healthcheck `/health`.
  - `api`: build từ `TimMach_api/Dockerfile`, phơi bày `8080`, gọi migrations bằng `entrypoint.sh` (goose up) trước khi chạy server.
- Chạy/dừng: `docker compose up -d --build` và `docker compose down` (thêm `-v` để xóa volume).

## Database (Postgres 16)
- Kết nối mặc định: host `localhost` (Compose: `db`), port `5432`, user/pass `postgres` / `postgres`, db `heartdb`.
- `DB_URL` mẫu: `postgres://postgres:postgres@localhost:5432/heartdb?sslmode=disable` (Compose dùng `@db:5432`).
- Migrations: `TimMach_api/db/migrations/20251121170000_init_schema.sql`; chạy tự động khi container `api` start hoặc thủ công bằng `make goose-up` trong `TimMach_api`.
- Lược đồ chính:
  - `users(id text, email unique, password_hash, created_at)`; id sinh từ sequence `user_id_seq` với format utils.FormatUserID.
  - `patients(id bigserial, user_id fk users, name, gender smallint, dob date, created_at)`.
  - `predictions(id bigserial, patient_id fk patients, probability double, risk_label text, raw_features jsonb, created_at)`.
  - `exercise_templates(id bigserial, name, intensity, description, duration_min, freq_per_week, target_risk_level, tags text[])`.
  - `exercise_recommendations(id bigserial, patient_id fk patients, prediction_id fk predictions, plan jsonb, created_at)`.
- Indexes đáng chú ý: `patients(user_id)`, `predictions(patient_id)`, `predictions(patient_id, created_at desc)`, `exercise_templates(target_risk_level)`, `exercise_recommendations(patient_id/prediction_id)`.

## Backend: `TimMach_api` (Go 1.23)
- Stack & thư viện: Gin, pgx + sqlc, JWT (golang-jwt), Resty (call ML, timeout 10s), zap logger, goose migrations.
- Khởi động: `entrypoint.sh` yêu cầu `DB_URL`, chạy `goose` áp dụng migration, sau đó chạy `server` (định nghĩa trong `main.go`).
- Env chính (config/config.go):
  - `DB_URL` mặc định `postgres://postgres:postgres@localhost:5432/heartdb?sslmode=disable` (Compose dùng host `db`).
  - `JWT_SECRET` mặc định `dev-secret`.
  - `ML_BASE_URL` mặc định `http://localhost:8000` (Compose: `http://ml:8000`).
  - `PORT` mặc định `8080`.
- Database schema (db/migrations/20251121170000_init_schema.sql):
  - `users` (id text từ sequence `user_id_seq`, email unique, password_hash).
  - `patients` (thuộc user, tên/giới tính/ngày sinh).
  - `predictions` (kết quả ML, xác suất, risk_label, raw_features JSONB).
  - `exercise_templates` (template gợi ý tập luyện theo risk_level) và `exercise_recommendations` (plan lưu theo prediction).
- API hiện có (tất cả mount dưới `/api`, auth JWT qua header Bearer):
  - `POST /users/register` → tạo user, trả `{token, user}`; `POST /users/login` → `{access_token, user}`; `GET /users/me` → profile.
  - Patients (yêu cầu JWT):
    - `POST /patients`, `GET /patients`, `GET/PUT/PATCH/DELETE /patients/:id` (validate sở hữu user).
    - `GET /patients` nhận thêm `risk` (low/medium/high/none) để lọc theo kết quả dự đoán mới nhất; mỗi patient trả thêm `latest_prediction` (probability, risk_label, created_at) nếu có.
  - Predictions (yêu cầu JWT): `POST /patients/:id/predict` gửi 11 input (age_years, gender, height, weight, ap_hi, ap_lo, cholesterol, gluc, smoke, alco, active) tới ML, lưu prediction + gợi ý tập (dùng template theo risk_label hoặc fallback). `GET /patients/:id/predictions` trả lịch sử (limit/offset).
  - Exercises:
    - Template: `GET/POST /exercise-templates` (tạo và xem danh sách template bài tập theo risk_level).
    - Recommendation: `GET /patients/:id/recommendations` trả các kế hoạch đã lưu (mỗi lần predict sẽ lưu vào `exercise_recommendations`).
  - Report PDF: `GET /patients/:id/report.pdf` (JWT) xuất hồ sơ bệnh nhân tim mạch (thông tin cơ bản, dự đoán gần nhất, gợi ý tập luyện, lịch sử).
  - Stats (yêu cầu JWT): `GET /stats` trả `total_patients` và `risk_counts` (high/medium/low/none) tính trên prediction mới nhất của từng bệnh nhân.
- Thư mục chính:
  - `modules/users|patients|predictions`: router + handler.
  - `middleware/auth.go`: decode JWT, set userID vào context.
  - `db/sqlc`: code generate; `db/migrations`: schema.
  - `makefile`: lệnh goose/sqlc (dùng `.env` nếu có).

## Frontend: `TimMach_client` (React + TS)
- Stack: React 18 + Vite, React Router v7, React Query v5, Axios, react-hook-form + Zod, Recharts, tailwind-like utility classes.
- Cấu hình: `VITE_API_URL` (mặc định `http://localhost:8080/api`), token lưu trong `localStorage`. `AuthContext` gắn Authorization header global và bảo vệ route.
- Luồng chính:
  - Auth: Login/Register trang riêng, sau khi thành công lưu token + user; logout xóa token.
  - Dashboard: hiển thị tổng bệnh nhân + phân bố nguy cơ (high/medium/low/none) dựa trên dự đoán mới nhất; danh sách 5 bệnh nhân gần nhất.
  - Patients: list có phân trang (limit/offset) + lọc theo `risk` (low/medium/high/none); hiển thị nguy cơ gần nhất (risk_label, probability%, thời gian dự đoán); tạo/sửa/xóa; form gồm name, gender (0/1/2), dob.
  - Patient detail: thông tin cơ bản + prediction gần nhất; nút Predict, Edit, History.
  - Predict page: form 11 trường input; gọi API tạo prediction, hiện thẻ kết quả (probability, risk_label) và thẻ khuyến nghị tập luyện.
  - History page: biểu đồ Recharts đường theo probability% và bảng lịch sử.
  - Profile page: xem thông tin user; Sidebar có Dashboard/Patients/Profile + Logout.
- Script: `npm run dev` (Vite), `build`, `preview`.

## ML service: `ml-python` (FastAPI)
- Endpoint:
  - `GET /health` (trả ok khi model load).
  - `POST /predict` nhận 11 feature (age_years, gender, height, weight, ap_hi, ap_lo, cholesterol, gluc, smoke, alco, active); tính thêm BMI + bp_ratio; trả `{probability, label, risk_level(low|medium|high)}`.
- Model & dữ liệu: RandomForest đã lưu ở `training/heart_disease_model.pkl`; dataset gốc `training/cardio_train.csv`; notebook `training/heart_disease_notebook.ipynb`.
- Phụ thuộc: FastAPI, scikit-learn, pandas, joblib (xem `requirements.txt`). Compose base URL: `http://ml:8000`.
