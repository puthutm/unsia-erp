package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	sharedhttpclient "github.com/unsia-erp/shared-httpclient"
	sharedrbac "github.com/unsia-erp/shared-rbac"
	"github.com/unsia-erp/unsia-crm-service/internal/domain"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type CampaignCreateRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Channel   string `json:"channel"`
	StartDate string `json:"start_date"` // format: YYYY-MM-DD
	EndDate   string `json:"end_date"`   // format: YYYY-MM-DD
}

type AgentCreateRequest struct {
	PersonID         string `json:"person_id" binding:"required"`
	AgentCode        string `json:"agent_code" binding:"required"`
	OrganizationName string `json:"organization_name"`
}

type AgentApproveRequest struct {
	ApprovalStatus string `json:"approval_status" binding:"required,oneof=approved rejected"`
	Status         string `json:"status" binding:"required,oneof=active suspended"`
}

type ReferralCreateRequest struct {
	ReferralType     string  `json:"referral_type" binding:"required,oneof=agent individual public"`
	ReferrerPersonID *string `json:"referrer_person_id"`
	AgentID          *string `json:"agent_id"`
	ReferralCode     string  `json:"referral_code" binding:"required"`
}

type LeadCreateRequest struct {
	PersonID       string  `json:"person_id" binding:"required"`
	StudyProgramID *string `json:"study_program_id"`
	LeadSourceID   *string `json:"lead_source_id"`
	CampaignID     *string `json:"campaign_id"`
	ReferralCode   *string `json:"referral_code"`
}

type LeadActivityCreateRequest struct {
	ActivityType string `json:"activity_type" binding:"required"`
	Note         string `json:"note"`
}

type LeadConvertRequest struct {
	PmbWaveID           *string `json:"pmb_wave_id"`
	AdmissionPathID     *string `json:"admission_path_id"`
	TargetEntryPeriodID *string `json:"target_entry_period_id"`
}

type CRMHandler struct {
	repo      *repository.CRMRepository
	db        *gorm.DB
	pmbClient *sharedhttpclient.Client
}

func NewCRMHandler(db *gorm.DB) *CRMHandler {
	pmbURL := os.Getenv("PMB_SERVICE_URL")
	if pmbURL == "" {
		pmbURL = "http://localhost:8004"
	}
	srvToken := os.Getenv("CRM_SERVICE_TOKEN")
	if srvToken == "" {
		srvToken = "crm_service_secret_token"
	}

	pmbClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      pmbURL,
		ServiceToken: srvToken,
		SourceName:   "crm-service",
		Timeout:      10 * time.Second,
	})

	return &CRMHandler{
		repo:      repository.NewCRMRepository(db),
		db:        db,
		pmbClient: pmbClient,
	}
}

func generateLeadNumber() string {
	now := time.Now().Format("20060102")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(900000))
	num := nBig.Int64() + 100000
	return fmt.Sprintf("L%s%d", now, num)
}

// ListCampaigns returns all campaigns
func (h *CRMHandler) ListCampaigns(c *gin.Context) {
	campaigns, err := h.repo.ListCampaigns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data campaign").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(campaigns).WithContext(c))
}

// CreateCampaign creates a new campaign
func (h *CRMHandler) CreateCampaign(c *gin.Context) {
	var req CampaignCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var startPtr, endPtr *time.Time
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "start_date format must be YYYY-MM-DD").WithContext(c))
			return
		}
		startPtr = &start
	}
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "end_date format must be YYYY-MM-DD").WithContext(c))
			return
		}
		endPtr = &end
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	campaign := domain.Campaign{
		Code:      req.Code,
		Name:      req.Name,
		Channel:   req.Channel,
		StartDate: startPtr,
		EndDate:   endPtr,
		Status:    "active",
		CreatedBy: actor,
	}

	if err := h.repo.CreateCampaign(&campaign); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan campaign").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.campaign.create",
		Module:       "crm",
		ResourceType: "campaign",
		ResourceID:   campaign.ID,
		NewValue:     campaign,
	})

	c.JSON(http.StatusCreated, sharederr.Success(campaign).WithContext(c))
}

// ListAgents returns all agents
func (h *CRMHandler) ListAgents(c *gin.Context) {
	agents, err := h.repo.ListAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data agent").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(agents).WithContext(c))
}

// CreateAgent creates a new agent
func (h *CRMHandler) CreateAgent(c *gin.Context) {
	var req AgentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	agent := domain.Agent{
		PersonID:         req.PersonID,
		AgentCode:        req.AgentCode,
		OrganizationName: req.OrganizationName,
		Status:           "inactive",
		ApprovalStatus:   "pending",
	}

	if err := h.repo.CreateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan agent").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.agent.register",
		Module:       "crm",
		ResourceType: "agent",
		ResourceID:   agent.ID,
		NewValue:     agent,
	})

	c.JSON(http.StatusCreated, sharederr.Success(agent).WithContext(c))
}

