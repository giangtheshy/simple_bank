CREATE TABLE public.users (
	username varchar NOT NULL,
	full_name varchar NOT NULL,
	hash_password varchar NOT NULL,
	email varchar NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	changed_password_at timestamptz NOT NULL DEFAULT '2023-06-29 21:13:17.733+07'::timestamp with time zone,
	CONSTRAINT users_pk PRIMARY KEY (username),
	CONSTRAINT users_un UNIQUE (email)
);

ALTER TABLE public.accounts ADD CONSTRAINT accounts_fk FOREIGN KEY ("owner") REFERENCES public.users(username) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.accounts ADD CONSTRAINT accounts_un UNIQUE ("owner",currency);
