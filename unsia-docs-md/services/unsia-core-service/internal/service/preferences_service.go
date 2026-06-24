package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettingsService manages user/system preferences
type SettingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) *SettingsService {
	return &SettingsService{db: db}
}

type UserSetting struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	Key       string    `gorm:"column:key"`
	Value     string    `gorm:"column:value;type:jsonb"`
	IsDefault bool      `gorm:"column:is_default"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}

// SystemSetting represents system settings
type SystemSetting struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Key         string    `gorm:"column:key;uniqueIndex"`
	Value       string    `gorm:"column:value;type:jsonb"`
	Description string    `gorm:"column:description"`
	Type       string    `gorm:"column:type"` // string, number, boolean, json
	IsPublic   bool      `gorm:"column:is_public"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}

// UserPreference represents user preferences
type UserPreference struct {
	UserID         string `json:"user_id"`
	Theme         string `json:"theme"` // light, dark, system
	Language      string `json:"language"` // en, id
	Timezone      string `json:"timezone"`
	EmailNotif    bool   `json:"email_notif"`
	PushNotif    bool   `json:"push_notif"`
	Digest       string `json:"digest"` // none, daily, weekly
	CompactView   bool   `json:"compact_view"`
	ShowTutor   bool   `json:"show_tutor"`
}

// GetUserPreference gets user preference
func (s *SettingsService) GetUserPreference(userID string) (*UserPreference, error) {
	defaults := UserPreference{
		UserID:      userID,
		Theme:       "system",
		Language:    "id",
		Timezone:    "Asia/Jakarta",
		EmailNotif:  true,
		PushNotif:  true,
		Digest:     "weekly",
		CompactView: false,
		ShowTutor:  true,
	}

	// Load from DB
	prefs := defaults
	rows, err := s.db.Table("user_settings").Where("user_id = ?", userID).Rows()
	if err != nil {
		return &defaults, nil
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		rows.Scan(&key, &value)
		switch key {
		case "theme":
			prefs.Theme = value
		case "language":
			prefs.Language = value
		case "timezone":
			prefs.Timezone = value
		case "email_notif":
			prefs.EmailNotif = value == "true"
		case "push_notif":
			prefs.PushNotif = value == "true"
		case "digest":
			prefs.Digest = value
		case "compact_view":
			prefs.CompactView = value == "true"
		case "show_tutor":
			prefs.ShowTutor = value == "true"
		}
	}

	return &prefs, nil
}

// SetUserPreference sets user preference
func (s *SettingsService) SetUserPreference(userID string, pref UserPreference) error {
	settings := map[string]string{
		"theme":        pref.Theme,
		"language":     pref.Language,
		"timezone":     pref.Timezone,
		"email_notif":  boolToString(pref.EmailNotif),
		"push_notif":   boolToString(pref.PushNotif),
		"digest":      pref.Digest,
		"compact_view": boolToString(pref.CompactView),
		"show_tutor":  boolToString(pref.ShowTutor),
	}

	for key, value := range settings {
		setting := UserSetting{
			ID:        uuid.New().String(),
			UserID:   userID,
			Key:      key,
			Value:    value,
			UpdatedAt: time.Now(),
		}

		s.db.Where("user_id = ? AND key = ?", userID, key).FirstOrCreate(&setting, setting)
	}

	return nil
}

// GetSystemSetting gets system setting
func (s *SettingsService) GetSystemSetting(key string) (string, bool) {
	var setting SystemSetting
	if err := s.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return "", false
	}
	return setting.Value, true
}

// SetSystemSetting sets system setting
func (s *SettingsService) SetSystemSetting(key, value, description, settingType string) error {
	setting := SystemSetting{
		Key:         key,
		Value:       value,
		Description: description,
		Type:        settingType,
		UpdatedAt:  time.Now(),
	}

	return s.db.Where("key = ?", key).FirstOrCreate(&setting, setting).Error
}

// ThemeService manages themes
type ThemeService struct{}

func NewThemeService() *ThemeService {
	return &ThemeService{}
}

type Theme struct {
	ID        string `json:"id"`
	Name     string `json:"name"`
	Primary  string `json:"primary"`
	Secondary string `json:"secondary"`
	Accent   string `json:"accent"`
	Mode     string `json:"mode"` // light, dark
}

