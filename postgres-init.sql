
# create database portfolio;

# \c portfolio;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE portfolio (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    ticker VARCHAR(512) NOT NULL unique,
    quantity  float NOT NULL,
    price float NOT NULL,
    update_time TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE OR REPLACE FUNCTION portfolio_notify_func() RETURNS trigger as $$
DECLARE
  payload text;
BEGIN
	IF TG_OP = 'DELETE' THEN
    payload := row_to_json(tmp)::text FROM (
			SELECT
				OLD.id as id,
				TG_OP as _op,
				TG_TABLE_NAME as _tablename
		) tmp;
	ELSE
		payload := row_to_json(tmp)::text FROM (
			SELECT
				NEW.*,
				TG_TABLE_NAME as _tablename,
				TG_OP as _op
		) tmp;
		IF octet_length( payload ) > 8000 THEN
			-- payload is too big for a pg_notify.
			payload := row_to_json(tmp)::text FROM (
				SELECT
					NEW.id as id,
					'payload length > 8000 bytes' as error,
					TG_TABLE_NAME as _tablename,
					TG_OP as _op
			) tmp;
		END IF;
	END IF;
  PERFORM pg_notify('portfolio_update'::text, payload);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

create trigger portfolio_notify_trig after insert or update or delete on portfolio for each row execute procedure portfolio_notify_func();


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


insert into portfolio (ticker, quantity, price) values ('FB', 230, 130);
insert into portfolio (ticker, quantity, price) values ('AAPL', 100, 110);
insert into portfolio (ticker, quantity, price) values ('GOOG', 410, 650);
insert into portfolio (ticker, quantity, price) values ('AMZN', 310, 700);
