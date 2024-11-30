import { Link } from "react-router";
import BaseLayout from "../components/layouts/BaseLayout";

export default function ErrorPage({ error }: { error?: Error }) {
    console.error(error)

    return (
        <BaseLayout>
            <h1 className="text-8xl font-bold">¡Ocurrió un error!</h1>
            <p className="text-lg text-stone-400">
                Parece ser que algo no funcionó como se esperaba.
                Por favor vuelve al <Link to="/">inicio</Link>.
            </p>
        </BaseLayout>
    )
}
