ALTER TABLE executions DROP COLUMN execution_error_output IF EXISTS, DROP COLUMN pod_name IF EXISTS, DROP COLUMN pod_cleaned IF EXISTS;