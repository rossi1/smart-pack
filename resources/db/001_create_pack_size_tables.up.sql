CREATE TABLE smartpack (
    id SERIAL PRIMARY KEY,
    size INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

INSERT INTO smartpack (size) VALUES (250), (500), (1000), (2000), (5000);

CREATE INDEX idx_smartpack_size ON smartpack (size);

CREATE OR REPLACE FUNCTION check_positive_smartpack_size()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.size <= 0 THEN
        RAISE EXCEPTION 'Size must be greater than 0';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_smartpack_size_check
BEFORE INSERT OR UPDATE ON smartpack
FOR EACH ROW
EXECUTE FUNCTION check_positive_smartpack_size();
