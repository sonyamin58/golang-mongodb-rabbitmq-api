-- Oracle Database Schema Initialization
-- Mini Bank Application
-- Compatible with Oracle XE 21c+

-- Drop existing tables (in reverse dependency order)
BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE transactions PURGE';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/
BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE accounts PURGE';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/
BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE users PURGE';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/
BEGIN
   EXECUTE IMMEDIATE 'DROP SEQUENCE users_seq';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/
BEGIN
   EXECUTE IMMEDIATE 'DROP SEQUENCE accounts_seq';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/
BEGIN
   EXECUTE IMMEDIATE 'DROP SEQUENCE transactions_seq';
EXCEPTION WHEN OTHERS THEN NULL;
END;
/

-- Create sequences
CREATE SEQUENCE users_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
CREATE SEQUENCE accounts_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
CREATE SEQUENCE transactions_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;

-- ============================================================
-- TABLE: users
-- ============================================================
CREATE TABLE users (
    id               NUMBER(20)      DEFAULT users_seq.NEXTVAL PRIMARY KEY,
    email            VARCHAR2(255)   NOT NULL UNIQUE,
    password_hash    VARCHAR2(255)   NOT NULL,
    full_name        VARCHAR2(255)   NOT NULL,
    phone_number     VARCHAR2(20)     NOT NULL UNIQUE,
    address          VARCHAR2(500),
    role             VARCHAR2(20)    DEFAULT 'user' NOT NULL,
    is_active        NUMBER(1)       DEFAULT 1 NOT NULL,
    created_at       TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at       TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT chk_user_role CHECK (role IN ('user', 'admin')),
    CONSTRAINT chk_user_active CHECK (is_active IN (0, 1))
);

COMMENT ON TABLE users IS 'User profile data for Mini Bank application';
COMMENT ON COLUMN users.id IS 'Primary key';
COMMENT ON COLUMN users.email IS 'User email address, must be unique';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.full_name IS 'User full name';
COMMENT ON COLUMN users.phone_number IS 'Phone number in format +62xxx';
COMMENT ON COLUMN users.role IS 'User role: user or admin';

-- ============================================================
-- TABLE: accounts
-- ============================================================
CREATE TABLE accounts (
    id               NUMBER(20)      DEFAULT accounts_seq.NEXTVAL PRIMARY KEY,
    user_id          NUMBER(20)      NOT NULL,
    account_number   VARCHAR2(50)    NOT NULL UNIQUE,
    balance          NUMBER(18,2)    DEFAULT 0 NOT NULL,
    currency        VARCHAR2(3)     DEFAULT 'IDR' NOT NULL,
    account_type     VARCHAR2(20)    DEFAULT 'savings' NOT NULL,
    is_active        NUMBER(1)       DEFAULT 1 NOT NULL,
    created_at       TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at       TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_account_type CHECK (account_type IN ('savings', 'checking')),
    CONSTRAINT chk_account_active CHECK (is_active IN (0, 1)),
    CONSTRAINT chk_balance_non_negative CHECK (balance >= 0)
);

COMMENT ON TABLE accounts IS 'Bank account data linked to users';
COMMENT ON COLUMN accounts.account_number IS 'Unique account number, 10 digits starting with 988';
COMMENT ON COLUMN accounts.balance IS 'Current account balance in IDR';

-- ============================================================
-- TABLE: transactions
-- ============================================================
CREATE TABLE transactions (
    id                  NUMBER(20)      DEFAULT transactions_seq.NEXTVAL PRIMARY KEY,
    user_id             NUMBER(20)      NOT NULL,
    account_id          NUMBER(20)      NOT NULL,
    transaction_type    VARCHAR2(20)    NOT NULL,
    amount              NUMBER(18,2)    NOT NULL,
    fee                 NUMBER(18,2)    DEFAULT 0 NOT NULL,
    total_amount        NUMBER(18,2)    NOT NULL,
    balance_before      NUMBER(18,2)   NOT NULL,
    balance_after       NUMBER(18,2)    NOT NULL,
    currency            VARCHAR2(3)     DEFAULT 'IDR' NOT NULL,
    
    -- Transfer specific
    to_account_id       NUMBER(20),
    to_account_number   VARCHAR2(50),
    to_account_name     VARCHAR2(255),
    
    -- Withdraw specific
    destination_bank    VARCHAR2(100),
    destination_account VARCHAR2(50),
    destination_name    VARCHAR2(255),
    
    -- Topup specific
    payment_method      VARCHAR2(50),
    sender_bank         VARCHAR2(100),
    sender_account      VARCHAR2(50),
    sender_name         VARCHAR2(255),
    
    reference_id        VARCHAR2(100)  UNIQUE,
    note                VARCHAR2(500),
    status              VARCHAR2(20)   DEFAULT 'pending' NOT NULL,
    celery_task_id      VARCHAR2(100),
    
    created_at          TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at          TIMESTAMP       DEFAULT SYSTIMESTAMP NOT NULL,
    
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_transactions_account FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT chk_tx_type CHECK (transaction_type IN ('TOPUP', 'WITHDRAW', 'TRANSFER_IN', 'TRANSFER_OUT')),
    CONSTRAINT chk_tx_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    CONSTRAINT chk_amount_positive CHECK (amount > 0)
);

COMMENT ON TABLE transactions IS 'All transaction history for audit trail';
COMMENT ON COLUMN transactions.celery_task_id IS 'Celery task ID for async processing tracking';

-- ============================================================
-- INDEXES
-- ============================================================

-- Users indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);

-- Accounts indexes
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_number ON accounts(account_number);
CREATE INDEX idx_accounts_active ON accounts(is_active);
CREATE INDEX idx_accounts_user_active ON accounts(user_id, is_active);

-- Transactions indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_reference ON transactions(reference_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_type_created ON transactions(transaction_type, created_at DESC);
CREATE INDEX idx_transactions_celery_task ON transactions(celery_task_id);
CREATE INDEX idx_transactions_to_account ON transactions(to_account_number);

-- ============================================================
-- SEED DATA (optional)
-- ============================================================

-- Admin user (password: Admin123!)
INSERT INTO users (id, email, password_hash, full_name, phone_number, role, is_active)
VALUES (users_seq.NEXTVAL, 'admin@minibank.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsW4WzOFbMB3dHI.Hu', 'System Admin', '+6280000000001', 'admin', 1);

COMMIT;
