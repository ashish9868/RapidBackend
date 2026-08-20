CREATE TABLE superadmins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT,
    last_name TEXT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email_verified_at DATETIME,
    is_active BOOLEAN DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    slug TEXT NOT NULL UNIQUE,
    settings JSON,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    project_id INTEGER,

    first_name TEXT,
    last_name TEXT,
    email TEXT NOT NULL,
    password TEXT NOT NULL,

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


CREATE TABLE project_collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    project_id INTEGER NOT NULL,

    name TEXT,
    type TEXT DEFAULT 'base',
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    collection_id INTEGER NOT NULL,

    name TEXT,
    type TEXT DEFAULT 'base',

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
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    collection_id INTEGER NOT NULL,

    data JSON,

    version INTEGER DEFAULT 1,

    name TEXT,
    type TEXT DEFAULT 'base',

    is_required BOOLEAN DEFAULT 0,
    is_indexed BOOLEAN DEFAULT 0,
    is_unique BOOLEAN DEFAULT 0,
    is_sortable BOOLEAN DEFAULT 0,
    is_filterable BOOLEAN DEFAULT 0,

    options JSON,

    created_by_id INTEGER NOT NULL,
    updated_by_id INTEGER NOT NULL,

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


CREATE TABLE access_key_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    collection_id INTEGER NOT NULL,
    collection TEXT NOT NULL,

    access_token TEXT,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (collection_id)
        REFERENCES project_collections(id)
        ON DELETE CASCADE
);


CREATE TABLE project_pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    project_id INTEGER NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    slug TEXT NOT NULL UNIQUE,

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
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name TEXT UNIQUE,
    description TEXT,

    is_system_template BOOLEAN DEFAULT 1,

    html_content TEXT,
    text_content TEXT,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    project_id INTEGER,

    value TEXT,

    created_by_id INTEGER NOT NULL,
    updated_by_id INTEGER NOT NULL,

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