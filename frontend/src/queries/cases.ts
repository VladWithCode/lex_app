import { useMutation, useQuery } from "@tanstack/react-query";
import { CreateCase, FindAllCases, FindCaseById, FindCaseWithAccords, UpdateCase } from "../../wailsjs/go/controllers/CaseController"
import { db } from "../../wailsjs/go/models";

export type CaseFilters = {
    caseType?: string;
    caseId?: string;
    caseNo?: string;
    caseYear?: string;
    alias?: string;
    limit: number;
}

const caseKeys = {
    all: ["cases"] as const,
    lists: () => [...caseKeys.all, "list"] as const,
    list: (filters?: CaseFilters) => [...caseKeys.lists(), filters] as const,
    //listWith: (filters: CaseFilters) => [...caseKeys.lists(), filters] as const,
    details: () => [...caseKeys.all, "detail"] as const,
    detail: (id: string) => [...caseKeys.details(), id] as const,

    detailsAndAccords: () => [...caseKeys.details(), "accords"] as const,
    detailAndAccords: (id: string, accordCount: number) => [...caseKeys.detailsAndAccords(), id, accordCount] as const
}

export function useCases(filters: CaseFilters) {
    return useQuery({
        queryKey: caseKeys.list(filters),
        queryFn: async () => {
            return await new Promise<db.lexCase[]>((resolve) => {
                setTimeout(async () => {
                    resolve(await FindAllCases())
                }, 1500)
            })
        },
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
