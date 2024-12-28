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

export const tribunalCategoryMap = {
    aux: "Auxiliar",
    civ: "Civil",
    fam: "Familiar",
    mer: "Mercantil",
    sec: "Secretaría Civil",
    cjm: "CJM",
    tri: "Tribunal Laboral",
    oth: "Otro",
}

export type TribunalCategoryOptions = Record<keyof typeof tribunalCategoryMap, { title: string, elements: { val: string, label: string }[] }>
export function getTribunalCategoryOptions(): TribunalCategoryOptions {
    let options: TribunalCategoryOptions = {
        aux: { title: tribunalCategoryMap.aux, elements: [] },
        civ: { title: tribunalCategoryMap.civ, elements: [] },
        fam: { title: tribunalCategoryMap.fam, elements: [] },
        mer: { title: tribunalCategoryMap.mer, elements: [] },
        sec: { title: tribunalCategoryMap.sec, elements: [] },
        cjm: { title: tribunalCategoryMap.cjm, elements: [] },
        tri: { title: tribunalCategoryMap.tri, elements: [] },
        oth: { title: tribunalCategoryMap.oth, elements: [] },
    }

    for (let [k, v] of Object.entries(caseTypeNameMap)) {
        let trib = k.slice(0, 3) as keyof typeof tribunalCategoryMap
        if (tribunalCategoryMap[trib]) {
            options[trib].elements.push({ val: k, label: v });
        } else {
            options.oth.elements.push({ val: k, label: v });
        }
    }
    return options
}
