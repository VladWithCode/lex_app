import { AppWindow, AppWindowMac, Minus, Plus } from "lucide-react"
import { Button } from "./ui/button"
import { Quit, WindowIsMaximised, WindowMaximise, WindowMinimise, WindowUnmaximise } from "../../wailsjs/runtime/runtime"

export default function AppControls() {
    return (
        <div className="flex items-center ml-auto gap-px justify-between pt-1">
            <MinimizeControl />
            <MaximizeControl />
            <QuitControl />
        </div>
    )
}

function MinimizeControl() {
    const minimizeApp =  () => {
        WindowMinimise()
    }
    return (
        <Button
            onClick={minimizeApp}
            className="relative hover:bg-stone-500 size-7 [&_svg]:size-6"
            variant="ghost">
            <Minus className="absolute -bottom-1" />
        </Button>
    )
}

function MaximizeControl() {
    const toggleMaximize = async () => {
        const isMaximized = await WindowIsMaximised()

        if (isMaximized) {
            WindowUnmaximise()
        } else {
            WindowMaximise()
        }
    }
    return (
        <Button
            onClick={toggleMaximize}
            className="hover:bg-stone-500 size-7 [&_svg]:size-6"
            variant="ghost">
            <AppWindowMac />
        </Button>
    )
}

function QuitControl() {
    const quitApp = () => {
        Quit()
    }
    return (
        <Button
            onClick={quitApp}
            className="hover:bg-destructive size-7 [&_svg]:size-8"
            variant="ghost" >
            <Plus className="rotate-45" />
        </Button>
    )
}
