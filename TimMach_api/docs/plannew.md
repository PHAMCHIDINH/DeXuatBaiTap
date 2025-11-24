# Plan: Implement Reports Feature

## Overview
Add database table and API endpoints to store patient report metadata and enable report management (create, list, download, email).

---

## Current State vs New State

### Current Logic (Old)
```
POST /patients/:id/report/email
    ↓
[1] Generate PDF on-the-fly
[2] Send email immediately
[3] NO database record
    ↓
Response: "Email sent"
```

**Issues:**
- Must regenerate PDF every time
- No report history
- Cannot resend old reports

### New Logic (With reports table)

**Workflow 1: Create Report**
```
POST /api/patients/:id/reports
    ↓
[1] Validate user ownership
[2] Generate PDF using buildPatientReportPDF()
[3] Save file to disk: tmp/reports/:patient_id/:filename.pdf
[4] Insert record to DB with file_url
    ↓
Response: Report object with id, filename, file_url
```

**Workflow 2: List Reports**
```
GET /api/patients/:id/reports?limit=10&offset=0
    ↓
[1] Validate ownership
[2] Query: ListReportsByPatient()
    ↓
Response: Array of reports with recipients tracking
```

**Workflow 3: Download Report**
```
GET /api/reports/:report_id/download
    ↓
[1] Get report from DB
[2] Validate ownership
[3] Read file from disk
    ↓
Response: PDF file (binary)
```

**Workflow 4: Send Report via Email**
```
POST /api/reports/:report_id/email
Body: {
  "email": "patient@example.com",
  "subject": "Your Health Report",
  "message": "..."
}
    ↓
[1] Get report from DB
[2] Validate ownership
[3] Read PDF from disk
[4] Send email via mailer service
[5] Update recipients JSONB in DB
    ↓
Response: {message, report_id, sent_to}
```

**Workflow 5: Delete Report**
```
DELETE /api/reports/:report_id
    ↓
[1] Get report from DB
[2] Validate ownership
[3] Delete file from disk
[4] DELETE FROM reports
    ↓
Response: 204 No Content
```

---

## Database Schema

### reports table
```sql
CREATE TABLE reports (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    file_url TEXT NOT NULL,
    recipients JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_patient_id ON reports(patient_id);
CREATE INDEX idx_reports_recipients ON reports USING GIN (recipients);
```

### recipients JSONB structure
```json
[
  {
    "email": "patient@example.com",
    "sent_at": "2025-11-24T10:35:00Z",
    "status": "sent"
  },
  {
    "email": "doctor@clinic.com",
    "sent_at": "2025-11-24T11:00:00Z",
    "status": "sent"
  }
]
```

---

## SQL Queries (reports.sql)

```sql
-- name: CreateReport :one
INSERT INTO reports (patient_id, filename, file_url, recipients)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReportByID :one
SELECT * FROM reports
WHERE id = $1
LIMIT 1;

-- name: ListReportsByPatient :many
SELECT * FROM reports
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReportRecipients :one
UPDATE reports
SET recipients = $2
WHERE id = $1
RETURNING *;

-- name: DeleteReport :exec
DELETE FROM reports
WHERE id = $1;

-- name: CountReportsByPatient :one
SELECT COUNT(*) FROM reports
WHERE patient_id = $1;
```

---

## API Endpoints

### New Routes (patients/routers.go)
```go
func RegisterRoutes(r *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
    p := r.Group("/patients")
    p.Use(authMW)
    {
        // Existing endpoints...
        
        // NEW: Report management
        p.POST("/:id/reports", h.CreateReport)      // Create report
        p.GET("/:id/reports", h.ListReports)        // List reports
    }
    
    // NEW: Report actions
    reports := r.Group("/reports")
    reports.Use(authMW)
    {
        reports.GET("/:id/download", h.DownloadReport)   // Download PDF
        reports.POST("/:id/email", h.SendReportEmail)    // Send email
        reports.DELETE("/:id", h.DeleteReport)           // Delete report
    }
}
```

---

## Data Types (types.go)

