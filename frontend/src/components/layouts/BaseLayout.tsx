import { Outlet } from "react-router";
import { SidebarProvider } from "../ui/sidebar";
import AppSidebar from "../sidebar/AppSidebar";
import { PropsWithChildren } from "react";
import { useSidebarStore } from "../../stores/useSidebarStore";

export default function BaseLayout({ children }: PropsWithChildren) {
    const sidebarStore = useSidebarStore()
    return (
        <SidebarProvider open={sidebarStore.open} defaultOpen={false}>
            <AppSidebar />
            <div className="relative flex flex-col w-full h-full z-0 overflow-scroll p-2">
                {children !== undefined ? children : <Outlet />}
            </div>
        </SidebarProvider>
    )
}
