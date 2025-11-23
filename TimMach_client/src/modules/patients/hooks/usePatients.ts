import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import * as api from '../api';
import {
  CreatePatientRequest,
  ListPatientsParams,
  ListPatientsResponse,
  PatientResponse,
  UpdatePatientRequest,
} from '../../../types/api';

export function usePatientsList(params: ListPatientsParams = {}) {
  return useQuery<ListPatientsResponse>({
    queryKey: ['patients', params.limit ?? null, params.offset ?? null, params.risk ?? null],
    queryFn: () => api.listPatients(params),
  });
}

export function usePatient(id?: string) {
  return useQuery<PatientResponse>({
    queryKey: ['patient', id],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return api.getPatient(id);
    },
    enabled: !!id,
  });
}

export function useCreatePatient() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreatePatientRequest) => api.createPatient(payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['patients'] });
    },
  });
}

export function useUpdatePatient(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdatePatientRequest) => api.updatePatient(id, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['patients'] });
      qc.invalidateQueries({ queryKey: ['patient', id] });
    },
  });
}

export function useDeletePatient() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deletePatient(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['patients'] });
    },
  });
}
