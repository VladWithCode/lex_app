import { useState } from "react"
import { getTribunalCategoryOptions, TribunalCategoryOptions } from "../../lib/caseTypeNames"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Select, SelectItem, SelectLabel, SelectContent, SelectGroup, SelectTrigger, SelectValue } from "../ui/select"

export type CaseFiltersParams = React.PropsWithChildren & {
    filters: { search: string, caseNo: string, caseYear: string, caseType: string };
    setFilter: (key: string, value: string) => void;
}

const TribunalOptions = getTribunalCategoryOptions()

export default function CaseFilters({ setFilter, filters }: CaseFiltersParams) {
    const onFieldChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setFilter(e.target.name, e.target.value)
    }

    return (
        <div className="space-y-1.5 max-w-full p-px overflow-hidden">
            <div className="flex gap-2">
                <h4 className="text-xl font-medium">Filtros</h4>
            </div>
            <div className="flex gap-3">
                <div className="basis-80 shrink space-y-1">
                    <Label htmlFor="case-filter-search">Buscar</Label>
                    <Input
                        id="case-filter-search"
                        name="search"
                        type="search"
                        defaultValue={filters.search}
                        onChange={onFieldChange}
                        placeholder="Busca palabras clave, no. de expediente, alias, contenido del acuerdo"
                        title="Busca palabras clave, no. de expediente, alias, contenido del acuerdo" />
                </div>
                <div className="shrink grow max-w-56 space-y-1">
                    <Label htmlFor="case-filter-year">No. Expediente</Label>
                    <Input
                        id="case-filter-case-no"
                        name="caseNo"
                        type="text"
                        defaultValue={filters.caseNo}
                        onChange={onFieldChange}
                        placeholder="El no. de expediente (ej. 84, 678, 00264)"
                        title="El no. de expediente (ej. 84, 678, 00264)" />
                </div>
                <div className="shrink grow max-w-56 space-y-1">
                    <Label htmlFor="case-filter-case-year">Año</Label>
                    <Input
                        id="case-filter-case-year"
                        name="caseYear"
                        type="number"
                        defaultValue={filters.caseYear}
                        onChange={onFieldChange}
                        placeholder="El año en el no. de expediente (ej. 2003, 1999)"
                        title="El año en el no. de expediente (ej. 2003, 1999)" />
                </div>
                <div className="shrink grow max-w-56 space-y-1">
                    <Label htmlFor="case-fitler-case-type">Juzgado</Label>
                    <TribunalSelect
                        tribunals={TribunalOptions}
                        onChange={val => {
                            setFilter("caseType", val)
                        }}
                        defaultValue={filters.caseType} />
                </div>
                {/* TODO: Implement Advanced Filters
                    <div className="flex gap-3 items-center shrink grow-0 max-w-32 ml-auto mt-auto">
                    <Button
                        className="bg-stone-900 p-2 [&_svg]:size-5 text-stone-50 hover:text-stone-900"
                        title="Filtros avanzados"
                        type="button">
                        <div className="text-lg font-medium">Todos</div>
                        <LucideFilter className="stroke-current" />
                    </Button>
                </div> */}
            </div>
        </div>
    )
}

function TribunalSelect(
    {
        tribunals,
        onChange,
        defaultValue,
    }: {
        tribunals: TribunalCategoryOptions,
        onChange: (val: string) => void,
        defaultValue: string
    }
) {
    const [selected, setSelected] = useState<string>("unset")

    return (
        <Select value={selected} onValueChange={val => {
            if (val === "unset") {
                setSelected("")
                onChange("")
            } else {
                setSelected(val)
                onChange(val)
            }
        }} defaultValue={defaultValue}>
            <SelectTrigger title="El juzgado del expediente">
                <SelectValue placeholder="Elige un juzgado..." />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="unset">Todos</SelectItem>
                {Object.entries(tribunals).map(([k, v]) => (
                    <TribunalGroup key={k} title={v.title} elements={v.elements} />
                ))}
            </SelectContent>
        </Select>
    )
}

function TribunalGroup({ elements, title }: { elements: { val: string, label: string }[], title: string }) {
    return (
        <SelectGroup>
            <SelectLabel className="text-lg">{title}</SelectLabel>
            {elements.map(el => (<SelectItem className="text-base" key={title + "." + el.val + "." + el.label} value={el.val}>{el.label}</SelectItem>))}
        </SelectGroup>
    )
}
