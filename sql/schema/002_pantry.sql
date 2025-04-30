-- gooseUp
CREATE TABLE pantry (
    id uuid NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    name text NOT NULL,
    added_at timestamp NOT NULL DEFAULT now(),
    expiry_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- gooseDown
DROP TABLE pantry;