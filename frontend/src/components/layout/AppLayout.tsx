import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Header from "./Header";

export default function AppLayout() {
    return (
        <div className="d-flex">
            <Sidebar />

            <div className="flex-grow-1 min-vh-100 overflow-hidden">
                <Header />

                <main className="p-4">
                    <Outlet />
                </main>
            </div>
        </div>
    );
}