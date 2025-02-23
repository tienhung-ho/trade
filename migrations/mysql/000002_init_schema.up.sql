CREATE TABLE `category` (
                          category_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(200) NOT NULL,
                          description TEXT,
                          status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                          created_by CHAR(30) DEFAULT NULL,
                          updated_by CHAR(30) DEFAULT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          deleted_at DATETIME DEFAULT NULL
);

CREATE TABLE image (
                       image_id BIGINT AUTO_INCREMENT PRIMARY KEY,
                       url VARCHAR(300) NOT NULL,
                       alt_text VARCHAR(255),

    -- resource_id: Khoá tham chiếu tới ID của thực thể
                       resource_id BIGINT NOT NULL,

    -- resource_type: Xác định loại thực thể (product, category, account, blog_post, ...)
                       resource_type VARCHAR(50) NOT NULL,

                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
