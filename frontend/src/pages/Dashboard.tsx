import { Link } from "react-router";
import RecentCases from "../components/cases/RecentCases";
import { Separator } from "../components/ui/separator";

export default function Dashboard() {
    return (
        <>
            <h1 className="text-6xl font-semibold">lexApp</h1>
            <p className="text-lg text-stone-500 pt-2">Bienvendo, Usuario</p>
            <Separator className="my-2" />
            <div className="py-2 px-4">
                <h2 className="text-3xl">
                    <span>Ultimas Actualizaciones</span>
                    <Link to="/cases/recent" className="relative bottom-1 text-sm mt-auto ml-3 underline">Ver Todas&nbsp;&gt;&gt;</Link>
                </h2>
                <div className="max-w-full overflow-auto py-1">
                    <RecentCases />
                </div>
            </div>
        </>
    )
}
