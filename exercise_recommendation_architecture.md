
# Kiến Trúc Module “Exercise Recommendation”

Tài liệu này mô tả kiến trúc mở rộng của hệ thống dự đoán bệnh tim mạch để **đưa ra gợi ý tập luyện thể dục** dựa trên:
- Hồ sơ bệnh nhân (patients)
- Các chỉ số sức khỏe nhập vào
- Kết quả dự đoán nguy cơ tim mạch (predictions từ ML service)

---

## 1. Mục tiêu chức năng

### 1.1. Chức năng chính

- Sau khi chạy dự đoán tim mạch cho một bệnh nhân, hệ thống:
  - Nhận được **risk_level** (low / medium / high) và probability từ ML service.
  - Kết hợp với thông tin sức khỏe cơ bản (tuổi, BMI, mức độ vận động hiện tại, hút thuốc, uống rượu, …).
  - Sinh ra **kế hoạch tập luyện gợi ý (Exercise Plan)**:
    - Loại bài tập (đi bộ, đạp xe, yoga, kéo giãn, …)
    - Cường độ (intensity: low / medium / high)
    - Thời lượng (phút/buổi)
    - Tần suất (số buổi/tuần)
    - Ghi chú lưu ý (notes)

- Lưu lại kế hoạch tập luyện này trong DB để:
  - Xem lại theo từng lần dự đoán.
  - Mở rộng về sau: theo dõi lịch sử luyện tập, thu thập feedback, làm recommender ML.

### 1.2. Phạm vi giai đoạn 1

- Logic gợi ý dựa trên **rule-based** (luật đơn giản).
- ML service (FastAPI) **chỉ lo predict tim mạch**, không thay đổi nhiều.
- Logic recommendation nằm chủ yếu ở backend Go.

---

## 2. Kiến trúc tổng thể

### 2.1. Sơ đồ khái quát

```text
Client (React) 
    ↓
Go Backend (modules: users, patients, predictions, recommendations)
    ↓
ML Service (FastAPI, model .pkl)
    ↓
PostgreSQL (users, patients, predictions, exercise_templates, exercise_recommendations)
```

Luồng chính khi “Predict + Recommendation”:

1. Client gọi `POST /api/patients/:id/predict`.
2. Go backend:
   - Kiểm tra quyền user với bệnh nhân.
   - Gửi request sang ML service `/predict`.
   - Nhận lại `probability`, `label`, `risk_level`.
   - Áp dụng **rule-based** để chọn kế hoạch tập luyện phù hợp.
   - Lưu `prediction` + `exercise_recommendation` vào DB.
   - Trả về response gồm:
     - `prediction`
     - `recommendation`

---

## 3. Thiết kế database (mở rộng)

### 3.1. Bảng `exercise_templates`

Thư viện bài tập mẫu, dùng để build kế hoạch recommendation.

| Cột             | Kiểu           | Mô tả                                          |
|-----------------|----------------|-----------------------------------------------|
| id              | UUID (PK)      | Khóa chính                                    |
| name            | text           | Tên bài tập (VD: Đi bộ nhanh)                |
| intensity       | text           | low / medium / high                          |
| description     | text           | Mô tả chi tiết bài tập                        |
| duration_min    | int            | Thời lượng khuyến nghị (phút/buổi)           |
| freq_per_week   | int            | Số buổi/tuần khuyến nghị                      |
| target_risk_min | text (optional)| Nhóm nguy cơ tối thiểu (VD: low)             |
| target_risk_max | text (optional)| Nhóm nguy cơ tối đa (VD: medium)             |
| tags            | text[] / jsonb | Nhãn: cardio, flexibility, strength, senior… |

> Giai đoạn 1 có thể chỉ cần `name`, `intensity`, `duration_min`, `freq_per_week`, `description`.

---

### 3.2. Bảng `exercise_recommendations`

