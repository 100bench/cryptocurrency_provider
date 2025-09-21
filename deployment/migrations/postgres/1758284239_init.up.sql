-- Валюты и их тип
BEGIN;
CREATE TABLE symbols (
    code text PRIMARY KEY,                      -- 'BTC','ETH','USD'
    kind text NOT NULL CHECK (kind IN ('crypto','fiat'))
);

-- Котировки base→USD
CREATE TABLE rates (
    base_code   text NOT NULL REFERENCES symbols(code),   -- 'BTC','ETH' связать
    ts          timestamptz NOT NULL,                     -- момент котировки (UTC)
    price       numeric(20,8) NOT NULL,                   -- цена base в USD
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (base_code, ts)                                 -- идемпотентность по времени
);

-- Добавляем базовые валюты
INSERT INTO symbols (code, kind) VALUES 
    ('BTC', 'crypto'),
    ('ETH', 'crypto'),
    ('USD', 'fiat')
ON CONFLICT (code) DO NOTHING;

END;