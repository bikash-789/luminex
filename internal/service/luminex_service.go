package service

import (
	"context"

	pb "github.com/bikash-789/comm-protos/luminex/v1"
	"github.com/bikash-789/comm-protos/luminex/v1/request"
	"github.com/bikash-789/comm-protos/luminex/v1/response"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"luminex-service/internal/biz"
	gh "luminex-service/internal/biz/github"
)

type LuminexService struct {
	pb.UnimplementedLuminexServer
	handler       biz.ILuminexServiceHandler
	githubHandler gh.IGithubHandler
	log           *log.Helper
}

func NewLuminexService(handler biz.ILuminexServiceHandler, githubHandler gh.IGithubHandler, logger log.Logger) *LuminexService {
	return &LuminexService{
		UnimplementedLuminexServer: pb.UnimplementedLuminexServer{},
		handler:                    handler,
		githubHandler:              githubHandler,
		log:                        log.NewHelper(logger),
	}
}

func (s *LuminexService) GetHealth(ctx context.Context, _ *emptypb.Empty) (*response.HealthResponse, error) {
	s.log.WithContext(ctx).Info("API call: GetHealth")
	return &response.HealthResponse{
		Status: "OK 👍",
	}, nil
}

func (s *LuminexService) GetPRMetrics(ctx context.Context, req *request.RepositoryRequest) (*response.PRMetricsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetPRMetrics, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetPRMetrics(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get PR metrics: %v", err)
		return nil, err
	}
	return stats, nil
}

func (s *LuminexService) GetMonthlyStats(ctx context.Context, req *request.RepositoryRequest) (*response.MonthlyStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetMonthlyStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetMonthlyStats(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get monthly stats: %v", err)
		return nil, err
	}

	data := make([]*response.MonthData, 0, len(stats.Data))
	for _, item := range stats.Data {
		data = append(data, &response.MonthData{
			Month:     item.Month,
			OpenPrs:   item.OpenPrs,
			MergedPrs: item.MergedPrs,
			Issues:    item.Issues,
		})
	}

	return &response.MonthlyStatsResponse{
		Data: data,
	}, nil
}

func (s *LuminexService) GetRepoStats(ctx context.Context, req *request.RepositoryRequest) (*response.RepoStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetRepoStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetRepoStats(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get repo stats: %v", err)
		return nil, err
	}
	return stats, nil
}

func (s *LuminexService) GetContributorStats(ctx context.Context, req *request.RepositoryRequest) (*response.ContributorStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetContributorStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetContributorStats(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get contributor stats: %v", err)
		return nil, err
	}

	contributors := make([]*response.ContributorData, 0, len(stats.TopContributors))
	for _, c := range stats.TopContributors {
		contributors = append(contributors, &response.ContributorData{
			Username:      c.Username,
			Contributions: c.Contributions,
			AvatarUrl:     c.AvatarUrl,
		})
	}
	return stats, nil
}

func (s *LuminexService) GetIssueStats(ctx context.Context, req *request.RepositoryRequest) (*response.IssueStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetIssueStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetIssueStats(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get issue stats: %v", err)
		return nil, err
	}
	return stats, nil
}

func (s *LuminexService) GetDetailedPRStats(ctx context.Context, req *request.RepositoryRequest) (*response.DetailedPRStatsResponse, error) {
	s.log.WithContext(ctx).Infof("API call: GetDetailedPRStats, repo: %s/%s", req.Owner, req.Repo)
	stats, err := s.githubHandler.GetDetailedPRMetrics(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get detailed PR stats: %v", err)
		return nil, err
	}
	return stats, nil
}
