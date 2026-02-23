document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".tab").forEach((btn) => {
    btn.addEventListener("click", () => {
      document.querySelectorAll(".tab").forEach((b) => b.classList.remove("active"));
      document.querySelectorAll(".tab-content").forEach((s) => (s.style.display = "none"));
      btn.classList.add("active");
      const section = document.getElementById(btn.dataset.tab);
      if (section) {
        section.style.display = "";
        section.classList.add("active");
      }

      if (btn.dataset.tab === "dashboard") initDashboard();
      if (btn.dataset.tab === "reports") loadReportsList();
    });
  });

  mermaid.initialize({
    startOnLoad: false,
    theme: "base",
    themeVariables: {
      darkMode: true,
      background: "#1e1e2e",
      primaryColor: "#89b4fa",
      primaryTextColor: "#1e1e2e",
      primaryBorderColor: "#45475a",
      secondaryColor: "#cba6f7",
      secondaryTextColor: "#1e1e2e",
      secondaryBorderColor: "#45475a",
      tertiaryColor: "#a6e3a1",
      tertiaryTextColor: "#1e1e2e",
      tertiaryBorderColor: "#45475a",
      lineColor: "#a6adc8",
      textColor: "#cdd6f4",
      mainBkg: "#313244",
      nodeBorder: "#45475a",
      clusterBkg: "#181825",
      clusterBorder: "#45475a",
      titleColor: "#cdd6f4",
      edgeLabelBackground: "#313244",
      nodeTextColor: "#cdd6f4",
    },
  });

  initDashboard();
});

async function api(path, options = {}) {
  const res = await fetch(`/api${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  if (res.status === 204) return null;
  const ct = res.headers.get("Content-Type") || "";
  if (ct.includes("text/markdown")) return res.text();
  return res.json();
}

function formatDate(isoString) {
  return new Date(isoString).toLocaleString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatDateShort(isoString) {
  return new Date(isoString).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

function escapeHtml(str) {
  const div = document.createElement("div");
  div.textContent = str || "";
  return div.innerHTML;
}
