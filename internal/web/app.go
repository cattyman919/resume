package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/cattyman919/autocv/internal/core"
	"github.com/cattyman919/autocv/internal/core/config"
	"github.com/cattyman919/autocv/internal/core/domain"
)

//go:embed static
var staticFiles embed.FS

type App struct {
	core          *core.App
	router        *chi.Mux
	typstMu       sync.Map
	envFieldCache map[string]string
}

func NewApp() (*App, error) {
	coreApp, err := core.NewApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize core app: %w", err)
	}

	app := &App{
		core:          coreApp,
		envFieldCache: make(map[string]string),
	}

	app.loadEnvFieldCache()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", app.handleIndex)
	router.Get("/cv/{typeName}", app.handleCVEditor)

	router.Get("/partials/editor/{typeName}", app.handleEditorPartial)
	router.Get("/partials/personal-info", app.handlePersonalInfoPartial)
	router.Get("/partials/settings", app.handleSettingsPartial)
	router.Get("/partials/experiences/{typeName}", app.handleExperiencesPartial)
	router.Get("/partials/projects/{typeName}", app.handleProjectsPartial)
	router.Get("/partials/skills", app.handleSkillsPartial)
	router.Get("/partials/education", app.handleEducationPartial)
	router.Get("/partials/certificates", app.handleCertificatesPartial)
	router.Get("/partials/awards", app.handleAwardsPartial)
	router.Get("/partials/layout/{typeName}", app.handleLayoutPartial)
	router.Get("/partials/description/{typeName}", app.handleDescriptionPartial)

	router.Post("/api/personal-info", app.handleSavePersonalInfo)
	router.Post("/api/settings", app.handleSaveSettings)
	router.Post("/api/experiences/{typeName}", app.handleSaveExperiences)
	router.Post("/api/experience/{typeName}/{index}", app.handleSaveExperience)
	router.Post("/api/experience/{typeName}/add", app.handleAddExperience)
	router.Post("/api/experience/{typeName}/{index}/delete", app.handleDeleteExperience)
	router.Post("/api/projects/{typeName}", app.handleSaveProjects)
	router.Post("/api/project/{typeName}/{index}", app.handleSaveProject)
	router.Post("/api/project/{typeName}/add", app.handleAddProject)
	router.Post("/api/project/{typeName}/{index}/delete", app.handleDeleteProject)
	router.Post("/api/skills", app.handleSaveSkills)
	router.Post("/api/education", app.handleSaveEducation)
	router.Post("/api/certificates", app.handleSaveCertificates)
	router.Post("/api/awards", app.handleSaveAwards)
	router.Post("/api/layout/{typeName}", app.handleSaveLayout)
	router.Post("/api/description/{typeName}", app.handleSaveDescription)
	router.Post("/api/generate/{typeName}", app.handleGeneratePDF)
	router.Get("/api/pdf/{typeName}", app.handleServePDF)
	router.Head("/api/pdf/{typeName}", app.handleServePDF)

	router.Post("/api/cv-type/create", app.handleCreateCVType)
	router.Post("/api/cv-type/rename", app.handleRenameCVType)
	router.Post("/api/cv-type/delete", app.handleDeleteCVType)

	router.Get("/api/env-resolve", app.handleEnvResolve)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to create static file system: %w", err)
	}
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	app.router = router

	return app, nil
}

func (a *App) loadEnvFieldCache() {
	rawContent, err := os.ReadFile(filepath.Join("config", "general.yaml"))
	if err != nil {
		slog.Warn("Failed to read raw general.yaml for env tracking", "err", err)
		return
	}
	a.envFieldCache = config.TrackEnvFields(string(rawContent))
}

func (a *App) ListenAndServe(addr string) error {
	slog.Info("Starting web server", "addr", addr)
	return http.ListenAndServe(addr, a.router)
}

