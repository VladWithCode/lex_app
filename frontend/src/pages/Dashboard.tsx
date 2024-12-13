import CaseListing from "../components/cases/CaseListing";
import { Separator } from "../components/ui/separator";

export default function Dashboard() {
    return (
        <>
            <h1 className="text-6xl font-semibold">lexApp</h1>
            <p className="text-lg text-stone-500 pt-2">Bienvendo, Usuario</p>
            <Separator className="my-2" />
            <div className="py-2 px-4">
                <h2 className="text-3xl">Actividad reciente</h2>
                <p className="text-stone-400 pb-2">Ultimos casos en los que trabajaste</p>
                <div className="max-w-full overflow-auto py-1">
                    <CaseListing />
                </div>
            </div>
        </>
    )
}
