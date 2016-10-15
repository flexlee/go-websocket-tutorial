
CREATE OR REPLACE FUNCTION update_change_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.update_time = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_change_timestamp_trig BEFORE UPDATE
    ON portfolio FOR EACH ROW EXECUTE PROCEDURE
    update_change_timestamp();
