CREATE TABLE mvp_signal_log
(
id SERIAL UNIQUE ,
signal VARCHAR(255),
chain VARCHAR(255),
trade_time TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE,
CONSTRAINT mvp_signal_log_pkey PRIMARY KEY (id)
);

CREATE TABLE public.mvp_order
(
id VARCHAR(26) NOT NULL,
client_order_id VARCHAR(255),
strategy VARCHAR(255),
status VARCHAR(255),
currency_pair VARCHAR(255),
avg_price VARCHAR(255),
executed_qty VARCHAR(255),
finished_at TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE,
CONSTRAINT mvp_order_pkey PRIMARY KEY (id)
);

ALTER TABLE public.mvp_signal_log
ADD strategy VARCHAR(255);

CREATE TABLE public.mvp_user
(
id VARCHAR(26) NOT NULL,
name VARCHAR(255) UNIQUE ,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE,
CONSTRAINT mvp_user_pkey PRIMARY KEY (id)
);

CREATE TABLE public.mvp_signal_source
(
id VARCHAR(26) NOT NULL,
enabled BOOL,
type VARCHAR(255),
ip VARCHAR(255),
signal_version INTEGER,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE,
CONSTRAINT mvp_signal_source_pkey PRIMARY KEY (id)
);

CREATE TABLE public.mvp_strategy
(
id VARCHAR(26) NOT NULL,
enabled BOOL,
user_id VARCHAR(26),
signal_source_id VARCHAR(26),
type VARCHAR(255),
name VARCHAR(255),
exchange VARCHAR(255),
margin VARCHAR(25),
leverage FLOAT8,
fixed_trade_amount FLOAT8,
trade_strategy VARCHAR(255),
currency_pair VARCHAR(255),
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE,
CONSTRAINT mvp_strategy_pkey PRIMARY KEY (id)
);

ALTER TABLE public.mvp_order
ADD COLUMN side VARCHAR,
ADD COLUMN coin VARCHAR,
ADD COLUMN signal_id VARCHAR,
ADD COLUMN signal VARCHAR;
ALTER TABLE public.mvp_order ADD FOREIGN KEY (signal_id) REFERENCES mvp_signal_source (id);

-- TODO: Remove one of these relations as it's not normalized
ALTER TABLE public.mvp_strategy ADD FOREIGN KEY (user_id) REFERENCES mvp_user (id);
ALTER TABLE public.mvp_strategy ADD FOREIGN KEY (signal_source_id) REFERENCES public.mvp_signal_source (id);

-- Add Account Leverage column to strategy
ALTER TABLE public.mvp_strategy ADD COLUMN account_leverage FLOAT8;

-- DROP
DROP TABLE IF EXISTS mvp_signal_log;
DROP TABLE IF EXISTS mvp_trade;
DROP TABLE IF EXISTS mvp_signal_source;
DROP TABLE IF EXISTS mvp_strategy;
DROP TABLE IF EXISTS mvp_user;
