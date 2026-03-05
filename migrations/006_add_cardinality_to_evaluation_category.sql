-- ========================================================================
-- Add Cardinality to EvaluationCategory
-- ========================================================================
-- This migration adds the cardinality field to EvaluationCategory table.
-- Cardinality represents the number of evaluations expected for each
-- category. When calculating final grades, missing evaluations will be
-- counted as 0 unless the student is marked as excused (is_excused=true).
-- ========================================================================

ALTER TABLE `EvaluationCategory` 
ADD COLUMN `cardinality` INT NOT NULL DEFAULT 1 AFTER `max_score`;

COMMIT;