Lưu kế hoạch tập luyện được generate cho từng lần dự đoán.

| Cột           | Kiểu           | Mô tả                                        |
|--------------|----------------|---------------------------------------------|
| id           | UUID (PK)      | Khóa chính                                  |
| patient_id   | UUID (FK)      | Liên kết với `patients.id`                  |
| prediction_id| UUID (FK)      | Liên kết với `predictions.id`               |
| plan         | jsonb          | Nội dung chi tiết kế hoạch gợi ý            |
| created_at   | timestamptz    | Thời điểm tạo gợi ý                         |

Ví dụ nội dung `plan` (json):

```json
{
  "summary": "Bạn nên tập thể dục cường độ nhẹ đến vừa 5 buổi/tuần.",
  "items": [
    {
      "name": "Đi bộ nhanh",
      "intensity": "low",
      "duration_min": 30,
      "freq_per_week": 5,
      "notes": "Theo dõi nhịp tim, không để vượt quá 120 bpm."
    },
    {
      "name": "Yoga hít thở",
      "intensity": "low",
      "duration_min": 20,
      "freq_per_week": 3,
      "notes": "Tập vào buổi sáng hoặc tối, tránh sau ăn no."
    }
  ]
}
```

---

## 4. Thiết kế module backend

### 4.1. Vị trí module

Có hai lựa chọn:

1. **Tạo module mới**: `modules/recommendations`
2. **Mở rộng module `predictions`** để xử lý luôn recommendation

Giai đoạn đầu có thể gộp chung trong `predictions`, sau tách ra `recommendations` nếu logic phức tạp.

### 4.2. Trách nhiệm chính

- **Module predictions**:
  - Gọi ML service để dự đoán.
  - Lưu bản ghi vào bảng `predictions`.
  - Gọi logic recommendation.

- **Module recommendations** (hiện có thể chỉ là một service/internal package):
  - Nhận input: risk_level, age, BMI, active, smoke, alco, …
  - Áp dụng rule-based để chọn các bài tập từ `exercise_templates`.
  - Build object `plan` (struct/JSON).
  - Lưu `exercise_recommendations` vào DB.
  - Trả `plan` cho controller.

---

## 5. Thiết kế API

### 5.1. Mở rộng API `POST /api/patients/:id/predict`

**Request**: giữ nguyên như hiện tại (các chỉ số sức khỏe cần cho ML).

**Response (mới):**

```json
{
  "prediction": {
    "id": "uuid",
    "probability": 0.82,
    "risk_level": "high",
    "label": 1,
    "created_at": "2025-01-01T10:00:00Z"
  },
  "recommendation": {
    "summary": "Bạn nên tập cường độ nhẹ 5 buổi/tuần.",
    "items": [
      {
        "name": "Đi bộ nhanh",
        "intensity": "low",
        "duration_min": 30,
        "freq_per_week": 5,
        "notes": "Theo dõi nhịp tim, không quá 120 bpm."
      },
      {
        "name": "Yoga thư giãn",
        "intensity": "low",
        "duration_min": 20,
        "freq_per_week": 3,
        "notes": "Hít thở sâu, chậm, đều."
      }
    ]
  }
}
```

> Lưu ý: `prediction.id` có thể dùng để liên kết với `exercise_recommendations.prediction_id`.

---

### 5.2. API lịch sử gợi ý (tùy chọn giai đoạn 1)

- `GET /api/patients/:id/recommendations`
  - Trả danh sách các recommendation đã sinh ra (kèm prediction, created_at).
- `GET /api/patients/:id/predictions` (hiện đã có)
  - Có thể mở rộng để trả kèm tóm tắt recommendation mới nhất.

---

## 6. Luồng xử lý chi tiết (Sequence)

### 6.1. Luồng: Predict + Recommendation

1. **Client** gửi `POST /api/patients/:id/predict` với các chỉ số sức khỏe.
2. **Auth middleware**:
   - Decode JWT.
   - Lấy `userID`.
