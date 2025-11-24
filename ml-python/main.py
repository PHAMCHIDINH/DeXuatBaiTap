import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
import joblib
import pandas as pd
import traceback
from typing import List, Optional

app = FastAPI(
    title="Heart Disease Risk Prediction API",
    description="API dự đoán nguy cơ bệnh tim mạch sử dụng mô hình RandomForest đã train và lưu trong heart_disease_model.pkl",
    version="1.0.0",
)

# Đường dẫn tới file model đã lưu từ notebook (dùng đường dẫn tuyệt đối để tránh lỗi khi đổi cwd)
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
MODEL_PATH = os.path.join(BASE_DIR, "training", "heart_disease_model.pkl")

# Load model khi khởi động server
try:
    model = joblib.load(MODEL_PATH)
except Exception as e:
    print(f"Không thể load model từ {MODEL_PATH}: {e}")
    model = None


class HeartInput(BaseModel):
    """
    Đầu vào dự đoán nguy cơ bệnh tim mạch.

    Lưu ý:
    - age_years: tuổi theo NĂM (ví dụ 45.5)
    - gender: 1 = Nữ, 2 = Nam (giống dataset gốc)
    - cholesterol: 1 (tốt), 2 (trung bình), 3 (xấu)
    - gluc: 1 (bình thường), 2 (cao), 3 (rất cao)
    - smoke, alco, active: 0 = Không, 1 = Có
    """
    age_years: float = Field(..., example=54)
    gender: int = Field(..., ge=1, le=2, example=2)
    height: float = Field(..., example=165)  # cm
    weight: float = Field(..., example=70)   # kg
    ap_hi: int = Field(..., example=130)     # huyết áp tâm thu
    ap_lo: int = Field(..., example=80)      # huyết áp tâm trương
    cholesterol: int = Field(..., ge=1, le=3, example=2)
    gluc: int = Field(..., ge=1, le=3, example=1)
    smoke: int = Field(..., ge=0, le=1, example=0)
    alco: int = Field(..., ge=0, le=1, example=0)
    active: int = Field(..., ge=0, le=1, example=1)


class RiskFactor(BaseModel):
    field: str
    status: str
    message: str
    contribution: Optional[float] = None


class HeartPrediction(BaseModel):
    probability: float
    label: int
    risk_level: str
    factors: List[RiskFactor] = []


@app.get("/health")
def health_check():
    """
    Endpoint kiểm tra trạng thái service.
    """
    if model is None:
        return {"status": "error", "detail": "Model not loaded"}
    return {"status": "ok"}


@app.post("/predict", response_model=HeartPrediction)
def predict_heart_disease(data: HeartInput):
    """
    Dự đoán nguy cơ bệnh tim mạch.

    Trả về:
    - probability: Xác suất (0-1) bị bệnh (cardio=1)
    - label: 0 hoặc 1
    - risk_level: low / medium / high
    - factors: Danh sách các yếu tố nguy cơ kèm mô tả
    """
    if model is None:
        raise HTTPException(status_code=500, detail="Model chưa được load. Kiểm tra lại file heart_disease_model.pkl")

    try:
        # Chuyển dữ liệu đầu vào thành dict
        input_dict = data.dict()

        # Tính thêm các feature giống như lúc train:
        # bmi = weight / (height/100)^2
        height_m = input_dict["height"] / 100.0
        bmi = input_dict["weight"] / (height_m ** 2)

        # bp_ratio = ap_hi / ap_lo
        if input_dict["ap_lo"] == 0:
            # tránh chia cho 0 – có thể raise lỗi hoặc set giá trị mặc định
            raise ValueError("ap_lo (huyết áp tâm trương) không được bằng 0.")
        bp_ratio = input_dict["ap_hi"] / input_dict["ap_lo"]

        # Tạo row với đúng các cột đã dùng khi train trong notebook:
        # ['age_years','gender','height','weight','ap_hi','ap_lo',
        #  'cholesterol','gluc','smoke','alco','active','bmi','bp_ratio']
        row = {
            "age_years": input_dict["age_years"],
            "gender": input_dict["gender"],
            "height": input_dict["height"],
            "weight": input_dict["weight"],
            "ap_hi": input_dict["ap_hi"],
            "ap_lo": input_dict["ap_lo"],
            "cholesterol": input_dict["cholesterol"],
            "gluc": input_dict["gluc"],
            "smoke": input_dict["smoke"],
            "alco": input_dict["alco"],
            "active": input_dict["active"],
            "bmi": bmi,
            "bp_ratio": bp_ratio,
        }

        # Tạo DataFrame 1 dòng
        X = pd.DataFrame([row])

        # Dự đoán
        proba = model.predict_proba(X)[0, 1]
        label = int(model.predict(X)[0])

        # Map risk level theo xác suất (tuỳ bạn chỉnh lại ngưỡng)
        if proba < 0.33:
            risk_level = "low"
        elif proba < 0.66:
            risk_level = "medium"
        else:
            risk_level = "high"

        factors = []
        if input_dict["cholesterol"] >= 3:
            factors.append(RiskFactor(field="cholesterol", status="high", message="Cholesterol cao (>= 3)"))
        if input_dict["gluc"] >= 3:
            factors.append(RiskFactor(field="gluc", status="high", message="Đường huyết cao (>= 3)"))
        if input_dict["ap_hi"] > 140 or input_dict["ap_lo"] > 90:
            factors.append(RiskFactor(field="blood_pressure", status="high", message="Huyết áp cao (>140/90)"))
        if bmi >= 30:
            factors.append(RiskFactor(field="bmi", status="high", message="BMI cao (>= 30)"))
        if input_dict["smoke"] == 1:
            factors.append(RiskFactor(field="smoke", status="yes", message="Hút thuốc"))
        if input_dict["alco"] == 1:
            factors.append(RiskFactor(field="alco", status="yes", message="Uống rượu"))
        if input_dict["active"] == 0:
            factors.append(RiskFactor(field="active", status="low", message="Ít vận động"))

        return HeartPrediction(
            probability=float(round(proba, 4)),
            label=label,
            risk_level=risk_level,
            factors=factors,
        )

    except HTTPException:
        # Bắn thẳng ra ngoài nếu đã là HTTPException
        raise
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Prediction error: {e}")


# Để chạy server: uvicorn main:app --reload
if __name__ == "__main__":
    import uvicorn

    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
