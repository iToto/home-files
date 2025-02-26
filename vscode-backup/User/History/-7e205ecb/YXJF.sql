CREATE TABLE social_url
(
    id VARCHAR(26) NOT NULL UNIQUE ,
    url_slug_id VARCHAR(26) UNIQUE ,
    social_urls JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT social_url_pkey PRIMARY KEY (id)
);

CREATE TABLE url_slug
(
    id VARCHAR(26) NOT NULL UNIQUE ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT url_slug_pkey PRIMARY KEY (id)
);

ALTER TABLE url_slug ADD FOREIGN KEY (id) REFERENCES social_url (url_slug_id);
