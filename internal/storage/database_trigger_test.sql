-- `database_schema_test.sql` is generated using prisma's schema dump functionality. This functionality does not support dumping out functions or triggers. So they are manually added in this file. Details on generating and updating `database_schema_test.sql` can be found in the frontend README. This file is manually updated as and when needed.

-- CreateUpdateAtFunction
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Candidate updated_at trigger
CREATE TRIGGER update_candidate_updated_at BEFORE UPDATE ON candidates FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

-- FileUpload updated_at trigger
CREATE TRIGGER update_file_upload_updated_at BEFORE UPDATE ON file_uploads FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();
