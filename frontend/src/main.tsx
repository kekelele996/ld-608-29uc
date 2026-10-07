import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { routes } from "./router/routes";
import { useSessionStore } from "./stores/SessionStore";
import { LoginModal } from "./components/common/LoginModal";
import { RoleText, type Role } from "./constants/GroundTaskStatus";
import "./styles.css";

function Shell() {
  const [active, setActive] = useState<string>(routes[0]?.route ?? "/dashboard");
  const [showLogin, setShowLogin] = useState(false);
  const user = useSessionStore((s) => s.user);
  const clear = useSessionStore((s) => s.clear);
  const current = routes.find((route) => route.route === active) ?? routes[0];
  const Page = current?.page;

  useEffect(() => {
    if (!user) setShowLogin(true);
  }, [user]);

  return (
    <div className="shell">
      <aside>
        <div className="brand">航空地勤周转<br />保障平台 <span>ground-turn</span></div>
        <nav>
          {routes.map((route) => (
            <button
              key={route.route}
              className={active === route.route ? "active" : ""}
              onClick={() => setActive(route.route)}
            >
              {route.name}
            </button>
          ))}
        </nav>
        <div className="session-box">
          {user ? (
            <>
              <div className="session-user">{user.display_name}</div>
              <div className="session-role">{RoleText[user.role as Role] ?? user.role} · {user.team_code}</div>
              <button className="btn small ghost" onClick={clear}>退出登录</button>
            </>
          ) : (
            <button className="btn small" onClick={() => setShowLogin(true)}>登录</button>
          )}
        </div>
      </aside>
      {Page ? <Page /> : null}
      {showLogin ? <LoginModal onClose={() => setShowLogin(false)} /> : null}
    </div>
  );
}

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Shell />
  </React.StrictMode>
);
