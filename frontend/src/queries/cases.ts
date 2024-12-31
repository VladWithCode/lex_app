import { useMutation, useQuery } from "@tanstack/react-query";
import { CreateCase, FindCaseById, FindCases, FindCaseWithAccords, UpdateCase } from "../../wailsjs/go/controllers/CaseController"
import { db } from "../../wailsjs/go/models";

export type FindCaseOptions = Partial<db.FindCaseOptions> & {
    search?: string;
}

const caseQueryKeys = {
    all: ["cases"] as const,
    lists: () => [...caseQueryKeys.all, "list"] as const,
    list: (filters?: FindCaseOptions) => [...caseQueryKeys.lists(), filters] as const,
    //listWith: (filters: CaseFilters) => [...caseKeys.lists(), filters] as const,
    details: () => [...caseQueryKeys.all, "detail"] as const,
    detail: (id: string) => [...caseQueryKeys.details(), id] as const,

    detailsAndAccords: () => [...caseQueryKeys.details(), "accords"] as const,
    detailAndAccords: (id: string, accordCount: number) => [...caseQueryKeys.detailsAndAccords(), id, accordCount] as const
}

export function useCases(filters: FindCaseOptions) {
    return useQuery({
        queryKey: caseQueryKeys.list(filters),
        queryFn: async () => {
            return await FindCases(filters as db.FindCaseOptions)
        }
    })
}

export function useCase(id: string) {
    return useQuery({
        queryKey: caseQueryKeys.detail(id),
        queryFn: async () => {
            return await FindCaseById(id)
        }
    })
}

export function useCaseWithAccords(id: string, accordCount: number) {
    return useQuery({
        queryKey: caseQueryKeys.detailAndAccords(id, accordCount),
        queryFn: async ({ queryKey }) => {
            const [_key, _, __, id, count] = queryKey
            return await FindCaseWithAccords(id, count)
        }
    })
}

type CreateCaseParams = {
    caseId: string;
    caseType: string;
    alias?: string;
}
export function useCreateCase() {
    return useMutation({
        mutationFn: ({ caseId, caseType, alias }: CreateCaseParams) => {
            return CreateCase(caseId, caseType, alias || "")
        }
    })
}

type UpdateCaseParams = {
    id: string;
    caseData: Partial<db.LexCase>;
}
export function useUpdateCase() {
    return useMutation({
        mutationFn: ({ caseData, id }: UpdateCaseParams) => {
            return UpdateCase(id, new db.LexCase(caseData))
        }
    })
}
