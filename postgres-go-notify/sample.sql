create table audit(
	log text
);
create table audit_log(
	username TEXT,
	add_time TIMESTAMP,
	table_name TEXT,
	operation TEXT,
	row_before JSON,
	row_after JSON
);


CREATE OR REPLACE FUNCTION fn_audit_trigger()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS
$$
	DECLARE
		old_row json = null;
		new_row json = null;
		notifyPayload text = null;
		
	BEGIN
		IF TG_OP IN ('DELETE', 'UPDATE') THEN
			old_row = row_to_json(OLD);
		END IF;
		
		IF TG_OP IN ('INSERT', 'UPDATE') THEN
			new_row = row_to_json(NEW);
		END IF;
		
		INSERT INTO audit_log
		(
			username,
			add_time,
			table_name,
			operation,
			row_before,
			row_after
		)
		VALUES
		(
			SESSION_USER,
			CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
			TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME,
			TG_OP,
			old_row,
			new_row
		);
		notifyPayload = json_build_object(
			'action', TG_OP,
            'table_name',  TG_TABLE_NAME,
			'new', row_to_json(NEW),
			'old',  row_to_json(OLD)
            )::text;
		PERFORM pg_notify('myChannel', notifyPayload);
		RETURN NEW;
	END;
$$


CREATE TRIGGER tg_audit_log
AFTER INSERT OR DELETE OR UPDATE
ON audit
FOR EACH ROW
EXECUTE PROCEDURE fn_audit_trigger();


select * from audit;
INSERT INTO audit (log) values ('log-1');
select * from audit_log;

UPDATE audit set log = 'log-100' where log = 'log-1';

DELETE from audit where log = 'log-100';

LISTEN "myChannel";
UNLISTEN "myChannel";
