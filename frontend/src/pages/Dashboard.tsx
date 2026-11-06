import { useEffect, useState } from "react";
import api from "../services/api";
import type { ProductPaginationResponse } from "../types/product";

export default function Dashboard() {
    const [totalProducts, setTotalProducts] = useState(0);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const response = await api.get<ProductPaginationResponse>(
                    "/products",
                    {
                        params: { page: 1, limit: 1 },
                    },
                );

                setTotalProducts(response.data.total_items);
            } catch (error) {
                console.error("Failed to fetch products:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, []);

    return (
        <div>
            <h4 className="mb-4">Dashboard</h4>

            <div className="row mb-4">
                <div className="col-md-4">
                    <div className="card border-0 shadow-sm">
                        <div className="card-body">
                            <p className="text-secondary mb-2">Total Products</p>
                            <h3 className="fw-bold mb-0">
                                {loading ? "..." : totalProducts}
                            </h3>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}