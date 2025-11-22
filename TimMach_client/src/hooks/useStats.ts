import { useQuery } from '@tanstack/react-query';
import { getStats } from '../api/stats';
import { StatsResponse } from '../types/api';

export function useStats() {
  return useQuery<StatsResponse>({
    queryKey: ['stats'],
    queryFn: getStats,
  });
}
