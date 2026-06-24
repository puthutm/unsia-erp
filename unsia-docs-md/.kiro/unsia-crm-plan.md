# CRM Service Implementation Plan

## Service Overview
**Service Name**: unsia-crm-service (Customer Relationship Management)
**Port**: 8002
**Database**: crm_db (PostgreSQL)

## Core Features

### 1. Lead Management
- Create & manage prospects/camaba
- Lead status (new, contacted, qualified, converted)
- lead sources

### 2. Campaign Management
- PMB campaigns
- Campaign tracking

### 3. Activities
- Follow-up activities
- Call logs, meetings

---

## Domain Models

```
Lead
├── ID (UUID)
├── Name
├── Email
├── Phone
├── Source (FACEBOOK, INSTAGRAM, WEBSITE, REFERRAL)
├── Status (NEW, CONTACTED, QUALIFIED, CONVERTED, LOST)
├── StudyProgramInterestID (FK)
├── CampaignID (FK)
├── CreatedAt
└── UpdatedAt

Activity
├── LeadID (FK)
├── Type (CALL, MEETING, WHATSAPP)
├── Description
├── DueDate
└── CompletedAt
```

---

## Implementation Timeline: 6-8 weeks
