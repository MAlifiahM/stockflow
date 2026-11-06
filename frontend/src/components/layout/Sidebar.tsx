import { useState } from "react";
import { NavLink } from "react-router-dom";

const menus = [
    {
        label: "Overview",
        items: [{ name: "Dashboard", path: "/" }],
    },
    {
        label: "Inventory",
        items: [{ name: "Products", path: "/products" }],
    },
    {
        label: "Sales",
        items: [{ name: "Orders", path: "/orders" }],
    },
];

export default function Sidebar() {
    const [collapsed, setCollapsed] = useState(false);

    return (
        <aside
            className="border-end bg-white p-3 flex-shrink-0"
            style={{
                width: collapsed ? 72 : 240,
                minHeight: "100vh",
                transition: "width 0.25s ease",
            }}
        >
            <div className="d-flex align-items-center justify-content-between mb-4">
                {!collapsed && <h4 className="mb-0">StockFlow</h4>}

                <button
                    type="button"
                    className="btn btn-light border"
                    onClick={() => setCollapsed(!collapsed)}
                    aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
                    title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
                >
                    {collapsed ? "»" : "«"}
                </button>
            </div>

            {menus.map((group) => (
                <div className="mb-4" key={group.label}>
                    {!collapsed && (
                        <small className="text-secondary fw-semibold">
                            {group.label}
                        </small>
                    )}

                    <nav className="nav flex-column mt-2 gap-1">
                        {group.items.map((item) => (
                            <NavLink
                                key={item.path}
                                to={item.path}
                                end={item.path === "/"}
                                title={item.name}
                                className={({ isActive }) =>
                                    `nav-link rounded ${
                                        isActive
                                            ? "active bg-primary text-white"
                                            : "text-dark"
                                    }`
                                }
                            >
                                {collapsed
                                    ? item.name.charAt(0)
                                    : item.name}
                            </NavLink>
                        ))}
                    </nav>
                </div>
            ))}
        </aside>
    );
}