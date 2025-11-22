export interface User {
  id: string;
  email: string;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  user: User;
}

export interface RegisterResponse {
  token: string;
  user: User;
}

export interface CreatePatientRequest {
  name: string;
  gender: number;
  dob: string; // YYYY-MM-DD
}

export interface UpdatePatientRequest {
  name?: string;
  gender?: number;
  dob?: string;
}

export interface PatientResponse {
  id: string;
  user_id: string;
  name: string;
  gender: number;
  dob: string;
  created_at: string;
  latest_prediction?: PatientLatestPrediction | null;
}

export interface ListPatientsResponse {
  patients: PatientResponse[];
}

export interface ListPatientsParams {
  limit?: number;
  offset?: number;
  risk?: string;
}

export interface CreatePredictionRequest {
  age_years: number;
  gender: number;
  height: number;
  weight: number;
  ap_hi: number;
  ap_lo: number;
  cholesterol: number;
  gluc: number;
  smoke: number;
  alco: number;
  active: number;
}

export interface PredictionResponse {
  id: string;
  patient_id: string;
  probability: number;
  risk_label: string;
  model_version?: string;
  raw_features?: Record<string, unknown> | null;
  created_at: string;
}

export interface ListPredictionsResponse {
  predictions: PredictionResponse[];
}

export interface ListPredictionsParams {
  limit?: number;
  offset?: number;
}

export interface RecommendationItem {
  name: string;
  intensity: string;
  duration_min: number;
  freq_per_week: number;
  notes?: string;
}

export interface RecommendationPlan {
  summary: string;
  items: RecommendationItem[];
}

export interface CreatePredictionResponse {
  prediction: PredictionResponse;
  recommendation: RecommendationPlan;
}

export interface PatientLatestPrediction {
  probability: number;
  risk_label: string;
  created_at: string;
}

export interface StatsRiskCount {
  risk_label: string;
  count: number;
}

export interface StatsResponse {
  total_patients: number;
  risk_counts: StatsRiskCount[];
}

export interface TemplateResponse {
  id: number;
  name: string;
  intensity: string;
  description: string;
  duration_min: number;
  freq_per_week: number;
  target_risk_level: string;
  tags: string[];
}

export interface CreateTemplateRequest {
  name: string;
  intensity: string;
  description: string;
  duration_min: number;
  freq_per_week: number;
  target_risk_level: string;
  tags?: string[];
}

export interface ListTemplatesResponse {
  templates: TemplateResponse[];
}

export interface RecommendationResponse {
  id: number;
  patient_id: string;
  prediction_id: string;
  plan?: RecommendationPlan | null;
  created_at: string;
}

export interface ListRecommendationsResponse {
  recommendations: RecommendationResponse[];
}
