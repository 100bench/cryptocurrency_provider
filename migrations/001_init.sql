-- Валюты и их тип
CREATE TABLE symbols (
    code text PRIMARY KEY,                      -- 'BTC','ETH','USD'
    kind text NOT NULL CHECK (kind IN ('crypto','fiat'))
);

-- Котировки base→USD
CREATE TABLE rates (
    id          bigserial PRIMARY KEY,
    base_code   text NOT NULL REFERENCES symbols(code),   -- 'BTC','ETH'
    ts          timestamptz NOT NULL,                     -- момент котировки (UTC)
    price       numeric(20,8) NOT NULL,                   -- цена base в USD
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (base_code, ts)                                 -- идемпотентность по времени
);

