document.getElementById("reports-sidebar-toggle").addEventListener("click", () => {
  const sidebar = document.getElementById("reports-sidebar");
  sidebar.classList.add("open");
  document.getElementById("reports-sidebar-close").classList.remove("hidden");
});

document.getElementById("reports-sidebar-close").addEventListener("click", () => {
  const sidebar = document.getElementById("reports-sidebar");
  sidebar.classList.remove("open");
  document.getElementById("reports-sidebar-close").classList.add("hidden");
});

document.getElementById("btn-upload-report").addEventListener("click", () => {
  document.getElementById("report-file-input").click();
});

document.getElementById("report-file-input").addEventListener("change", async (e) => {
  const file = e.target.files[0];
  if (!file) return;
  e.target.value = "";

  const content = await file.text();
  const filename = file.name.endsWith(".md") ? file.name : file.name + ".md";

  try {
    await api("/reports", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ filename, content }),
    });
    await loadReportsList();
    viewReport(filename, document.querySelector(`[data-filename="${CSS.escape(filename)}"]`));
  } catch (err) {
    alert("Upload failed: " + err.message);
  }
});

async function loadReportsList() {
  const container = document.getElementById("reports-list");

  try {
    const names = await api("/reports");
    if (!names.length) {
      container.innerHTML = '<div class="empty-state py-8">No reports yet.</div>';
      return;
    }
    container.innerHTML = names
      .map(
        (name) =>
          `<div class="report-item" data-filename="${escapeHtml(name)}" onclick="viewReport('${escapeHtml(name)}', this)">${escapeHtml(name.replace(/\.md$/, ""))}</div>`
      )
      .join("");
  } catch (err) {
    container.innerHTML = '<div class="empty-state py-8">Failed to load reports.</div>';
  }
}

async function viewReport(filename, clickedEl) {
  document.querySelectorAll("#reports-list .report-item").forEach((el) => {
    el.classList.remove("active");
  });
  if (clickedEl) clickedEl.classList.add("active");

  const sidebar = document.getElementById("reports-sidebar");
  sidebar.classList.remove("open");
  document.getElementById("reports-sidebar-close").classList.add("hidden");

  const content = document.getElementById("report-content");
  content.innerHTML = `<div class="flex items-center justify-center h-32 text-overlay0 text-sm">Loading...</div>`;

  try {
    const md = await api(`/reports/${encodeURIComponent(filename)}`);

    const renderer = new globalThis.marked.Marked(
      globalThis.markedHighlight.markedHighlight({
        langPrefix: "hljs language-",
        highlight(code, lang) {
          if (lang && hljs.getLanguage(lang)) {
            return hljs.highlight(code, { language: lang }).value;
          }
          return hljs.highlightAuto(code).value;
        },
      })
    );

    const html = renderer.parse(md);
    content.innerHTML = html;

    content.querySelectorAll("pre code.language-mermaid").forEach((block) => {
      const pre = block.parentElement;
      const div = document.createElement("div");
      div.className = "mermaid";
      div.textContent = block.textContent;
      pre.replaceWith(div);
    });

    await mermaid.run({ nodes: content.querySelectorAll(".mermaid") });
  } catch (err) {
    content.innerHTML = `<div class="empty-state">Failed to load report: ${escapeHtml(err.message)}</div>`;
  }
}
