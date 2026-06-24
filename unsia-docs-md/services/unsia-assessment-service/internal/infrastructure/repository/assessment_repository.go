package repository

import (
	"errors"

	"github.com/unsia-erp/unsia-assessment-service/internal/domain"
	"gorm.io/gorm"
)

type AssessmentRepository struct {
	db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) *AssessmentRepository {
	return &AssessmentRepository{db: db}
}

func (r *AssessmentRepository) CreateSession(s *domain.AssessmentSession) error {
	return r.db.Create(s).Error
}

func (r *AssessmentRepository) GetSessionByID(id string) (*domain.AssessmentSession, error) {
	var s domain.AssessmentSession
	err := r.db.Preload("Participants").Where("id = ?", id).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *AssessmentRepository) RegisterParticipant(p *domain.AssessmentParticipant) error {
	return r.db.Create(p).Error
}

func (r *AssessmentRepository) GetParticipantByID(id string) (*domain.AssessmentParticipant, error) {
	var p domain.AssessmentParticipant
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *AssessmentRepository) CreateAttempt(a *domain.AssessmentAttempt) error {
	return r.db.Create(a).Error
}

func (r *AssessmentRepository) GetAttemptByID(id string) (*domain.AssessmentAttempt, error) {
	var a domain.AssessmentAttempt
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AssessmentRepository) UpdateAttemptScore(id string, score float64, status string) error {
	updates := map[string]interface{}{
		"total_score": &score,
		"status":      status,
	}
	return r.db.Model(&domain.AssessmentAttempt{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AssessmentRepository) CreateQuestionBank(qb *domain.QuestionBank) error {
	return r.db.Create(qb).Error
}

func (r *AssessmentRepository) CreateQuestion(q *domain.Question) error {
	return r.db.Create(q).Error
}

func (r *AssessmentRepository) GetQuestionByID(id string) (*domain.Question, error) {
	var q domain.Question
	err := r.db.Preload("Options").Where("id = ?", id).First(&q).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *AssessmentRepository) CreateQuestionVersion(qv *domain.QuestionVersion) error {
	return r.db.Create(qv).Error
}

func (r *AssessmentRepository) GetLatestQuestionVersion(questionID string) (int, error) {
	var maxVal int
	row := r.db.Table("question_versions").
		Select("COALESCE(MAX(version_number), 0)").
		Where("question_id = ?", questionID).
		Row()
	err := row.Scan(&maxVal)
	return maxVal, err
}

func (r *AssessmentRepository) CreateQuestionOption(qo *domain.QuestionOption) error {
	return r.db.Create(qo).Error
}

func (r *AssessmentRepository) SaveAnswer(ans *domain.AssessmentAnswer) error {
	var existing domain.AssessmentAnswer
	err := r.db.Where("attempt_id = ? AND question_id = ?", ans.AttemptID, ans.QuestionID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(ans).Error
		}
		return err
	}

	existing.SelectedOptionID = ans.SelectedOptionID
	existing.AnswerText = ans.AnswerText
	existing.Score = ans.Score
	existing.GradedBy = ans.GradedBy
	existing.GradedAt = ans.GradedAt
	return r.db.Save(&existing).Error
}

func (r *AssessmentRepository) GetAnswersByAttemptID(attemptID string) ([]domain.AssessmentAnswer, error) {
	var answers []domain.AssessmentAnswer
	err := r.db.Where("attempt_id = ?", attemptID).Find(&answers).Error
	return answers, err
}
