CREATE TABLE social-url
(
    id VARCHAR(26) NOT NULL UNIQUE ,
    url_slug_id VARCHAR(26) UNIQUE ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT 'CURRENT_TIMESTAMP' NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT social-url_pkey PRIMARY KEY (id)
);

CREATE TABLE url-slug
(
    id VARCHAR(26) NOT NULL UNIQUE ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT 'CURRENT_TIMESTAMP' NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT url-slug_pkey PRIMARY KEY (id)
);

ALTER TABLE url-slug ADD FOREIGN KEY (id) REFERENCES social-url (url_slug_id);
REFERENCES social-url (url_slug_id);
