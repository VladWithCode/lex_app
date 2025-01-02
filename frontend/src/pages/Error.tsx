import { Link } from "react-router";
import BaseLayout from "../components/layouts/BaseLayout";
import BasePageHeader from "@/components/layouts/BasePageHeader";

export default function ErrorPage({ error }: { error?: Error }) {
    console.error(error)

    return (
        <BaseLayout>
            <BasePageHeader title="¡Ocurrió un error!" description="" />
            <p className="text-lg text-stone-400">
                Parece ser que algo no funcionó como se esperaba.
                Por favor vuelve al <Link to="/">inicio</Link>.
            </p>
        </BaseLayout>
    )
}
