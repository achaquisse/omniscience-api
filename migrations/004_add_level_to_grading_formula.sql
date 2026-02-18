-- ========================================================================
-- Add level_id to GradingFormula table
-- ========================================================================
-- This migration adds level_id to GradingFormula to support multiple
-- grading formulas per course for different levels.
-- ========================================================================

-- Add level_id column to EvaluationGradingFormula
ALTER TABLE `EvaluationGradingFormula` 
ADD COLUMN `level_id` bigint(20) DEFAULT NULL AFTER `course_id`,
ADD KEY `idx_level_id` (`level_id`);

-- Add unique constraint to prevent duplicate formulas for same course+level combination
ALTER TABLE `EvaluationGradingFormula` 
ADD UNIQUE KEY `idx_unique_course_level` (`course_id`, `level_id`, `is_active`);

-- Add foreign key constraint if Level table exists
SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Level') > 0,
  'ALTER TABLE `EvaluationGradingFormula` ADD CONSTRAINT `fk_grading_formula_level` FOREIGN KEY (`level_id`) REFERENCES `Level` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
