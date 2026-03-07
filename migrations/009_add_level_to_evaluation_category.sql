-- ========================================================================
-- Add level_id to EvaluationCategory
-- ========================================================================
-- Different levels may require different cardinalities and display_order
-- for the same course categories. This migration adds level_id so that
-- EvaluationCategory can be scoped per course+level combination.
-- ========================================================================

ALTER TABLE `EvaluationCategory`
ADD COLUMN `level_id` bigint(20) DEFAULT NULL AFTER `course_id`,
ADD KEY `idx_level_id` (`level_id`);

ALTER TABLE `EvaluationCategory`
ADD UNIQUE KEY `idx_unique_course_level_name` (`course_id`, `level_id`, `name`);

SET @s = (SELECT IF(
  (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='Level') > 0,
  'ALTER TABLE `EvaluationCategory` ADD CONSTRAINT `fk_eval_category_level` FOREIGN KEY (`level_id`) REFERENCES `Level` (`id`) ON DELETE CASCADE',
  'SELECT 1'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