// ApproveAgent approves or rejects agent registration
func (h *CRMHandler) ApproveAgent(c *gin.Context) {
	agentID := c.Param("id")
	var req AgentApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	agent, err := h.repo.GetAgentByID(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Agent tidak ditemukan").WithContext(c))
		return
	}

	oldAgent := *agent

	err = h.repo.UpdateAgentApproval(agentID, req.ApprovalStatus, actor, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate persetujuan agent").WithContext(c))
		return
	}

	updatedAgent, _ := h.repo.GetAgentByID(agentID)

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.agent.approve",
		Module:       "crm",
		ResourceType: "agent",
		ResourceID:   agentID,
		OldValue:     oldAgent,
		NewValue:     updatedAgent,
	})

	c.JSON(http.StatusOK, sharederr.Success(updatedAgent).WithContext(c))
}

// ListReferrals returns all referrals
func (h *CRMHandler) ListReferrals(c *gin.Context) {
	referrals, err := h.repo.ListReferrals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data referral").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(referrals).WithContext(c))
}

// CreateReferral creates a new referral code
func (h *CRMHandler) CreateReferral(c *gin.Context) {
	var req ReferralCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	referral := domain.Referral{
		ReferralType:     req.ReferralType,
		ReferrerPersonID: req.ReferrerPersonID,
		AgentID:          req.AgentID,
		ReferralCode:     req.ReferralCode,
		IsValid:          true,
	}

	if err := h.repo.CreateReferral(&referral); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan referral").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.referral.create",
		Module:       "crm",
		ResourceType: "referral",
		ResourceID:   referral.ID,
		NewValue:     referral,
	})

	c.JSON(http.StatusCreated, sharederr.Success(referral).WithContext(c))
}

// ListLeads lists leads applying own_lead scope checks
func (h *CRMHandler) ListLeads(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	scopeStr := sharedrbac.ResolveDataScope(claims)
	userScope := sharedrbac.ParseScope(scopeStr, userID)

	leads, err := h.repo.ListLeads(userScope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data lead").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(leads).WithContext(c))
}

// CreateLead creates a new lead
func (h *CRMHandler) CreateLead(c *gin.Context) {
	var req LeadCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	var referralID *string
	if req.ReferralCode != nil && *req.ReferralCode != "" {
		ref, err := h.repo.GetReferralByCode(*req.ReferralCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REFERRAL", "Kode referral tidak valid atau kedaluwarsa").WithContext(c))
			return
		}
		referralID = &ref.ID
	}

	lead := domain.Lead{
		PersonID:       req.PersonID,
		StudyProgramID: req.StudyProgramID,
		LeadSourceID:   req.LeadSourceID,
		CampaignID:     req.CampaignID,
		ReferralID:     referralID,
		LeadNumber:     generateLeadNumber(),
		Status:         "new",
		OwnerUserID:    &userID,
	}

	if err := h.repo.CreateLead(&lead); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan lead").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.lead.create",
		Module:       "crm",
		ResourceType: "lead",
		ResourceID:   lead.ID,
		NewValue:     lead,
	})

	c.JSON(http.StatusCreated, sharederr.Success(lead).WithContext(c))
}

// CreateLeadActivity records an activity on a lead
func (h *CRMHandler) CreateLeadActivity(c *gin.Context) {
	leadID := c.Param("id")
	var req LeadActivityCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	// Verify lead exists
	lead, err := h.repo.GetLeadByID(leadID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Lead tidak ditemukan").WithContext(c))
		return
	}

	activity := domain.LeadActivity{
		LeadID:       lead.ID,
		UserID:       &userID,
		ActivityType: req.ActivityType,
		Note:         req.Note,
		ActivityAt:   time.Now(),
	}

	if err := h.repo.CreateLeadActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan aktivitas lead").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(activity).WithContext(c))
}

// PMBResponse matches the response format of PMB service
type PMBResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
}