3. **Backend (predictions controller)**:
   - Check quyền sở hữu bệnh nhân: `patients.user_id == userID`.
   - Build payload gửi sang ML service.

4. **Gọi ML service (FastAPI)**:
   - `POST http://ml-service:8000/predict`
   - Nhận `probability`, `label`, `risk_level`.

5. **Lưu prediction**:
   - Gọi DB tạo bản ghi trong `predictions`.

6. **Gọi recommendation service**:
   - Input:
     - risk_level
     - tuổi (age_years)
     - BMI
     - active, smoke, alco, …
   - Rule-based logic chọn bài tập từ `exercise_templates`.
   - Build `plan` JSON.
   - Lưu vào `exercise_recommendations` (patient_id, prediction_id, plan).

7. **Trả response** về client:
   - `prediction` + `recommendation`.

---

## 7. Rule-based recommendation (Giai đoạn 1)

### 7.1. Các input chính

- `age_years`
- `bmi`
- `risk_level` (low / medium / high)
- `active` (0/1)
- `smoke` (0/1)
- `alco` (0/1)

### 7.2. Ví dụ luật đơn giản

- Nếu `risk_level = high` hoặc `age_years > 60`:
  - Chỉ chọn bài tập `intensity = low`.
  - Giới hạn duration 20–30 phút/buổi.
  - Khuyến cáo “tham khảo thêm ý kiến bác sĩ”.

- Nếu `risk_level = medium` & BMI > 25:
  - Chọn kết hợp bài tập `low` và một phần `medium`.
  - Tập trung vào **cardio nhẹ**: đi bộ nhanh, đạp xe nhẹ.

- Nếu `risk_level = low` & active = 1:
  - Có thể chọn thêm `intensity = medium`.
  - Gợi ý thêm chạy bộ nhẹ, bơi.

### 7.3. Triển khai

- Định nghĩa các rule trong code (Go) dưới dạng:
  - Cấu trúc `if/else` rõ ràng, hoặc
  - Một bảng rule (slice/map) và logic matching.
- Lấy danh sách `exercise_templates` từ DB hoặc hard-code (giai đoạn đầu).
- Lắp ghép thành `plan`.

---

## 8. Mở rộng trong tương lai

### 8.1. Thu thập feedback

- Cho phép người dùng:
  - Đánh dấu “Đã tập / Chưa tập / Quá nặng / Quá nhẹ”.
  - Lưu feedback vào bảng riêng (vd: `exercise_feedback`).

### 8.2. Recommender ML

- Dùng data lịch sử:
  - hồ sơ bệnh nhân,
  - lịch sử prediction,
  - lịch sử kế hoạch tập luyện,
  - feedback,
- Train model ML để:
  - Gợi ý kế hoạch tập luyện cá nhân hóa hơn.
  - Đánh giá hiệu quả (giảm nguy cơ tim mạch sau thời gian).

- Thêm một ML service mới:
  - `/ml/recommend` → trả về exercise plan.
- Backend Go:
  - Chuyển từ rule-based sang gọi ML recommender (hoặc hybrid).

---

## 9. Tóm tắt

Module “Exercise Recommendation” giúp hệ thống:

- Chuyển từ **dự đoán thụ động** sang **gợi ý hành động cụ thể** cho người dùng.
- Giai đoạn đầu dùng logic rule-based đơn giản, dễ triển khai.
- Kiến trúc được thiết kế để sẵn sàng:
  - Lưu lịch sử kế hoạch tập luyện.
  - Thu thập feedback.
  - Nâng cấp lên mô hình Machine Learning recommender trong tương lai.

Tài liệu này nhằm:
- Làm rõ **schema DB** mới.
- Mô tả **luồng API** và **module backend** liên quan.
- Làm nền tảng để implement trong Go backend, React frontend, và mở rộng ML service.
