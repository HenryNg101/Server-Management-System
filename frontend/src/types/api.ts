export type Role = "admin" | "user";

export interface Server {
  id: number;
  name: string;
  ipv4: string;
  status: "UP" | "DOWN" | "NO_DATA" | string;
  created_at: string;
  updated_at: string;
}

export interface ServerListResponse {
  servers?: Server[];
  total?: number;
  page?: number;
  page_size?: number;
  total_pages?: number;
}

export interface ReportServer {
  serverID: number;
  uptime: number;
  cpu_usage_avg: number;
  memory_usage_avg: number;
  read_bps_avg: number;
  write_bps_avg: number;
}

export interface Report {
  total_servers: number;
  servers_up: number;
  servers_down: number;
  servers_stats: ReportServer[];
}

export interface ImportJob {
  id: string;
  status: string;
  processed_rows: number;
  success_rows_count: number;
  failed_rows_count: number;
  error?: string;
  failures_file_url?: string;
}

export interface User {
  name: string;
  email: string;
  role: Role;
  created_at: string;
}
