import { useState} from "react";
import {useLocation, useNavigate, useParams} from "react-router-dom";
import api from "../services/api";
import type {Product} from "../types/product.ts";

export default function ProductForm() {
    const navigate = useNavigate();
    const { id } = useParams();
    const isEdit = Boolean(id);
    const location = useLocation();
    const product = location.state?.product as Product | undefined;

    const [sku, setSku] = useState(product?.sku ?? "");
    const [name, setName] = useState(product?.name ?? "");
    const [price, setPrice] = useState(product?.price ?? 0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    if (isEdit && !product) {
        navigate("/products", { replace: true });
        return null;
    }

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        setError("");

        const payload = {
            sku: sku.trim(),
            name: name.trim(),
            price: Number(price),
        };

        if (!payload.sku || !payload.name || !Number.isInteger(payload.price) || payload.price < 0) {
            setError("Please provide valid product information.");
            return;
        }

        setLoading(true);

        try {
            if (isEdit) {
                await api.put(`/products/${id}`, payload);
            } else {
                await api.post("/products", payload);
            }

            navigate("/products");
        } catch (error) {
            console.error("Failed to save product:", error);
            setError("Failed to save product. Please try again.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4">
                <h4>{isEdit ? "Edit Product" : "Add Product"}</h4>
                <button
                    type="button"
                    className="btn btn-outline-secondary"
                    onClick={() => navigate("/products")}
                    disabled={loading}
                >
                    Back
                </button>
            </div>

            <div className="card border-0 shadow-sm">
                <div className="card-body p-4">
                    <form onSubmit={handleSubmit}>
                        <div className="mb-3">
                            <label className="form-label">SKU</label>
                            <input
                                type="text"
                                className="form-control"
                                placeholder="e.g. EL-013"
                                required
                                defaultValue={sku}
                                disabled={isEdit}
                                onChange={(e) => setSku(e.target.value)}
                            />
                        </div>

                        <div className="mb-3">
                            <label className="form-label">Product Name</label>
                            <input
                                type="text"
                                className="form-control"
                                placeholder="Enter product name"
                                required
                                defaultValue={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                        </div>

                        <div className="mb-4">
                            <label className="form-label">Price (IDR)</label>
                            <input
                                type="number"
                                className="form-control"
                                placeholder="Enter price"
                                min="0"
                                step="1"
                                required
                                onChange={(e) => setPrice(Number(e.target.value))}
                                defaultValue={price}
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
                                disabled={loading}
                            >
                                Cancel
                            </button>

                            <button
                                type="submit"
                                className="btn btn-primary"
                                disabled={loading}
                            >
                                {loading
                                    ? "Saving..."
                                    : isEdit
                                        ? "Update Product"
                                        : "Save Product"}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}