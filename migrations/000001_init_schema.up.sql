CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,                
    url TEXT NOT NULL UNIQUE,                
    title TEXT NOT NULL,                    
    current_price DECIMAL(12, 2) DEFAULT 0,  
    target_price DECIMAL(12, 2) DEFAULT 0,   
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)

CREATE INDEX IF NOT EXISTS idx_products_url ON products (url);