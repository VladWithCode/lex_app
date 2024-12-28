import { HashRouter, Route, Routes } from "react-router";
import BaseLayout from "./components/layouts/BaseLayout";
import ErrorPage from "./pages/Error";
import Dashboard from "./pages/Dashboard";
import CasesPage from "./pages/cases/CasesPage";

export default function Router() {
    return (
        <HashRouter>
            <Routes>
                <Route path="/" errorElement={<ErrorPage />} element={<BaseLayout />}>
                    <Route index element={<Dashboard />} />
                    <Route path="/casos" element={<CasesPage />} />
                </Route>

                <Route path="*" element={<ErrorPage error={new Error("Not found")} />} />
            </Routes>
        </HashRouter>
    )
}
