CREATE TABLE product (
    product_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,  -- FK: người bán (user_id)
    category_id BIGINT  DEFAULT NULL,  -- FK: danh mục, cho phép NULL
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    stock INT NOT NULL DEFAULT 0,
    status ENUM('Active', 'Inactive', 'Pending') NOT NULL DEFAULT 'Pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    created_by CHAR(30) DEFAULT NULL,
    updated_by CHAR(30) DEFAULT NULL,


    CONSTRAINT fk_product_seller 
      FOREIGN KEY (user_id) REFERENCES `user`(user_id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
      
    CONSTRAINT fk_product_category 
      FOREIGN KEY (category_id) REFERENCES category(category_id)
      ON DELETE SET NULL
      ON UPDATE CASCADE
) ENGINE=InnoDB;

