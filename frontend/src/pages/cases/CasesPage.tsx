import { Link, useLocation } from "react-router";
import CaseFilters from "../../components/cases/CaseFilters";
import { Separator } from "../../components/ui/separator";
import { useState, useRef, useEffect } from "react";
import { useCases } from "../../queries/cases";
import { LucideLoader } from "lucide-react";
import CaseCard from "../../components/cases/CaseCard";
import { CardContent } from "../../components/ui/card";
import { formatDateToShortReadable } from "../../lib/formatUtils";
import { db } from "../../../wailsjs/go/models";
import { Button } from "../../components/ui/button";
import BasePageHeader from "@/components/layouts/BasePageHeader";
import GeneralUpdatesDialog from "@/components/cases/GeneralUpdatesDialog";

export default function CasesPage() {
    const { params, setParam } = useCasesSearchParams()
    const blockAction = false

    return (
        <>
            <BasePageHeader title="Casos" description="Busca actualizaciones entre los casos registrados." />
            <Separator className="my-2" />
            <div className="flex gap-4 items-center">
                <Button size="lg" className="text-base font-bold" asChild>
                    <Link to="/casos/nuevo">Registrar Caso</Link>
                </Button>
                <GeneralUpdatesDialog blockAction={blockAction} filters={params} />
            </div>
            <Separator className="my-2" />
            <CaseFilters setFilter={setParam} filters={params} />
            <Separator className="my-2" />
            <Listing filters={params} />
        </>
    )
}

function Listing(
    {
        filters
    }: {
        filters: {
            search: string;
            caseNo: string;
            caseYear: string;
            caseType: string;
        }
    } & React.PropsWithChildren
) {
    const { data, isSuccess, isLoading, isError } = useCases({
        IncludeAccords: true,
        MaxAccords: 1,
        CaseNo: filters.caseNo,
        CaseYear: filters.caseYear,
        CaseType: filters.caseType,
        search: filters.search,
    })

    if (isError || !isSuccess) {
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

    return <div className="flex-1 grid grid-cols-3 auto-rows-[14rem] gap-4 max-h-full overflow-auto">
        {data.length == 0
            ? <div className="col-span-full self-center justify-self-center py-24">
                <p className="text-stone-400 text-2xl font-semibold">No se encontraron casos para la búsqueda ingresada</p>
            </div>
            : data.map(c => (<ListingCase key={c.id} c={c} />))}
    </div>
}

function ListingCase({ c }: { c: db.LexCase }) {
    return (
        <CaseCard caseData={c}>
            <CardContent className="flex flex-col px-2 py-1 gap-0.5 grow shrink">
                {
                    c.accords.length > 0
                        ? (
                            <>
                                <p className="text-stone-200 line-clamp-2 overflow-clip text-ellipsis">
                                    {c.accords[0].content}
                                </p>
                                <p className="text-stone-400 font-semibold mt-auto text-sm">
                                    {c.nature.toLowerCase()}
                                </p>
                                <p className="text-stone-400 font-semibold mt-auto">
                                    Fecha de Acuerdo: <span className="capitalize">{
                                        formatDateToShortReadable(new Date(c.accords[0].dateStr))
                                    }</span>
                                </p>
                            </>
                        )
                        : (
                            <>
                                <p className="text-stone-400 m-auto">SIN ACUERDO</p>
                                <p className="text-stone-400 font-semibold mt-auto text-sm">
                                    {c.nature.toLowerCase()}
                                </p>
                            </>
                        )
                }
            </CardContent>
        </CaseCard>
    )
}

function useCasesSearchParams() {
    const location = useLocation()
    const [params, setParams] = useState(new URLSearchParams(location.search))
    const debounceId = useRef<Timer | null>(null)

    const setParam = (key: string, value: string) => {
        if (debounceId.current) {
            clearTimeout(debounceId.current)
        }
        debounceId.current = setTimeout(() => {
            params.set(key, value)
            history.pushState(
                { params: params.toString() },
                "",
                "/#" + location.pathname + "?" + params.toString(),
            )
            setParams(new URLSearchParams(params))
        }, 250)
    }

    useEffect(() => {
        const handler = () => {
            let s = window.location.hash.substring(window.location.hash.indexOf("?"))
            setParams(new URLSearchParams(s))
        }
        window.addEventListener("popstate", handler)
        return () => window.removeEventListener("popstate", handler)
    }, [])

    return {
        params: {
            search: params.get("search") || "",
            caseNo: params.get("caseNo") || "",
            caseYear: params.get("caseYear") || "",
            caseType: params.get("caseType") || "",
        },
        setParam,
    }
}