func (a *App) getTypeMutex(typeName string) *sync.Mutex {
	val, _ := a.typstMu.LoadOrStore(typeName, &sync.Mutex{})
	return val.(*sync.Mutex)
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	typeNames := a.core.GetCVTypeNames()
	if len(typeNames) == 0 {
		http.Error(w, "no CV types found", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cv/"+typeNames[0], http.StatusFound)
}

func (a *App) handleCVEditor(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	typeNames := a.core.GetCVTypeNames()
	component := EditorPage(typeName, typeNames, &a.core.CVConfig.GeneralCfg, &a.core.CVConfig.SettingsCfg, cvType, a.core.TypstAvailable, a.envFieldCache)
	component.Render(r.Context(), w)
}

func (a *App) handleEditorPartial(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	component := EditorPanel(typeName, &a.core.CVConfig.GeneralCfg, &a.core.CVConfig.SettingsCfg, cvType, a.envFieldCache)
	component.Render(r.Context(), w)
}

func (a *App) handlePersonalInfoPartial(w http.ResponseWriter, r *http.Request) {
	component := PersonalInfoSection(&a.core.CVConfig.GeneralCfg, a.envFieldCache)
	component.Render(r.Context(), w)
}

func (a *App) handleSettingsPartial(w http.ResponseWriter, r *http.Request) {
	component := SettingsSection(&a.core.CVConfig.SettingsCfg)
	component.Render(r.Context(), w)
}

func (a *App) handleExperiencesPartial(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	component := ExperiencesSection(typeName, cvType.Experiences, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleProjectsPartial(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	component := ProjectsSection(typeName, cvType.Projects, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSkillsPartial(w http.ResponseWriter, r *http.Request) {
	component := SkillsSection(a.core.CVConfig.GeneralCfg.Skills, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleEducationPartial(w http.ResponseWriter, r *http.Request) {
	component := EducationSection(a.core.CVConfig.GeneralCfg.Educations, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleCertificatesPartial(w http.ResponseWriter, r *http.Request) {
	component := CertificatesSection(a.core.CVConfig.GeneralCfg.Certifications, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleAwardsPartial(w http.ResponseWriter, r *http.Request) {
	component := AwardsSection(a.core.CVConfig.GeneralCfg.Awards, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleLayoutPartial(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	component := LayoutSection(typeName, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleDescriptionPartial(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}
	component := DescriptionSection(typeName, cvType.Description, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSavePersonalInfo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.GeneralCfg.PersonalInfo = domain.PersonalInfo{
		Name:           r.FormValue("name"),
		Email:          r.FormValue("email"),
		Phone:          r.FormValue("phone"),
		Website:        r.FormValue("website"),
		Linkedin:       r.FormValue("linkedin"),
		LinkedinHandle: r.FormValue("linkedin_handle"),
		Github:         r.FormValue("github"),
		GithubHandle:   r.FormValue("github_handle"),
		ProfilePic:     r.FormValue("profile_pic"),
		Location:       r.FormValue("location"),
	}
	if err := a.core.SaveGeneralConfig(); err != nil {
		slog.Error("Failed to save general config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := PersonalInfoForm(&a.core.CVConfig.GeneralCfg, a.envFieldCache)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveSettings(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.SettingsCfg = domain.CVSettings{
		Font:        r.FormValue("font"),
		AccentColor: r.FormValue("accent_color"),
	}
	if err := a.core.SaveSettingsConfig(); err != nil {
		slog.Error("Failed to save settings config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := SettingsForm(&a.core.CVConfig.SettingsCfg)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveExperiences(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var experiences []domain.Experience
	if err := json.Unmarshal([]byte(r.FormValue("data")), &experiences); err != nil {
		http.Error(w, "failed to parse experiences", http.StatusBadRequest)
		return
	}
	cvType.Experiences = experiences

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ExperiencesSection(typeName, cvType.Experiences, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveExperience(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	index := chi.URLParam(r, "index")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var idx int
	if _, err := fmt.Sscanf(index, "%d", &idx); err != nil || idx < 0 || idx >= len(cvType.Experiences) {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	cvType.Experiences[idx].Role = r.FormValue("role")
	cvType.Experiences[idx].JobType = r.FormValue("job_type")
	cvType.Experiences[idx].Company = r.FormValue("company")
	cvType.Experiences[idx].Location = r.FormValue("location")
	cvType.Experiences[idx].Dates = r.FormValue("dates")

	pointsStr := r.FormValue("points")
	cvType.Experiences[idx].Points = []string{}
	if pointsStr != "" {
		for _, p := range strings.Split(pointsStr, "\n") {
			p = strings.TrimSpace(p)
			if p != "" {
				cvType.Experiences[idx].Points = append(cvType.Experiences[idx].Points, p)
			}
		}
	}

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "pdf-regenerate")
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleAddExperience(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	cvType.Experiences = append(cvType.Experiences, domain.Experience{
		Role:    "New Role",
		Company: "Company",
		Dates:   "Date",
		Points:  []string{"Description point"},
	})

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ExperiencesSection(typeName, cvType.Experiences, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleDeleteExperience(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	index := chi.URLParam(r, "index")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var idx int
	if _, err := fmt.Sscanf(index, "%d", &idx); err != nil || idx < 0 || idx >= len(cvType.Experiences) {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}

	cvType.Experiences = append(cvType.Experiences[:idx], cvType.Experiences[idx+1:]...)

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ExperiencesSection(typeName, cvType.Experiences, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveProjects(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var projects []domain.Project
	if err := json.Unmarshal([]byte(r.FormValue("data")), &projects); err != nil {
		http.Error(w, "failed to parse projects", http.StatusBadRequest)
		return
	}
	cvType.Projects = projects

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ProjectsSection(typeName, cvType.Projects, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveProject(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	index := chi.URLParam(r, "index")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var idx int
	if _, err := fmt.Sscanf(index, "%d", &idx); err != nil || idx < 0 || idx >= len(cvType.Projects) {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	cvType.Projects[idx].Name = r.FormValue("name")
	cvType.Projects[idx].Github = r.FormValue("github")
	cvType.Projects[idx].GithubHandle = r.FormValue("github_handle")

	pointsStr := r.FormValue("points")
	cvType.Projects[idx].Points = []string{}
	if pointsStr != "" {
		for _, p := range strings.Split(pointsStr, "\n") {
			p = strings.TrimSpace(p)
			if p != "" {
				cvType.Projects[idx].Points = append(cvType.Projects[idx].Points, p)
			}
		}
	}

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "pdf-regenerate")
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleAddProject(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	cvType.Projects = append(cvType.Projects, domain.Project{
		Name:   "New Project",
		Points: []string{"Description point"},
	})

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ProjectsSection(typeName, cvType.Projects, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	index := chi.URLParam(r, "index")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var idx int
	if _, err := fmt.Sscanf(index, "%d", &idx); err != nil || idx < 0 || idx >= len(cvType.Projects) {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}

	cvType.Projects = append(cvType.Projects[:idx], cvType.Projects[idx+1:]...)

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ProjectsSection(typeName, cvType.Projects, cvType.Layouts)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveSkills(w http.ResponseWriter, r *http.Request) {
	var skills []domain.SkillGroup
	if err := json.Unmarshal([]byte(r.FormValue("data")), &skills); err != nil {
		http.Error(w, "failed to parse skills", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.GeneralCfg.Skills = skills

	if err := a.core.SaveGeneralConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := SkillsSection(a.core.CVConfig.GeneralCfg.Skills, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveEducation(w http.ResponseWriter, r *http.Request) {
	var educations []domain.Education
	if err := json.Unmarshal([]byte(r.FormValue("data")), &educations); err != nil {
		http.Error(w, "failed to parse education", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.GeneralCfg.Educations = educations

	if err := a.core.SaveGeneralConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := EducationSection(a.core.CVConfig.GeneralCfg.Educations, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveCertificates(w http.ResponseWriter, r *http.Request) {
	var certs []domain.Certificate
	if err := json.Unmarshal([]byte(r.FormValue("data")), &certs); err != nil {
		http.Error(w, "failed to parse certificates", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.GeneralCfg.Certifications = certs

	if err := a.core.SaveGeneralConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := CertificatesSection(a.core.CVConfig.GeneralCfg.Certifications, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveAwards(w http.ResponseWriter, r *http.Request) {
	var awards []domain.Award
	if err := json.Unmarshal([]byte(r.FormValue("data")), &awards); err != nil {
		http.Error(w, "failed to parse awards", http.StatusBadRequest)
		return
	}
	a.core.CVConfig.GeneralCfg.Awards = awards

	if err := a.core.SaveGeneralConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := AwardsSection(a.core.CVConfig.GeneralCfg.Awards, nil)
	component.Render(r.Context(), w)
}

func (a *App) handleSaveLayout(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	var layouts []domain.Layout
	if err := json.Unmarshal([]byte(r.FormValue("data")), &layouts); err != nil {
		http.Error(w, "failed to parse layout", http.StatusBadRequest)
		return
	}
	cvType.Layouts = layouts

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "pdf-regenerate")
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleSaveDescription(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	cvType.Description = r.FormValue("description")

	if err := a.core.SaveCVTypeConfig(cvType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "pdf-regenerate")
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleGeneratePDF(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")

	if !a.core.TypstAvailable {
		http.Error(w, "Typst is not installed — PDF generation unavailable", http.StatusServiceUnavailable)
		return
	}

	cvType := a.core.GetCVType(typeName)
	if cvType == nil {
		http.Error(w, "CV type not found", http.StatusNotFound)
		return
	}

	mu := a.getTypeMutex(typeName)
	if !mu.TryLock() {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"status": "already_generating"}`)
		return
	}

	err := a.core.GenerateCV(cvType)
	mu.Unlock()

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "error", "message": %q}`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status": "done"}`)
}

func (a *App) handleServePDF(w http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "typeName")
	pdfPath := filepath.Join("out", a.core.GetPDFPath(typeName))

	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.Error(w, "PDF not found — generate one first", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, pdfPath)
}

func (a *App) handleCreateCVType(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	typeName := r.FormValue("type_name")
	cloneFrom := r.FormValue("clone_from")

	if err := a.core.CreateCVType(typeName, cloneFrom); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.core.ReloadConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	typeNames := a.core.GetCVTypeNames()
	component := CVTypeSelector(typeName, typeNames)
	component.Render(r.Context(), w)
}

func (a *App) handleRenameCVType(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	oldName := r.FormValue("old_name")
	newName := r.FormValue("new_name")

	if err := a.core.RenameCVType(oldName, newName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.core.ReloadConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	typeNames := a.core.GetCVTypeNames()
	component := CVTypeSelector(newName, typeNames)
	component.Render(r.Context(), w)
}

func (a *App) handleDeleteCVType(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	typeName := r.FormValue("type_name")

	if err := a.core.DeleteCVType(typeName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.core.ReloadConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	typeNames := a.core.GetCVTypeNames()
	newSelection := ""
	if len(typeNames) > 0 {
		newSelection = typeNames[0]
	}
	component := CVTypeSelector(newSelection, typeNames)
	component.Render(r.Context(), w)
}

func (a *App) handleEnvResolve(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	if field == "" {
		http.Error(w, "missing field parameter", http.StatusBadRequest)
		return
	}

	if config.HasEnvTemplate(field) {
		varName := config.ExtractEnvVarName(field)
		resolved := os.Getenv(varName)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"is_env":    "true",
			"env_var":   varName,
			"resolved":  resolved,
			"raw":       field,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"is_env":   "false",
		"resolved": field,
		"raw":      field,
	})
}
