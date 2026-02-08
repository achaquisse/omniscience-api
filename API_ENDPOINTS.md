# Evaluation System API Endpoints

## Overview
This document describes the REST API endpoints for the evaluation/grading system implemented in the Omniscience API.

## Base URL
All endpoints require authentication via the `AuthMiddleware`.

---

## Evaluation Categories

### List Evaluation Categories
**GET** `/evaluation-categories?course_id={course_id}&include_inactive={true|false}`

Lists all evaluation categories for a course.

**Query Parameters:**
- `course_id` (required): ID of the course
- `include_inactive` (optional): Include inactive categories (default: false)

**Response:** Array of EvaluationCategory objects

### Get Evaluation Category
**GET** `/evaluation-categories/:id`

Get a specific evaluation category by ID.

**Response:** EvaluationCategory object

### Create Evaluation Category
**POST** `/evaluation-categories`

Create a new evaluation category.

**Request Body:**
```json
{
  "course_id": 1,
  "name": "Quizzes",
  "description": "Weekly quizzes",
  "weight": 25.00,
  "max_score": 20.00,
  "drop_lowest": 1,
  "is_extra_credit": false,
  "display_order": 1
}
```

**Response:** Created EvaluationCategory object

### Update Evaluation Category
**PUT** `/evaluation-categories/:id`

Update an existing evaluation category.

**Request Body:** Partial or full EvaluationCategory object

**Response:** Updated EvaluationCategory object

### Delete Evaluation Category
**DELETE** `/evaluation-categories/:id`

Delete an evaluation category (cascades to evaluation items and grades).

**Response:** 204 No Content

---

## Grading Formulas

### Get Grading Formula
**GET** `/grading-formulas?course_id={course_id}`

Get the active grading formula for a course.

**Query Parameters:**
- `course_id` (required): ID of the course

**Response:** GradingFormula object

### Create Grading Formula
**POST** `/grading-formulas`

Create a new grading formula for a course.

**Request Body:**
```json
{
  "course_id": 1,
  "formula_type": "WEIGHTED_AVERAGE",
  "passing_percentage": 60.00,
  "grading_scale": {
    "A": {"min": 90, "max": 100},
    "B": {"min": 80, "max": 89},
    "C": {"min": 70, "max": 79},
    "D": {"min": 60, "max": 69},
    "F": {"min": 0, "max": 59}
  }
}
```

**Supported Formula Types:**
- `WEIGHTED_AVERAGE`: Calculate final grade using weighted average of category scores
- `POINTS_BASED`: Sum total points earned
- `PASS_FAIL`: Simple pass/fail based on threshold
- `CUSTOM`: Custom formula defined in formula_config

**Response:** Created GradingFormula object

### Update Grading Formula
**PUT** `/grading-formulas/:id`

Update an existing grading formula.

**Request Body:** Partial or full GradingFormula object

**Response:** Updated GradingFormula object

---

## Evaluation Items

### List Evaluation Items
**GET** `/evaluation-items?category_id={category_id}&student_class_id={student_class_id}`

List evaluation items (e.g., Quiz 1, Midterm Exam).

**Query Parameters:**
- `category_id` (optional): Filter by category
- `student_class_id` (optional): Filter by class

**Response:** Array of EvaluationItem objects

### Get Evaluation Item
**GET** `/evaluation-items/:id`

Get a specific evaluation item.

**Response:** EvaluationItem object

### Create Evaluation Item
**POST** `/evaluation-items`

Create a new evaluation item.

**Request Body:**
```json
{
  "category_id": 1,
  "student_class_id": 10,
  "name": "Quiz 1",
  "description": "Chapter 1-3",
  "date": "2026-02-15",
  "due_date": "2026-02-15T23:59:59Z",
  "max_score": 20.00,
  "status": "DRAFT"
}
```

**Response:** Created EvaluationItem object

### Update Evaluation Item
**PUT** `/evaluation-items/:id`

Update an evaluation item.

**Request Body:** Partial or full EvaluationItem object

**Response:** Updated EvaluationItem object

### Delete Evaluation Item
**DELETE** `/evaluation-items/:id`

Delete an evaluation item (cascades to grades).

**Response:** 204 No Content

---

## Grades

### List Grades
**GET** `/grades?evaluation_item_id={id}&registration_id={id}&student_class_id={id}`

List grades with optional filters.

**Query Parameters:**
- `evaluation_item_id` (optional): Filter by evaluation item
- `registration_id` (optional): Filter by student registration
- `student_class_id` (optional): Filter by class

**Response:** Array of Grade objects

### Get Grade
**GET** `/grades/:id`

Get a specific grade.

**Response:** Grade object

### Create Grade
**POST** `/grades`

Create a single grade.

**Request Body:**
```json
{
  "evaluation_item_id": 1,
  "registration_id": 100,
  "score": 18.5,
  "is_late": false,
  "comments": "Excellent work",
  "graded_by": "teacher@example.com",
  "status": "DRAFT"
}
```

**Note:** Percentage is automatically calculated based on score and max_score.

**Response:** Created Grade object

### Update Grade
**PUT** `/grades/:id`

Update a grade. Creates audit history record.

**Request Body:** Partial or full Grade object

**Response:** Updated Grade object

### Batch Create Grades
**POST** `/grades/batch`

Create multiple grades at once for an evaluation item.

