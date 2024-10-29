CREATE TABLE share_requests (
    id UUID NOT NULL PRIMARY KEY,
    request_by UUID NOT NULL,
    request_to UUID NOT NULL,
    rsa_public_key TEXT NOT NULL,
    status VARCHAR(255) NOT NULL,
    user_profile_json TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE share_request_files (
    share_request_id UUID NOT NULL,
    file_id UUID NOT NULL,
    PRIMARY KEY (share_request_id, file_id),    
    CONSTRAINT fk_share_request_file_share_request_id FOREIGN KEY (share_request_id) REFERENCES share_requests(id)
);

ALTER TABLE share_requests
ADD CONSTRAINT fk_share_requests_request_by FOREIGN KEY (request_by) REFERENCES users(id);

ALTER TABLE share_requests
ADD CONSTRAINT fk_share_requests_request_to FOREIGN KEY (request_to) REFERENCES users(id);

-- Attach the trigger function to the `share_requests` table
CREATE TRIGGER trigger_update_share_requests_updated_at
BEFORE UPDATE ON share_requests
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();