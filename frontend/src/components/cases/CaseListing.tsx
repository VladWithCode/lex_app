import CaseCard from "./CaseCard"
import { useCases } from "../../queries/cases"
import { LucideLoader } from "lucide-react"

export default function CaseListing() {
    const { data, isSuccess, isLoading, isError } = useCases({ Limit: 5 })

    if (isError) {
        return <div className="flex h-36 items-center justify-center">
            <p className="text-stone-200 text-xl font-semibold">Ocurrio un error al recuperar los casos</p>
        </div>
    }

    if (isLoading) {
        return <div className="flex flex-col h-36 items-center justify-center">
            <LucideLoader className="animate-spin" />
            <p className="text-lg text-stone-400 uppercase font-semibold">Cargando Casos Recientes...</p>
        </div>
    }

    return <div className="flex gap-4 max-w-full overflow-visible mb-2">
        {isSuccess && data.map(c => <CaseCard key={c.id} caseData={c} />)}
    </div>
}
