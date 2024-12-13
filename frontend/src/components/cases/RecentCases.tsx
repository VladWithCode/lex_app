import CaseCard from "./CaseCard"
import { useCases } from "../../queries/cases"
import { LucideLoader } from "lucide-react"
import { CardContent } from "../ui/card"
import { formatDateToShortReadable } from "../../lib/formatUtils"

export default function RecentCases() {
    const { data, isSuccess, isLoading, isError } = useCases({ Limit: 5, IncludeAccords: true, MaxAccords: 1 })

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
        {isSuccess && data.map(c => (
            <CaseCard key={c.id} caseData={c}>
                <CardContent className="p-2 space-y-2">
                    <p className="text-xl text-stone-200 line-clamp-3 overflow-clip text-ellipsis">
                        {c.accords[0].content}
                    </p>
                    <p className="text-stone-400 font-semibold">
                        Fecha de Acuerdo: <span className="capitalize">{
                            formatDateToShortReadable(new Date(c.accords[0].date))
                        }</span>
                    </p>
                </CardContent>
            </CaseCard>
        ))}
    </div>
}
