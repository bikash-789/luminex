package github

import (
	"context"
	"github.com/bikash-789/comm-protos/luminex/v1/request"
	"github.com/bikash-789/comm-protos/luminex/v1/response"
)

type IGithubHandler interface {
	GetPRMetrics(ctx context.Context, req *request.RepositoryRequest) (*response.PRMetricsResponse, error)
	GetMonthlyStats(ctx context.Context, req *request.RepositoryRequest) (*response.MonthlyStatsResponse, error)
	GetRepoStats(ctx context.Context, req *request.RepositoryRequest) (*response.RepoStatsResponse, error)
	GetContributorStats(ctx context.Context, req *request.RepositoryRequest) (*response.ContributorStatsResponse, error)
	GetIssueStats(ctx context.Context, req *request.RepositoryRequest) (*response.IssueStatsResponse, error)
	GetDetailedPRMetrics(ctx context.Context, req *request.RepositoryRequest) (*response.DetailedPRStatsResponse, error)
}
