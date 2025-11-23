export interface StatsRiskCount {
  risk_label: string;
  count: number;
}

export interface StatsResponse {
  total_patients: number;
  risk_counts: StatsRiskCount[];
}
