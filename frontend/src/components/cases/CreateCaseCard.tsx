import { LucideLoader } from "lucide-react"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "../ui/card"
import { Separator } from "../ui/separator"
import { Label } from "../ui/label"
import { Input } from "../ui/input"
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "../ui/select"
import { Button } from "../ui/button"
import { Link } from "react-router"
import { cn } from "../../lib/utils"
import { useState } from "react"
import { useCreateCase } from "../../queries/cases"
import { CaseType, caseTypeToName, getTribunalCategoryOptions } from "../../lib/caseTypeNames"

const TribunalOptions = getTribunalCategoryOptions()

export default function CreateCaseCard() {
    const [fields, setFields] = useState<{
        caseNo: string,
        caseType: string,
        alias: string
    }>({ caseNo: "", caseType: "", alias: "" })
    const onFieldChange = (field: keyof typeof fields, value: string) => {
        setFields(prev => ({ ...prev, [field]: value }))
    }

    const [isLoading, setIsLoading] = useState(false)
    const [displayErr, setDisplayErr] = useState("")
    const [caseUUID, setCaseUUID] = useState("")

    const createCase = useCreateCase()
    const submitRegister: React.MouseEventHandler<HTMLButtonElement> = (e) => {
        e.preventDefault()
        setIsLoading(true)

        setTimeout(() => {
            createCase.mutate({
                caseId: fields.caseNo,
                caseType: fields.caseType,
                alias: fields.alias,
            }, {
                onSuccess: ({ id }) => {
                    setCaseUUID(id)
                },
                onError: (err) => {
                    let errMsg = String(err)
                    if (errMsg.includes("UNIQUE")) {
                        setDisplayErr(
                            "Ya existe un registro para el caso "
                                + fields.caseNo
                                + " | "
                                + caseTypeToName(fields.caseType as CaseType)
                        )
                    } else if (errMsg.includes("caseId invalid format")) {
                        setDisplayErr(
                            "El no. de expediente "
                            + fields.caseNo
                            + " no es válido. El formato debe ser '123/2024[-I]'"
                        )
                    } else {
                        setDisplayErr("Ocurrió un error al registrar el caso")
                    }
                },
                onSettled: () => {
                    setIsLoading(false)
                }
            })
        }, 100)
    }
    const onReset = () => {
        setDisplayErr("")
        setCaseUUID("")
        setFields({ caseNo: "", caseType: "", alias: "" })
        createCase.reset()
    }

    return (
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
                <div className={createCase.isError ? "text-sm text-red-500/80 pt-2" : "opacity-0"}>
                    {displayErr || "error"}
                </div>
            </CardContent>
            <Separator className="my-2" />
            <CardFooter className="justify-end items-end gap-2">
                <FooterActions
                    uuid={caseUUID}
                    show={createCase.isSuccess ? "success" : "default"}
                    onReset={onReset}
                    onSubmit={submitRegister} />
            </CardFooter>
        </Card>
    )
}

function FooterActions({
    uuid,
    show,
    onReset,
    onSubmit,
}: {
    uuid: string,
    show: "success" | "default",
    onReset: () => void,
    onSubmit: React.MouseEventHandler<HTMLButtonElement>
}) {
    if (show === "success") {
        return (
            <>
                <Button
                    type="button"
                    variant="secondary"
                    size="lg"
                    className="text-base"
                    onClick={onReset}>Limpiar</Button>
                <Button
                    type="button"
                    size="lg"
                    className="text-base font-semibold"
                    asChild>
                    <Link to={"/casos/" + uuid}>Ver Caso</Link>
                </Button>
            </>
        )
    }

    return (
        <Button
            type="button"
            variant="secondary"
            size="lg"
            className="text-base"
            onClick={onSubmit}>Registrar</Button>
    )
}
