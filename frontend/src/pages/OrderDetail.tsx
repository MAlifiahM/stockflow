import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../services/api";
import type { Order } from "../types/order";

export default function OrderDetail() {
    const { id } = useParams();
    const navigate = useNavigate();

    const [order, setOrder] = useState<Order | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        const fetchOrder = async () => {
            setLoading(true);
            setError("");

            try {
                const response = await api.get<{ data: Order }>(
                    `/orders/${id}`
                );

                setOrder(response.data.data);
            } catch (error) {
                console.error("Failed to fetch order:", error);
                setError("Failed to fetch order details.");
            } finally {
                setLoading(false);
            }
        };

        if (id) {
            fetchOrder();
        }
    }, [id]);

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4">
                <div>
                    <h4 className="mb-1">Order Detail</h4>
                    <p className="text-muted mb-0">
                        View sales order information.
                    </p>
                </div>

                <button
                    className="btn btn-outline-secondary"
                    onClick={() => navigate("/orders")}
                >
                    Back to Orders
                </button>
            </div>

            {loading && <p>Loading order...</p>}

            {error && (
                <div className="alert alert-danger">{error}</div>
            )}

            {order && (
                <>
                    <div className="card border-0 shadow-sm mb-4">
                        <div className="card-header bg-white py-3">
                            <h5 className="mb-0">Order Information</h5>
                        </div>

                        <div className="card-body">
                            <div className="row g-3">
                                <div className="col-md-6">
                                    <label className="text-muted small">Order ID</label>
                                    <div className="fw-semibold">#{order.id}</div>
                                </div>

                                <div className="col-md-6">
                                    <label className="text-muted small">Customer</label>
                                    <div className="fw-semibold">{order.customer}</div>
                                </div>

                                <div className="col-md-6">
                                    <label className="text-muted small">Status</label>
                                    <div>
                                        <span className={`badge ${
                                            order.status === "CREATED"
                                                ? "text-bg-success"
                                                : "text-bg-secondary"
                                        }`}>
                                            {order.status}
                                        </span>
                                    </div>
                                </div>

                                <div className="col-md-6">
                                    <label className="text-muted small">Created At</label>
                                    <div>
                                        {new Date(order.created_at).toLocaleString("id-ID")}
                                    </div>
                                </div>

                                <div className="col-md-6">
                                    <label className="text-muted small">Total</label>
                                    <div>
                                        {new Intl.NumberFormat("id-ID", {
                                            style: "currency",
                                            currency: "IDR",
                                            maximumFractionDigits: 0,
                                        }).format(order.total)}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="card border-0 shadow-sm mb-4">
                        <div className="card-header bg-white py-3">
                            <h5 className="mb-0">Order Items</h5>
                        </div>

                        <div className="table-responsive">
                            <table className="table table-hover align-middle mb-0">
                                <thead className="table-light">
                                <tr>
                                    <th>SKU</th>
                                    <th>Quantity</th>
                                    <th>Unit Price</th>
                                    <th>Subtotal</th>
                                </tr>
                                </thead>

                                <tbody>
                                {order.items.map((item, index) => (
                                    <tr key={`${item.sku}-${index}`}>
                                        <td className="fw-semibold">{item.sku}</td>
                                        <td>{item.qty}</td>
                                        <td>
                                            {new Intl.NumberFormat("id-ID", {
                                                style: "currency",
                                                currency: "IDR",
                                                maximumFractionDigits: 0,
                                            }).format(item.unit_price)}
                                        </td>
                                        <td>
                                            {new Intl.NumberFormat("id-ID", {
                                                style: "currency",
                                                currency: "IDR",
                                                maximumFractionDigits: 0,
                                            }).format(item.qty * item.unit_price)}
                                        </td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </>

            )}
        </div>
    );
}