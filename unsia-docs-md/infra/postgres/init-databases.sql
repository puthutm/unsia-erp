-- Init databases for UNSIA ERP

SELECT 'Creating databases...' AS progress;

CREATE DATABASE core_db;
CREATE DATABASE reference_db;
CREATE DATABASE crm_db;
CREATE DATABASE pmb_db;
CREATE DATABASE finance_db;
CREATE DATABASE academic_db;
CREATE DATABASE hris_db;
CREATE DATABASE lms_db;
CREATE DATABASE assessment_db;
CREATE DATABASE portal_db;

SELECT 'Enabling pgcrypto extension on all databases...' AS progress;

\c core_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c reference_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c crm_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c pmb_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c finance_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c academic_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c hris_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c lms_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c assessment_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c portal_db
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SELECT 'Database initialization completed successfully.' AS progress;