**Request Body:**
```json
{
  "evaluation_item_id": 1,
  "grades": [
    {
      "registration_id": 100,
      "score": 18.5,
      "comments": "Good work"
    },
    {
      "registration_id": 101,
      "score": 15.0
    }
  ]
}
```

**Response:** Summary with created grades

### Publish Grades
**POST** `/grades/publish/:evaluation_item_id`

Publish all draft grades for an evaluation item.

**Response:** 
```json
{
  "message": "grades published successfully",
  "published": 25
}
```

---

## Reports & Final Grades

### Get Student Report
**GET** `/reports/student/:registration_id`

Get detailed grade report for a student including all published grades grouped by category.

**Response:**
```json
{
  "registration_id": 100,
  "student_class_id": 10,
  "category_scores": {
    "Quizzes": {
      "category_id": 1,
      "category_name": "Quizzes",
      "weight": 25.00,
      "grades": [...]
    }
  },
  "total_grades": 15,
  "final_grade": {...}
}
```

### Get Class Report
**GET** `/reports/class/:student_class_id`

Get performance statistics and analytics for an entire class.

**Response:**
```json
{
  "student_class_id": 10,
  "class_name": "French A1 - 2026",
  "course_name": "French",
  "evaluation_stats": [
    {
      "evaluation_item_id": 1,
      "evaluation_name": "Quiz 1",
      "category_name": "Quizzes",
      "total_graded": 25,
      "average_percentage": 82.5,
      "min_percentage": 45.0,
      "max_percentage": 100.0,
      "passing_count": 22,
      "failing_count": 3
    }
  ],
  "total_students": 25,
  "passing_students": 20,
  "average_final_grade": 78.3,
  "final_grades": [...]
}
```

### Get Final Grade
**GET** `/final-grades/:registration_id`

Get calculated final grade for a student.

**Response:** FinalGrade object with calculated scores and letter grade

### Calculate Final Grades
**POST** `/final-grades/calculate/:student_class_id`

Trigger recalculation of final grades for all students in a class.

**Response:**
```json
{
  "message": "final grades calculated successfully",
  "student_class_id": 10
}
```

---

## Data Models

### EvaluationCategory
```json
{
  "id": 1,
  "course_id": 1,
  "name": "Quizzes",
  "description": "Weekly quizzes",
  "weight": 25.00,
  "max_score": 20.00,
  "drop_lowest": 1,
  "is_extra_credit": false,
  "display_order": 1,
  "is_active": true,
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

### GradingFormula
```json
{
  "id": 1,
  "course_id": 1,
  "formula_type": "WEIGHTED_AVERAGE",
  "passing_percentage": 60.00,
  "grading_scale": {
    "A": {"min": 90, "max": 100},
    "B": {"min": 80, "max": 89}
  },
  "is_active": true
}
```

### EvaluationItem
```json
{
  "id": 1,
  "category_id": 1,
  "student_class_id": 10,
  "name": "Quiz 1",
  "date": "2026-02-15",
  "max_score": 20.00,
  "status": "PUBLISHED"
}
```

### Grade
```json
{
  "id": 1,
  "evaluation_item_id": 1,
  "registration_id": 100,
  "score": 18.5,
  "percentage": 92.50,
  "letter_grade": "A",
  "is_excused": false,
  "is_late": false,
  "status": "PUBLISHED",
  "graded_at": "2026-02-15T14:30:00Z"
}
```

### FinalGrade
```json
{
  "id": 1,
  "registration_id": 100,
  "student_class_id": 10,
  "calculated_percentage": 85.50,
  "letter_grade": "B",
  "is_passing": true,
  "category_scores": {
    "Quizzes": {
      "category_id": 1,
      "score": 88.0,
      "weight": 25.0
    }
  },
  "status": "PUBLISHED",
  "calculation_date": "2026-02-20T10:00:00Z"
}
```

---

## Features

### Automatic Calculations
- **Percentage**: Automatically calculated from score and max_score
- **Letter Grade**: Automatically assigned based on grading scale
- **Final Grade**: Calculated using weighted average with support for drop lowest scores
- **Category Scores**: Aggregated per category with weight application

### Grade Status Workflow
1. **DRAFT**: Grade is saved but not visible to students
2. **PUBLISHED**: Grade is finalized and visible

### Audit Trail
- All grade changes are logged in `GradeHistory` table
- Tracks old/new values and who made the change

### Drop Lowest Scores
- Configure `drop_lowest` on evaluation categories
- Automatically excludes lowest N scores from category average

### Extra Credit
- Mark categories as extra credit
- Adds to total without counting in denominator

### Late Penalties
- Track late submissions with `is_late` flag
- Apply percentage penalty with `late_penalty` field

---

## Database Schema

The schema is defined in `.dev-db/2026-02.sql` and includes:

### Core Tables
- `EvaluationCategory`: Configurable evaluation categories per course
- `GradingFormula`: Grading formulas and scales per course
- `EvaluationItem`: Individual evaluation instances
- `Grade`: Student grades
- `GradeHistory`: Audit trail
- `FinalGrade`: Calculated final grades
- `GradeRubric`: Detailed rubric criteria
- `GradeRubricScore`: Rubric scores per student

### Database Views
- `vw_StudentGradeSummary`: Published grades with course/category info
- `vw_ClassPerformanceStats`: Aggregated class statistics
- `vw_FinalGradeReport`: Final grades with student details
