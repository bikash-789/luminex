package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github-insights-dashboard/models"
)

type GithubRepo interface {
	GetPRMetrics(ctx context.Context, owner, repo string) (*models.PRStats, error)
	GetMonthlyStats(ctx context.Context, owner, repo string) (*models.MonthlyStats, error)
	GetRepoStats(ctx context.Context, owner, repo string) (*models.RepoStats, error)
	GetContributorStats(ctx context.Context, owner, repo string) (*models.ContributorStats, error)
	GetIssueStats(ctx context.Context, owner, repo string) (*models.IssueStats, error)
	GetDetailedPRMetrics(ctx context.Context, owner, repo string) (*models.DetailedPRStats, error)
}

type GithubApp struct {
	repo GithubRepo
	log  *log.Helper
}

func NewGithubApp(repo GithubRepo, logger log.Logger) *GithubApp {
	return &GithubApp{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (g *GithubApp) GetPRMetrics(ctx context.Context, owner, repo string) (*models.PRStats, error) {
	g.log.WithContext(ctx).Infof("GetPRMetrics: owner=%s, repo=%s", owner, repo)
	return g.repo.GetPRMetrics(ctx, owner, repo)
}

func (g *GithubApp) GetMonthlyStats(ctx context.Context, owner, repo string) (*models.MonthlyStats, error) {
	g.log.WithContext(ctx).Infof("GetMonthlyStats: owner=%s, repo=%s", owner, repo)
	return g.repo.GetMonthlyStats(ctx, owner, repo)
}

func (g *GithubApp) GetRepoStats(ctx context.Context, owner, repo string) (*models.RepoStats, error) {
	g.log.WithContext(ctx).Infof("GetRepoStats: owner=%s, repo=%s", owner, repo)
	return g.repo.GetRepoStats(ctx, owner, repo)
}

func (g *GithubApp) GetContributorStats(ctx context.Context, owner, repo string) (*models.ContributorStats, error) {
	g.log.WithContext(ctx).Infof("GetContributorStats: owner=%s, repo=%s", owner, repo)
	return g.repo.GetContributorStats(ctx, owner, repo)
}

func (g *GithubApp) GetIssueStats(ctx context.Context, owner, repo string) (*models.IssueStats, error) {
	g.log.WithContext(ctx).Infof("GetIssueStats: owner=%s, repo=%s", owner, repo)
	return g.repo.GetIssueStats(ctx, owner, repo)
}

func (g *GithubApp) GetDetailedPRMetrics(ctx context.Context, owner, repo string) (*models.DetailedPRStats, error) {
	g.log.WithContext(ctx).Infof("GetDetailedPRMetrics: owner=%s, repo=%s", owner, repo)
	return g.repo.GetDetailedPRMetrics(ctx, owner, repo)
} 