CREATE TABLE public.accounts (
    account_id character varying,
    label character varying NOT NULL,
    type character varying NOT NULL,
    currency character varying NOT NULL,
    CONSTRAINT accounts_pkey PRIMARY KEY (account_id)
);