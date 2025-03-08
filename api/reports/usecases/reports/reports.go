package reports

import "reports/usecases/reports/dto"

func GetReports() []dto.Report {
	return []dto.Report{
		{"Monthly Report", "This is the content of the monthly report."},
		{"Weekly Summary", "Summary of this week's activities."},
		{"Annual Review", "Detailed review of the entire year's performance."},
		{"Incident Report", "Description of a specific incident that occurred."},
		{"Project Update", "Current status and updates on the project."},
	}
}
