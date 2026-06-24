-- hris_db: 000006_create_attendances.up.sql

CREATE TABLE attendances (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id      UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    attendance_date  DATE NOT NULL DEFAULT CURRENT_DATE,
    check_in         TIME,
    check_out        TIME,
    status           VARCHAR(50) NOT NULL DEFAULT 'present'
);

CREATE INDEX idx_attendances_employee_date ON attendances (employee_id, attendance_date);
