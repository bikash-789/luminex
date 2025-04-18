package models

type PRStats struct {
	AvgMergeTime string `json:"avg_merge_time"`
	OpenPRs      int    `json:"open_prs"`
	MergedLast7  int    `json:"merged_last_week"`
}

type MonthData struct {
	Month     string `json:"month"`
	OpenPRs   int    `json:"open_prs"`
	MergedPRs int    `json:"merged_prs"`
	Issues    int    `json:"issues"`
}

type MonthlyStats struct {
	Data []MonthData `json:"data"`
}

type RepoStats struct {
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Watchers    int    `json:"watchers"`
	Size        int    `json:"size_kb"`
	LastUpdated string `json:"last_updated"`
	Language    string `json:"language"`
}

type ContributorStats struct {
	TotalContributors int             `json:"total_contributors"`
	TopContributors   []ContributorData `json:"top_contributors"`
	CommitsLast30Days int             `json:"commits_last_30_days"`
	AvgCommitsPerDay  float64         `json:"avg_commits_per_day"`
}

type ContributorData struct {
	Username     string `json:"username"`
	Contributions int    `json:"contributions"`
	AvatarURL    string `json:"avatar_url"`
}

type IssueStats struct {
	OpenIssues        int    `json:"open_issues"`
	ClosedIssues      int    `json:"closed_issues"`
	AvgResolutionTime string `json:"avg_resolution_time"`
	OldestOpenIssue   string `json:"oldest_open_issue"`
	IssuesLast30Days  int    `json:"issues_last_30_days"`
}

type DetailedPRStats struct {
	PRStats
	SmallPRs         int `json:"small_prs"`
	MediumPRs        int `json:"medium_prs"`
	LargePRs         int `json:"large_prs"`
	AvgComments      int `json:"avg_comments"`
	PRsWithoutReview int `json:"prs_without_review"`
} 