// Available themes
var availableThemes = map[string]Theme{
	"default": {
		ID:         "default",
		Name:       "Default",
		Primary:   "#2563eb",
		Secondary: "#64748b",
		Accent:    "#06b6d4",
		Mode:      "light",
	},
	"dark": {
		ID:         "dark",
		Name:       "Dark",
		Primary:   "#3b82f6",
		Secondary: "#475569",
		Accent:    "#22d3ee",
		Mode:      "dark",
	},
	"ocean": {
		ID:         "ocean",
		Name:       "Ocean",
		Primary:   "#0ea5e9",
		Secondary: "#64748b",
		Accent:    "#14b8a6",
		Mode:      "light",
	},
	"forest": {
		ID:         "forest",
		Name:       "Forest",
		Primary:   "#22c55e",
		Secondary: "#78716c",
		Accent:    "#84cc16",
		Mode:      "light",
	},
}

// GetTheme gets theme by ID
func (s *ThemeService) GetTheme(themeID string) *Theme {
	if theme, ok := availableThemes[themeID]; ok {
		return &theme
	}
	defaultTheme := availableThemes["default"]
	return &defaultTheme
}

// ListThemes lists all available themes
func (s *ThemeService) ListThemes() []Theme {
	themes := make([]Theme, 0, len(availableThemes))
	for _, theme := range availableThemes {
		themes = append(themes, theme)
	}
	return themes
}

// LanguageService manages translations
type LanguageService struct {
	translations map[string]map[string]string
}

func NewLanguageService() *LanguageService {
	s := &LanguageService{
		translations: make(map[string]map[string]string),
	}
	s.loadTranslations()
	return s
}

func (s *LanguageService) loadTranslations() {
	// Indonesian translations
	s.translations["id"] = map[string]string{
		"welcome": "Selamat Datang",
		"login":   "Masuk",
		"logout":  "Keluar",
		"home":    "Beranda",
		"profile": "Profil",
		"settings": "Pengaturan",
		"admin":  "Admin",
		"user":   "Pengguna",
		"save":   "Simpan",
		"cancel": "Batal",
		"delete": "Hapus",
		"edit":   "Edit",
		"add":    "Tambah",
		"search": "Cari",
		"filter": "Filter",
		"export":  "Ekspor",
		"import": "Impor",
		"refresh": "Segarkan",
		"loading": "Memuat...",
		"error":  "Kesalahan",
		"success": "Berhasil",
		"warning": "Peringatan",
		"info":   "Informasi",
		"confirm": "Konfirmasi",
		"yes":   "Ya",
		"no":   "Tidak",
		"submit": "Kirim",
		"close": "Tutup",
	}

	// English translations
	s.translations["en"] = map[string]string{
		"welcome": "Welcome",
		"login":  "Login",
		"logout": "Logout",
		"home":  "Home",
		"profile": "Profile",
		"settings": "Settings",
		"admin": "Admin",
		"user":  "User",
		"save":  "Save",
		"cancel": "Cancel",
		"delete": "Delete",
		"edit":  "Edit",
		"add":   "Add",
		"search": "Search",
		"filter": "Filter",
		"export": "Export",
		"import": "Import",
		"refresh": "Refresh",
		"loading": "Loading...",
		"error": "Error",
		"success": "Success",
		"warning": "Warning",
		"info": "Information",
		"confirm": "Confirm",
		"yes": "Yes",
		"no": "No",
		"submit": "Submit",
		"close": "Close",
	}
}

// GetTranslation gets translation
func (s *LanguageService) GetTranslation(lang, key string) string {
	if trans, ok := s.translations[lang]; ok {
		if val, ok := trans[key]; ok {
			return val
		}
	}
	// Fallback to English
	if trans, ok := s.translations["en"]; ok {
		if val, ok := trans[key]; ok {
			return val
		}
	}
	return key
}

// GetTranslations gets all translations for a language
func (s *LanguageService) GetTranslations(lang string) map[string]string {
	if trans, ok := s.translations[lang]; ok {
		return trans
	}
	return s.translations["en"]
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
