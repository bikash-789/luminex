package github

import (
	"context"
	"fmt"
	"github.com/bikash-789/comm-protos/luminex/v1/request"
	"github.com/bikash-789/comm-protos/luminex/v1/response"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"luminex-service/internal/interfaces/entity"
	"time"
)

type GithubClient struct {
	client *github.Client
	ctx    context.Context
}

func NewGithubClient(config entity.GithubConfig) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token})
	tc := oauth2.NewClient(ctx, ts)

	return &GithubClient{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

func (g *GithubClient) GetPRMetrics(req *request.RepositoryRequest) (*response.PRMetricsResponse, error) {
	opts := &github.PullRequestListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
	owner := req.Owner
	repo := req.Repo
	prs, _, err := g.client.PullRequests.List(g.ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	var totalMergeTime time.Duration
	var mergedCount int
	var openCount int
	var mergedLast7Days int
	now := time.Now()

	for _, pr := range prs {
		if pr.State != nil && *pr.State == "open" {
			openCount++
		}
		if pr.MergedAt != nil && pr.CreatedAt != nil {
			mergeTime := pr.MergedAt.Time.Sub(pr.CreatedAt.Time)
			totalMergeTime += mergeTime
			mergedCount++

			if now.Sub(pr.MergedAt.Time).Hours() < 24*7 {
				mergedLast7Days++
			}
		}
	}

	avg := "N/A"
	if mergedCount > 0 {
		avg = (totalMergeTime / time.Duration(mergedCount)).String()
	}

	return &response.PRMetricsResponse{
		AvgMergeTime: avg,
		OpenPrs:      int32(openCount),
		MergedLast_7: int32(mergedLast7Days),
	}, nil
}

func (g *GithubClient) GetMonthlyStats(req *request.RepositoryRequest) (*response.MonthlyStatsResponse, error) {
	owner := req.Owner
	repo := req.Repo
	data := make([]*response.MonthData, 12)
	currentMonth := time.Now()

	for i := 0; i < 12; i++ {
		monthTime := currentMonth.AddDate(0, -i, 0)
		data[11-i] = &response.MonthData{
			Month: monthTime.Format("Jan 2006"),
		}
	}

	prOpts := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	prs, _, err := g.client.PullRequests.List(g.ctx, owner, repo, prOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs: %w", err)
	}

	for _, pr := range prs {
		if pr.CreatedAt == nil {
			continue
		}

		prCreatedTime := pr.CreatedAt.Time

		for i, item := range data {
			monthTime, _ := time.Parse("Jan 2006", item.Month)

			if prCreatedTime.Year() == monthTime.Year() && prCreatedTime.Month() == monthTime.Month() {
				if pr.State != nil && *pr.State == "open" {
					data[i].OpenPrs++
				}

				if pr.MergedAt != nil {
					data[i].MergedPrs++
				}

				break
			}
		}
	}

	issueOpts := &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, _, err := g.client.Issues.ListByRepo(g.ctx, owner, repo, issueOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}

	for _, issue := range issues {
		if issue.PullRequestLinks != nil || issue.CreatedAt == nil {
			continue
		}

		issueCreatedTime := issue.CreatedAt.Time

		for i, item := range data {
			monthTime, _ := time.Parse("Jan 2006", item.Month)

			if issueCreatedTime.Year() == monthTime.Year() && issueCreatedTime.Month() == monthTime.Month() {
				data[i].Issues++
				break
			}
		}
	}

	return &response.MonthlyStatsResponse{Data: data}, nil
}

func (g *GithubClient) GetRepoStats(req *request.RepositoryRequest) (*response.RepoStatsResponse, error) {
	owner := req.Owner
	repo := req.Repo
	repository, _, err := g.client.Repositories.Get(g.ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository data: %w", err)
	}

	return &response.RepoStatsResponse{
		Stars:       int32(repository.GetStargazersCount()),
		Forks:       int32(repository.GetForksCount()),
		Watchers:    int32(repository.GetWatchersCount()),
		SizeKb:      int32(repository.GetSize()),
		LastUpdated: repository.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Language:    repository.GetLanguage(),
	}, nil
}

func (g *GithubClient) GetContributorStats(req *request.RepositoryRequest) (*response.ContributorStatsResponse, error) {
	owner := req.Owner
	repo := req.Repo
	contributors, _, err := g.client.Repositories.ListContributors(g.ctx, owner, repo, &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributors: %w", err)
	}

	result := &response.ContributorStatsResponse{
		TotalContributors: int32(len(contributors)),
		TopContributors:   make([]*response.ContributorData, 0),
	}

	limit := 5
	if len(contributors) < limit {
		limit = len(contributors)
	}

	for i := 0; i < limit; i++ {
		contributor := contributors[i]
		result.TopContributors = append(result.TopContributors, &response.ContributorData{
			Username:      contributor.GetLogin(),
			Contributions: int32(contributor.GetContributions()),
			AvatarUrl:     contributor.GetAvatarURL(),
		})
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	commitOpts := &github.CommitsListOptions{
		Since:       thirtyDaysAgo,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	commits, _, err := g.client.Repositories.ListCommits(g.ctx, owner, repo, commitOpts)
	if err != nil {
		return result, nil
	}

	result.CommitsLast_30Days = int32(len(commits))
	result.AvgCommitsPerDay = float32(float64(len(commits)) / 30.0)

	return result, nil
}

func (g *GithubClient) GetIssueStats(req *request.RepositoryRequest) (*response.IssueStatsResponse, error) {
	owner := req.Owner
	repo := req.Repo
	issueOpts := &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	issues, _, err := g.client.Issues.ListByRepo(g.ctx, owner, repo, issueOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}

	var openIssues, closedIssues int
	var totalResolutionTime time.Duration
	var resolutionCount int
	var oldestOpenIssue *github.Issue
	var issuesLast30Days int

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, issue := range issues {
		if issue.PullRequestLinks != nil {
			continue
		}

		if issue.GetState() == "open" {
			openIssues++

			if oldestOpenIssue == nil || issue.CreatedAt.Time.Before(oldestOpenIssue.CreatedAt.Time) {
				oldestOpenIssue = issue
			}
		} else {
			closedIssues++

			if issue.CreatedAt != nil && issue.ClosedAt != nil {
				resolutionTime := issue.ClosedAt.Time.Sub(issue.CreatedAt.Time)
				totalResolutionTime += resolutionTime
				resolutionCount++
			}
		}

		if issue.CreatedAt != nil && issue.CreatedAt.Time.After(thirtyDaysAgo) {
			issuesLast30Days++
		}
	}

	result := &response.IssueStatsResponse{
		OpenIssues:        int32(openIssues),
		ClosedIssues:      int32(closedIssues),
		IssuesLast_30Days: int32(issuesLast30Days),
	}

	if resolutionCount > 0 {
		avgResolutionTime := totalResolutionTime / time.Duration(resolutionCount)
		result.AvgResolutionTime = avgResolutionTime.String()
	} else {
		result.AvgResolutionTime = "N/A"
	}

	if oldestOpenIssue != nil {
		result.OldestOpenIssue = oldestOpenIssue.CreatedAt.Time.Format("2006-01-02")
	} else {
		result.OldestOpenIssue = "N/A"
	}

	return result, nil
}

func (g *GithubClient) GetDetailedPRMetrics(req *request.RepositoryRequest) (*response.DetailedPRStatsResponse, error) {
	owner := req.Owner
	repo := req.Repo
	basicStatsResp, err := g.GetPRMetrics(req)
	if err != nil {
		return nil, err
	}

	result := &response.DetailedPRStatsResponse{
		AvgMergeTime: basicStatsResp.AvgMergeTime,
		OpenPrs:      basicStatsResp.OpenPrs,
		MergedLast_7: basicStatsResp.MergedLast_7,
	}

	opts := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	prs, _, err := g.client.PullRequests.List(g.ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs: %w", err)
	}

	var totalComments int
	var prsWithComments int

	for _, pr := range prs {
		if pr.ChangedFiles == nil {
			continue
		}

		changedFiles := *pr.ChangedFiles
		if changedFiles < 10 {
			result.SmallPrs++
		} else if changedFiles <= 30 {
			result.MediumPrs++
		} else {
			result.LargePrs++
		}

		if pr.ReviewComments != nil && *pr.ReviewComments == 0 && pr.State != nil && *pr.State == "closed" {
			result.PrsWithoutReview++
		}

		if pr.Comments != nil && *pr.Comments > 0 {
			totalComments += *pr.Comments
			prsWithComments++
		}
	}

	if prsWithComments > 0 {
		result.AvgComments = int32(totalComments / prsWithComments)
	}

	return result, nil
}
