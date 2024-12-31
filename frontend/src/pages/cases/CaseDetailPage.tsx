import CaseAccordCard from "@/components/cases/CaseAccordCard";
import { Separator } from "@/components/ui/separator";
import { CaseType, caseTypeToName } from "@/lib/caseTypeNames";
import { cn } from "@/lib/utils";
import { useCaseWithAccords } from "@/queries/cases";
import { LucideLoader } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";
import { db } from "wailsjs/go/models";

export default function CaseDetailPage() {
    const { caseUUID } = useParams()
    const { data, status } = useCaseWithAccords(String(caseUUID), 15)

    if (status === "pending") {
        return <div className="flex flex-col items-center justify-center gap-4 grow basis-full">
            <LucideLoader size="96" className="animate-spin" />
            <p className="text-xl text-stone-400 uppercase font-semibold">Cargando información del caso...</p>
        </div>
    }

    if (status === "error") {
        return <div className="flex items-center justify-center grow basis-full">
            <p className="text-stone-200 text-4xl font-semibold">Ocurrio un error al recuperar los casos</p>
        </div>
    }

    return (
        <>
            <h1 className="text-6xl">Caso {data.caseId} - {caseTypeToName(data.caseType as CaseType)} | lexApp</h1>
            {data.alias === ""
                || <p className="text-2xl text-stone-200 font-medium py-1"> {data.alias} </p>
            }
            <p className="text-lg text-stone-400 pt-2">Acuerdos y detalles del caso No. {data.caseId} del juzgado {caseTypeToName(data.caseType as CaseType)}</p>
            <Separator className="my-2" />
            <CaseDetails data={data} />
            <Separator className="my-2" />
            <div className="flex flex-col shrink-0 grow gap-2 max-h-full overflow-auto">
                <h2 className="text-2xl text-stone-200">Acuerdos</h2>
                <div className="grid grid-cols-3 gap-4 max-h-full overflow-auto">
                    {
                        data.accords == null || data.accords.length == 0
                            ? <div className="col-span-full self-center justify-self-center py-24">
                                <p className="text-stone-400 text-2xl font-semibold">No se encontraron acuerdos para el caso</p>
                            </div>
                            : data.accords.map(c => (<CaseAccordCard key={c.id} accord={c} />))
                    }
                </div>
            </div>
        </>
    )
}

function CaseDetails({ data }: { data: db.LexCase }) {
    const [showFullId, setShowFullId] = useState(false)
    return (
        <div className="flex flex-col shrink grow-0 gap-2">
            <h2 className="text-2xl text-stone-200">Detalles</h2>
            <div className="flex gap-2 text-base">
                <div
                    className="grow-0 text-stone-300 rounded-sm cursor-pointer hover:bg-stone-800/70"
                    onClick={() => setShowFullId(!showFullId)}>
                    <p className="font-bold my-0">UUID</p>
                    <p className={cn(
                        "text-xl max-w-[10ch] overflow-clip text-ellipsis line-clamp-1",
                        showFullId ? "max-w-full" : ""
                    )}>{data.id}</p>
                </div>
                <Separator className="mx-2" orientation="vertical" />
                <div className="grow-0 text-stone-300">
                    <p className="font-bold my-0">No.</p>
                    <p className="text-xl">{data.caseNo}</p>
                </div>
                <Separator className="mx-2" orientation="vertical" />
                <div className="grow-0 text-stone-300">
                    <p className="font-semibold my-0">Año</p>
                    <p className="text-xl">{data.caseYear}</p>
                </div>
                <Separator className="mx-2" orientation="vertical" />
                <div className="grow-0 text-stone-300">
                    <p className="font-semibold my-0">Alias</p>
                    <p className="text-xl">
                        {
                            data.alias
                            || <span className="block text-center">-</span>
                        }
                    </p>
                </div>
                <Separator className="mx-2" orientation="vertical" />
                <div className="grow-0 text-stone-300">
                    <p className="font-semibold my-0">Numeros de expediente relacionados</p>
                    <p className="text-xl">
                        {
                            data.otherIds && data.otherIds.length > 0
                                ? data.otherIds.join(", ")
                                : <span className="block text-center">-</span>
                        }
                    </p>
                </div>
            </div>
        </div>
    )
}
