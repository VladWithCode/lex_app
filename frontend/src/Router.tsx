import { HashRouter, Route, Routes } from "react-router";
import BaseLayout from "./components/layouts/BaseLayout";
import ErrorPage from "./pages/Error";
import Dashboard from "./pages/Dashboard";
import CasesPage from "./pages/cases/CasesPage";
import NewCasePage from "./pages/cases/NewCasePage";
import CaseDetailPage from "./pages/cases/CaseDetailPage";

export default function Router() {
    return (
        <HashRouter>
            <Routes>
                <Route path="/" errorElement={<ErrorPage />} element={<BaseLayout />}>
                    <Route index element={<Dashboard />} />
                    <Route path="/casos" element={<CasesPage />} />
                    <Route path="/casos/nuevo" element={<NewCasePage />} />
                    <Route path="/casos/:caseUUID" element={<CaseDetailPage />} />
                </Route>

                <Route path="*" element={<ErrorPage error={new Error("Not found")} />} />
            </Routes>
        </HashRouter>
    )
}
