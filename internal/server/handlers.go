package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github-insights-dashboard/api/github/v1"
	"github-insights-dashboard/internal/service"
)

type APIHandler struct {
	githubHandler *GitHubHandler
	log           *log.Helper
}

type GitHubHandler struct {
	svc *service.LuminexService
	log *log.Helper
}

func NewLuminexHandler(svc *service.LuminexService, logger log.Logger) http.Handler {
	
	githubHandler := &GitHubHandler{
		svc: svc,
		log: log.NewHelper(logger),
	}
	
	
	h := &APIHandler{
		githubHandler: githubHandler,
		log:          log.NewHelper(logger),
	}

	
	mux := http.NewServeMux()
	
	
	mux.HandleFunc("/api/health", h.Health)
	
	
	mux.HandleFunc("/api/metrics", githubHandler.GetPRMetrics)
	mux.HandleFunc("/api/monthly-stats", githubHandler.GetMonthlyStats)
	mux.HandleFunc("/api/repo-stats", githubHandler.GetRepoStats)
	mux.HandleFunc("/api/contributor-stats", githubHandler.GetContributorStats)
	mux.HandleFunc("/api/issue-stats", githubHandler.GetIssueStats)
	mux.HandleFunc("/api/detailed-pr-stats", githubHandler.GetDetailedPRStats)

	return mux
}

func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func validateRequest(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return "", "", false
	}

	org := r.URL.Query().Get("org")
	repo := r.URL.Query().Get("repo")
	if org == "" || repo == "" {
		http.Error(w, "Missing required parameters: org, repo", http.StatusBadRequest)
		return "", "", false
	}

	return org, repo, true
}

func handleResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *GitHubHandler) GetPRMetrics(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetPRMetrics(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get PR metrics: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
}

func (h *GitHubHandler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetMonthlyStats(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get monthly stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
}

func (h *GitHubHandler) GetRepoStats(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetRepoStats(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get repo stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
}

func (h *GitHubHandler) GetContributorStats(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetContributorStats(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get contributor stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
}

func (h *GitHubHandler) GetIssueStats(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetIssueStats(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get issue stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
}

func (h *GitHubHandler) GetDetailedPRStats(w http.ResponseWriter, r *http.Request) {
	org, repo, valid := validateRequest(w, r)
	if !valid {
		return
	}

	ctx := r.Context()
	result, err := h.svc.GetDetailedPRStats(ctx, &v1.RepositoryRequest{
		Owner: org,
		Repo:  repo,
	})
	if err != nil {
		h.log.Errorf("Failed to get detailed PR stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handleResponse(w, result)
} 