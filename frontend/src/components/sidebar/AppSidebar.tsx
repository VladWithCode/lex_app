import { Home, LucideFolder, SearchX } from "lucide-react";
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarTrigger,
    useSidebar,
} from "../ui/sidebar";
import { useNavigate } from "react-router";
import { cn } from "../../lib/utils";
import { Separator } from "../ui/separator";
import { useSidebarStore } from "../../stores/useSidebarStore";

const items = [
    {
        title: "Inicio",
        url: "/",
        icon: Home,
    },
    {
        title: "Mis Casos",
        url: "/casos",
        icon: LucideFolder,
    },
    {
        title: "Buscar",
        url: "/buscador",
        icon: SearchX,
    },
]

export default function AppSidebar() {
    const { toggleSidebar } = useSidebarStore()
    return (
        <Sidebar className="h-[calc(100vh-2rem)] mt-auto" variant="floating" collapsible="icon">
            <SidebarHeader>
                <SidebarTrigger onClick={() => toggleSidebar()} className="ml-auto"  />
            </SidebarHeader>
            <Separator className="w-11/12 mx-auto" ></Separator>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel children="Menu" />
                    <NavMain />
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    )
}

function NavMain() {
    const { open } = useSidebar()
    const navigate = useNavigate()
    return (
        <SidebarMenu>
            {items.map(i => (
                <SidebarMenuItem key={i.title}>
                    <SidebarMenuButton
                        onClick={() => navigate(i.url)}
                        className={cn("[&_svg]:size-4 transition-[height,_width]", open && "[&_svg]:size-6")}>
                        <i.icon />
                        {i.title}
                    </SidebarMenuButton>
                </SidebarMenuItem>
            ))}
        </SidebarMenu>
    )
}
