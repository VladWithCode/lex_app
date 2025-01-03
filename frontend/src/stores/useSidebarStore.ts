import { create } from "zustand";

interface SidebarState {
    open: boolean
    toggleSidebar: () => void
}

export const useSidebarStore = create<SidebarState>(set => ({
    open: true,

    toggleSidebar: () => set(state => ({ open: !state.open }))
}))
