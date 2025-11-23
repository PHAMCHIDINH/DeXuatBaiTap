import client from '../../api/client';
import { User } from '../../types/api';

export async function me(): Promise<User> {
  const { data } = await client.get<User>('/users/me');
  return data;
}
