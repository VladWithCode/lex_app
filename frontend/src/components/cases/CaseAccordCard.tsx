import { db } from "wailsjs/go/models"
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card"
import { useState } from "react"
import { cn } from "@/lib/utils"
import { Separator } from "../ui/separator"
import { formatDateToShortReadable } from "@/lib/formatUtils"

export default function CaseAccordCard({ accord, className }: {
    accord: db.Accord
    className?: string
}) {
    const [showFullContent, setShowFullContent] = useState(false)
    return (
        <Card className={cn(
            "bg-zinc-900",
            className
        )}>
            <CardHeader className="p-4">
                <CardTitle className="text-lg">
                    Fecha: {
                        accord.dateStr && accord.dateStr.length > 0
                            ? formatDateToShortReadable(new Date(accord.dateStr))
                            : "Sin Fecha"
                    }
                </CardTitle>
            </CardHeader>
            <Separator className="my-2" />
            <CardContent className="p-4">
                <p className={cn(
                    "text-stone-200 text-lg line-clamp-5 teel-ellipsis overflow-clip",
                    showFullContent && "line-clamp-none"
                )}>{accord.content}</p>
            </CardContent>
            {/* <Separator className="my-2" />
            <CardFooter className="p-4">
                <Button
                    size="lg"
                    variant="destructive"
                    className="ml-auto text-base">Eliminar</Button>
            </CardFooter> */}
        </Card>
    )
}