### Request DTOs
```go
type SendReportEmailRequest struct {
    Email   string `json:"email" binding:"required,email"`
    Subject string `json:"subject"`
    Message string `json:"message"`
}

type ListReportsRequest struct {
    Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
    Offset int32 `form:"offset" binding:"omitempty,min=0"`
}
```

### Response DTOs
```go
type ReportResponse struct {
    ID         int64                `json:"id"`
    PatientID  int64                `json:"patient_id"`
    Filename   string               `json:"filename"`
    FileURL    string               `json:"file_url"`
    Recipients []ReportRecipient    `json:"recipients"`
    CreatedAt  time.Time            `json:"created_at"`
}

type ReportRecipient struct {
    Email  string    `json:"email"`
    SentAt time.Time `json:"sent_at"`
    Status string    `json:"status"`
}

type ListReportsResponse struct {
    Reports []ReportResponse `json:"reports"`
    Total   int64            `json:"total"`
}
```

---

## Handler Functions (report.go)

### 1. CreateReport
```go
func (h *Handler) CreateReport(c *gin.Context) {
    // [1] Get userID from context
    // [2] Get patient_id from URL param
    // [3] Validate patient belongs to user
    // [4] Generate PDF using buildPatientReportPDF()
    // [5] Save file to disk: tmp/reports/{patient_id}/{filename}.pdf
    // [6] Create file_url: /api/reports/{report_id}/download
    // [7] Insert to DB: CreateReport()
    // [8] Return ReportResponse
}
```

### 2. ListReports
```go
func (h *Handler) ListReports(c *gin.Context) {
    // [1] Get userID from context
    // [2] Get patient_id from URL param
    // [3] Validate patient belongs to user
    // [4] Parse limit/offset from query params
    // [5] Query: ListReportsByPatient()
    // [6] Query: CountReportsByPatient()
    // [7] Convert db.Report[] to ReportResponse[]
    // [8] Return ListReportsResponse
}
```

### 3. DownloadReport
```go
func (h *Handler) DownloadReport(c *gin.Context) {
    // [1] Get userID from context
    // [2] Get report_id from URL param
    // [3] Get report from DB: GetReportByID()
    // [4] Get patient to validate ownership
    // [5] Read file from disk using file_url
    // [6] Set headers: Content-Type, Content-Disposition
    // [7] Return PDF bytes
}
```

### 4. SendReportEmail
```go
func (h *Handler) SendReportEmail(c *gin.Context) {
    // [1] Get userID from context
    // [2] Get report_id from URL param
    // [3] Bind SendReportEmailRequest
    // [4] Get report from DB
    // [5] Validate ownership via patient
    // [6] Read PDF file from disk
    // [7] Send email via mailer.Send()
    // [8] Update recipients JSONB in DB
    // [9] Return success response
}
```

### 5. DeleteReport
```go
func (h *Handler) DeleteReport(c *gin.Context) {
    // [1] Get userID from context
    // [2] Get report_id from URL param
    // [3] Get report from DB
    // [4] Validate ownership
    // [5] Delete file from disk
    // [6] DeleteReport from DB
    // [7] Return 204 No Content
}
```

---

## Converter Functions (converter.go)

```go
func mapReportResponse(r db.Report) (ReportResponse, error) {
    // Parse recipients JSONB to []ReportRecipient
    var recipients []ReportRecipient
    if err := json.Unmarshal(r.Recipients, &recipients); err != nil {
        return ReportResponse{}, err
    }
    
    return ReportResponse{
        ID:         r.ID,
        PatientID:  r.PatientID,
        Filename:   r.Filename,
        FileURL:    r.FileURL,
        Recipients: recipients,
        CreatedAt:  r.CreatedAt.Time,
    }, nil
}
```

---

## File Storage Strategy

### Option 1: Local Filesystem (Recommended for MVP)
```
tmp/reports/
  ├── 123/                    # patient_id
  │   ├── report_123_20251124_103045.pdf
  │   └── report_123_20251124_110000.pdf
  └── 456/
      └── report_456_20251124_120000.pdf
```

**Pros:**
- Simple implementation
- No external dependencies
- Good for development

**Cons:**
- Files lost if container restarts (unless volume mounted)
- Not scalable for multiple servers
- Need cleanup job for old files

