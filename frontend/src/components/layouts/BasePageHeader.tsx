import { useLocation, useNavigate } from "react-router";
import { LucideChevronUp } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "../ui/button";

export default function BasePageHeader({
    title,
    description
}: {
    title: string | React.ReactElement;
    description: string | React.ReactElement;
}) {
    const navigate = useNavigate()
    const pathname = useLocation().pathname
    const isHome = pathname === "/"

    return (
        <div className="flex items-center gap-2 mb-2">
            <Button
                type="button"
                className={cn(
                    "rounded-lg [&_svg]:size-8 bg-transparent text-stone-50 gap-0 hover:bg-primary hover:text-primary-foreground shrink grow-0",
                    isHome && "hidden"
                )}
                onClick={() => !isHome && navigate(-1)}>
                <LucideChevronUp className="-rotate-90 stroke-current" />
                <span className="text-lg font-semibold -ml-0.5">Volver</span>
            </Button>
            <div className="grow basis-full">
                <h1 className="text-6xl font-semibold">{title}</h1>
                <p className="text-stone-400 pt-2">{description}</p>
            </div>
        </div>
    )
}
