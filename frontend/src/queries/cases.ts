import { useMutation, useQuery } from "@tanstack/react-query";
import { CreateCase, FindCaseById, FindCases, FindCaseWithAccords, UpdateCase } from "../../wailsjs/go/controllers/CaseController"
import { db } from "../../wailsjs/go/models";

export type FindCaseOptions = Partial<db.lexFindCaseOptions>

const caseKeys = {
    all: ["cases"] as const,
    lists: () => [...caseKeys.all, "list"] as const,
    list: (filters?: FindCaseOptions) => [...caseKeys.lists(), filters] as const,
    //listWith: (filters: CaseFilters) => [...caseKeys.lists(), filters] as const,
    details: () => [...caseKeys.all, "detail"] as const,
    detail: (id: string) => [...caseKeys.details(), id] as const,

    detailsAndAccords: () => [...caseKeys.details(), "accords"] as const,
    detailAndAccords: (id: string, accordCount: number) => [...caseKeys.detailsAndAccords(), id, accordCount] as const
}

export function useCases(filters: FindCaseOptions) {
    return useQuery({
        queryKey: caseKeys.list(filters),
        queryFn: async () => {
            return await FindCases(filters as db.lexFindCaseOptions)
        }
    })
}

export function useCase(id: string) {
    return useQuery({
        queryKey: caseKeys.detail(id),
        queryFn: async () => {
            return await FindCaseById(id)
        }
    })
}

export function useCaseWithAccords(id: string, accordCount: number) {
    return useQuery({
        queryKey: caseKeys.detailAndAccords(id, accordCount),
        queryFn: async ({ queryKey }) => {
            const [_key, _, __, id, count] = queryKey
            return await FindCaseWithAccords(id, count)
        }
    })
}

type CreateCaseParams = {
    caseId: string;
    caseType: string;
}
export function useCreateCase() {
    return useMutation({
        mutationFn: ({ caseId, caseType }: CreateCaseParams) => {
            return CreateCase(caseId, caseType)
        }
    })
}

type UpdateCaseParams = {
    id: string;
    caseData: Partial<db.lexCase>;
}
export function useUpdateCase() {
    return useMutation({
        mutationFn: ({ caseData, id }: UpdateCaseParams) => {
            return UpdateCase(id, caseData as db.lexCase)
        }
    })
}
