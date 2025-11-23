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

export interface PatientLatestPrediction {
  probability: number;
  risk_label: string;
  created_at: string;
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

export interface ListPatientsParams {
  limit?: number;
  offset?: number;
  risk?: string;
}
