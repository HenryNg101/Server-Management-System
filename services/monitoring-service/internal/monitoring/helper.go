package monitoring

import (
	"fmt"
	"strings"
	"time"
)

func buildReportHTML(report *Report, startTime, endTime time.Time) string {
	var b strings.Builder

	b.WriteString("<h2>Server Report</h2>")
	b.WriteString(fmt.Sprintf(
		"<p><b>Period:</b> %s → %s</p>",
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
	))

	// Summary
	b.WriteString("<h3>Summary</h3>")
	b.WriteString(fmt.Sprintf(`
	<ul>
		<li>Total Servers: %d</li>
		<li>Up: %d</li>
		<li>Down: %d</li>
	</ul>
	`, report.TotalServers, report.ServersUp, report.ServersDown))

	// Table report
	b.WriteString("<h3>Servers stats</h3>")
	b.WriteString(`<table border="1" cellpadding="5" cellspacing="0">
	<tr>
		<th>Server ID</th>
		<th>Uptime</th>
		<th>CPU Avg</th>
		<th>Memory Avg</th>
	</tr>`)

	for _, s := range report.Stats {
		b.WriteString(fmt.Sprintf(
			`<tr>
				<td>%d</td>
				<td>%.2f%%</td>
				<td>%.2f</td>
				<td>%.2f</td>
			</tr>`,
			s.ServerID,
			s.Uptime*100,
			s.CPUUsageAvg,
			s.MemoryUsageAvg,
		))
	}

	b.WriteString("</table>")

	return b.String()
}
