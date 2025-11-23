import client from '../../api/client';
import { StatsResponse } from './types';

export async function getStats(): Promise<StatsResponse> {
  const { data } = await client.get<StatsResponse>('/stats');
  return data;
}
