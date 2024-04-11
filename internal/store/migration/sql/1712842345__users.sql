
DROP TYPE IF EXISTS public.user_status;

CREATE TYPE  public.user_status AS ENUM ( 'pending', 'active' );
ALTER TYPE public.user_status OWNER TO postgres;

CREATE SEQUENCE IF NOT EXISTS public.user_id_seq AS INTEGER START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE IF NOT EXISTS public.user
(
    id            INTEGER                DEFAULT nextval('public.user_id_seq'::regclass) NOT NULL,
    name          CHARACTER VARYING(45)                                                  NOT NULL,
    email         CHARACTER VARYING(60)                                                  NOT NULL,
    password      CHARACTER VARYING(60)                                                  NOT NULL,
    status public.user_status DEFAULT 'pending',
    meta_data     json                   DEFAULT NULL,
    created_at    timestamp              DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp              DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_email_ukey UNIQUE ("email"),
    CONSTRAINT user_id_pkey PRIMARY KEY ("id")
);
ALTER SEQUENCE public.user_id_seq OWNED BY public.user.id;
CREATE INDEX if not exists user_email_index_key ON public.user(email);
