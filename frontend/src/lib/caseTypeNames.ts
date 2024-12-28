export type CaseType = keyof typeof caseTypeNameMap
export const caseTypeNameMap = {
    aux1: "Auxiliar Primero",
    aux2: "Auxiliar Segundo",
    civ2: "Civil Segundo",
    civ3: "Civil Tercero",
    civ4: "Civil Cuarto",
    fam1: "Familiar Primero",
    fam2: "Familiar Segundo",
    fam3: "Familiar Tercero",
    fam4: "Familiar Cuarto",
    fam5: "Familiar Quinto",
    mer1: "Mercantil Primero",
    mer2: "Mercantil Segundo",
    mer3: "Mercantil Tercero",
    mer4: "Mercantil Cuarto",
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
