export type CaseType = keyof typeof caseTypeNameMap
export const caseTypeNameMap = {
    aux1: "Auxiliar 1",
    aux2: "Auxiliar 2",
    civ2: "Civil 2",
    civ3: "Civil 3",
    civ4: "Civil 4",
    fam1: "Familiar 1",
    fam2: "Familiar 2",
    fam3: "Familiar 3",
    fam4: "Familiar 4",
    fam5: "Familiar 5",
    mer1: "Mercantil 1",
    mer2: "Mercantil 2",
    mer3: "Mercantil 3",
    mer4: "Mercantil 4",
    seccc: "Civil Colegiada",
    seccu: "Civil Unitaria",
    cjmf: "CJM",
    cjmf2: "CJM 2",
    tribl: "Tribunal Laboral",
} as const

export function caseTypeToName(ct: CaseType): string {
    let name = caseTypeNameMap[ct]

    return name || "Otro"
}
