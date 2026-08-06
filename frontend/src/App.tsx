import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import { useAuth } from "./contexts/auth";
import AppShell from "./components/layout/AppShell";
import Login from "./pages/Login";
import Servers from "./pages/Servers";
import ServerDetail from "./pages/ServerDetails";
import Reports from "./pages/Reports";
import Imports from "./pages/Imports";
import Users from "./pages/Users";

function Guard({ children, admin = false }: { children: React.ReactNode; admin?: boolean }) { const { isAuthenticated, role } = useAuth(); if (!isAuthenticated) return <Navigate to="/" replace />; return admin && role !== "admin" ? <Navigate to="/servers" replace /> : <>{children}</>; }
function LoginRoute() { const { isAuthenticated } = useAuth(); return isAuthenticated ? <Navigate to="/servers" replace /> : <Login />; }
export default function App() { return <BrowserRouter><AuthProvider><Routes><Route path="/" element={<LoginRoute />} /><Route element={<Guard><AppShell /></Guard>}><Route path="/servers" element={<Servers />} /><Route path="/servers/:id" element={<ServerDetail />} /><Route path="/reports" element={<Guard admin><Reports /></Guard>} /><Route path="/imports" element={<Guard admin><Imports /></Guard>} /><Route path="/users" element={<Guard admin><Users /></Guard>} /></Route><Route path="*" element={<Navigate to="/servers" replace />} /></Routes></AuthProvider></BrowserRouter>; }
