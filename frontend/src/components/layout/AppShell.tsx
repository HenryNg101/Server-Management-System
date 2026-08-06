import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/auth";

const nav = [{ to: "/servers", label: "Fleet", icon: "▦" }, { to: "/reports", label: "Reporting", icon: "◔", admin: true }, { to: "/imports", label: "Import jobs", icon: "⇧", admin: true }, { to: "/users", label: "Users", icon: "♙", admin: true }];
export default function AppShell() {
  const { role, logout } = useAuth(); const navigate = useNavigate();
  const handleLogout = async () => { await logout(); navigate("/"); };
  return <div className="app-shell"><aside className="sidebar"><NavLink to="/servers" className="brand"><span className="brand-mark">S</span><span>SMS<span className="brand-accent">-X</span></span></NavLink><span className="product-label">SERVER OPERATIONS</span><nav>{nav.filter((item) => !item.admin || role === "admin").map((item) => <NavLink key={item.to} to={item.to} className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}><span>{item.icon}</span>{item.label}</NavLink>)}</nav><div className="sidebar-bottom"><div className="user-card"><span className="avatar">{role === "admin" ? "A" : "U"}</span><div><strong>{role === "admin" ? "Administrator" : "Operator"}</strong><small>{role} access</small></div></div><button className="nav-link logout" onClick={handleLogout}><span>↪</span>Sign out</button></div></aside><main className="main-content"><header className="topbar"><div><span className="eyebrow">CONTROL CENTER</span><h1>Server management</h1></div><div className="live-state"><i /> System online</div></header><Outlet /></main></div>;
}
