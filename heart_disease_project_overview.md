# Dự Án Dự Đoán Nguy Cơ Bệnh Tim Mạch

## 1. Tổng quan dự án
Hệ thống dự đoán nguy cơ mắc bệnh tim mạch dựa trên các chỉ số sức khỏe của bệnh nhân. Mục tiêu là xây dựng một backend mạnh mẽ, dễ mở rộng, chia thành hai phần chính:
- **Backend chính (Golang)**: CRUD, quản lý người dùng, bệnh nhân, lưu lịch sử dự đoán, expose API cho frontend.
- **ML Service (Python)**: Load mô hình `.pkl`, chuẩn hóa dữ liệu đầu vào, thực hiện dự đoán, trả kết quả cho backend.

Kiến trúc hướng microservice giúp dễ scale và bảo trì.

---

## 2. Các thành phần chính của hệ thống
### 2.1 Frontend
Có thể là React/Vue/Angular hoặc giao diện đơn giản. Nhiệm vụ:
- Cho phép nhập thông tin bệnh nhân và các chỉ số y tế.
- Gửi request đến backend (Go).
- Hiển thị kết quả dự đoán: xác suất, mức độ nguy cơ.
- Hiển thị lịch sử dự đoán.

### 2.2 Backend chính (Golang)
Đảm nhiệm:
- **Xác thực người dùng (Auth)**: login/signup bằng JWT.
- **CRUD bệnh nhân**.
- **CRUD người dùng** (tùy thiết kế).
- **Gọi Python ML service** khi cần dự đoán.
- **Lưu prediction vào cơ sở dữ liệu**.
- Cung cấp API để frontend tương tác.

### 2.3 ML Service (Python – FastAPI)
- Load file mô hình **`model.pkl`**.
- Chuẩn hóa dữ liệu đầu vào.
- Predict và trả về probability + label.
- Expose endpoint:
  - `POST /predict`
  - `GET /health` (optional)

---

## 3. Kiến trúc tổng thể
```
Client (Web/Mobile)
        |
        v
Golang Backend  <---->  Database (PostgreSQL)
        |
        v
Python ML Service (FastAPI, load model.pkl)
```

Backend Go gọi nội bộ đến Python qua HTTP, ví dụ:
```
POST http://ml-service:8000/predict
```

---

## 4. Thiết kế API
### 4.1 Backend Go
**Auth:**
- `POST /api/auth/login`
- `POST /api/auth/register`

**Patients:**
- `POST /api/patients`
- `GET /api/patients/:id`
- `PUT /api/patients/:id`
- `DELETE /api/patients/:id`

**Prediction:**
- `POST /api/patients/:id/predict`
- `GET /api/patients/:id/predictions`

### 4.2 ML Service (Python)
**Prediction endpoint:**
- `POST /predict`

Request body:
```json
{
  "age": 54,
  "sex": 1,
  "cholesterol": 245,
  "trestbps": 140,
  "thalach": 150,
  "cp": 2,
  "fbs": 0,
  "restecg": 1,
  "exang": 0,
  "oldpeak": 1.4,
  "slope": 2,
  "ca": 0,
  "thal": 2
}
```

Response:
```json
{
  "probability": 0.82,
  "label": "high_risk"
}
```

---

## 5. Thiết kế Database
Sử dụng PostgreSQL.

### Bảng `users`
| Trường | Kiểu | Mô tả |
|--------|------|--------|
| id | PK | UUID |
| email | text | Email đăng nhập |
| password_hash | text | Mật khẩu mã hóa |
| created_at | timestamp | Thời gian tạo |

### Bảng `patients`
| Trường | Kiểu | Mô tả |
|--------|------|--------|
| id | PK | UUID |
| user_id | FK | Chủ sở hữu |
| name | text | Tên bệnh nhân |
| gender | int | 0=female,1=male |
| dob | date | Ngày sinh |

### Bảng `predictions`
| Trường | Kiểu | Mô tả |
|--------|------|--------|
| id | PK | UUID |
| patient_id | FK | Liên kết bệnh nhân |
| probability | float | Xác suất mắc bệnh |
| risk_label | text | mức độ nguy cơ |
| raw_features | JSON | Input gốc |
| model_version | text | để quản lý model |
| created_at | timestamp | thời điểm dự đoán |

---

## 6. Cấu trúc thư mục dự án
```
heart-disease-risk/
├── backend-go/
│   ├── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── models/
│   │   └── db/
│   ├── go.mod
│   ├── Dockerfile
│
├── ml-python/
│   ├── main.py
│   ├── heart_model.pkl
│   ├── preprocess.py
│   ├── requirements.txt
│   ├── schema.py
│   ├── Dockerfile
│
├── docker-compose.yml
└── README.md
```

---

## 7. Quy trình hoạt động End-to-End
1. Người dùng nhập chỉ số và gửi yêu cầu dự đoán.
2. Frontend gửi request đến backend Go.
3. Backend Go validate dữ liệu → gọi ML Python.
4. Python ML service load model và dự đoán.
5. Python trả kết quả về cho Go.
6. Go lưu vào DB và trả lại kết quả cho frontend.
7. Người dùng xem kết quả hoặc lịch sử.

---

## 8. Chiến lược triển khai (Deployment)
### Option 1: 2 container tách biệt (khuyến khích)
- `go-backend` chạy API
- `python-ml` chạy FastAPI
- Docker Compose chạy cả hai

### Option 2: 1 container (demo, MVP)
- Startup script chạy Python backend & Go service chung 1 container.

---

## 9. Mở rộng trong tương lai
- Thêm dashboard thống kê.
- Cho phép nhiều phiên bản mô hình chạy song song.
- Logging/monitoring bằng Prometheus - Grafana.
- Triển khai Kubernetes khi hệ thống lớn.

---

## 10. Kết luận
Tài liệu này mô tả toàn bộ kiến trúc và cách tổ chức dự án "Dự đoán nguy cơ bệnh tim mạch" theo hướng microservice, gồm Go backend và Python ML service. Cách chia tách này giúp dễ phát triển, bảo trì, nâng cấp và mở rộng khi cần.

