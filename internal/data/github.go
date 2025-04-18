package data

import (
	"context"
	"fmt"
	"time"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"github-insights-dashboard/models"
)

type githubRepo struct {
	client *github.Client
	ctx    context.Context
	log    *log.Helper
	data   *Data
}


func newGithubRepo(token string, logger log.Logger, data *Data) *githubRepo {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &githubRepo{
		client: github.NewClient(tc),
		ctx:    ctx,
		log:    log.NewHelper(logger),
		data:   data,
	}
}


func (g *githubRepo) GetPRMetrics(ctx context.Context, owner, repo string) (*models.PRStats, error) {
	g.log.Infof("GetPRMetrics: owner=%s, repo=%s", owner, repo)
	
	opts := &github.PullRequestListOptions{
		State: "all", 
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	prs, _, err := g.client.PullRequests.List(ctx, owner, repo, opts)
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

	var avgMergeTime string
	if mergedCount > 0 {
		avgMergeTime = (totalMergeTime / time.Duration(mergedCount)).String()
	} else {
		avgMergeTime = "N/A"
	}

	return &models.PRStats{
		AvgMergeTime: avgMergeTime,
		OpenPRs:      openCount,
		MergedLast7:  mergedLast7Days,
	}, nil
}


func (g *githubRepo) GetMonthlyStats(ctx context.Context, owner, repo string) (*models.MonthlyStats, error) {
	g.log.Infof("GetMonthlyStats: owner=%s, repo=%s", owner, repo)
	
	
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
	prs, _, err := g.client.PullRequests.List(ctx, owner, repo, prOpts)
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
	issues, _, err := g.client.Issues.ListByRepo(ctx, owner, repo, issueOpts)
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


func (g *githubRepo) GetRepoStats(ctx context.Context, owner, repo string) (*models.RepoStats, error) {
	g.log.Infof("GetRepoStats: owner=%s, repo=%s", owner, repo)
	
	repository, _, err := g.client.Repositories.Get(ctx, owner, repo)
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


func (g *githubRepo) GetContributorStats(ctx context.Context, owner, repo string) (*models.ContributorStats, error) {
	g.log.Infof("GetContributorStats: owner=%s, repo=%s", owner, repo)
	
	
	contributors, _, err := g.client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
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
	
	commits, _, err := g.client.Repositories.ListCommits(ctx, owner, repo, commitOpts)
	if err != nil {
		
		return stats, nil
	}
	
	
	stats.CommitsLast30Days = len(commits)
	stats.AvgCommitsPerDay = float64(len(commits)) / 30.0

	return stats, nil
}


func (g *githubRepo) GetIssueStats(ctx context.Context, owner, repo string) (*models.IssueStats, error) {
	g.log.Infof("GetIssueStats: owner=%s, repo=%s", owner, repo)
	
	issueOpts := &github.IssueListByRepoOptions{
		State: "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	issues, _, err := g.client.Issues.ListByRepo(ctx, owner, repo, issueOpts)
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


func (g *githubRepo) GetDetailedPRMetrics(ctx context.Context, owner, repo string) (*models.DetailedPRStats, error) {
	g.log.Infof("GetDetailedPRMetrics: owner=%s, repo=%s", owner, repo)
	
	
	basicStats, err := g.GetPRMetrics(ctx, owner, repo)
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
	
	prs, _, err := g.client.PullRequests.List(ctx, owner, repo, opts)
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

func (r *githubRepo) GetClient() *github.Client {
	return r.client
} 