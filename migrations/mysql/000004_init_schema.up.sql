-- Bảng đơn hàng
CREATE TABLE `order` (
    order_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,  -- FK: người mua (user_id)
    seller_id BIGINT UNSIGNED NOT NULL,  -- <--- thêm cột này

    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    shipping_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    final_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    
    shipping_address TEXT NOT NULL,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_phone VARCHAR(20) NOT NULL,
    
--    payment_method ENUM('COD', 'BankTransfer', 'CreditCard', 'Wallet') NOT NULL,
--    payment_status ENUM('Pending', 'Paid', 'Failed', 'Refunded') NOT NULL DEFAULT 'Pending',
    
--    order_status ENUM('Pending', 'Confirmed', 'Shipping', 'Delivered', 'Cancelled', 'Refunded') NOT NULL DEFAULT 'Pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    created_by CHAR(30) DEFAULT NULL,
    updated_by CHAR(30) DEFAULT NULL,
    
    CONSTRAINT fk_order_user 
      FOREIGN KEY (user_id) REFERENCES `user`(user_id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
    CONSTRAINT fk_order_item_seller
      FOREIGN KEY (seller_id) REFERENCES `user`(user_id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE

) ENGINE=InnoDB;

-- Bảng chi tiết đơn hàng
CREATE TABLE order_item (
    order_item_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,  -- FK: đơn hàng
    product_id BIGINT UNSIGNED NOT NULL,  -- FK: sản phẩm
    
    quantity INT UNSIGNED NOT NULL,
    unit_price DECIMAL(15, 2) NOT NULL,  -- Giá tại thời điểm mua
    total_price DECIMAL(15, 2) NOT NULL,  -- Tổng giá của item (unit_price * quantity)
    
--    item_status ENUM('Normal', 'Cancelled', 'Returned') NOT NULL DEFAULT 'Normal',
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_order_item_order 
      FOREIGN KEY (order_id) REFERENCES `order`(order_id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
      
    CONSTRAINT fk_order_item_product 
      FOREIGN KEY (product_id) REFERENCES product(product_id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Index cho bảng order
CREATE INDEX idx_order_user ON `order` (user_id);
CREATE INDEX idx_order_created ON `order` (created_at);

-- Index cho bảng order_item
CREATE INDEX idx_order_item_product ON order_item (product_id);
CREATE INDEX idx_order_item_order ON order_item (order_id);

ALTER TABLE product 
ADD COLUMN price DECIMAL(15, 2) NOT NULL DEFAULT 0.00 
AFTER description;
