import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../services/api";
import type { Product, ProductPaginationResponse } from "../types/product";
import type { OrderItemRequest } from "../types/order";

export default function CreateOrder() {
    const navigate = useNavigate();
    const [customer, setCustomer] = useState("");
    const [items, setItems] = useState<OrderItemRequest[]>([
        { sku: "", qty: 1 },
    ]);
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const response = await api.get<ProductPaginationResponse>(
                    "/products",
                    { params: { page: 1, limit: 100 } }
                );

                setProducts(response.data.data.filter(
                    (product) => product.is_active
                ));
            } catch (error) {
                console.error("Failed to fetch products:", error);
                setError("Failed to load products.");
            }
        };

        fetchProducts();
    }, []);

    const addItem = () => {
        setItems((prev) => [...prev, { sku: "", qty: 1 }]);
    };

    const removeItem = (index: number) => {
        setItems((prev) => prev.filter((_, i) => i !== index));
    };

    const updateItem = (
        index: number,
        field: "sku" | "qty",
        value: string | number
    ) => {
        setItems((prev) =>
            prev.map((item, i) =>
                i === index ? { ...item, [field]: value } : item
            )
        );
    };

    const total = items.reduce((sum, item) => {
        const product = products.find((p) => p.sku === item.sku);
        return sum + (product ? product.price * item.qty : 0);
    }, 0);

    const totalQuantity = items.reduce(
        (sum, item) => sum + item.qty,
        0
    );

    const handleSubmit = async () => {
        setError("");

        if (!customer.trim()) {
            setError("Customer name is required.");
            return;
        }

        if (items.length === 0) {
            setError("At least one item is required.");
            return;
        }

        if (items.some((item) => !item.sku || item.qty <= 0)) {
            setError("Please select a product and enter a valid quantity.");
            return;
        }

        const skus = items.map((item) => item.sku);

        if (new Set(skus).size !== skus.length) {
            setError("Duplicate products are not allowed.");
            return;
        }

        setLoading(true);

        try {
            await api.post("/orders", {
                customer: customer.trim(),
                items,
            });

            navigate(`/orders`);
        } catch (error) {
            console.error("Failed to create order:", error);

            setError("Failed to create order. Please check product stock and try again.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <div className="mb-4">
                <div className="d-flex justify-content-between align-items-center mb-4">
                    <div>
                        <h4 className="mb-1">Create Sales Order</h4>
                        <p className="text-muted mb-0">
                            Create a new sales order.
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

            <div className="card border-0 shadow-sm mb-4">
                <div className="card-header bg-white py-3">
                    <h5 className="mb-0">Customer Information</h5>
                </div>

                <div className="card-body">
                    <label className="form-label">Customer Name</label>
                    <input
                        type="text"
                        className="form-control"
                        placeholder="Enter customer name"
                        value={customer}
                        onChange={(e) => setCustomer(e.target.value)}
                        disabled={loading}
                    />
                </div>
            </div>

            <div className="card border-0 shadow-sm mb-4">
                <div className="card-header bg-white py-3 d-flex justify-content-between align-items-center">
                    <h5 className="mb-0">Order Items</h5>

                    <button
                        type="button"
                        className="btn btn-sm btn-outline-primary"
                        onClick={addItem}
                        disabled={loading}
                    >
                        Add Item
                    </button>
                </div>

                <div className="card-body">
                    <div className="table-responsive">
                        <table className="table align-middle">
                            <thead className="table-light">
                            <tr>
                                <th>Product</th>
                                <th>SKU</th>
                                <th>Available Stock</th>
                                <th style={{ width: "130px" }}>Quantity</th>
                                <th className="text-end">Action</th>
                            </tr>
                            </thead>

                            <tbody>
                            {items.map((item, index) => {
                                const product = products.find(
                                    (p) => p.sku === item.sku
                                );

                                return (
                                    <tr key={index}>
                                        <td style={{ minWidth: "220px" }}>
                                            <select
                                                className="form-select"
                                                value={item.sku}
                                                onChange={(e) =>
                                                    updateItem(index, "sku", e.target.value)
                                                }
                                                disabled={loading}
                                            >
                                                <option value="">Select Product</option>

                                                {products.map((product) => (
                                                    <option
                                                        key={product.id}
                                                        value={product.sku}
                                                    >
                                                        {product.name}
                                                    </option>
                                                ))}
                                            </select>
                                        </td>

                                        <td>{product?.sku ?? "-"}</td>

                                        <td>{product?.stock ?? "-"}</td>

                                        <td>
                                            <input
                                                type="number"
                                                className="form-control"
                                                min="1"
                                                value={item.qty}
                                                onChange={(e) =>
                                                    updateItem(
                                                        index,
                                                        "qty",
                                                        Number(e.target.value)
                                                    )
                                                }
                                                disabled={loading}
                                            />
                                        </td>

                                        <td className="text-end">
                                            <button
                                                type="button"
                                                className="btn btn-sm btn-outline-danger"
                                                onClick={() => removeItem(index)}
                                                disabled={loading || items.length === 1}
                                            >
                                                Remove
                                            </button>
                                        </td>
                                    </tr>
                                );
                            })}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <div className="card border-0 shadow-sm mb-4">
                <div className="card-body">
                    <h5 className="mb-3">Order Summary</h5>

                    <div className="d-flex justify-content-between mb-2">
                        <span className="text-muted">Total Quantity</span>
                        <strong>{totalQuantity}</strong>
                    </div>

                    <div className="d-flex justify-content-between">
                        <span className="text-muted">Total Amount</span>
                        <strong>
                            {new Intl.NumberFormat("id-ID", {
                                style: "currency",
                                currency: "IDR",
                                maximumFractionDigits: 0,
                            }).format(total)}
                        </strong>
                    </div>
                </div>
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
                    onClick={() => navigate("/orders")}
                    disabled={loading}
                >
                    Cancel
                </button>

                <button
                    type="button"
                    className="btn btn-primary"
                    onClick={handleSubmit}
                    disabled={loading || products.length === 0}
                >
                    {loading ? "Creating..." : "Create Order"}
                </button>
            </div>
        </div>
    );
}