### Option 2: S3/MinIO (For Production)
- Upload to S3 after generation
- file_url = presigned S3 URL
- Automatic cleanup via lifecycle policy

---

## Implementation Checklist

### Phase 1: Database & Queries
- [x] Create migration: `20251124023816_add_table_report.sql`
- [x] Create queries: `db/queries/reports.sql`
- [ ] Run migration: `make goose-up`
- [ ] Generate sqlc: `make sqlc`

### Phase 2: Data Types
- [ ] Add request/response DTOs to `modules/patients/types.go`
- [ ] Add converter functions to `modules/patients/converter.go`

### Phase 3: Handlers
- [ ] Implement `CreateReport()` in `modules/patients/report.go`
- [ ] Implement `ListReports()` in `modules/patients/report.go`
- [ ] Implement `DownloadReport()` in `modules/patients/report.go`
- [ ] Implement `SendReportEmail()` in `modules/patients/report.go`
- [ ] Implement `DeleteReport()` in `modules/patients/report.go`

### Phase 4: Routes
- [ ] Register new routes in `modules/patients/routers.go`

### Phase 5: File Storage
- [ ] Create directory: `tmp/reports/`
- [ ] Add cleanup utility (optional)
- [ ] Update .gitignore to exclude tmp/reports/

### Phase 6: Testing
- [ ] Test create report endpoint
- [ ] Test list reports endpoint
- [ ] Test download report endpoint
- [ ] Test send email endpoint
- [ ] Test delete report endpoint
- [ ] Test ownership validation
- [ ] Test file storage/retrieval

---

## Example Usage Scenario

**Doctor creates and sends report:**
```
1. View patient info
   GET /api/patients/123

2. Create report
   POST /api/patients/123/reports
   Response: {id: 5, filename: "report_123_20251124.pdf", ...}

3. View reports
   GET /api/patients/123/reports
   Response: {reports: [{id: 5, ...}, {id: 4, ...}], total: 2}

4. Send to patient
   POST /api/reports/5/email
   Body: {email: "patient@gmail.com", subject: "Your Report"}
   DB updated: recipients += [{"email":"patient@gmail.com","sent_at":"..."}]

5. Send to another doctor
   POST /api/reports/5/email
   Body: {email: "doctor@clinic.com"}
   DB updated: recipients += [{"email":"doctor@clinic.com","sent_at":"..."}]

6. Patient downloads report
   GET /api/reports/5/download
   Response: PDF file
```

---

## Key Differences: Old vs New

| Feature | Old Logic | New Logic |
|---------|-----------|-----------|
| **Generate PDF** | Every email send | Once at report creation |
| **Storage** | None | File + DB record |
| **History** | No | Yes (reports table) |
| **Resend** | Must regenerate | Use existing file |
| **Tracking** | No | Yes (recipients JSONB) |
| **Performance** | Slow (regenerate each time) | Fast (reuse file) |
| **Audit Trail** | No | Yes (created_at, recipients) |

---

## Security Considerations

1. **Ownership Validation**
   - Always verify patient belongs to authenticated user
   - Check both at patient level and report level

2. **File Access Control**
   - Don't expose file system paths in URLs
   - Use report_id for download endpoint
   - Validate ownership before file download

3. **Email Validation**
   - Validate email format
   - Optional: Rate limiting on email endpoint
   - Optional: Email verification

4. **File Cleanup**
   - Consider adding job to delete old reports
   - Or implement soft delete with retention policy

---

## Future Enhancements

- [ ] Bulk email sending (multiple recipients at once)
- [ ] Report templates (different report types)
- [ ] Scheduled reports (weekly/monthly)
- [ ] Email delivery status tracking (bounced, opened)
- [ ] Report versioning
- [ ] Internationalization (English/Vietnamese)
- [ ] PDF customization (clinic branding)
- [ ] S3 integration for production
- [ ] Report sharing via public link (with expiry)
- [ ] Report comments/annotations

---

## Notes

- Keep existing `/patients/:id/report.pdf` and `/patients/:id/report/email` endpoints for backward compatibility (mark as deprecated)
- PDF generation logic remains unchanged (buildPatientReportPDF)
- Mailer service remains unchanged
- Template remains unchanged (patient_report.html)
