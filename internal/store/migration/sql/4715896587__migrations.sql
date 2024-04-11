
CREATE SEQUENCE IF NOT EXISTS public.db_migrations_id_seq AS INTEGER START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE IF NOT EXISTS public.db_migrations
(
    id         INTEGER                  DEFAULT nextval('public.db_migrations_id_seq'::regclass),
    name       CHARACTER VARYING(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT db_migrations_name UNIQUE ("name"),
    CONSTRAINT "db_migrations_pkey" PRIMARY KEY ("id")
);
ALTER SEQUENCE public.db_migrations_id_seq OWNED BY public.db_migrations.id;