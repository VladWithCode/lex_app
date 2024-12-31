import { Outlet } from "react-router";
import { SidebarProvider } from "../ui/sidebar";
import AppSidebar from "../sidebar/AppSidebar";
import { PropsWithChildren } from "react";
import { useSidebarStore } from "../../stores/useSidebarStore";
import { Toaster } from "../ui/sonner";

export default function BaseLayout({ children }: PropsWithChildren) {
    const sidebarStore = useSidebarStore()
    return (
        <SidebarProvider className="h-[calc(100vh-2rem)] min-h-0" open={sidebarStore.open} defaultOpen={true}>
            <AppSidebar />
            <div className="relative flex flex-col w-full h-[calc(100vh-2rem)] z-0 p-2">
                {children !== undefined ? children : <Outlet />}
            </div>
            <Toaster closeButton richColors duration={8000} />
        </SidebarProvider>
    )
}
