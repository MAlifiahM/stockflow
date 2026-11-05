INSERT INTO products (sku, name, price, stock) VALUES
('EL-001', 'Mechanical Keyboard', 750000, 25),
('EL-002', 'Wireless Mouse', 250000, 40),
('EL-003', 'USB-C Hub', 450000, 15),
('EL-004', 'Monitor 24 Inch', 175000, 12),
('EL-005', 'Laptop Stand', 175000, 30),
('EL-006', 'Webcam HD', 550000, 18),
('EL-007', 'Headset Gaming', 625000, 20),
('EL-008', 'USB Flash Drive 64GB', 95000, 50),
('EL-009', 'SSD External 1TB', 1250000, 10),
('EL-010', 'Mouse Pad XL', 120000, 35),
('EL-011', 'USB-C Cable', 75000, 60),
('EL-012', 'Bluetooth Speaker', 380000, 14)
ON CONFLICT (sku) DO NOTHING;