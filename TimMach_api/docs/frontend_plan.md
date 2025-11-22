# Kế hoạch FE (React + TS) bám theo API TimMach_api

## 0) Stack & setup nhanh
- Stack: Vite + React 18/19 + TS, React Router v7, React Query, Axios, Tailwind + shadcn/ui, react-hook-form + Zod, Recharts.
- Khởi tạo: `npm create vite@latest my-app -- --template react-ts`.
- Cài lib chính: `npm i react-router-dom@7 axios @tanstack/react-query tailwindcss postcss autoprefixer class-variance-authority tailwind-merge lucide-react`.
- Cài dev: `npm i -D @types/node`.
- Tailwind: `npx tailwindcss init -p`; chỉnh `tailwind.config.cjs` content: `["./index.html","./src/**/*.{ts,tsx}"]`; thêm plugin từ shadcn nếu dùng.
- shadcn/ui: cài CLI, generate các component base (Button, Input, Card, Table, Dialog...) để dùng lại.
- React Query: wrap `QueryClientProvider` tại `main.tsx`, set `retry` phù hợp (thường 1–2 lần) và `staleTime` ngắn.
- Axios: `api/client.ts` có interceptor gắn `Authorization: Bearer <token>` từ `AuthContext/storage`; interceptor response bắt 401 → logout.
- react-hook-form + Zod: `npm i react-hook-form zod @hookform/resolvers`; mỗi form dùng `zodResolver(schema)`.
- Recharts: `npm i recharts` để vẽ history chart.

## 1) Trang (pages) & API cần gọi
- Login `/login`
  - Form: `email`, `password`
  - POST `/api/users/login` → lưu `access_token` vào `localStorage`, set vào AuthContext, redirect `/dashboard`.
- Register `/register` (nếu mở đăng ký)
  - Form: `email`, `password`, `confirmPassword`
  - POST `/api/users/register` → nhận `token` (same JWT) → lưu & chuyển hướng (tự login hoặc sang `/login`).
- Dashboard `/dashboard`
  - Overview: số bệnh nhân, số lần predict gần đây, danh sách 5 bệnh nhân mới.
  - API: GET `/api/patients?limit=5` (backend chưa có `/stats`, có thể thêm sau).
- Patients list `/patients`
  - Bảng: name, gender, dob → tính `age`, số lần dự đoán (hiện backend chưa trả count; muốn có thì cần tính thêm từ history).
  - Search/filter (client-side hoặc thêm query nếu backend bổ sung), pagination bằng `limit`/`offset`.
  - API: GET `/api/patients?limit=&offset=`.
  - Nút “Create patient”.
- Patient form
  - Tạo mới: `/patients/new` → POST `/api/patients`
  - Sửa: `/patients/:id/edit` → PUT/PATCH `/api/patients/:id`
  - Fields theo backend: `name` (string), `gender` (int16), `dob` (`yyyy-mm-dd`).
- Patient detail `/patients/:id`
  - Thông tin bệnh nhân: GET `/api/patients/:id`
  - Latest prediction: GET `/api/patients/:id/predictions?limit=1`
  - Buttons: “Predict now” → `/patients/:id/predict`; “View history” → `/patients/:id/history`.
- Predict `/patients/:id/predict`
  - Form fields (trùng `CreatePredictionRequest`):
    - `age_years` (float, có thể auto tính từ dob), `gender` (int), `height`, `weight`, `ap_hi`, `ap_lo`, `cholesterol`, `gluc`, `smoke`, `alco`, `active`.
  - Submit → POST `/api/patients/:id/predict`
  - Hiển thị: `probability`, `risk_level`, và recommendation; màu low/medium/high.
- History `/patients/:id/history`
  - GET `/api/patients/:id/predictions?limit=&offset=`
  - Bảng/line chart: thời gian – probability, màu theo `risk_label`.
- Profile `/profile` (optional)
  - GET `/api/users/me` → show `email`, `created_at`.
  - Đổi mật khẩu cần backend bổ sung.
- Logout
  - Xóa token khỏi storage + context, chuyển `/login`.

## 2) Components nên tách
- Layout: `components/layout/MainLayout`, `Header` (logo, avatar, dropdown logout), `Sidebar` (Dashboard, Patients, Profile).
- Auth: `components/auth/LoginForm`, `RegisterForm`.
- Patients: `PatientsTable`, `PatientForm`, `PatientSummaryCard`.
- Predictions: `PredictForm`, `PredictionResultCard`, `PredictionHistoryTable`, `PredictionChart`.
- UI chung: `components/ui/{Button,Input,Select,Card,Modal,...}`.
- Context: `context/AuthContext` giữ `token`, `user`, `login`, `logout`; persist token bằng `utils/storage.ts`.

## 3) Cấu trúc thư mục gợi ý (Vite + React Router + TS)
```
src/
  api/
    client.ts          // fetch/axios base, inject Bearer token
    users.ts           // /api/users/*
    patients.ts        // /api/patients/*
    predictions.ts     // /api/patients/:id/predict, /predictions
  pages/
    auth/{LoginPage.tsx, RegisterPage.tsx}
    dashboard/DashboardPage.tsx
    patients/{PatientsListPage.tsx, PatientDetailPage.tsx, PatientFormPage.tsx, PatientPredictPage.tsx, PatientHistoryPage.tsx}
    profile/ProfilePage.tsx
  components/
    layout/{MainLayout.tsx, Sidebar.tsx, Header.tsx}
    ui/{...}
    patients/{PatientsTable.tsx, PatientForm.tsx, PatientSummaryCard.tsx}
    predictions/{PredictForm.tsx, PredictionResultCard.tsx, PredictionHistoryTable.tsx, PredictionChart.tsx}
  hooks/{useAuth.ts, usePatients.ts, usePredictions.ts}
  context/AuthContext.tsx
  routes/index.tsx
  types/{api.ts, common.ts}
  utils/{storage.ts, format.ts}
  App.tsx
  main.tsx
```

## 4) Model & payload (cho `types/api.ts`)
- User
  - Login request: `{email, password}` → response `{access_token, user: {id, email, created_at}}`
  - Register response: `{user, token}`
  - Me response: `{id, email, created_at}`
- Patient
  - Create: `{name, gender, dob: "YYYY-MM-DD"}`
  - Update: cùng field trên (optional) via PUT/PATCH.
  - Response: `{id, user_id, name, gender, dob, created_at}`
  - List response: `{patients: PatientResponse[]}`
- Prediction
  - Create request: `{age_years, gender, height, weight, ap_hi, ap_lo, cholesterol, gluc, smoke, alco, active}`
  - Create response: `{prediction: PredictionResponse, recommendation: RecommendationPlan}`
  - List response: `{predictions: PredictionResponse[]}`

## 5) Luồng auth chuẩn
- Sau `login/register`: lưu token vào `localStorage`, set `Authorization: Bearer <token>` cho mọi request (axios interceptor hoặc fetch wrapper).
- Khi 401: logout và chuyển tới `/login`.
- Không để token trong query string; dùng HTTPS trong môi trường thật.

## 6) Ghi chú triển khai nhanh
- Tính tuổi: từ `dob` → hiển thị `age_years` và auto-fill vào predict form (cho phép chỉnh).
- Gender/enum: backend dùng int; FE map 0/1/2 thành label (ví dụ 0=unknown, 1=male, 2=female nếu theo dữ liệu ML).
- Risk color: low=green, medium=amber, high=red.
- Dashboard số liệu: hiện chỉ có patients list/history; nếu cần thống kê thêm, bổ sung API `/api/stats` hoặc tính client-side từ list.
