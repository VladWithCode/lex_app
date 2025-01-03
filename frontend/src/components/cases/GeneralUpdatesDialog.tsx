import { useState } from "react";
import { Button } from "../ui/button";
import { LucideLoader } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Input } from "@/components/ui/input";
import { Separator } from "../ui/separator";
import { Label } from "../ui/label";
import { useFindAndUpdateCaseAccords } from "@/queries/cases";
import { toast } from "sonner";
import { CaseFilters } from "./CaseFilters";

type SearchParams = {
    fromDate: string;
    untilDate: string;
    daysBack: number;
    exhaustSearch: boolean;
}
type SPKey = keyof SearchParams
type SPValue = SearchParams[SPKey]

export default function GeneralUpdatesDialog({
    blockAction,
    filters,
}: {
    blockAction: boolean;
    filters: CaseFilters;
}) {
    const [isOpen, setIsOpen] = useState(false)
    const [searchParams, setSearchParams] = useState<SearchParams>({
        fromDate: new Date().toISOString().split('T')[0],
        untilDate: "",
        daysBack: 0,
        exhaustSearch: false,
    })
    const setField = (field: SPKey, value: SPValue) => {
        setSearchParams(prev => ({ ...prev, [field]: value }))
    }
    const findAndUpdate = useFindAndUpdateCaseAccords()

    return (
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
                <Button
                    size="lg"
                    className="text-base font-bold active:scale-95 transition-transform duration-150 disabled:opacity-50"
                    disabled={blockAction}>Buscar actualizaciones</Button>
            </DialogTrigger>
            <DialogContent className="w-full max-w-lg">
                <DialogHeader>
                    <DialogTitle className="text-lg">Buscar actualizaciones</DialogTitle>
                    <DialogDescription asChild>
                        <div>
                            <p className="text-base pb-1">
                                Buscando actualizaciones con los filtros especificados.
                            </p>
                            <div className="flex gap-2 items-center">
                                <p className="text-base">
                                    <span className="font-bold">Buscar:</span>
                                    <span className="text-stone-100"> {filters.search || "N/A"} </span>
                                </p>
                                <p className="text-base">
                                    <span className="font-bold">Numero</span>
                                    <span className="text-stone-100"> {filters.caseNo || "N/A"} </span>
                                </p>
                                <p className="text-base">
                                    <span className="font-bold">Año:</span>
                                    <span className="text-stone-100"> {filters.caseYear || "N/A"} </span>
                                </p>
                                <p className="text-base">
                                    <span className="font-bold">Juzgado:</span>
                                    <span className="text-stone-100"> {filters.caseType || "N/A"} </span>
                                </p>
                            </div>
                        </div>
                    </DialogDescription>
                </DialogHeader>
                <div>
                    <Separator className="mb-2" />
                    <h3 className="col-span-full text-lg">Parametros</h3>
                    <div className="grid grid-cols-2 gap-4 py-1">
                        <div className="col-span-1">
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Label htmlFor="searchupdates-from-date">Desde (fecha):</Label>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p className="text-base">Determina la fecha de inicio para la busqueda.</p>
                                        <p className="text-base">Por defecto es la fecha actual.</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                            <Input
                                id="searchupdates-from-date"
                                name="fromDate"
                                type="date"
                                value={searchParams.fromDate}
                                onChange={(e) => setField("fromDate", e.target.value)} />
                        </div>
                        {/* TODO: Implement until date
                            <div className="space-y-2 grow">
                                <TooltipProvider>
                                    <Tooltip>
                                        <TooltipTrigger asChild>
                                            <Label htmlFor="searchupdates-until-date">Hasta (fecha): </Label>
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p className="text-base">Determina la fecha de fin para la busqueda.</p>
                                            <p className="text-base">Esta debe ocurrir antes de la fecha de inicio.</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </TooltipProvider>
                                <Input
                                    id="searchupdates-until-date"
                                    name="untilDate"
                                    type="date"
                                    value={searchParams.untilDate}
                                    onChange={(e) => setField("untilDate", e.target.value)} />
                            </div> */}
                        <div className="col-span-1">
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Label htmlFor="searchupdates-from-date">No. de Días de busqueda</Label>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p className="text-base">Alternativo a la fecha de fin. Especifica un número de días (en el pasado) a buscar.</p>
                                        <p className="text-base">Por defecto es 0 (o solo buscar 1 día) y es ignorado si se especifica una fecha de fin.</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                            <Input
                                id="searchupdates-days-back"
                                name="daysBack"
                                type="number"
                                min="0"
                                value={searchParams.daysBack}
                                onChange={(e) => {
                                    setField("daysBack", Number(e.target.value))
                                }} />
                        </div>

                        <TooltipProvider>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <div className="col-span-full flex items-center gap-1 w-fit">
                                        <Checkbox
                                            id="searchupdates-exhaust-search"
                                            name="exhaustSearch"
                                            checked={searchParams.exhaustSearch}
                                            onCheckedChange={(chkd) => {
                                                let isChecked = chkd !== 'indeterminate' && Boolean(chkd)
                                                setField("exhaustSearch", isChecked)
                                            }} />
                                        <Label htmlFor="searchupdates-exhaust-search">Busqueda exhaustiva</Label>
                                    </div>
                                </TooltipTrigger>
                                <TooltipContent>
                                    <p className="text-base max-w-[60ch]">Si se activa, la busqueda continuará hasta recuperar todas las actualizaciones disponibles en el rango de las fechas especificadas.</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>

                    </div>
                    <Separator className="mt-4" />
                </div>
                <DialogFooter>
                    <Button
                        onClick={() => {
                            console.log("Finding and updating")
                            findAndUpdate.mutate({
                                searchStartDate: new Date(searchParams.fromDate),
                                maxSearchBack: searchParams.daysBack,
                                exhaustSearch: searchParams.exhaustSearch,
                                findOpts: {
                                    CaseType: filters.caseType,
                                    CaseYear: filters.caseYear,
                                    CaseNo: filters.caseNo,
                                    Search: filters.search,
                                    IncludeAccords: false,
                                }
                            }, {
                                onSuccess: () => {
                                    toast.success("Actualizado")
                                },
                                onError: (err) => {
                                    if (String(err).includes("no updates")) {
                                        toast.error("No hay actualizaciones para el caso")
                                    } else {
                                        toast.error("Error al actualizar")
                                    }
                                }
                            })
                        }}
                        className="text-base"
                        disabled={blockAction || findAndUpdate.isPending}>
                        {
                            findAndUpdate.status === "pending"
                                ? <LucideLoader className="animate-spin" />
                                : "Buscar actualizaciones"

                        }
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
