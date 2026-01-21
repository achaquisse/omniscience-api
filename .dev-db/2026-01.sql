
-- FEV 2026 Changes
ALTER TABLE `Course` ADD COLUMN `teacher_email` VARCHAR(500);
UPDATE `Course` SET `teacher_email` = 'animake.co.mz@gmail.com' WHERE `id` IN (1, 3, 4);

CREATE TABLE `Attendance` (
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `registration_id` bigint(20) NOT NULL,
                              `date` date NOT NULL,
                              `status` varchar(20) NOT NULL,
                              `remarks` text DEFAULT NULL,
                              `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
                              `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `unique_registration_date` (`registration_id`, `date`),
                              KEY `idx_registration_id` (`registration_id`),
                              KEY `idx_date` (`date`),
                              KEY `idx_status` (`status`),
                              CONSTRAINT `fk_attendance_registration` FOREIGN KEY (`registration_id`) REFERENCES `Registration` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE StudentClass MODIFY COLUMN disabled TINYINT(1) NOT NULL DEFAULT 0;"