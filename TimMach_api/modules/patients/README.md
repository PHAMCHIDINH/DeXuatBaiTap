# Module `patients`

Quản lý hồ sơ bệnh nhân theo từng user (bác sĩ hoặc chính bệnh nhân).

## Routes (đều cần Auth middleware)
- `POST /api/patients` – tạo bệnh nhân
- `GET /api/patients` – danh sách của user (có `limit/offset`)
- `GET /api/patients/:id` – chi tiết
- `PUT|PATCH /api/patients/:id` – cập nhật
- `DELETE /api/patients/:id` – xoá

## Request/Response
- Create: `{"name","gender","dob"}` → 201 `PatientResponse`
- Update: các field optional → 200 `PatientResponse`
- List: query `limit,offset` → 200 `{"patients":[PatientResponse]}`

`PatientResponse`:
```json
{"id":"uuid","user_id":"uuid","name":"...","gender":1,"dob":"yyyy-mm-dd","created_at":"..."}
```

## Logic (controllers)
- Lấy `userID` từ context (middleware).
- Create: validate body → parse `dob` → `CreatePatient`.
- List: `ListPatientsByUser(userID, limit, offset)`.
- Detail/Update/Delete: parse `patientID` → `GetPatientByID` → kiểm tra sở hữu `patient.user_id == userID` → thao tác tiếp (`UpdatePatient`, `DeletePatient`).

## Phụ thuộc
- DB queries: `CreatePatient`, `ListPatientsByUser`, `GetPatientByID`, `UpdatePatient`, `DeletePatient`.
- Utils: lấy user từ context, parse UUID/date, trả JSON lỗi.

## Lỗi chuẩn
`{"error": "missing user in context"}`, `{"error": "invalid patient id"}`, `{"error": "patient does not belong to user"}`, ...
