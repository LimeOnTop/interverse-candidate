CREATE TABLE IF NOT EXISTS candidates (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    position VARCHAR(255),
    experience BIGINT NOT NULL DEFAULT 0,
    skills TEXT,
    resume_url VARCHAR(500),
    linkedin_url VARCHAR(500),
    github_url VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    interviewer_id BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS candidates_email_unique
    ON candidates (email)
    WHERE email IS NOT NULL;
