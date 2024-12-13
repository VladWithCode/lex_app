import { MouseEvent } from "react";
import AppControls from "./AppControls";
import { WindowIsMaximised, WindowMaximise, WindowUnmaximise } from "../../wailsjs/runtime/runtime";

async function onDbClick(e: MouseEvent<HTMLDivElement>) {
    e.preventDefault()
    if (await WindowIsMaximised()) {
        WindowUnmaximise()
    } else {
        WindowMaximise()
    }
}

export default function TitleBar() {
    return (
        <div
            style={{ "--wails-draggable": "drag" } as React.CSSProperties}
            className="relative top-0 left-0 right-0 w-full h-8 flex px-2 hover:cursor-grab focus:cursor-grabbing active:cursor-grabbing select-none"
            onDoubleClick={onDbClick}>
            <div className="flex items-center">
                <p className="ml-2 text-stone-400 font-semibold">lexApp</p>
            </div>
            <AppControls />
        </div>
    )
}
