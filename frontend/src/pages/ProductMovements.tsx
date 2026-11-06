import { useEffect, useState } from "react";
import {useLocation, useNavigate, useParams} from "react-router-dom";
import api from "../services/api";
import type {Product, StockMovement} from "../types/product";

export default function ProductMovements() {
    const { id } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const product = location.state?.product as Product | undefined;
    const [movements, setMovements] = useState<StockMovement[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        if (!id) return;

        const fetchMovements = async () => {
            setLoading(true);
            setError("");

            try {
                const response = await api.get<{ data: StockMovement[] }>(
                    `/products/${id}/movements`
                );

                setMovements(response.data.data);
            } catch (error) {
                console.error("Failed to fetch movements:", error);
                setError("Failed to fetch product movements.");
            } finally {
                setLoading(false);
            }
        };

        fetchMovements();
    }, [id]);

    return (
        <div>
            <div className="mb-4">
                <div className="d-flex justify-content-between align-items-center mb-4">
                    <div>
                        <h4 className="mb-1">Product Movements</h4>
                        <p className="text-muted mb-0">
                            Stock movement history for this product.
                        </p>
                    </div>
                    <button
                        type="button"
                        className="btn btn-outline-secondary"
                        onClick={() => navigate("/products")}
                        disabled={loading}
                    >
                        Back
                    </button>
                </div>
            </div>

            {/* Product Details */}
            <div className="card border-0 shadow-sm mb-4">
                <div className="card-header bg-white py-3">
                    <h6 className="mb-0">Product Details</h6>
                </div>

                <div className="card-body">
                    <div className="row g-3">
                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">Product ID</small>
                            <div className="fw-semibold">{product?.id ?? id ?? "-"}</div>
                        </div>

                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">SKU</small>
                            <div className="fw-semibold">{product?.sku ?? "-"}</div>
                        </div>

                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">Product Name</small>
                            <div className="fw-semibold">{product?.name ?? "-"}</div>
                        </div>

                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">Price</small>
                            <div className="fw-semibold">
                                {product
                                    ? new Intl.NumberFormat("id-ID", {
                                        style: "currency",
                                        currency: "IDR",
                                        maximumFractionDigits: 0,
                                    }).format(product.price)
                                    : "-"}
                            </div>
                        </div>

                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">Current Stock</small>
                            <div className="fw-semibold">{product?.stock ?? "-"}</div>
                        </div>

                        <div className="col-md-4 col-sm-6">
                            <small className="text-muted">Status</small>
                            <div>
                                {product ? (
                                    <span
                                        className={`badge ${
                                            product.is_active
                                                ? "text-bg-success"
                                                : "text-bg-secondary"
                                        }`}
                                    >
                                    {product.is_active ? "Active" : "Inactive"}
                                </span>
                                ) : (
                                    "-"
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Movement History */}
            <div className="mb-3">
                <h5 className="mb-0">Movement History</h5>
            </div>

            {loading && <p>Loading movements...</p>}

            {error && <div className="alert alert-danger">{error}</div>}

            <div className="card border-0 shadow-sm">
                <div className="table-responsive">
                    <table className="table table-hover align-middle mb-0">
                        <thead className="table-light">
                        <tr>
                            <th>ID</th>
                            <th>Type</th>
                            <th>Quantity</th>
                            <th>Note</th>
                            <th>Date</th>
                        </tr>
                        </thead>

                        <tbody>
                        {movements.length === 0 && !loading ? (
                            <tr>
                                <td colSpan={5} className="text-center py-4 text-muted">
                                    No movements found
                                </td>
                            </tr>
                        ) : (
                            movements.map((movement) => (
                                <tr key={movement.id}>
                                    <td>{movement.id}</td>
                                    <td>
                                        <span
                                            className={`badge ${
                                                movement.type === "in"
                                                    ? "text-bg-success"
                                                    : "text-bg-danger"
                                            }`}
                                        >
                                            {movement.type.toUpperCase()}
                                        </span>
                                    </td>
                                    <td>{movement.qty}</td>
                                    <td>{movement.note || "-"}</td>
                                    <td>
                                        {new Date(movement.created_at).toLocaleString("id-ID")}
                                    </td>
                                </tr>
                            ))
                        )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}