package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"github-insights-dashboard/models"
)

type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

func NewGitHubClient(token string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

func (g *GitHubClient) GetPRMetrics(owner, repo string) (*models.PRStats, error) {
	opts := &github.PullRequestListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
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

	return &models.PRStats{
		AvgMergeTime: avg,
		OpenPRs:      openCount,
		MergedLast7:  mergedLast7Days,
	}, nil
}

func (g *GitHubClient) GetMonthlyStats(owner, repo string) (*models.MonthlyStats, error) {
	monthlyData := make([]models.MonthData, 12)
	currentMonth := time.Now()
	
	for i := 0; i < 12; i++ {
		monthTime := currentMonth.AddDate(0, -i, 0)
		monthlyData[11-i] = models.MonthData{
			Month: monthTime.Format("Jan 2006"),
		}
	}
	
	prOpts := &github.PullRequestListOptions{
		State: "all",
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
		
		for i, data := range monthlyData {
			monthTime, _ := time.Parse("Jan 2006", data.Month)
			
			if prCreatedTime.Year() == monthTime.Year() && prCreatedTime.Month() == monthTime.Month() {
				if pr.State != nil && *pr.State == "open" {
					monthlyData[i].OpenPRs++
				}
				
				if pr.MergedAt != nil {
					monthlyData[i].MergedPRs++
				}
				
				break
			}
		}
	}
	
	issueOpts := &github.IssueListByRepoOptions{
		State: "all",
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
		
		for i, data := range monthlyData {
			monthTime, _ := time.Parse("Jan 2006", data.Month)
			
			if issueCreatedTime.Year() == monthTime.Year() && issueCreatedTime.Month() == monthTime.Month() {
				monthlyData[i].Issues++
				break
			}
		}
	}
	
	return &models.MonthlyStats{Data: monthlyData}, nil
}

func (g *GitHubClient) GetRepoStats(owner, repo string) (*models.RepoStats, error) {
	repository, _, err := g.client.Repositories.Get(g.ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository data: %w", err)
	}

	stats := &models.RepoStats{
		Stars:       repository.GetStargazersCount(),
		Forks:       repository.GetForksCount(),
		Watchers:    repository.GetWatchersCount(),
		Size:        repository.GetSize(),
		Language:    repository.GetLanguage(),
	}

	if repository.UpdatedAt != nil {
		stats.LastUpdated = repository.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	return stats, nil
}

func (g *GitHubClient) GetContributorStats(owner, repo string) (*models.ContributorStats, error) {
	contributors, _, err := g.client.Repositories.ListContributors(g.ctx, owner, repo, &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributors: %w", err)
	}

	stats := &models.ContributorStats{
		TotalContributors: len(contributors),
		TopContributors:   make([]models.ContributorData, 0),
	}

	limit := 5
	if len(contributors) < limit {
		limit = len(contributors)
	}
	
	for i := 0; i < limit; i++ {
		contributor := contributors[i]
		stats.TopContributors = append(stats.TopContributors, models.ContributorData{
			Username:      contributor.GetLogin(),
			Contributions: contributor.GetContributions(),
			AvatarURL:     contributor.GetAvatarURL(),
		})
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	commitOpts := &github.CommitsListOptions{
		Since: thirtyDaysAgo,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	commits, _, err := g.client.Repositories.ListCommits(g.ctx, owner, repo, commitOpts)
	if err != nil {
		return stats, nil
	}
	
	stats.CommitsLast30Days = len(commits)
	stats.AvgCommitsPerDay = float64(len(commits)) / 30.0

	return stats, nil
}

func (g *GitHubClient) GetIssueStats(owner, repo string) (*models.IssueStats, error) {
	issueOpts := &github.IssueListByRepoOptions{
		State: "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	issues, _, err := g.client.Issues.ListByRepo(g.ctx, owner, repo, issueOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}
	
	stats := &models.IssueStats{}
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
	
	stats.OpenIssues = openIssues
	stats.ClosedIssues = closedIssues
	stats.IssuesLast30Days = issuesLast30Days
	
	if resolutionCount > 0 {
		avgResolutionTime := totalResolutionTime / time.Duration(resolutionCount)
		stats.AvgResolutionTime = avgResolutionTime.String()
	} else {
		stats.AvgResolutionTime = "N/A"
	}
	
	if oldestOpenIssue != nil {
		stats.OldestOpenIssue = oldestOpenIssue.CreatedAt.Time.Format("2006-01-02")
	} else {
		stats.OldestOpenIssue = "N/A"
	}
	
	return stats, nil
}

func (g *GitHubClient) GetDetailedPRMetrics(owner, repo string) (*models.DetailedPRStats, error) {
	basicStats, err := g.GetPRMetrics(owner, repo)
	if err != nil {
		return nil, err
	}
	
	detailedStats := &models.DetailedPRStats{
		PRStats: *basicStats,
	}
	
	opts := &github.PullRequestListOptions{
		State: "all",
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
			detailedStats.SmallPRs++
		} else if changedFiles <= 30 {
			detailedStats.MediumPRs++
		} else {
			detailedStats.LargePRs++
		}
		
		if pr.ReviewComments != nil && *pr.ReviewComments == 0 && pr.State != nil && *pr.State == "closed" {
			detailedStats.PRsWithoutReview++
		}
		
		if pr.Comments != nil && *pr.Comments > 0 {
			totalComments += *pr.Comments
			prsWithComments++
		}
	}
	
	if prsWithComments > 0 {
		detailedStats.AvgComments = totalComments / prsWithComments
	}
	
	return detailedStats, nil
} 