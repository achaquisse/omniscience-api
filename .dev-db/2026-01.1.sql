
-- Add created_by and updated_by columns to Attendance table
ALTER TABLE `Attendance` ADD COLUMN `created_by` VARCHAR(500) DEFAULT NULL;
ALTER TABLE `Attendance` ADD COLUMN `updated_by` VARCHAR(500) DEFAULT NULL;
