-- Create the `user` table
CREATE TABLE `user` (
  user_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  
  -- Thông tin đăng nhập và định danh
  fullname VARCHAR(300) DEFAULT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(20) DEFAULT NULL,
  
  -- Trạng thái và vai trò
  status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
  gender ENUM('Male', 'Female', 'Other') DEFAULT NULL,
  
  -- Thông tin hồ sơ mở rộng dạng JSON
  profile JSON DEFAULT NULL,
  
  -- Thông tin liên quan đến đăng nhập
  last_login DATETIME DEFAULT NULL,
  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,   -- Thêm trường deleted_at cho Soft Delete
  created_by CHAR(30) DEFAULT NULL,
  updated_by CHAR(30) DEFAULT NULL
) ENGINE=InnoDB;

-- Create the `user_wallets` table
CREATE TABLE `user_wallets` (
  wallet_id INT AUTO_INCREMENT PRIMARY KEY,
  
  -- Liên kết với người dùng
  user_id BIGINT UNSIGNED NOT NULL,
  
  -- Thông tin ví blockchain
  wallet_address TEXT NOT NULL,      -- Ví dụ: địa chỉ Ethereum có 42 ký tự
  encrypted_mnemonic TEXT DEFAULT NULL,       -- Mnemonic đã được mã hóa (không lưu plaintext)
  wallet_type ENUM('ethereum','bitcoin','cosmos','solana','citcoin','other') DEFAULT 'other',
  
  -- Số dư ví (được cập nhật định kỳ)
  balance DECIMAL(65,18) DEFAULT 0,
  
  -- Thông tin bổ sung dạng JSON (ví dụ: cấu hình, metadata)
  metadata JSON DEFAULT NULL,
  
  -- Thời gian tạo và cập nhật
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  -- Các cột cho Soft Delete và audit
  deleted_at DATETIME DEFAULT NULL,
  status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
  created_by CHAR(30) DEFAULT NULL,
  updated_by CHAR(30) DEFAULT NULL,
  
  -- Ràng buộc liên kết với bảng `user`
  FOREIGN KEY (user_id) REFERENCES `user`(user_id) ON DELETE CASCADE
) ENGINE=InnoDB;

