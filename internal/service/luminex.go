package service

import (
	"context"

	v1 "github-insights-dashboard/api/github/v1"
	"github-insights-dashboard/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ProviderSet = wire.NewSet(
	NewLuminexService,
)

type LuminexService struct {
	v1.UnimplementedLuminexServer
	githubHandler *biz.GithubApp
	log          *log.Helper
}

func NewLuminexService(githubHandler *biz.GithubApp, logger log.Logger) *LuminexService {
	return &LuminexService{
		UnimplementedLuminexServer: v1.UnimplementedLuminexServer{},
		githubHandler:             githubHandler,
		log:                      log.NewHelper(logger),
	}
}

func (s *LuminexService) GetHealth(ctx context.Context, _ *emptypb.Empty) (*v1.HealthResponse, error) {
	s.log.WithContext(ctx).Info("API call: GetHealth")
	return &v1.HealthResponse{
		Status: "ok",
	}, nil
}

func (s *LuminexService) GetPRMetrics(ctx context.Context, req *v1.RepositoryRequest) (*v1.PRMetricsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetPRMetrics, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetPRMetrics(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get PR metrics: %v", err)
		return nil, err
	}

	return &v1.PRMetricsResponse{
		AvgMergeTime: stats.AvgMergeTime,
		OpenPrs:      int32(stats.OpenPRs),
		MergedLast_7: int32(stats.MergedLast7),
	}, nil
}

func (s *LuminexService) GetMonthlyStats(ctx context.Context, req *v1.RepositoryRequest) (*v1.MonthlyStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetMonthlyStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetMonthlyStats(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get monthly stats: %v", err)
		return nil, err
	}

	data := make([]*v1.MonthData, 0, len(stats.Data))
	for _, item := range stats.Data {
		data = append(data, &v1.MonthData{
			Month:     item.Month,
			OpenPrs:   int32(item.OpenPRs),
			MergedPrs: int32(item.MergedPRs),
			Issues:    int32(item.Issues),
		})
	}

	return &v1.MonthlyStatsResponse{
		Data: data,
	}, nil
}

func (s *LuminexService) GetRepoStats(ctx context.Context, req *v1.RepositoryRequest) (*v1.RepoStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetRepoStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetRepoStats(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get repo stats: %v", err)
		return nil, err
	}

	return &v1.RepoStatsResponse{
		Stars:       int32(stats.Stars),
		Forks:       int32(stats.Forks),
		Watchers:    int32(stats.Watchers),
		SizeKb:      int32(stats.Size),
		LastUpdated: stats.LastUpdated,
		Language:    stats.Language,
	}, nil
}

func (s *LuminexService) GetContributorStats(ctx context.Context, req *v1.RepositoryRequest) (*v1.ContributorStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetContributorStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetContributorStats(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get contributor stats: %v", err)
		return nil, err
	}

	contributors := make([]*v1.ContributorData, 0, len(stats.TopContributors))
	for _, c := range stats.TopContributors {
		contributors = append(contributors, &v1.ContributorData{
			Username:      c.Username,
			Contributions: int32(c.Contributions),
			AvatarUrl:     c.AvatarURL,
		})
	}

	return &v1.ContributorStatsResponse{
		TotalContributors:  int32(stats.TotalContributors),
		TopContributors:    contributors,
		CommitsLast_30Days: int32(stats.CommitsLast30Days),
		AvgCommitsPerDay:   float32(stats.AvgCommitsPerDay),
	}, nil
}

func (s *LuminexService) GetIssueStats(ctx context.Context, req *v1.RepositoryRequest) (*v1.IssueStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetIssueStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetIssueStats(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get issue stats: %v", err)
		return nil, err
	}

	return &v1.IssueStatsResponse{
		OpenIssues:        int32(stats.OpenIssues),
		ClosedIssues:      int32(stats.ClosedIssues),
		AvgResolutionTime: stats.AvgResolutionTime,
		OldestOpenIssue:   stats.OldestOpenIssue,
		IssuesLast_30Days: int32(stats.IssuesLast30Days),
	}, nil
}

func (s *LuminexService) GetDetailedPRStats(ctx context.Context, req *v1.RepositoryRequest) (*v1.DetailedPRStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetDetailedPRStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetDetailedPRMetrics(ctx, req.Owner, req.Repo)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get detailed PR stats: %v", err)
		return nil, err
	}

	return &v1.DetailedPRStatsResponse{
		AvgMergeTime:     stats.AvgMergeTime,
		OpenPrs:          int32(stats.OpenPRs),
		MergedLast_7:     int32(stats.MergedLast7),
		SmallPrs:         int32(stats.SmallPRs),
		MediumPrs:        int32(stats.MediumPRs),
		LargePrs:         int32(stats.LargePRs),
		AvgComments:      int32(stats.AvgComments),
		PrsWithoutReview: int32(stats.PRsWithoutReview),
	}, nil
} 