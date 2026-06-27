package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"gorm.io/gorm"
)

type AcademicRepository struct {
	db *gorm.DB
}

func NewAcademicRepository(db *gorm.DB) *AcademicRepository {
	return &AcademicRepository{db: db}
}

// Student operations
func (r *AcademicRepository) CreateStudent(s *domain.Student) error {
	return r.db.Create(s).Error
}

func (r *AcademicRepository) GetStudentByID(id string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// UpdateStudentAdvisor assigns/removes academic advisor (PA) for a student
func (r *AcademicRepository) UpdateStudentAdvisor(studentID string, advisorID *string) error {
	updates := map[string]interface{}{
		"advisor_id": advisorID,
		"updated_at": time.Now(),
	}
	return r.db.Model(&domain.Student{}).Where("id = ?", studentID).Updates(updates).Error
}

// GetStudentsByAdvisor retrieves students assigned to a specific PA
func (r *AcademicRepository) GetStudentsByAdvisor(advisorID string, page, limit int) ([]domain.Student, int64, error) {
	var students []domain.Student
	var total int64

	query := r.db.Model(&domain.Student{}).Where("advisor_id = ?", advisorID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("nim asc").Limit(limit).Offset(offset).Find(&students).Error
	return students, total, err
}

func (r *AcademicRepository) GetStudentByApplicantID(appID string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.Where("applicant_id = ?", appID).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *AcademicRepository) ListStudents(studyProgramID string, page, limit int) ([]domain.Student, int64, error) {
	var students []domain.Student
	var total int64

	query := r.db.Model(&domain.Student{})
	if studyProgramID != "" {
		query = query.Where("study_program_id = ?", studyProgramID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("nim asc").Limit(limit).Offset(offset).Find(&students).Error
	return students, total, err
}

// NIM Sequence Generation (with lock)
func (r *AcademicRepository) GenerateNIM(tx *gorm.DB, studyProgramID string, periodID string, year string) (string, error) {
	var seq domain.NimSequence

	// Try finding the sequence
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("study_program_id = ? AND entry_period_id = ? AND sequence_year = ?", studyProgramID, periodID, year).
		First(&seq).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create first sequence record
			seq = domain.NimSequence{
				StudyProgramID: studyProgramID,
				EntryPeriodID:  periodID,
				SequenceYear:   year,
				LastNumber:     1,
			}
			if err := tx.Create(&seq).Error; err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		// Increment
		seq.LastNumber++
		if err := tx.Model(&seq).Update("last_number", seq.LastNumber).Error; err != nil {
			return "", err
		}
	}

	// Format program code (take first 4 chars of UUID or mapping, let's keep it simple)
	// We'll generate a NIM like: [Year_Suffix][Prog_Suffix][Sequence_Formatted]
	// E.g., 202611020005
	progPrefix := studyProgramID
	if len(progPrefix) > 4 {
		progPrefix = progPrefix[:4]
	}
	nim := fmt.Sprintf("%s%s%04d", year, progPrefix, seq.LastNumber)
	return nim, nil
}

// Class operations
func (r *AcademicRepository) CreateClass(c *domain.Class) error {
	return r.db.Create(c).Error
}

func (r *AcademicRepository) GetClassByID(id string) (*domain.Class, error) {
	var c domain.Class
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// KRS operations
func (r *AcademicRepository) CreateKRS(krs *domain.KRS) error {
	return r.db.Create(krs).Error
}

func (r *AcademicRepository) GetKRSByID(id string) (*domain.KRS, error) {
	var krs domain.KRS
	err := r.db.Preload("Items").Where("id = ?", id).First(&krs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &krs, nil
}

func (r *AcademicRepository) UpdateKRSStatus(id string, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "submitted" {
		updates["submitted_at"] = &now
	} else if status == "approved" {
		updates["approved_at"] = &now
	}
	return r.db.Model(&domain.KRS{}).Where("id = ?", id).Updates(updates).Error
}

// KrsItem operations
func (r *AcademicRepository) GetKrsItemByID(id string) (*domain.KrsItem, error) {
	var item domain.KrsItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *AcademicRepository) UpdateKrsItemStatus(id string, status string) error {
	return r.db.Model(&domain.KrsItem{}).Where("id = ?", id).Update("status", status).Error
}

// Grade operations
func (r *AcademicRepository) CreateGrade(gr *domain.Grade) error {
	return r.db.Create(gr).Error
}

func (r *AcademicRepository) GetGradeByKrsItem(krsItemID string) (*domain.Grade, error) {
	var gr domain.Grade
	err := r.db.Where("krs_item_id = ?", krsItemID).First(&gr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gr, nil
}

func (r *AcademicRepository) UpdateGrade(id string, numGrade float64, letGrade string, ptGrade float64, source string, editor string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"numeric_grade": &numGrade,
		"letter_grade":  letGrade,
		"grade_point":   &ptGrade,
		"source":        source,
		"submitted_at":  &now,
		"submitted_by":  &editor,
	}
	return r.db.Model(&domain.Grade{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AcademicRepository) CreateCurriculum(c *domain.Curriculum) error {
	return r.db.Create(c).Error
}

func (r *AcademicRepository) GetCurriculumByID(id string) (*domain.Curriculum, error) {
	var c domain.Curriculum
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *AcademicRepository) CreateCourse(co *domain.Course) error {
	return r.db.Create(co).Error
}

func (r *AcademicRepository) GetCourseByID(id string) (*domain.Course, error) {
	var c domain.Course
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *AcademicRepository) CreateCurriculumCourse(cc *domain.CurriculumCourse) error {
	return r.db.Create(cc).Error
}

func (r *AcademicRepository) CreateCourseOffering(co *domain.CourseOffering) error {
	return r.db.Create(co).Error
}

func (r *AcademicRepository) GetCourseOfferingByID(id string) (*domain.CourseOffering, error) {
	var co domain.CourseOffering
	err := r.db.Where("id = ?", id).First(&co).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &co, nil
}

func (r *AcademicRepository) CreateClassLecturer(cl *domain.ClassLecturer) error {
	return r.db.Create(cl).Error
}

func (r *AcademicRepository) GetGradeByID(id string) (*domain.Grade, error) {
	var g domain.Grade
	err := r.db.Where("id = ?", id).First(&g).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *AcademicRepository) CreateGradeHistory(gh *domain.GradeHistory) error {
	return r.db.Create(gh).Error
}

// Student grades retrieval for GPA / transcript / KHS
func (r *AcademicRepository) GetStudentGrades(studentID string, academicPeriodID string) ([]domain.Grade, error) {
	var grades []domain.Grade
	query := r.db.Joins("JOIN krs_items ON krs_items.id = grades.krs_item_id").
		Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Where("krs.student_id = ? AND krs.status = 'approved'", studentID)
	
	if academicPeriodID != "" {
		query = query.Where("krs.academic_period_id = ?", academicPeriodID)
	}

	err := query.Find(&grades).Error
	return grades, err
}

func (r *AcademicRepository) GetSksByKrsItemID(krsItemID string) (int, error) {
	var sks int
	err := r.db.Table("krs_items").
		Select("courses.sks").
		Joins("JOIN classes ON classes.id = krs_items.class_id").
		Joins("JOIN course_offerings ON course_offerings.id = classes.course_offering_id").
		Joins("JOIN courses ON courses.id = course_offerings.course_id").
		Where("krs_items.id = ?", krsItemID).
		Row().Scan(&sks)
	return sks, err
}

// Grade Conversion operations
func (r *AcademicRepository) GetGradeConversionByLetter(letter string) (*domain.GradeConversion, error) {
	var gc domain.GradeConversion
	err := r.db.Where("grade_letter = ?", letter).First(&gc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gc, nil
}

func (r *AcademicRepository) CreateGradeConversion(gc *domain.GradeConversion) error {
	return r.db.Create(gc).Error
}

func (r *AcademicRepository) UpdateGradeConversion(gc *domain.GradeConversion) error {
	return r.db.Model(&domain.GradeConversion{}).Where("id = ?", gc.ID).Updates(map[string]interface{}{
		"min_score": gc.MinScore,
		"max_score": gc.MaxScore,
		"grade_points": gc.GradePoints,
	}).Error
}

func (r *AcademicRepository) GetAllGradeConversions() ([]domain.GradeConversion, error) {
	var gcs []domain.GradeConversion
	err := r.db.Order("min_score DESC").Find(&gcs).Error
	return gcs, err
}

// Grade Entry operations
func (r *AcademicRepository) CreateGradeEntry(ge *domain.GradeEntry) error {
	return r.db.Create(ge).Error
}

func (r *AcademicRepository) UpdateGradeEntry(id string, numGrade *float64, letGrade *string, ptGrade *float64) error {
	updates := map[string]interface{}{
		"numeric_grade": numGrade,
		"letter_grade": letGrade,
		"grade_point": ptGrade,
	}
	return r.db.Model(&domain.GradeEntry{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AcademicRepository) UpdateGradeStatus(id string, status string) error {
	return r.db.Model(&domain.Grade{}).Where("id = ?", id).Update("status", status).Error
}

// Transcript generation
func (r *AcademicRepository) GenerateTranscript(studentID string, academicPeriodID string) (map[string]interface{}, error) {
	grades, err := r.GetStudentGrades(studentID, academicPeriodID)
	if err != nil {
		return nil, err
	}

	var totalSks int
	var totalPoints float64
	var courses []map[string]interface{}

	for _, g := range grades {
		sks, _ := r.GetSksByKrsItemID(g.KrsItemID)
		gradePoint := 0.0
		if g.GradePoint != nil {
			gradePoint = *g.GradePoint
		}
		totalPoints += gradePoint * float64(sks)
		totalSks += sks

		courses = append(courses, map[string]interface{}{
			"krs_item_id":  g.KrsItemID,
			"numeric":     g.NumericGrade,
			"letter":      g.LetterGrade,
			"grade_point": gradePoint,
			"sks":        sks,
		})
	}

	ipk := 0.0
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	result := map[string]interface{}{
		"student_id":  studentID,
		"courses":   courses,
		"total_sks": totalSks,
		"ipk":      ipk,
	}

	return result, nil
}

// IPK calculation
func (r *AcademicRepository) CalculateIPK(studentID string) (map[string]interface{}, error) {
	grades, err := r.GetStudentGrades(studentID, "")
	if err != nil {
		return nil, err
	}

	var totalSks int
	var totalPoints float64

	for _, g := range grades {
		sks, _ := r.GetSksByKrsItemID(g.KrsItemID)
		gradePoint := 0.0
		if g.GradePoint != nil {
			gradePoint = *g.GradePoint
		}
		totalPoints += gradePoint * float64(sks)
		totalSks += sks
	}

	ipk := 0.0
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	return map[string]interface{}{
		"student_id": studentID,
		"ipk":     ipk,
		"total_sks": totalSks,
	}, nil
}

// IPS calculation
func (r *AcademicRepository) CalculateIPS(studentID string, academicPeriodID string) (map[string]interface{}, error) {
	grades, err := r.GetStudentGrades(studentID, academicPeriodID)
	if err != nil {
		return nil, err
	}

	var totalSks int
	var totalPoints float64

	for _, g := range grades {
		sks, _ := r.GetSksByKrsItemID(g.KrsItemID)
		gradePoint := 0.0
		if g.GradePoint != nil {
			gradePoint = *g.GradePoint
		}
		totalPoints += gradePoint * float64(sks)
		totalSks += sks
	}

	ips := 0.0
	if totalSks > 0 {
		ips = totalPoints / float64(totalSks)
	}

	return map[string]interface{}{
		"student_id":         studentID,
		"academic_period_id": academicPeriodID,
		"ips":               ips,
		"total_sks":         totalSks,
	}, nil
}

// KHS operations
func (r *AcademicRepository) CreateKHS(khs *domain.KHS) error {
	return r.db.Create(khs).Error
}

func (r *AcademicRepository) GetKHSByStudentAndPeriod(studentID string, academicPeriodID string) (*domain.KHS, error) {
	var khs domain.KHS
	err := r.db.Where("student_id = ? AND academic_period_id = ?", studentID, academicPeriodID).First(&khs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &khs, nil
}

// Grade Component operations
func (r *AcademicRepository) CreateGradeComponent(gc *domain.GradeComponent) error {
	return r.db.Create(gc).Error
}

func (r *AcademicRepository) GetGradeComponentsByGradeID(gradeID string) ([]domain.GradeComponent, error) {
	var components []domain.GradeComponent
	err := r.db.Where("grade_id = ?", gradeID).Order("weight DESC").Find(&components).Error
	return components, err
}

func (r *AcademicRepository) UpdateGradeComponent(id string, weight float64, isActive bool) error {
	return r.db.Model(&domain.GradeComponent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"weight": weight,
		"is_active": isActive,
	}).Error
}

// Grade Component Score operations
func (r *AcademicRepository) CreateGradeComponentScore(gcs *domain.GradeComponentScore) error {
	return r.db.Create(gcs).Error
}

func (r *AcademicRepository) GetGradeComponentScores(gradeEntryID string) ([]domain.GradeComponentScore, error) {
	var scores []domain.GradeComponentScore
	err := r.db.Where("grade_entry_id = ?", gradeEntryID).Find(&scores).Error
	return scores, err
}

func (r *AcademicRepository) UpdateGradeComponentScore(id string, score float64) error {
	return r.db.Model(&domain.GradeComponentScore{}).Where("id = ?", id).Update("score", score).Error
}

// Grade History operations (track grade changes)
func (r *AcademicRepository) CreateGradeWithHistory(gr *domain.Grade, editor string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(gr).Error; err != nil {
			return err
		}
		// Log history - serialize to JSON
		oldVal, _ := json.Marshal(nil)
		newVal, _ := json.Marshal(gr.NumericGrade)
		gh := domain.GradeHistory{
			GradeID:   gr.ID,
			OldValue:  string(oldVal),
			NewValue:  string(newVal),
			ChangedBy: &editor,
			Reason:   "initial_grade",
			ChangedAt: time.Now(),
		}
		return tx.Create(&gh).Error
	})
}

func (r *AcademicRepository) UpdateGradeWithHistory(id string, numGrade float64, letGrade string, ptGrade float64, editor string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get old grade
		var oldGrade domain.Grade
		if err := tx.Where("id = ?", id).First(&oldGrade).Error; err != nil {
			return err
		}
		// Serialize old values
		oldVal, _ := json.Marshal(oldGrade.NumericGrade)
		newVal, _ := json.Marshal(numGrade)
		// Update
		now := time.Now()
		updates := map[string]interface{}{
			"numeric_grade": &numGrade,
			"letter_grade":  letGrade,
			"grade_point":   &ptGrade,
			"updated_at":  now,
		}
		if err := tx.Model(&domain.Grade{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		// Log history
		gh := domain.GradeHistory{
			GradeID:   id,
			OldValue:  string(oldVal),
			NewValue:  string(newVal),
			ChangedBy: &editor,
			Reason:   "grade_update",
			ChangedAt: now,
		}
		return tx.Create(&gh).Error
	})
}

// Graduation Record operations
func (r *AcademicRepository) CreateGraduationRecord(gr *domain.GraduationRecord) error {
	return r.db.Create(gr).Error
}

func (r *AcademicRepository) GetGraduationRecordByStudentID(studentID string) (*domain.GraduationRecord, error) {
	var gr domain.GraduationRecord
	err := r.db.Where("student_id = ?", studentID).First(&gr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gr, nil
}

func (r *AcademicRepository) UpdateGraduationRecord(id string, status string, certNumber string, note string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if certNumber != "" {
		updates["certificate_number"] = certNumber
	}
	if note != "" {
		updates["note"] = note
	}
	if status == "approved" {
		now := time.Now()
		updates["graduated_at"] = now
	}
	return r.db.Model(&domain.GraduationRecord{}).Where("id = ?", id).Updates(updates).Error
}

// Grade Entry With GradeID, StudentID operations
func (r *AcademicRepository) GetGradeEntryByID(id string) (*domain.GradeEntry, error) {
	var ge domain.GradeEntry
	err := r.db.Where("id = ?", id).First(&ge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ge, nil
}

func (r *AcademicRepository) GetGradeEntriesByGradeID(gradeID string) ([]domain.GradeEntry, error) {
	var entries []domain.GradeEntry
	err := r.db.Where("grade_id = ?", gradeID).Find(&entries).Error
	return entries, err
}

func (r *AcademicRepository) GetGradeEntriesByStudentID(studentID string) ([]domain.GradeEntry, error) {
	var entries []domain.GradeEntry
	err := r.db.Where("student_id = ?", studentID).Find(&entries).Error
	return entries, err
}

// Grade Appeal operations
func (r *AcademicRepository) CreateGradeAppeal(ga *domain.GradeAppeal) error {
	return r.db.Create(ga).Error
}

func (r *AcademicRepository) GetGradeAppealByID(id string) (*domain.GradeAppeal, error) {
	var ga domain.GradeAppeal
	err := r.db.Where("id = ?", id).First(&ga).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ga, nil
}

func (r *AcademicRepository) GetGradeAppealsByStudentID(studentID string) ([]domain.GradeAppeal, error) {
	var appeals []domain.GradeAppeal
	err := r.db.Where("student_id = ?", studentID).Order("created_at DESC").Find(&appeals).Error
	return appeals, err
}

func (r *AcademicRepository) UpdateGradeAppeal(id string, status string, reviewerID string, reviewNote string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
		"reviewed_at": now,
	}
	if reviewerID != "" {
		updates["reviewed_by"] = reviewerID
	}
	if reviewNote != "" {
		updates["review_note"] = reviewNote
	}
	return r.db.Model(&domain.GradeAppeal{}).Where("id = ?", id).Updates(updates).Error
}

// Academic Period operations
func (r *AcademicRepository) CreateAcademicPeriod(ap *domain.AcademicPeriod) error {
	return r.db.Create(ap).Error
}

func (r *AcademicRepository) GetAcademicPeriodByID(id string) (*domain.AcademicPeriod, error) {
	var ap domain.AcademicPeriod
	err := r.db.Where("id = ?", id).First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ap, nil
}

func (r *AcademicRepository) GetActiveAcademicPeriod() (*domain.AcademicPeriod, error) {
	var ap domain.AcademicPeriod
	err := r.db.Where("is_active = ?", true).First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ap, nil
}

func (r *AcademicRepository) ListAcademicPeriods(academicYearID string, page, limit int) ([]domain.AcademicPeriod, int64, error) {
	var periods []domain.AcademicPeriod
	var total int64

	query := r.db.Model(&domain.AcademicPeriod{})
	if academicYearID != "" {
		query = query.Where("academic_year_id = ?", academicYearID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("start_date DESC").Limit(limit).Offset(offset).Find(&periods).Error
	return periods, total, err
}

// Study Program operations
func (r *AcademicRepository) CreateStudyProgram(sp *domain.StudyProgram) error {
	return r.db.Create(sp).Error
}

func (r *AcademicRepository) GetStudyProgramByID(id string) (*domain.StudyProgram, error) {
	var sp domain.StudyProgram
	err := r.db.Where("id = ?", id).First(&sp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sp, nil
}

func (r *AcademicRepository) ListStudyPrograms(page, limit int) ([]domain.StudyProgram, int64, error) {
	var programs []domain.StudyProgram
	var total int64

	err := r.db.Model(&domain.StudyProgram{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.Where("is_active = ?", true).Order("code ASC").Limit(limit).Offset(offset).Find(&programs).Error
	return programs, total, err
}

// StudentAdvisor operations (PA Assignment)
func (r *AcademicRepository) CreateStudentAdvisor(sa *domain.StudentAdvisor) error {
	return r.db.Create(sa).Error
}

func (r *AcademicRepository) GetStudentAdvisor(studentID, academicYear string, semester int) (*domain.StudentAdvisor, error) {
	var sa domain.StudentAdvisor
	err := r.db.Where("student_id = ? AND academic_year = ? AND semester = ? AND is_active = ?", studentID, academicYear, semester, true).First(&sa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sa, nil
}

func (r *AcademicRepository) GetStudentAdvisorsByAdvisor(advisorID, academicYear string, semester int) ([]domain.StudentAdvisor, error) {
	var sas []domain.StudentAdvisor
	query := r.db.Where("advisor_id = ? AND is_active = ?", advisorID, true)
	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	if semester > 0 {
		query = query.Where("semester = ?", semester)
	}
	err := query.Order("assigned_at DESC").Find(&sas).Error
	return sas, err
}

func (r *AcademicRepository) DeactivateStudentAdvisor(studentID string) error {
	return r.db.Model(&domain.StudentAdvisor{}).Where("student_id = ? AND is_active = ?", studentID, true).Update("is_active", false).Error
}

// Lecturer operations
func (r *AcademicRepository) GetLecturerByID(id string) (*domain.Lecturer, error) {
	var l domain.Lecturer
	err := r.db.Where("id = ?", id).First(&l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (r *AcademicRepository) GetLecturerByNik(nik string) (*domain.Lecturer, error) {
	var l domain.Lecturer
	err := r.db.Where("nik = ?", nik).First(&l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (r *AcademicRepository) ListLecturers(studyProgramID string, page, limit int) ([]domain.Lecturer, int64, error) {
	var lecturers []domain.Lecturer
	var total int64

	query := r.db.Model(&domain.Lecturer{})
	if studyProgramID != "" {
		query = query.Where("study_program_id = ?", studyProgramID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Where("is_active = ?", true).Order("nik ASC").Limit(limit).Offset(offset).Find(&lecturers).Error
	return lecturers, total, err
}

// Grade retrieval for transcript
func (r *AcademicRepository) GetGradesByStudent(studentID string) ([]domain.Grade, error) {
	var grades []domain.Grade
	err := r.db.Joins("JOIN krs_items ON krs_items.id = grades.krs_item_id").
		Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Where("krs.student_id = ?", studentID).
		Order("krs.academic_period_id ASC").
		Find(&grades).Error
	return grades, err
}

func (r *AcademicRepository) GetGradesByCourse(courseID string) ([]domain.Grade, error) {
	var grades []domain.Grade
	err := r.db.Joins("JOIN krs_items ON krs_items.id = grades.krs_item_id").
		Joins("JOIN classes ON classes.id = krs_items.class_id").
		Joins("JOIN course_offerings ON course_offerings.id = classes.course_offering_id").
		Where("course_offerings.course_id = ?", courseID).
		Find(&grades).Error
	return grades, err
}

// GradeEntry retrieval for transcript
func (r *AcademicRepository) GetGradeEntriesByStudent(studentID string) ([]domain.GradeEntry, error) {
	var entries []domain.GradeEntry
	err := r.db.Where("student_id = ? AND status = 'published'", studentID).
		Order("entered_at ASC").
		Find(&entries).Error
	return entries, err
}

func (r *AcademicRepository) GetGradeEntryByStudentAndKrsItem(studentID, krsItemID string) (*domain.GradeEntry, error) {
	var ge domain.GradeEntry
	err := r.db.Where("student_id = ? AND krs_item_id = ?", studentID, krsItemID).First(&ge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ge, nil
}

// ============ Schedule Repository Methods ============

// CreateSchedule creates a new class schedule
func (r *AcademicRepository) CreateSchedule(s *domain.ClassSchedule) error {
	return r.db.Create(s).Error
}

// GetScheduleByID retrieves a schedule by ID
func (r *AcademicRepository) GetScheduleByID(id string) (*domain.ClassSchedule, error) {
	var s domain.ClassSchedule
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// GetSchedulesByClassID retrieves all schedules for a class
func (r *AcademicRepository) GetSchedulesByClassID(classID string) ([]domain.ClassSchedule, error) {
	var schedules []domain.ClassSchedule
	err := r.db.Where("class_id = ?", classID).Order("day_of_week ASC, start_time ASC").Find(&schedules).Error
	return schedules, err
}

// UpdateSchedule updates an existing schedule
func (r *AcademicRepository) UpdateSchedule(s *domain.ClassSchedule) error {
	return r.db.Save(s).Error
}

// DeleteSchedule deletes a schedule
func (r *AcademicRepository) DeleteSchedule(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.ClassSchedule{}).Error
}

// GetWeeklySchedule retrieves weekly schedule by study program and academic period
func (r *AcademicRepository) GetWeeklySchedule(studyProgramID, academicPeriodID string) ([]domain.ClassSchedule, error) {
	var schedules []domain.ClassSchedule
	query := r.db.Joins("JOIN classes ON classes.id = class_schedules.class_id").
		Joins("JOIN course_offerings ON course_offerings.id = classes.course_offering_id").
		Joins("JOIN study_programs ON study_programs.id = course_offerings.study_program_id").
		Where("1=1")

	if studyProgramID != "" {
		query = query.Where("study_programs.id = ?", studyProgramID)
	}
	if academicPeriodID != "" {
		query = query.Where("course_offerings.academic_period_id = ?", academicPeriodID)
	}

	err := query.Order("class_schedules.day_of_week ASC, class_schedules.start_time ASC").Find(&schedules).Error
	return schedules, err
}

// GetScheduleByLecturer retrieves schedules for a specific lecturer
func (r *AcademicRepository) GetScheduleByLecturer(lecturerID, dayOfWeek string) ([]domain.ClassSchedule, error) {
	var schedules []domain.ClassSchedule
	query := r.db.Joins("JOIN class_lecturers ON class_lecturers.class_id = class_schedules.class_id").
		Where("class_lecturers.lecturer_id = ?", lecturerID)

	if dayOfWeek != "" {
		query = query.Where("class_schedules.day_of_week = ?", dayOfWeek)
	}

	err := query.Order("class_schedules.day_of_week ASC, class_schedules.start_time ASC").Find(&schedules).Error
	return schedules, err
}

// GetStudentSchedule retrieves schedule for a student based on their enrolled classes
func (r *AcademicRepository) GetStudentSchedule(studentID, dayOfWeek string) ([]domain.ClassSchedule, error) {
	var schedules []domain.ClassSchedule
	query := r.db.Joins("JOIN class_enrollments ON class_enrollments.class_id = class_schedules.class_id").
		Joins("JOIN krs ON krs.id = class_enrollments.krs_id").
		Where("krs.student_id = ? AND krs.status = 'approved'", studentID)

	if dayOfWeek != "" {
		query = query.Where("class_schedules.day_of_week = ?", dayOfWeek)
	}

	err := query.Order("class_schedules.day_of_week ASC, class_schedules.start_time ASC").Find(&schedules).Error
	return schedules, err
}

// ============ Attendance Repository Methods ============

// CreateAttendance creates a new attendance record
func (r *AcademicRepository) CreateAttendance(a *domain.StudentAttendance) error {
	return r.db.Create(a).Error
}

// UpdateAttendance updates an existing attendance record
func (r *AcademicRepository) UpdateAttendance(a *domain.StudentAttendance) error {
	return r.db.Save(a).Error
}

// GetAttendance retrieves attendance by student, class, and date
func (r *AcademicRepository) GetAttendance(studentID, classID string, sessionDate time.Time) (*domain.StudentAttendance, error) {
	var a domain.StudentAttendance
	err := r.db.Where("student_id = ? AND class_id = ? AND session_date = ?", studentID, classID, sessionDate).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// GetStudentAttendance retrieves all attendance records for a student
func (r *AcademicRepository) GetStudentAttendance(studentID, classID, academicPeriodID string) ([]domain.StudentAttendance, error) {
	var attendances []domain.StudentAttendance
	query := r.db.Where("student_id = ?", studentID)

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}

	err := query.Order("session_date DESC").Find(&attendances).Error
	return attendances, err
}

// GetClassAttendance retrieves all attendance records for a class
func (r *AcademicRepository) GetClassAttendance(classID, sessionDate string) ([]domain.StudentAttendance, error) {
	var attendances []domain.StudentAttendance
	query := r.db.Where("class_id = ?", classID)

	if sessionDate != "" {
		date, _ := time.Parse("2006-01-02", sessionDate)
		query = query.Where("session_date = ?", date)
	}

	err := query.Order("student_id ASC").Find(&attendances).Error
	return attendances, err
}

// GetAttendanceStats calculates attendance statistics
func (r *AcademicRepository) GetAttendanceStats(studentID, classID string) (map[string]interface{}, error) {
	query := r.db.Model(&domain.StudentAttendance{})

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}

	var total int64
	var present, absent, excused, sick int64

	query.Count(&total)
	r.db.Model(&domain.StudentAttendance{}).Where("student_id = ? AND status = 'present'", studentID).Count(&present)
	r.db.Model(&domain.StudentAttendance{}).Where("student_id = ? AND status = 'absent'", studentID).Count(&absent)
r.db.Model(&domain.StudentAttendance{}).Where("student_id = ? AND status = 'excused'", studentID).Count(&excused)
	r.db.Model(&domain.StudentAttendance{}).Where("student_id = ? AND status = 'sick'", studentID).Count(&sick)

	percentage := 0.0
	if total > 0 {
		percentage = float64(present) / float64(total) * 100
	}

	stats := map[string]interface{}{
		"total":     total,
		"present":   present,
		"absent":    absent,
		"excused":  excused,
		"sick":     sick,
		"percentage": percentage,
	}

	return stats, nil
}

// IsStudentEnrolledInClass checks if a student is enrolled in a class
func (r *AcademicRepository) IsStudentEnrolledInClass(studentID, classID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.KrsItem{}).
		Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Where("krs.student_id = ? AND krs_items.class_id = ? AND krs.status = 'approved'", studentID, classID).
		Count(&count).Error

	return count > 0, err
}

func (r *AcademicRepository) UpdateStudent(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.Student{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AcademicRepository) GetKrsByStudentID(studentID string, academicPeriodID string, status string) ([]domain.KRS, error) {
	var list []domain.KRS
	query := r.db.Model(&domain.KRS{}).Where("student_id = ?", studentID)
	if academicPeriodID != "" {
		query = query.Where("academic_period_id = ?", academicPeriodID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Preload("Items").Find(&list).Error
	return list, err
}

func (r *AcademicRepository) GetGradesByStudentID(studentID string, academicPeriodID string) ([]domain.Grade, error) {
	return r.GetStudentGrades(studentID, academicPeriodID)
}

func (r *AcademicRepository) GetCourseByKrsItemID(krsItemID string) (*domain.Course, error) {
	var course domain.Course
	err := r.db.Table("courses").
		Joins("JOIN course_offerings ON course_offerings.course_id = courses.id").
		Joins("JOIN classes ON classes.course_offering_id = course_offerings.id").
		Joins("JOIN krs_items ON krs_items.class_id = classes.id").
		Where("krs_items.id = ?", krsItemID).
		First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *AcademicRepository) GetPeriodNameByKrsItemID(krsItemID string) (string, error) {
	var periodName string
	err := r.db.Table("academic_periods").
		Select("academic_periods.period_name").
		Joins("JOIN krs ON krs.academic_period_id = academic_periods.id").
		Joins("JOIN krs_items ON krs_items.krs_id = krs.id").
		Where("krs_items.id = ?", krsItemID).
		Row().Scan(&periodName)
	return periodName, err
}
