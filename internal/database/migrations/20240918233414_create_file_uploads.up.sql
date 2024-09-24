CREATE TABLE file_uploads (
    id UUID PRIMARY KEY,  
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Attach the trigger function to the `file_uploads` table
CREATE TRIGGER trigger_update_file_uploads_updated_at
BEFORE UPDATE ON file_uploads
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();