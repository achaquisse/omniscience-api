-- ========================================================================
-- OMNISCIENCE API - CONSOLIDATED GRADING SYSTEM SCHEMA
-- ========================================================================
-- This is a consolidated migration that combines the original 2026-02.sql
-- schema with all subsequent migrations (001-004).
--
-- Changes incorporated:
-- - 001: Removed status fields, added is_active boolean
-- - 002: Merged EvaluationItem into EvaluationGrade table
-- - 003: Removed weight/drop_lowest/is_extra_credit from EvaluationCategory
-- - 004: Added level field to Registration for level-based grading
--
-- Date: February 2026
-- ========================================================================

-- --------------------------------------------------------
-- Table: EvaluationCategory
-- Defines evaluation categories for each EvaluationGrade
-- Simplified from original - removed weight, drop_lowest, is_extra_credit
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationCategory` (
                                                    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `name` varchar(100) NOT NULL,
    `description` text DEFAULT NULL,
    `max_score` decimal(10,2) NOT NULL DEFAULT 20.00 COMMENT 'Maximum possible score',
    `display_order` int(11) DEFAULT 0,
    `is_active` tinyint(1) DEFAULT 1,
    `created_at` datetime DEFAULT current_timestamp(),
    `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    `created_by` varchar(500) DEFAULT NULL,
    `updated_by` varchar(500) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_is_active` (`is_active`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraint if Course table exists
SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Course') > 0,
  'ALTER TABLE `EvaluationCategory` ADD CONSTRAINT `fk_eval_category_course` FOREIGN KEY (`course_id`) REFERENCES `Course` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------
-- Table: EvaluationGradingFormula
-- Stores grading formula configurations for courses
-- Supports CUSTOM formulas with JSON configuration for complex calculations
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationGradingFormula` (
                                                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `formula_type` varchar(50) NOT NULL COMMENT 'WEIGHTED_AVERAGE, POINTS_BASED, PASS_FAIL, CUSTOM',
    `formula_config` json DEFAULT NULL COMMENT 'JSON configuration for CUSTOM formulas (multi-stage, level-based, etc.)',
    `passing_percentage` decimal(5,2) DEFAULT NULL COMMENT 'Minimum percentage to pass',
    `grading_scale` json DEFAULT NULL COMMENT 'JSON defining grade ranges (A, B, C, etc.)',
    `is_active` tinyint(1) DEFAULT 1,
    `created_at` datetime DEFAULT current_timestamp(),
    `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    `created_by` varchar(500) DEFAULT NULL,
    `updated_by` varchar(500) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_is_active` (`is_active`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraint if Course table exists
SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Course') > 0,
  'ALTER TABLE `EvaluationGradingFormula` ADD CONSTRAINT `fk_grading_formula_course` FOREIGN KEY (`course_id`) REFERENCES `Course` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------
-- Table: EvaluationGrade
-- Individual student grades for evaluations
-- Merged with EvaluationItem - now contains both item and grade data
-- Removed status field (simplified workflow)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationGrade` (
                                                 `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `category_id` bigint(20) NOT NULL,
    `student_class_id` bigint(20) NOT NULL,
    `registration_id` bigint(20) NOT NULL,
    `name` varchar(200) NOT NULL COMMENT 'Name of the evaluation (e.g., Quiz 1, Midterm)',
    `description` text DEFAULT NULL,
    `date` date DEFAULT NULL,
    `due_date` datetime DEFAULT NULL,
    `max_score` decimal(10,2) NOT NULL,
    `score` decimal(10,2) DEFAULT NULL,
    `percentage` decimal(5,2) DEFAULT NULL COMMENT 'Calculated percentage',
    `letter_grade` varchar(5) DEFAULT NULL COMMENT 'Letter grade (A, B, C, etc.)',
    `is_excused` tinyint(1) DEFAULT 0 COMMENT 'Excused from this evaluation',
    `is_late` tinyint(1) DEFAULT 0 COMMENT 'Submitted late',
    `late_penalty` decimal(5,2) DEFAULT 0.00 COMMENT 'Penalty percentage for late submission',
    `comments` text DEFAULT NULL,
    `graded_by` varchar(500) DEFAULT NULL,
    `graded_at` datetime DEFAULT NULL,
    `created_at` datetime DEFAULT current_timestamp(),
    `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    `created_by` varchar(500) DEFAULT NULL,
    `updated_by` varchar(500) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_student_class_id` (`student_class_id`),
    KEY `idx_registration_id` (`registration_id`),
    KEY `idx_date` (`date`),
    KEY `idx_graded_at` (`graded_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraints
ALTER TABLE `EvaluationGrade` ADD CONSTRAINT `fk_grade_category` FOREIGN KEY (`category_id`) REFERENCES `EvaluationCategory` (`id`) ON DELETE CASCADE;

SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='StudentClass') > 0,
  'ALTER TABLE `EvaluationGrade` ADD CONSTRAINT `fk_grade_student_class` FOREIGN KEY (`student_class_id`) REFERENCES `StudentClass` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Registration') > 0,
  'ALTER TABLE `EvaluationGrade` ADD CONSTRAINT `fk_grade_registration` FOREIGN KEY (`registration_id`) REFERENCES `Registration` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------
-- Table: EvaluationGradeHistory
-- Audit trail for grade changes
-- Removed status fields (old_status, new_status)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationGradeHistory` (
                                                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `grade_id` bigint(20) NOT NULL,
    `old_score` decimal(10,2) DEFAULT NULL,
    `new_score` decimal(10,2) DEFAULT NULL,
    `change_reason` text DEFAULT NULL,
    `changed_by` varchar(500) NOT NULL,
    `changed_at` datetime DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    KEY `idx_grade_id` (`grade_id`),
    KEY `idx_changed_at` (`changed_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraint
ALTER TABLE `EvaluationGradeHistory` ADD CONSTRAINT `fk_grade_history_grade` FOREIGN KEY (`grade_id`) REFERENCES `EvaluationGrade` (`id`) ON DELETE CASCADE;

-- --------------------------------------------------------
-- Table: EvaluationFinalGrade
-- Calculated final grades for students in courses
-- Removed status field (simplified workflow)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationFinalGrade` (
                                                      `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `registration_id` bigint(20) NOT NULL,
    `student_class_id` bigint(20) NOT NULL,
    `calculated_score` decimal(10,2) DEFAULT NULL COMMENT 'Calculated total score',
    `calculated_percentage` decimal(5,2) DEFAULT NULL,
    `letter_grade` varchar(5) DEFAULT NULL,
    `is_passing` tinyint(1) DEFAULT NULL,
    `category_scores` json DEFAULT NULL COMMENT 'JSON with scores per category',
    `calculation_date` datetime DEFAULT NULL,
    `override_score` decimal(10,2) DEFAULT NULL COMMENT 'Manual override if needed',
    `override_reason` text DEFAULT NULL,
    `created_at` datetime DEFAULT current_timestamp(),
    `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    `created_by` varchar(500) DEFAULT NULL,
    `updated_by` varchar(500) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_unique_final_grade` (`registration_id`, `student_class_id`),
    KEY `idx_student_class_id` (`student_class_id`),
    KEY `idx_calculation_date` (`calculation_date`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraints
SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Registration') > 0,
  'ALTER TABLE `EvaluationFinalGrade` ADD CONSTRAINT `fk_final_grade_registration` FOREIGN KEY (`registration_id`) REFERENCES `Registration` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='StudentClass') > 0,
  'ALTER TABLE `EvaluationFinalGrade` ADD CONSTRAINT `fk_final_grade_class` FOREIGN KEY (`student_class_id`) REFERENCES `StudentClass` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------
-- Table: EvaluationGradeRubric
-- Rubric criteria for detailed evaluation
-- Updated to reference Grade directly (no EvaluationItem)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationGradeRubric` (
                                                       `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `evaluation_category_id` bigint(20) DEFAULT NULL,
    `grade_id` bigint(20) DEFAULT NULL,
    `criteria_name` varchar(200) NOT NULL,
    `description` text DEFAULT NULL,
    `max_points` decimal(10,2) NOT NULL,
    `display_order` int(11) DEFAULT 0,
    `created_at` datetime DEFAULT current_timestamp(),
    `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`evaluation_category_id`),
    KEY `idx_grade_id` (`grade_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraints
ALTER TABLE `EvaluationGradeRubric` ADD CONSTRAINT `fk_rubric_category` FOREIGN KEY (`evaluation_category_id`) REFERENCES `EvaluationCategory` (`id`) ON DELETE CASCADE;
ALTER TABLE `EvaluationGradeRubric` ADD CONSTRAINT `fk_rubric_grade` FOREIGN KEY (`grade_id`) REFERENCES `EvaluationGrade` (`id`) ON DELETE CASCADE;

-- --------------------------------------------------------
-- Table: EvaluationGradeRubricScore
-- Individual rubric scores for students
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `EvaluationGradeRubricScore` (
                                                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `grade_id` bigint(20) NOT NULL,
    `rubric_id` bigint(20) NOT NULL,
    `score` decimal(10,2) NOT NULL,
    `comments` text DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_unique_rubric_score` (`grade_id`, `rubric_id`),
    KEY `idx_rubric_id` (`rubric_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Add foreign key constraints
ALTER TABLE `EvaluationGradeRubricScore` ADD CONSTRAINT `fk_rubric_score_grade` FOREIGN KEY (`grade_id`) REFERENCES `EvaluationGrade` (`id`) ON DELETE CASCADE;
ALTER TABLE `EvaluationGradeRubricScore` ADD CONSTRAINT `fk_rubric_score_rubric` FOREIGN KEY (`rubric_id`) REFERENCES `EvaluationGradeRubric` (`id`) ON DELETE CASCADE;

-- --------------------------------------------------------
-- Indexes for Performance
-- --------------------------------------------------------

CREATE INDEX IF NOT EXISTS `idx_grade_category_student` ON `EvaluationGrade` (`category_id`, `student_class_id`, `registration_id`);
CREATE INDEX IF NOT EXISTS `idx_final_grade_calculated` ON `EvaluationFinalGrade` (`calculation_date`);

-- ========================================================================
-- END OF CONSOLIDATED SCHEMA
-- ========================================================================

COMMIT;
