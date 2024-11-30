import { create } from "zustand";

interface SidebarState {
    open: boolean
    toggleSidebar: () => void
}

export const useSidebarStore = create<SidebarState>(set => ({
    open: false,

    toggleSidebar: () => set(state => ({ open: !state.open }))
}))
