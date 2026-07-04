-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users.users (
    id UUID DEFAULT gen_random_uuid(),
    name varchar(64) NOT NULL,
    primary_email varchar(128) NOT NULL,
    is_deleted boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT uq_users_primary_email UNIQUE (primary_email)
);

CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users.users
    FOR EACH ROW
    EXECUTE FUNCTION users.trigger_set_updated_at();

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
