import { useEffect, useState, type FormEvent } from "react";
import api from "../api/client";
import type { User } from "../types/api";

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [show, setShow] = useState(false);
  const [error, setError] = useState("");
  const load = async () => {
    try {
      const { data } = await api.get<User[]>("/users");
      setUsers(data);
      setError("");
    } catch {
      setError("Unable to load users. Check that the users service is running and your account is an administrator.");
    }
  };
  // This initial fetch synchronizes remote user data with this page.
  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => { void load(); }, []);
  return <><section className="page-intro split-intro"><div><span className="eyebrow">ACCESS CONTROL</span><h2>Users</h2><p>Manage people with access to the server operations workspace.</p></div><button className="primary-button" onClick={() => setShow(true)}>+ Add user</button></section><section className="panel">{error && <div className="error-message">{error}</div>}<div className="table-wrap"><table><thead><tr><th>USER</th><th>ROLE</th><th>ADDED</th></tr></thead><tbody>{users.length === 0 ? <tr><td colSpan={3} className="empty-cell">No users found.</td></tr> : users.map((user) => <tr key={user.email}><td><strong>{user.name}</strong><small>{user.email}</small></td><td><span className={`role-chip ${user.role}`}>{user.role}</span></td><td>{new Date(user.created_at).toLocaleDateString()}</td></tr>)}</tbody></table></div></section>{show && <CreateUser close={() => setShow(false)} saved={() => { setShow(false); void load(); }} />}</>;
}

function CreateUser({ close, saved }: { close(): void; saved(): void }) {
  const [name, setName] = useState(""); const [email, setEmail] = useState(""); const [password, setPassword] = useState(""); const [role, setRole] = useState("user"); const [error, setError] = useState("");
  const submit = async (e: FormEvent) => { e.preventDefault(); try { await api.post("/users/", { name, email, password, role }); saved(); } catch { setError("The user could not be created. Check the values and try again."); } };
  return <div className="modal-backdrop"><form className="modal" onSubmit={submit}><button className="close-button" type="button" onClick={close}>×</button><span className="eyebrow">NEW USER</span><h2>Add workspace user</h2>{error && <div className="error-message">{error}</div>}<label>FULL NAME<input value={name} onChange={(e) => setName(e.target.value)} required /></label><label>EMAIL<input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required /></label><label>TEMPORARY PASSWORD<input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required /></label><label>ROLE<select value={role} onChange={(e) => setRole(e.target.value)}><option value="user">User</option><option value="admin">Administrator</option></select></label><div className="modal-actions"><button type="button" className="secondary-button" onClick={close}>Cancel</button><button className="primary-button">Create user</button></div></form></div>;
}
