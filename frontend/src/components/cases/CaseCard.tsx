import { LucideEye } from "lucide-react";
import { Button } from "../ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { Separator } from "../ui/separator";
import { db } from "../../../wailsjs/go/models"
import { CaseType, caseTypeToName } from "../../lib/caseTypeNames";

export type CaseCardProps = React.PropsWithChildren & {
    caseData: db.LexCase
}

export default function CaseCard({ caseData, children }: CaseCardProps) {
    return (
        <Card className="bg-zinc-900 basis-1/3 flex-grow-0 flex-shrink-0 overflow-hidden">
            <CardHeader className="p-3">
                <CardTitle className="text-lg font-medium">Expediente {caseData.caseId} | {caseTypeToName(caseData.caseType as CaseType)}</CardTitle>
            </CardHeader>
            <Separator />
            {
                children ? children : <DefaultCardContent caseData={caseData} />
            }
            <Separator />
            <CardFooter className="p-0">
                <Button
                    className="relative z-0 flex text-stone-50 bg-zinc-900 hover:bg-stone-800 rounded-none justify-center flex-grow basis-1/3 max-w-1/2 py-1 [&_svg]:size-6 group">
                    <span className="absolute text-lg font-bold top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 transition-[opacity] duration-300 opacity-0 group-hover:opacity-100">Ver</span>
                    <LucideEye className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 stroke-current transition-[opacity] duration-300 opacity-100 group-hover:opacity-0" />
                </Button>
            </CardFooter>
        </Card>
    )
}

function DefaultCardContent({ caseData }: CaseCardProps) {
    return (
        <CardContent className="grow space-y-2 px-3 py-1.5 pb-3">
            <p className="text-stone-200 font-medium">Ultimo acceso: {String(caseData.lastAccessedAt)}</p>
        </CardContent>
    )
}
