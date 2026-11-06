import { useEffect, useState } from "react";
import api from "../services/api";
import type {
    Product,
    ProductPaginationResponse,
} from "../types/product";
import { useNavigate } from "react-router-dom";
import axios from "axios";


export default function Products() {
    const navigate = useNavigate();
    const [products, setProducts] = useState<Product[]>([]);
    const [search, setSearch] = useState("");
    const [page, setPage] = useState(1);
    const [limit] = useState(10);
    const [totalPages, setTotalPages] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [success, setSuccess] = useState("");
    const [refresh, setRefresh] = useState(0);

    useEffect(() => {
        const fetchProducts = async () => {
            setLoading(true);
            setError("");
            try {
                const response = await api.get<ProductPaginationResponse>(
                    "/products",
                    {
                        params: {
                            page: page,
                            limit: limit,
                            search: search,
                        }
                    },
                );

                setProducts(response.data.data);
                setTotalPages(response.data.total_pages);
            } catch (error) {
                console.error("failed to fetch products: ", error);
                setError("Failed to fetch products");
            } finally {
            setLoading(false);
            }
        };

        fetchProducts();
    }, [page, limit, search, refresh]);

    const handleDelete = async (id: number) => {
        setError("");
        setSuccess("");

        try {
            await api.delete(`/products/${id}`);

            setSelectedProduct(null);
            setSuccess("Product deleted successfully.");
            setRefresh((prev) => prev + 1);
        } catch (error) {
            console.error("Failed to delete product:", error);
            if (axios.isAxiosError(error)) {
                const errorCode = error.response?.data?.error?.code;

                if (errorCode === "PRODUCT_HAS_HISTORY") {
                    setError("Product has history.");
                } else {
                    console.log(errorCode);
                    setError("Something went wrong.");
                }
            }
        }
    };

    return (
        <div>
            <div className="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-4">
                <h4 className="mb-0">Products</h4>

                <button
                    className="btn btn-primary"
                    onClick={() => navigate("/products/new")}
                >
                    Add Product
                </button>
            </div>
            <input
                type="text"
                className="form-control mb-3"
                placeholder="Search products..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
            />

            {loading && <p>Loading...</p>}

            <div className="card border-0 shadow-sm">
                <div className="table-responsive">
                    <table className="table table-hover align-middle mb-0"  style={{ minWidth: "1100px" }}>
                        <thead className="table-light">
                        <tr>
                            <th>SKU</th>
                            <th>Product Name</th>
                            <th>Price</th>
                            <th>Stock</th>
                            <th>Status</th>
                            <th>Last Movement</th>
                            <th className="text-end">Actions</th>
                        </tr>
                        </thead>

                        <tbody>
                        {products.length === 0 ? (
                            <tr>
                                <td colSpan={7} className="text-center py-4 text-secondary">
                                    No products found
                                </td>
                            </tr>
                        ) : (
                            products.map((product) => (
                                <tr key={product.id}>
                                    <td className="fw-semibold">{product.sku}</td>

                                    <td>{product.name}</td>

                                    <td>
                                        {new Intl.NumberFormat("id-ID", {
                                            style: "currency",
                                            currency: "IDR",
                                            maximumFractionDigits: 0,
                                        }).format(product.price)}
                                    </td>

                                    <td>{product.stock}</td>

                                    <td>
                                <span
                                    className={`badge ${
                                        product.is_active
                                            ? "text-bg-success"
                                            : "text-bg-secondary"
                                    }`}
                                >
                                    {product.is_active ? "Active" : "Inactive"}
                                </span>
                                    </td>

                                    <td>
                                        {product.last_movement_at
                                            ? new Date(product.last_movement_at).toLocaleString("id-ID")
                                            : "-"}
                                    </td>

                                    <td>
                                        <div className="d-flex flex-nowrap justify-content-end gap-2">
                                            <button
                                                className="btn btn-sm btn-outline-success flex-shrink-0"
                                                onClick={() =>
                                                    navigate(`/products/${product.id}/movements`, {
                                                        state: { product },
                                                    })
                                                }
                                            >
                                                Movements
                                            </button>

                                            <button
                                                className="btn btn-sm btn-outline-primary flex-shrink-0"
                                                onClick={() =>
                                                    navigate(`/products/${product.id}/edit`, {
                                                        state: { product },
                                                    })
                                                }
                                            >
                                                Edit
                                            </button>

                                            <button
                                                className="btn btn-sm btn-outline-warning flex-shrink-0"
                                                onClick={() =>
                                                    navigate(`/products/${product.id}/stock`, {
                                                        state: { product },
                                                    })
                                                }
                                            >
                                                Edit Stock
                                            </button>

                                            <button
                                                className="btn btn-sm btn-outline-danger flex-shrink-0"
                                                onClick={() => setSelectedProduct(product)}
                                            >
                                                Delete
                                            </button>
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

            {(selectedProduct || success || error) && (
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
                                        ? "Delete Successful"
                                        : error
                                            ? "Delete Failed"
                                            : "Confirm Delete"}
                                </h5>
                                <button
                                    type="button"
                                    className="btn-close"
                                    onClick={() => {
                                        setSelectedProduct(null);
                                        setSuccess("");
                                        setError("");
                                    }}
                                />
                            </div>

                            <div className="modal-body">
                                {success ? (
                                    <div className="text-success">
                                        {success}
                                    </div>
                                ) : error ? (
                                    <div className="text-danger">
                                        {error}
                                    </div>
                                ) : (
                                    <>
                                        Are you sure you want to delete product{" "}
                                        <strong>{selectedProduct?.name}</strong>?
                                        <p className="text-muted mb-0 mt-2">
                                            This action cannot be undone.
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
                                            setSelectedProduct(null);
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
                                            onClick={() => setSelectedProduct(null)}
                                        >
                                            Cancel
                                        </button>

                                        <button
                                            type="button"
                                            className="btn btn-danger"
                                            onClick={() => handleDelete(selectedProduct!.id)}
                                        >
                                            Delete
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