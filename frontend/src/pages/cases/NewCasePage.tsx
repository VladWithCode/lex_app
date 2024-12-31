import { Select } from "@radix-ui/react-select";
import { Button } from "../../components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "../../components/ui/card";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";
import { Separator } from "../../components/ui/separator";
import { getTribunalCategoryOptions } from "../../lib/caseTypeNames";
import { SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "../../components/ui/select";
import { useState } from "react";
import { LucideLoader } from "lucide-react";
import { cn } from "../../lib/utils";

const TribunalOptions = getTribunalCategoryOptions()

export default function NewCasePage() {
    const [fields, setFields] = useState<{
        caseNo: string,
        caseType: string,
        alias: string
    }>({ caseNo: "", caseType: "", alias: "" })
    const onFieldChange = (field: keyof typeof fields, value: string) => {
        setFields(prev => ({ ...prev, [field]: value }))
    }
    const isLoading = false

    return (
        <>
            <h1 className="text-6xl font-semibold">Registrar Caso | lexApp</h1>
            <p className="text-lg text-stone-400 pt-2">Registra un nuevo caso para el que desees recibir actualizaciones.</p>
            <Separator className="my-2" />
            <div className="grid grid-cols-3 gap-4 max-h-full overflow-auto">
                <Card className="relative flex flex-col z-0">
                    <div className={cn(
                        "absolute inset-0 flex flex-col items-center justify-center gap-4 bg-stone-800/80 rounded-lg",
                        isLoading ? "opacity-100 visible" : "opacity-0 invisible pointer-events-none"
                    )}>
                        <LucideLoader size="48" className="animate-spin" />
                        <p className="text-lg text-stone-400 uppercase font-semibold">Registrando Caso...</p>
                    </div>
                    <CardHeader className="">
                        <CardTitle className="text-2xl font-semibold">Nuevo caso</CardTitle>
                    </CardHeader>
                    <Separator className="my-2" />
                    <CardContent className="grow shrink-0 space-y-2">
                        <div className="flex gap-2">
                            <div className="flex-1 space-y-2">
                                <Label htmlFor="new-case-caseno">No. de Expediente</Label>
                                <Input
                                    name="new-case-caseno"
                                    id="new-case-caseno"
                                    type="text"
                                    placeholder="Ej. 84/2003"
                                    value={fields.caseNo}
                                    onChange={e => onFieldChange("caseNo", e.target.value)} />
                            </div>
                            <div className="flex-1 space-y-2">
                                <Label>Juzgado</Label>
                                <Select
                                    name="new-case-casetype"
                                    value={fields.caseType}
                                    onValueChange={val => onFieldChange("caseType", val)}>
                                    <SelectTrigger title="El juzgado del expediente">
                                        <SelectValue placeholder="Elige un juzgado..." />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {Object.entries(TribunalOptions).map(([k, v]) => (
                                            <SelectGroup key={k}>
                                                <SelectLabel>{v.title}</SelectLabel>
                                                {v.elements.map(el => (
                                                    <SelectItem key={k + "." + el.val} value={el.val} title={el.val}>{el.label}</SelectItem>
                                                ))}
                                            </SelectGroup>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                        <div className="space-y-2">
                            <Label>Alias</Label>
                            <Input
                                name="caseNo"
                                id="new-case-caseno"
                                type="text"
                                placeholder="Ej. Caso Sr. Brambila"
                                value={fields.alias}
                                onChange={e => onFieldChange("alias", e.target.value)} />
                        </div>
                        <div className="text-sm text-red-500/80 pt-2 opacity-0">Error!</div>
                    </CardContent>
                    <Separator className="my-2" />
                    <CardFooter className="justify-end items-end">
                        <Button
                            type="button"
                            variant="secondary"
                            size="lg"
                            className="text-base"
                            onClick={() => alert(JSON.stringify(fields, null, 4))}>Registrar</Button>
                    </CardFooter>
                </Card>
            </div>
        </>
    )
}
