-- ========================================================================
-- Remove grading_scale and letter_grade columns
-- ========================================================================
-- This migration removes the grading_scale and letter_grade features
-- since the system now focuses on storing individual evaluation marks
-- rather than automatic grade conversion.
-- ========================================================================

-- Remove grading_scale column from EvaluationGradingFormula
ALTER TABLE `EvaluationGradingFormula` 
DROP COLUMN IF EXISTS `grading_scale`;

-- Remove letter_grade column from EvaluationGrade
ALTER TABLE `EvaluationGrade` 
DROP COLUMN IF EXISTS `letter_grade`;

-- Remove letter_grade column from EvaluationFinalGrade
ALTER TABLE `EvaluationFinalGrade` 
DROP COLUMN IF EXISTS `letter_grade`;

COMMIT;
