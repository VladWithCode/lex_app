import AppControls from "./AppControls";

export default function TitleBar() {
    return (
        <div
            style={{ "--wails-draggable": "drag" } as React.CSSProperties}
            className="relative top-0 left-0 right-0 w-full h-8 flex px-2 hover:cursor-grab focus:cursor-grabbing active:cursor-grabbing"
            >
            <div className="flex items-center">
                <p className="ml-2 text-stone-400 font-semibold">lexApp</p>
            </div>
            <AppControls />
        </div>
    )
}
