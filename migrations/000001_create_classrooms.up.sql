-- Create classrooms table
CREATE TABLE classrooms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    organization_name VARCHAR(255) NOT NULL,
    organization_id BIGINT NOT NULL,
    instructor_id BIGINT NOT NULL,
    instructor_login VARCHAR(255) NOT NULL,
    public BOOLEAN NOT NULL DEFAULT false,
    archived BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE UNIQUE INDEX idx_classrooms_slug ON classrooms (slug);
CREATE INDEX idx_classrooms_organization_name ON classrooms (organization_name);
CREATE INDEX idx_classrooms_instructor_id ON classrooms (instructor_id);
CREATE INDEX idx_classrooms_public ON classrooms (public) WHERE public = true;
CREATE INDEX idx_classrooms_archived ON classrooms (archived) WHERE archived = true;
CREATE INDEX idx_classrooms_created_at ON classrooms (created_at);

-- Add constraints
ALTER TABLE classrooms ADD CONSTRAINT chk_classrooms_slug_format
    CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$');
ALTER TABLE classrooms ADD CONSTRAINT chk_classrooms_name_length
    CHECK (char_length(name) >= 1 AND char_length(name) <= 255);
ALTER TABLE classrooms ADD CONSTRAINT chk_classrooms_slug_length
    CHECK (char_length(slug) >= 1 AND char_length(slug) <= 255);