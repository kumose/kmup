{
  _config+:: {
    local c = self,
    dashboardNamePrefix: 'Kmup',
    dashboardTags: ['kmup'],
    dashboardPeriod: 'now-1h',
    dashboardTimezone: 'default',
    dashboardRefresh: '1m',

    // Show issue by repository metrics with format kmup_issues_by_repository{repository="org/repo"} 5.
    // Requires Kmup 1.16.0 with ENABLED_ISSUE_BY_REPOSITORY set to true.
    showIssuesByRepository: true,
    // Show graphs for issue by label metrics with format kmup_issues_by_label{label="bug"} 2.
    // Requires Kmup 1.16.0 with ENABLED_ISSUE_BY_LABEL set to true.
    showIssuesByLabel: true,

    // Requires Kmup 1.16.0.
    showIssuesOpenClose: true,

    // add or remove metrics from dashboard
    kmupStatMetrics:
      [
        {
          name: 'kmup_organizations',
          description: 'Organizations',
        },
        {
          name: 'kmup_teams',
          description: 'Teams',
        },
        {
          name: 'kmup_users',
          description: 'Users',
        },
        {
          name: 'kmup_repositories',
          description: 'Repositories',
        },
        {
          name: 'kmup_milestones',
          description: 'Milestones',
        },
        {
          name: 'kmup_stars',
          description: 'Stars',
        },
        {
          name: 'kmup_releases',
          description: 'Releases',
        },
      ]
      +
      if c.showIssuesOpenClose then
        [
          {
            name: 'kmup_issues_open',
            description: 'Issues opened',
          },
          {
            name: 'kmup_issues_closed',
            description: 'Issues closed',
          },
        ] else
        [
          {
            name: 'kmup_issues',
            description: 'Issues',
          },
        ],
    //set this for using label colors on graphs
    issueLabels: [
      {
        label: 'bug',
        color: '#ee0701',
      },
      {
        label: 'duplicate',
        color: '#cccccc',
      },
      {
        label: 'invalid',
        color: '#e6e6e6',
      },
      {
        label: 'enhancement',
        color: '#84b6eb',
      },
      {
        label: 'help wanted',
        color: '#128a0c',
      },
      {
        label: 'question',
        color: '#cc317c',
      },
    ],
  },
}
