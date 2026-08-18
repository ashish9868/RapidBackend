CREATE TABLE superadmins (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email_verified_at DATETIME,
    is_active BOOLEAN DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    settings JSON,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255),

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    email_verified_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT 0,

    permissions_json JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unq_project_email
        UNIQUE (project_id, email),

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE access_key_tokens (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,
    collection VARCHAR(255) NOT NULL,

    access_token VARCHAR(255),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE project_collections (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255) NOT NULL,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',
    sort_order INTEGER DEFAULT 0,
    required BOOLEAN DEFAULT 0,

    rules JSON,
    options JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE project_collection_fields (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',

    is_required BOOLEAN DEFAULT 0,
    is_indexed BOOLEAN DEFAULT 0,
    is_unique BOOLEAN DEFAULT 0,
    is_sortable BOOLEAN DEFAULT 0,
    is_filterable BOOLEAN DEFAULT 0,

    options JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (collection_id)
        REFERENCES project_collections(id)
        ON DELETE CASCADE
);


CREATE TABLE project_collection_records (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,

    data JSON,

    version INTEGER DEFAULT 1,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',

    is_required BOOLEAN DEFAULT 0,
    is_indexed BOOLEAN DEFAULT 0,
    is_unique BOOLEAN DEFAULT 0,
    is_sortable BOOLEAN DEFAULT 0,
    is_filterable BOOLEAN DEFAULT 0,

    options JSON,

    created_by_id VARCHAR(255) NOT NULL,
    updated_by_id VARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (collection_id)
        REFERENCES project_collections(id)
        ON DELETE CASCADE,

    FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);


CREATE TABLE project_pages (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255) NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    slug VARCHAR(255) NOT NULL UNIQUE,

    settings JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unq_project_page
        UNIQUE (project_id, title),

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE email_templates (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    name VARCHAR(255) UNIQUE,
    description TEXT,

    is_system_template BOOLEAN DEFAULT 1,

    html_content TEXT,
    text_content TEXT,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE settings (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255),

    value TEXT,

    created_by_id VARCHAR(255) NOT NULL,
    updated_by_id VARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);