export interface Product {
    id: number;
    sku: string;
    name: string;
    price: number;
    stock: number;
    is_active: boolean;
    created_at: string;
    last_movement_at: string | null;
}

export interface ProductRequest {
    sku: string;
    name: string;
    price: number;
}

export interface ProductPaginationParams {
    search?: string;
    page?: number;
    limit?: number;
}

export interface ProductPaginationResponse {
    data: Product[];
    page: number;
    limit: number;
    total_items: number;
    total_pages: number;
}

export interface StockRequest {
    type: "in" | "out";
    qty: number;
    note?: string;
}

export interface StockMovement {
    id: number;
    product_id: number;
    type: "in" | "out";
    qty: number;
    note: string | null;
    created_at: string;
}