import {useEffect, useState} from "react";
import type {Order, OrderPaginationResponse} from "../types/order";

import {useNavigate} from "react-router-dom";
import api from "../services/api.ts";

export default function Orders() {
    const navigate = useNavigate();
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [page, setPage] = useState(1);
    const [limit] = useState(10);
    const [totalPages, setTotalPages] = useState(0);
    const [search, setSearch] = useState("");
    const [totalOrders, setTotalOrders] = useState(0);
    const [refresh, setRefresh] = useState(0);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [success, setSuccess] = useState("");

    useEffect(() => {
        const fetchOrders = async () => {
            setLoading(true);
            setError("");

            try {
                const response = await api.get<OrderPaginationResponse>(
                    "orders",
                    {
                        params: {
                            page: page,
                            limit: limit,
                            search: search,
                        }
                    });
                setOrders(response.data.data);
                setTotalPages(response.data.total_pages);
                setTotalOrders(response.data.total_items);
            } catch (error) {
                console.error("Failed to fetch orders:", error);
                setError("Failed to fetch orders");
            } finally {
                setLoading(false);
            }
        }

        fetchOrders();
    }, [page, limit, search, refresh]);

    const handleCancel = async (id: number) => {
        setError("");
        setSuccess("");

        try {
            await api.post(`/orders/${id}/cancel`);

            setSelectedOrder(null);
            setSuccess("Order cancelled successfully.");
            setRefresh((prev) => prev + 1);
        } catch (error) {
            console.error("Failed to cancel order:", error);
            setError("Failed to cancel order. Please try again.");
        }
    };

    return (
        <div>
            <div className="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-4">
                <h4 className="mb-0">Sales Orders</h4>

                <button
                    className="btn btn-primary"
                    onClick={() => navigate("/orders/new")}
                >
                    Create Order
                </button>
            </div>

            <div className="row mb-4">
                <div className="col-md-4">
                    <div className="card border-0 shadow-sm">
                        <div className="card-body">
                            <p className="text-secondary mb-2">Total Orders</p>
                            <h3 className="fw-bold mb-0">
                                {loading ? "..." : totalOrders}
                            </h3>
                        </div>
                    </div>
                </div>
            </div>

            <div className="card border-0 shadow-sm">
                {loading && <p>Loading...</p>}
                <input
                    type="text"
                    className="form-control mb-3"
                    placeholder="Search products..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                />


                <div className="table-responsive">
                    <table className="table table-hover align-middle mb-0">
                        <thead className="table-light">
                        <tr>
                            <th>ID</th>
                            <th>Customer</th>
                            <th>Status</th>
                            <th>Total</th>
                            <th>Created At</th>
                            <th>Total Items</th>
                            <th className="text-end">Actions</th>
                        </tr>
                        </thead>

                        <tbody>
                        {orders.length === 0 ? (
                            <tr>
                                <td colSpan={7} className="text-center py-4 text-muted">
                                    No orders found
                                </td>
                            </tr>
                        ) : (
                            orders.map((order) => (
                                <tr key={order.id}>
                                    <td className="fw-semibold">#{order.id}</td>
                                    <td>{order.customer}</td>
                                    <td>
                                <span
                                    className={`badge ${
                                        order.status === "CREATED"
                                            ? "text-bg-success"
                                            : "text-bg-secondary"
                                    }`}
                                >
                                    {order.status}
                                </span>
                                    </td>
                                    <td>
                                        {new Intl.NumberFormat("id-ID", {
                                            style: "currency",
                                            currency: "IDR",
                                            maximumFractionDigits: 0,
                                        }).format(order.total)}
                                    </td>
                                    <td>
                                        {new Date(order.created_at).toLocaleString("id-ID")}
                                    </td>
                                    <td>
                                        {(order.items ?? []).reduce(
                                            (total, item) => total + item.qty,
                                            0
                                        )}
                                    </td>
                                    <td>
                                        <div className="d-flex flex-nowrap justify-content-end gap-2">
                                            <button
                                                className="btn btn-sm btn-outline-primary"
                                                onClick={() => navigate(`/orders/${order.id}`)}
                                            >
                                                Detail
                                            </button>

                                            {order.status === "CREATED" && (
                                                <button
                                                    className="btn btn-sm btn-outline-danger"
                                                    onClick={() => setSelectedOrder(order)}
                                                >
                                                    Cancel
                                                </button>
                                            )}
                                        </div>
                                    </td>
                                </tr>
                            ))
                        )}
                        </tbody>
                    </table>
                </div>
            </div>

            <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mt-3">
                <span className="text-muted">
                    Page {page} of {totalPages}
                </span>

                <div className="d-flex gap-2">
                    <button
                        type="button"
                        className="btn btn-outline-secondary btn-sm"
                        disabled={page <= 1}
                        onClick={() => setPage((prev) => prev - 1)}
                    >
                        Previous
                    </button>

                    <button
                        type="button"
                        className="btn btn-outline-primary btn-sm"
                        disabled={page >= totalPages || totalPages === 0}
                        onClick={() => setPage((prev) => prev + 1)}
                    >
                        Next
                    </button>
                </div>
            </div>

            {(selectedOrder || success || error) && (
                <div
                    className="modal fade show d-block"
                    tabIndex={-1}
                    style={{ backgroundColor: "rgba(0,0,0,0.5)" }}
                >
                    <div className="modal-dialog modal-dialog-centered">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h5 className="modal-title">
                                    {success
                                        ? "Cancel Successful"
                                        : error
                                            ? "Cancel Failed"
                                            : "Confirm Cancel"}
                                </h5>

                                <button
                                    type="button"
                                    className="btn-close"
                                    onClick={() => {
                                        setSelectedOrder(null);
                                        setSuccess("");
                                        setError("");
                                    }}
                                />
                            </div>

                            <div className="modal-body">
                                {success ? (
                                    <div className="text-success">{success}</div>
                                ) : error ? (
                                    <div className="text-danger">{error}</div>
                                ) : (
                                    <>
                                        Are you sure you want to cancel order{" "}
                                        <strong>#{selectedOrder?.id}</strong> for{" "}
                                        <strong>{selectedOrder?.customer}</strong>?

                                        <p className="text-muted mb-0 mt-2">
                                            This action will restore the stock of all
                                            products in this order.
                                        </p>
                                    </>
                                )}
                            </div>

                            <div className="modal-footer">
                                {success || error ? (
                                    <button
                                        type="button"
                                        className="btn btn-primary"
                                        onClick={() => {
                                            setSelectedOrder(null);
                                            setSuccess("");
                                            setError("");
                                        }}
                                    >
                                        Close
                                    </button>
                                ) : (
                                    <>
                                        <button
                                            type="button"
                                            className="btn btn-secondary"
                                            onClick={() => setSelectedOrder(null)}
                                        >
                                            Keep Order
                                        </button>

                                        <button
                                            type="button"
                                            className="btn btn-danger"
                                            onClick={() => handleCancel(selectedOrder!.id)}
                                        >
                                            Confirm Cancel
                                        </button>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}