CREATE TABLE saga_state (
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    INDEX idx_saga_state_order_id (order_id),
    INDEX idx_saga_state_status (status)
) ENGINE=InnoDB;