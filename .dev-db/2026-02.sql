-- Evaluation System Enhancement
-- Migration created: February 2026
-- This migration enhances the evaluation system with configurable categories,
-- formulas, and comprehensive reporting capabilities

-- --------------------------------------------------------
-- Table: EvaluationCategory
-- Defines evaluation categories for each course (e.g., Quizzes, Assignments, Midterm, Final)
-- Replaces the EvaluationLevel approach with more flexible course-specific configuration
-- --------------------------------------------------------

CREATE TABLE `EvaluationCategory` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `course_id` bigint(20) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` text DEFAULT NULL,
  `weight` decimal(5,2) NOT NULL COMMENT 'Weight in percentage (0-100)',
  `max_score` decimal(10,2) NOT NULL DEFAULT 20.00 COMMENT 'Maximum possible score',
  `drop_lowest` int(11) DEFAULT 0 COMMENT 'Number of lowest scores to drop',
  `is_extra_credit` tinyint(1) DEFAULT 0 COMMENT 'Whether this is extra credit',
  `display_order` int(11) DEFAULT 0,
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` varchar(500) DEFAULT NULL,
  `updated_by` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_is_active` (`is_active`),
  CONSTRAINT `fk_eval_category_course` FOREIGN KEY (`course_id`) REFERENCES `Course` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: GradingFormula
-- Stores grading formula configurations for courses
-- --------------------------------------------------------

CREATE TABLE `GradingFormula` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `course_id` bigint(20) NOT NULL,
  `formula_type` varchar(50) NOT NULL COMMENT 'WEIGHTED_AVERAGE, POINTS_BASED, PASS_FAIL, CUSTOM',
  `formula_config` json DEFAULT NULL COMMENT 'JSON configuration for the formula',
  `passing_percentage` decimal(5,2) DEFAULT NULL COMMENT 'Minimum percentage to pass',
  `grading_scale` json DEFAULT NULL COMMENT 'JSON defining grade ranges (A, B, C, etc.)',
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` varchar(500) DEFAULT NULL,
  `updated_by` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  CONSTRAINT `fk_grading_formula_course` FOREIGN KEY (`course_id`) REFERENCES `Course` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: EvaluationItem
-- Individual evaluation instances within a category (e.g., Quiz 1, Quiz 2, Midterm Exam)
-- --------------------------------------------------------

