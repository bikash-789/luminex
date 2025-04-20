package github

import (
	"context"
	"github.com/bikash-789/comm-protos/luminex/v1/request"
	"github.com/bikash-789/comm-protos/luminex/v1/response"
	"github.com/go-kratos/kratos/v2/log"
	gh "luminex-service/internal/helpers/github"
	"luminex-service/internal/interfaces/entity"
)

type GithubHandler struct {
	githubConfig entity.GithubConfig
	githubHelper *gh.GithubClient
	log          *log.Helper
}

func NewGithubHandler(logger log.Logger, githubConfig entity.GithubConfig) *GithubHandler {
	return &GithubHandler{
		log:          log.NewHelper(logger),
		githubHelper: gh.NewGithubClient(githubConfig),
		githubConfig: githubConfig,
	}
}

func (g *GithubHandler) GetPRMetrics(ctx context.Context, req *request.RepositoryRequest) (*response.PRMetricsResponse, error) {
	g.log.WithContext(ctx).Infof("GetPRMetrics: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetPRMetrics(req)
}

func (g *GithubHandler) GetMonthlyStats(ctx context.Context, req *request.RepositoryRequest) (*response.MonthlyStatsResponse, error) {
	g.log.WithContext(ctx).Infof("GetMonthlyStats: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetMonthlyStats(req)
}

func (g *GithubHandler) GetRepoStats(ctx context.Context, req *request.RepositoryRequest) (*response.RepoStatsResponse, error) {
	g.log.WithContext(ctx).Infof("GetRepoStats: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetRepoStats(req)
}

func (g *GithubHandler) GetContributorStats(ctx context.Context, req *request.RepositoryRequest) (*response.ContributorStatsResponse, error) {
	g.log.WithContext(ctx).Infof("GetContributorStats: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetContributorStats(req)
}

func (g *GithubHandler) GetIssueStats(ctx context.Context, req *request.RepositoryRequest) (*response.IssueStatsResponse, error) {
	g.log.WithContext(ctx).Infof("GetIssueStats: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetIssueStats(req)
}

func (g *GithubHandler) GetDetailedPRMetrics(ctx context.Context, req *request.RepositoryRequest) (*response.DetailedPRStatsResponse, error) {
	g.log.WithContext(ctx).Infof("GetDetailedPRMetrics: owner=%s, repo=%s", req.Owner, req.Repo)
	return g.githubHelper.GetDetailedPRMetrics(req)
}
