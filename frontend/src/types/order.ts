export type OrderStatus = "CREATED" | "CANCELLED";

export interface OrderItem {
    id: number;
    order_id: number;
    sku: string;
    qty: number;
    unit_price: number;
}

export interface Order {
    id: number;
    customer: string;
    status: OrderStatus;
    total: number;
    created_at: string;
    cancelled_at: string | null;
    items: OrderItem[];
}

export interface OrderPaginationResponse {
    data: Order[];
    page: number;
    limit: number;
    total_items: number;
    total_pages: number;
}

export interface OrderItemRequest {
    sku: string;
    qty: number;
}

export interface OrderRequest {
    customer: string;
    items: OrderItemRequest[];
}