package monitoring

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

type ReportView struct {
	PeriodStart string
	PeriodEnd   string

	TotalServers int64
	ServersUp    int64
	ServersDown  int64

	Stats []*ServerOverview
}

func buildReportHTML(report *Report, startTime, endTime time.Time) string {
	view := ReportView{
		PeriodStart:  startTime.Format("2006-01-02 15:04:05"),
		PeriodEnd:    endTime.Format("2006-01-02 15:04:05"),
		TotalServers: report.TotalServers,
		ServersUp:    report.ServersUp,
		ServersDown:  report.ServersDown,
		Stats:        report.Stats,
	}

	const tpl = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style>
    body {
      margin: 0;
      padding: 24px;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif;
      background: #f6f8fb;
      color: #1f2937;
    }
    .container {
      max-width: 1200px;
      margin: 0 auto;
      background: #ffffff;
      border-radius: 16px;
      padding: 28px;
      box-shadow: 0 8px 30px rgba(15, 23, 42, 0.08);
    }
    .header h1 {
      margin: 0 0 8px 0;
      font-size: 28px;
      color: #0f172a;
    }
    .muted {
      color: #64748b;
      font-size: 14px;
    }
    .cards {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 14px;
      margin: 22px 0 28px;
    }
    .card {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 14px;
      padding: 16px;
    }
    .card .label {
      font-size: 13px;
      color: #64748b;
      margin-bottom: 8px;
    }
    .card .value {
      font-size: 28px;
      font-weight: 700;
      color: #0f172a;
    }
    .section-title {
      margin: 26px 0 12px;
      font-size: 18px;
      color: #0f172a;
    }
    .note {
      margin: 10px 0 18px;
      color: #475569;
      font-size: 13px;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      overflow: hidden;
      border-radius: 14px;
      border: 1px solid #e2e8f0;
    }
    thead th {
      background: #f1f5f9;
      text-align: left;
      font-size: 13px;
      color: #334155;
      padding: 12px 14px;
      border-bottom: 1px solid #e2e8f0;
      white-space: nowrap;
      vertical-align: bottom;
    }
    tbody td {
      padding: 12px 14px;
      border-bottom: 1px solid #e2e8f0;
      vertical-align: top;
      font-size: 14px;
    }
    tbody tr:nth-child(even) {
      background: #fbfdff;
    }
    .badge {
      display: inline-block;
      padding: 4px 10px;
      border-radius: 999px;
      font-size: 12px;
      font-weight: 700;
    }
    .good { background: #dcfce7; color: #166534; }
    .warn { background: #fef3c7; color: #92400e; }
    .bad  { background: #fee2e2; color: #991b1b; }
    .metric {
      font-weight: 600;
      color: #111827;
    }
    .small {
      font-size: 12px;
      color: #64748b;
      margin-top: 4px;
      line-height: 1.35;
    }
    .window {
      min-width: 210px;
    }
    @media (max-width: 900px) {
      .cards { grid-template-columns: repeat(2, minmax(0, 1fr)); }
    }
    @media (max-width: 640px) {
      .cards { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Server Monitoring Report</h1>
      <div class="muted">Period: {{.PeriodStart}} → {{.PeriodEnd}}</div>
    </div>

    <div class="cards">
      <div class="card">
        <div class="label">Total servers</div>
        <div class="value">{{.TotalServers}}</div>
      </div>
      <div class="card">
        <div class="label">Up</div>
        <div class="value">{{.ServersUp}}</div>
      </div>
      <div class="card">
        <div class="label">Down / no data</div>
        <div class="value">{{.ServersDown}}</div>
      </div>
      <div class="card">
        <div class="label">Reported servers</div>
        <div class="value">{{len .Stats}}</div>
      </div>
    </div>

    <div class="section-title">Servers by lowest uptime</div>
    <div class="note">
      Uptime is calculated from push telemetry only. For each server, the report window is clamped to the server’s own creation time so a newly created server is not penalized for time before it existed.
    </div>

    <table>
      <thead>
        <tr>
          <th>Server</th>
          <th>Actual monitoring window</th>
          <th>Uptime</th>
          <th>CPU avg</th>
          <th>Memory usage avg</th>
          <th>Working set avg</th>
          <th>RSS avg</th>
          <th>IO read avg</th>
          <th>IO write avg</th>
          <th>PIDs avg</th>
          <th>OOM events</th>
          <th>OOM kills</th>
        </tr>
      </thead>
      <tbody>
        {{range .Stats}}
        <tr>
          <td>
            <div class="metric">#{{.ServerID}}</div>
          </td>
          <td class="window">
            <div>{{formatTime .ActualStart}}</div>
            <div class="small">→ {{formatTime .ActualEnd}}</div>
            <div class="small">Duration: {{duration .ActualStart .ActualEnd}}</div>
          </td>
          <td>
            <span class="badge {{uptimeClass .Uptime}}">{{formatPct .Uptime}}</span>
            <div class="small">Availability in the observed window</div>
          </td>
          <td>
            <span class="metric">{{formatFloat .CPUUsageAvg}}</span>
            <div class="small">%</div>
          </td>
          <td>
            <span class="metric">{{formatBytes .MemoryUsageAvg}}</span>
            <div class="small">MiB</div>
          </td>
          <td>
            <span class="metric">{{formatBytes .MemoryWorkingSetAvg}}</span>
            <div class="small">MiB</div>
          </td>
          <td>
            <span class="metric">{{formatBytes .MemoryRSSAvg}}</span>
            <div class="small">MiB</div>
          </td>
          <td>
            <span class="metric">{{formatRate .ReadBPSAvg}}</span>
            <div class="small">B/s</div>
          </td>
          <td>
            <span class="metric">{{formatRate .WriteBPSAvg}}</span>
            <div class="small">B/s</div>
          </td>
          <td>
            <span class="metric">{{formatFloat .PIDsAvg}}</span>
            <div class="small">processes</div>
          </td>
          <td>
            <span class="metric">{{formatFloat .OOMEventsTotal}}</span>
          </td>
          <td>
            <span class="metric">{{formatFloat .OOMKillsTotal}}</span>
          </td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>
</body>
</html>`

	funcMap := template.FuncMap{
		"formatPct": func(v float64) string {
			return fmt.Sprintf("%.2f%%", v)
		},
		"formatFloat": func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
		"formatRate": func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
		"formatBytes": func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
		"formatTime": func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05")
		},
		"duration": func(start, end time.Time) string {
			if start.IsZero() || end.IsZero() || end.Before(start) {
				return "-"
			}
			return end.Sub(start).Truncate(time.Second).String()
		},
		"uptimeClass": func(v float64) string {
			switch {
			case v >= 95:
				return "good"
			case v >= 80:
				return "warn"
			default:
				return "bad"
			}
		},
	}

	t := template.Must(template.New("report").Funcs(funcMap).Parse(tpl))

	var b bytes.Buffer
	if err := t.Execute(&b, view); err != nil {
		return "<p>failed to render report</p>"
	}
	return b.String()
}
