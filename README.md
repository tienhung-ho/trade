# E-Commerce Platform with Private Blockchain Integration

## Overview
This project is a modern e-commerce platform integrated with a private blockchain for secure digital currency transactions. It supports wallet authentication, token staking, and face recognition-based identity verification. The system is built with a Go backend, Vue.js (Quasar) frontend, and a custom blockchain module using the Cosmos SDK.

## Features
- **Secure Transactions**: Leverages a private blockchain for token transfers and transaction validation.
- **Wallet Management**: Users can create and authenticate wallets using BIP39 mnemonic phrases.
- **Token Staking**: Supports staking functionality for internal tokens.
- **Face Recognition**: Implements FaceNet and OpenCV for secure user verification during login and transactions.
- **Responsive UI**: Built with Vue.js and Quasar for seamless product browsing, registration, and wallet management.
- **Scalable Backend**: Uses Go (Gin + GORM) for robust user, product, and order management.

## System Architecture
### Backend
- **Framework**: Go with Gin for RESTful APIs and GORM for MySQL database interactions.
- **Functionality**: Manages user authentication, product listings, order processing, and blockchain integration.
- **Database**: MySQL for persistent storage.

### Blockchain Layer
- **Framework**: Cosmos SDK with Ignite CLI for custom blockchain module development.
- **Features**:
  - Token balance management.
  - Wallet-based identity using BIP39 mnemonics.
  - Secure transaction validation.
- **Integration**: Communicates with the backend via gRPC for token transfers, wallet operations, and staking.

### Face Recognition System
- **Technology**: Python with FaceNet and OpenCV.
- **Functionality**: Provides identity verification for login and transaction authorization.
- **Integration**: Connected to the backend via gRPC.

### Frontend
- **Framework**: Vue.js with Quasar Framework for a responsive and modern UI.
- **Features**: Product browsing, user registration/login, wallet management, and token staking.

## Technologies Used
- **Backend**: Go, Gin, GORM, MySQL, Docker
- **Blockchain**: Cosmos SDK, Ignite CLI, gRPC
- **Face Recognition**: Python, FaceNet, OpenCV
- **Frontend**: Vue.js, Quasar Framework
- **Additional**: Redis, Supabase (for cloud storage)

## Prerequisites
- Go 1.20+
- Node.js 16+
- MySQL 8.0+
- Python 3.10+
- Docker
- Cosmos SDK and Ignite CLI
- Redis
- Supabase account for cloud storage

## Environment Setup
Create a `.env` file in the root directory (`tart-shop-manager`) based on the provided `.env.example`. The required environment variables are:

- **Server Configuration**:
  - `PORT`: Port for the backend server (e.g., `:3000`).
- **MySQL Database**:
  - `MYSQL_ROOT_PASSWORD`: Root password for MySQL.
  - `MYSQL_DATABASE`: Database name (e.g., `tradedb`).
  - `MYSQL_USER`: Database user (e.g., `admin`).
  - `MYSQL_PASSWORD`: Database user password.
  - `DB_URL`: Database connection URL for GORM.
  - `MYSQL_URL`: MySQL connection URL for migrations.
- **Redis**:
  - `REDIS_URL`: Redis connection URL (e.g., `redis://localhost:6379`).
- **Cryptography**:
  - `COST`: bcrypt cost factor for password hashing (e.g., `10`).
- **JWT Authentication**:
  - `ACCESS_SECRET_KEY`: Secret key for access tokens.
  - `REFRESH_SECRET_KEY`: Secret key for refresh tokens.
  - `WEB3_SECRET_KEY`: Secret key for Web3 access tokens.
  - `WEB3_REFRESH_TOKEN_KEY`: Secret key for Web3 refresh tokens.
- **gRPC**:
  - `END_POINT`: gRPC endpoint for blockchain and face recognition services (e.g., `localhost:9090`).
- **Blockchain**:
  - `ALICE`: Blockchain wallet address for admin (e.g., `cosmos15mkr0v0luhnqce5z07qs7uprgnns8fkxnr6zrx`).
  - `COIN_NAME`: Name of the custom token (e.g., `bitcoin`).
  - `ADMIN_NAME`: Admin account name (e.g., `alice`).
  - `WALLET_TYPE`: Wallet type for transactions (e.g., `bitcoin`).
- **Supabase Cloud Storage**:
  - `SUPABASE_ACCESS_KEY`: Access key for Supabase.
  - `SUPABASE_SECRET_KEY`: Secret key for Supabase.
  - `SUPABASE_ENDPOINT`: Supabase storage endpoint.
  - `SUPABASE_BUCKET`: Supabase bucket name (e.g., `trade`).
  - `SUPABASE_URL`: Supabase project URL.

A sample `.env.example` is provided in the repository for reference.

## Installation
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/tienhung-ho/tart-shop-manager.git
   cd tart-shop-manager
   ```

2. **Backend Setup**:
   - Copy the `.env.example` to `.env` and configure the variables:
     ```bash
     cp .env.example .env
     ```
   - Install dependencies:
     ```bash
     go mod tidy
     ```
   - Run the backend:
     ```bash
     go run ./cmd/app
     ```

5. **Frontend Setup**:
   - Navigate to the frontend directory:
     ```bash
     cd web/fe
     ```
   - Install dependencies:
     ```bash
     npm install
     ```
   - Run the development server:
     ```bash
     quasar dev
     ```

## License
This project is licensed under the MIT License.

## Contact
For questions or feedback, please open an issue on GitHub or contact the maintainers.