CREATE TABLE `EvaluationItem` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` bigint(20) NOT NULL,
  `student_class_id` bigint(20) NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` text DEFAULT NULL,
  `date` date DEFAULT NULL,
  `due_date` datetime DEFAULT NULL,
  `max_score` decimal(10,2) NOT NULL,
  `weight_override` decimal(5,2) DEFAULT NULL COMMENT 'Override category weight for this item',
  `status` varchar(20) DEFAULT 'DRAFT' COMMENT 'DRAFT, PUBLISHED, GRADED, ARCHIVED',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` varchar(500) DEFAULT NULL,
  `updated_by` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_student_class_id` (`student_class_id`),
  KEY `idx_date` (`date`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_eval_item_category` FOREIGN KEY (`category_id`) REFERENCES `EvaluationCategory` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_eval_item_class` FOREIGN KEY (`student_class_id`) REFERENCES `StudentClass` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: Grade
-- Individual student grades for evaluation items
-- --------------------------------------------------------

CREATE TABLE `Grade` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `evaluation_item_id` bigint(20) NOT NULL,
  `registration_id` bigint(20) NOT NULL,
  `score` decimal(10,2) DEFAULT NULL,
  `percentage` decimal(5,2) DEFAULT NULL COMMENT 'Calculated percentage',
  `letter_grade` varchar(5) DEFAULT NULL COMMENT 'Letter grade (A, B, C, etc.)',
  `is_excused` tinyint(1) DEFAULT 0 COMMENT 'Excused from this evaluation',
  `is_late` tinyint(1) DEFAULT 0 COMMENT 'Submitted late',
  `late_penalty` decimal(5,2) DEFAULT 0.00 COMMENT 'Penalty percentage for late submission',
  `comments` text DEFAULT NULL,
  `graded_by` varchar(500) DEFAULT NULL,
  `graded_at` datetime DEFAULT NULL,
  `status` varchar(20) DEFAULT 'DRAFT' COMMENT 'DRAFT, PUBLISHED',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` varchar(500) DEFAULT NULL,
  `updated_by` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_grade` (`evaluation_item_id`, `registration_id`),
  KEY `idx_registration_id` (`registration_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_grade_eval_item` FOREIGN KEY (`evaluation_item_id`) REFERENCES `EvaluationItem` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_grade_registration` FOREIGN KEY (`registration_id`) REFERENCES `Registration` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: GradeHistory
-- Audit trail for grade changes
-- --------------------------------------------------------

CREATE TABLE `GradeHistory` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `grade_id` bigint(20) NOT NULL,
  `old_score` decimal(10,2) DEFAULT NULL,
  `new_score` decimal(10,2) DEFAULT NULL,
  `old_status` varchar(20) DEFAULT NULL,
  `new_status` varchar(20) DEFAULT NULL,
  `change_reason` text DEFAULT NULL,
  `changed_by` varchar(500) NOT NULL,
  `changed_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_grade_id` (`grade_id`),
  KEY `idx_changed_at` (`changed_at`),
  CONSTRAINT `fk_grade_history_grade` FOREIGN KEY (`grade_id`) REFERENCES `Grade` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: FinalGrade
-- Calculated final grades for students in courses
-- --------------------------------------------------------

CREATE TABLE `FinalGrade` (
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
  `status` varchar(20) DEFAULT 'DRAFT' COMMENT 'DRAFT, PUBLISHED, FINAL',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` varchar(500) DEFAULT NULL,
  `updated_by` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_final_grade` (`registration_id`, `student_class_id`),
  KEY `idx_student_class_id` (`student_class_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_final_grade_registration` FOREIGN KEY (`registration_id`) REFERENCES `Registration` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_final_grade_class` FOREIGN KEY (`student_class_id`) REFERENCES `StudentClass` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: GradeRubric
-- Rubric criteria for detailed evaluation
-- --------------------------------------------------------

CREATE TABLE `GradeRubric` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `evaluation_category_id` bigint(20) DEFAULT NULL,
  `evaluation_item_id` bigint(20) DEFAULT NULL,
  `criteria_name` varchar(200) NOT NULL,
  `description` text DEFAULT NULL,
  `max_points` decimal(10,2) NOT NULL,
  `display_order` int(11) DEFAULT 0,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`evaluation_category_id`),
  KEY `idx_item_id` (`evaluation_item_id`),
  CONSTRAINT `fk_rubric_category` FOREIGN KEY (`evaluation_category_id`) REFERENCES `EvaluationCategory` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_rubric_item` FOREIGN KEY (`evaluation_item_id`) REFERENCES `EvaluationItem` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Table: GradeRubricScore
-- Individual rubric scores for students
-- --------------------------------------------------------

CREATE TABLE `GradeRubricScore` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `grade_id` bigint(20) NOT NULL,
  `rubric_id` bigint(20) NOT NULL,
  `score` decimal(10,2) NOT NULL,
  `comments` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_rubric_score` (`grade_id`, `rubric_id`),
  KEY `idx_rubric_id` (`rubric_id`),
  CONSTRAINT `fk_rubric_score_grade` FOREIGN KEY (`grade_id`) REFERENCES `Grade` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_rubric_score_rubric` FOREIGN KEY (`rubric_id`) REFERENCES `GradeRubric` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------
-- Sample Data: Default Grading Scales
-- --------------------------------------------------------

-- Example: Simple course configuration
-- Course 1: French (weighted average: periodic 60% + exam 40%)
INSERT INTO `EvaluationCategory` (`course_id`, `name`, `description`, `weight`, `max_score`, `drop_lowest`, `display_order`) VALUES
(1, 'Periodic Evaluations', 'Weekly/Monthly evaluations', 60.00, 20.00, 1, 1),
(1, 'Final Exam', 'End of period exam', 40.00, 20.00, 0, 2);

INSERT INTO `GradingFormula` (`course_id`, `formula_type`, `passing_percentage`, `grading_scale`) VALUES
(1, 'WEIGHTED_AVERAGE', 50.00, '{"A": {"min": 90, "max": 100}, "B": {"min": 80, "max": 89}, "C": {"min": 70, "max": 79}, "D": {"min": 60, "max": 69}, "F": {"min": 0, "max": 59}}');

-- Example: Complex course configuration
-- Course 2: English (multiple categories with different weights)
INSERT INTO `EvaluationCategory` (`course_id`, `name`, `description`, `weight`, `max_score`, `drop_lowest`, `display_order`) VALUES
(2, 'Homework', 'Weekly homework assignments', 20.00, 10.00, 2, 1),
(2, 'Quizzes', 'Periodic quizzes', 25.00, 20.00, 1, 2),
(2, 'Projects', 'Class projects', 25.00, 100.00, 0, 3),
(2, 'Final Exam', 'Final examination', 30.00, 100.00, 0, 4);

INSERT INTO `GradingFormula` (`course_id`, `formula_type`, `passing_percentage`, `grading_scale`) VALUES
(2, 'WEIGHTED_AVERAGE', 60.00, '{"A": {"min": 90, "max": 100}, "B": {"min": 80, "max": 89}, "C": {"min": 70, "max": 79}, "D": {"min": 60, "max": 69}, "F": {"min": 0, "max": 59}}');

-- --------------------------------------------------------
-- Indexes for Performance
-- --------------------------------------------------------

CREATE INDEX `idx_grade_published` ON `Grade` (`status`, `graded_at`);
CREATE INDEX `idx_final_grade_calculated` ON `FinalGrade` (`calculation_date`, `status`);
CREATE INDEX `idx_eval_item_class_date` ON `EvaluationItem` (`student_class_id`, `date`);

-- --------------------------------------------------------
-- Views for Reporting
-- --------------------------------------------------------

-- View: Student Grade Summary
CREATE OR REPLACE VIEW `vw_StudentGradeSummary` AS
SELECT 
    g.id AS grade_id,
    r.student_id,
    r.student_class_id,
    sc.course_id,
    c.name AS course_name,
    ec.name AS category_name,
    ei.name AS evaluation_name,
    ei.date AS evaluation_date,
    g.score,
    ei.max_score,
    g.percentage,
    g.letter_grade,
    g.status,
    g.graded_at,
    g.comments
FROM Grade g
INNER JOIN EvaluationItem ei ON g.evaluation_item_id = ei.id
INNER JOIN EvaluationCategory ec ON ei.category_id = ec.id
INNER JOIN Registration r ON g.registration_id = r.id
INNER JOIN StudentClass sc ON r.student_class_id = sc.id
INNER JOIN Course c ON sc.course_id = c.id
WHERE g.status = 'PUBLISHED';

-- View: Class Performance Statistics
CREATE OR REPLACE VIEW `vw_ClassPerformanceStats` AS
SELECT 
    sc.id AS student_class_id,
    sc.name AS class_name,
    c.id AS course_id,
    c.name AS course_name,
    ec.id AS category_id,
    ec.name AS category_name,
    ei.id AS evaluation_item_id,
    ei.name AS evaluation_name,
    COUNT(DISTINCT g.id) AS total_graded,
    AVG(g.percentage) AS average_percentage,
    MIN(g.percentage) AS min_percentage,
    MAX(g.percentage) AS max_percentage,
    STDDEV(g.percentage) AS std_deviation,
    SUM(CASE WHEN g.percentage >= 60 THEN 1 ELSE 0 END) AS passing_count,
    SUM(CASE WHEN g.percentage < 60 THEN 1 ELSE 0 END) AS failing_count
FROM Grade g
INNER JOIN EvaluationItem ei ON g.evaluation_item_id = ei.id
INNER JOIN EvaluationCategory ec ON ei.category_id = ec.id
INNER JOIN StudentClass sc ON ei.student_class_id = sc.id
INNER JOIN Course c ON sc.course_id = c.id
WHERE g.status = 'PUBLISHED' AND g.is_excused = 0
GROUP BY sc.id, c.id, ec.id, ei.id;

-- View: Final Grade Report
CREATE OR REPLACE VIEW `vw_FinalGradeReport` AS
SELECT 
    fg.id AS final_grade_id,
    r.student_id,
    s.firstName,
    s.lastName,
    r.student_class_id,
    sc.name AS class_name,
    c.id AS course_id,
    c.name AS course_name,
    c.teacher_email,
    fg.calculated_percentage,
    fg.letter_grade,
    fg.is_passing,
    fg.status,
    fg.calculation_date,
    p.start AS period_start,
    p.end AS period_end
FROM FinalGrade fg
INNER JOIN Registration r ON fg.registration_id = r.id
INNER JOIN Student s ON r.student_id = s.id
INNER JOIN StudentClass sc ON fg.student_class_id = sc.id
INNER JOIN Course c ON sc.course_id = c.id
INNER JOIN Period p ON sc.period_id = p.id
WHERE r.status = 'ACTIVE';

COMMIT;
