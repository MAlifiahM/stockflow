import { useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import api from "../services/api";
import type { Product } from "../types/product";

export default function Stock() {
    const navigate = useNavigate();
    const { id } = useParams();
    const location = useLocation();
    const product = location.state?.product as Product | undefined;

    const [type, setType] = useState<"in" | "out">("in");
    const [qty, setQty] = useState(1);
    const [note, setNote] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    if (!product) {
        navigate("/products", { replace: true });
        return null;
    }

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");

        if (!Number.isInteger(qty) || qty <= 0) {
            setError("Quantity must be a positive integer.");
            return;
        }

        setLoading(true);

        try {
            await api.post(`/products/${id}/stock`, {
                type,
                qty,
                note: note.trim() || undefined,
            });

            navigate("/products");
        } catch (error) {
            console.error("Failed to adjust stock:", error);

            if (
                typeof error === "object" &&
                error !== null &&
                "response" in error
            ) {
                const response = error.response as {
                    data?: {
                        error?: {
                            code?: string;
                            available?: number;
                        };
                    };
                };

                if (response.data?.error?.code === "INSUFFICIENT_STOCK") {
                    setError(
                        `Insufficient stock. Available stock: ${response.data.error.available ?? product.stock}.`,
                    );
                } else {
                    setError("Failed to adjust stock. Please try again.");
                }
            } else {
                setError("Failed to adjust stock. Please try again.");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4">
                <h4>Stock Adjustment</h4>

                <button
                    type="button"
                    className="btn btn-outline-secondary"
                    onClick={() => navigate("/products")}
                >
                    Back
                </button>
            </div>

            <div className="card border-0 shadow-sm">
                <div className="card-body p-4">
                    <form onSubmit={handleSubmit}>
                        <div className="mb-3">
                            <label className="form-label">Product</label>

                            <input
                                type="text"
                                className="form-control"
                                value={`${product.sku} - ${product.name}`}
                                disabled
                            />

                            <div className="form-text">
                                Current stock: {product.stock}
                            </div>
                        </div>

                        <div className="mb-3">
                            <label className="form-label">Movement Type</label>

                            <select
                                className="form-select"
                                value={type}
                                onChange={(e) =>
                                    setType(e.target.value as "in" | "out")
                                }
                                required
                            >
                                <option value="in">Stock In</option>
                                <option value="out">Stock Out</option>
                            </select>
                        </div>

                        <div className="mb-3">
                            <label className="form-label">Quantity</label>

                            <input
                                type="number"
                                className="form-control"
                                min="1"
                                step="1"
                                value={qty}
                                onChange={(e) => setQty(Number(e.target.value))}
                                required
                            />
                        </div>

                        <div className="mb-4">
                            <label className="form-label">
                                Note <span className="text-muted">(Optional)</span>
                            </label>

                            <textarea
                                className="form-control"
                                rows={3}
                                placeholder="Enter adjustment note"
                                value={note}
                                onChange={(e) => setNote(e.target.value)}
                            />
                        </div>

                        {error && (
                            <div className="alert alert-danger" role="alert">
                                {error}
                            </div>
                        )}

                        <div className="d-flex justify-content-end gap-2">
                            <button
                                type="button"
                                className="btn btn-outline-secondary"
                                onClick={() => navigate("/products")}
                            >
                                Cancel
                            </button>

                            <button
                                type="submit"
                                className="btn btn-primary"
                                disabled={loading}
                            >
                                {loading ? "Saving..." : "Submit Adjustment"}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}