// ConvertLeadToApplicant converts lead, calls PMB service, generates commissions
func (h *CRMHandler) ConvertLeadToApplicant(c *gin.Context) {
	leadID := c.Param("id")
	var req LeadConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	lead, err := h.repo.GetLeadByID(leadID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Lead tidak ditemukan").WithContext(c))
		return
	}

	if lead.Status == "converted" {
		c.JSON(http.StatusBadRequest, sharederr.Error("ALREADY_CONVERTED", "Lead ini sudah dikonversi sebelumnya").WithContext(c))
		return
	}

	// Prepare payload for PMB service
	pmbPayload := map[string]interface{}{
		"person_id":              lead.PersonID,
		"crm_lead_id":            lead.ID,
		"study_program_id":       lead.StudyProgramID,
		"pmb_wave_id":            req.PmbWaveID,
		"admission_path_id":      req.AdmissionPathID,
		"target_entry_period_id": req.TargetEntryPeriodID,
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	// Propagate context (which contains correlation ID)
	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	resp, err := h.pmbClient.Post(ctx, "/api/v1/pmb/applicants", pmbPayload)
	if err != nil {
		c.JSON(http.StatusBadGateway, sharederr.Error("PMB_SERVICE_UNAVAILABLE", fmt.Sprintf("Gagal menghubungi layanan PMB: %v", err)).WithContext(c))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusBadGateway, sharederr.Error("PMB_SERVICE_ERROR", fmt.Sprintf("Layanan PMB mengembalikan error code: %d", resp.StatusCode)).WithContext(c))
		return
	}

	var pmbResp PMBResponse
	if err := json.NewDecoder(resp.Body).Decode(&pmbResp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("PARSE_ERROR", "Gagal membaca response PMB").WithContext(c))
		return
	}

	// Begin GORM transaction in CRM database
	err = h.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 1. Update Lead Status
		oldStatus := lead.Status
		lead.Status = "converted"
		lead.ConvertedAt = &now
		lead.UpdatedAt = now
		if err := tx.Model(lead).Updates(map[string]interface{}{
			"status":       "converted",
			"converted_at": &now,
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}

		// 2. Create Lead Status History
		history := domain.LeadStatusHistory{
			LeadID:    lead.ID,
			OldStatus: &oldStatus,
			NewStatus: "converted",
			ChangedBy: &userID,
			Note:      "Lead converted to PMB Applicant ID: " + pmbResp.Data.ID,
			ChangedAt: now,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		// 3. Create Commission Record if Lead has Referral
		if lead.ReferralID != nil && *lead.ReferralID != "" {
			var ref domain.Referral
			if err := tx.Where("id = ?", *lead.ReferralID).First(&ref).Error; err == nil && ref.IsValid {
				// Get commission rule
				var rule domain.CommissionRule
				if err := tx.Where("referral_type = ? AND is_active = true", ref.ReferralType).First(&rule).Error; err == nil {
					commission := domain.CommissionRecord{
						LeadID:           lead.ID,
						CommissionRuleID: &rule.ID,
						ReferrerPersonID: ref.ReferrerPersonID,
						Amount:           rule.Amount,
						Status:           "draft",
					}
					if err := tx.Create(&commission).Error; err != nil {
						return err
					}
				}
			}
		}

		// 4. Publish Outbox Event
		envelope := sharedevent.EventEnvelope{
			EventName:        "crm.lead_converted",
			EventVersion:     "v1",
			PublisherService: "crm-service",
			AggregateType:    "lead",
			AggregateID:      lead.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"lead_id":      lead.ID,
				"applicant_id": pmbResp.Data.ID,
				"converted_at": now,
			},
		}

		conn := tx.Statement.ConnPool
		_, err = sharedevent.WriteOutbox(ctx, conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TRANSACTION_FAILED", "Gagal memproses konversi lead").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "crm.lead.convert",
		Module:       "crm",
		ResourceType: "lead",
		ResourceID:   lead.ID,
		NewValue:     lead,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"lead_id":      lead.ID,
		"applicant_id": pmbResp.Data.ID,
		"status":       "converted",
	}).WithContext(c))
}

func (h *CRMHandler) ListContacts(c *gin.Context) {
	var personIDs []string
	err := h.db.Model(&domain.Lead{}).Distinct("person_id").Pluck("person_id", &personIDs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kontak").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"contacts": personIDs,
	}).WithContext(c))
}

func (h *CRMHandler) CreateContact(c *gin.Context) {
	var req struct {
		PersonID string `json:"person_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}
	lead := domain.Lead{
		PersonID:   req.PersonID,
		LeadNumber: generateLeadNumber(),
		Status:     "new",
	}
	if err := h.db.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan kontak baru").WithContext(c))
		return
	}
	c.JSON(http.StatusCreated, sharederr.Success(lead).WithContext(c))
}

func (h *CRMHandler) ListOpportunities(c *gin.Context) {
	var opportunities []domain.Lead
	err := h.db.Where("status NOT IN (?)", []string{"converted", "lost"}).Order("created_at desc").Find(&opportunities).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data peluang").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(opportunities).WithContext(c))
}

func (h *CRMHandler) CreateOpportunity(c *gin.Context) {
	h.CreateLead(c)
}

