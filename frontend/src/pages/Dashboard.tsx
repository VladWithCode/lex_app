import { Link } from "react-router";
import RecentCases from "../components/cases/RecentCases";
import { Separator } from "../components/ui/separator";
import BasePageHeader from "@/components/layouts/BasePageHeader";

export default function Dashboard() {
    return (
        <>
            <BasePageHeader title="lexApp" description="Los ultimos acuerdos publicados para tus casos." />
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
