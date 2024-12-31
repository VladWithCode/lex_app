import { Separator } from "../../components/ui/separator";
import CreateCaseCard from "../../components/cases/CreateCaseCard";

export default function NewCasePage() {

    return (
        <>
            <h1 className="text-6xl font-semibold">Registrar Caso | lexApp</h1>
            <p className="text-lg text-stone-400 pt-2">Registra un nuevo caso para el que desees recibir actualizaciones.</p>
            <Separator className="my-2" />
            <div className="grid grid-cols-3 gap-4 max-h-full overflow-auto">
                <CreateCaseCard />
            </div>
        </>
    )
}
