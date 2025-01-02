import { Outlet, useNavigate } from "react-router";
import { SidebarProvider } from "../ui/sidebar";
import AppSidebar from "../sidebar/AppSidebar";
import { PropsWithChildren } from "react";
import { useSidebarStore } from "../../stores/useSidebarStore";
import { Toaster } from "../ui/sonner";
import { Button } from "../ui/button";
import { LucideChevronUp } from "lucide-react";

export default function BaseLayout({ children }: PropsWithChildren) {
    const sidebarStore = useSidebarStore()
    const navigate = useNavigate()
    return (
        <SidebarProvider className="h-[calc(100vh-2rem)] min-h-0" open={sidebarStore.open} defaultOpen={true}>
            <AppSidebar />
            <div className="relative flex flex-col w-full h-[calc(100vh-2rem)] z-0 p-2">
                <div className="flex items-center grow-0 mb-2">
                    <Button
                        type="button"
                        className="rounded-lg [&_svg]:size-8 bg-transparent text-stone-50 gap-0 hover:bg-primary hover:text-primary-foreground"
                        onClick={() => navigate(-1)}>
                        <LucideChevronUp className="-rotate-90 stroke-current" />
                        <span className="text-lg font-semibold -ml-0.5">Volver</span>
                    </Button>
                </div>
                {children !== undefined ? children : <Outlet />}
            </div>
            <Toaster closeButton richColors duration={8000} />
        </SidebarProvider>
    )
}
