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
        <Card className="bg-zinc-900 basis-1/3 flex-grow-0 flex-shrink-0 flex flex-col overflow-hidden">
            <CardHeader className="p-3">
                <CardTitle className="text-lg font-medium">
                    Expediente {caseData.caseId}
                    <br />
                    {caseTypeToName(caseData.caseType as CaseType)}
                </CardTitle>
            </CardHeader>
            <Separator />
            {
                children ? children : <DefaultCardContent caseData={caseData} />
            }
            <Separator />
            <CardFooter className="p-0 shrink-0">
                <Button
                    className="flex text-stone-50 bg-zinc-900 hover:bg-stone-800 rounded-none justify-center flex-grow basis-1/3 max-w-1/2 py-1">
                    <span className="mx-auto text-lg font-bold">Ver</span>
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
