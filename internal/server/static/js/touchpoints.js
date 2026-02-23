let allTouchpoints = [];
let metadata = { categories: [], tags: [] };
let editingId = null;

const CATEGORY_COLORS = [
  { bg: "#89b4fa", text: "#1e1e2e" },
  { bg: "#cba6f7", text: "#1e1e2e" },
  { bg: "#a6e3a1", text: "#1e1e2e" },
  { bg: "#fab387", text: "#1e1e2e" },
  { bg: "#f38ba8", text: "#1e1e2e" },
  { bg: "#f9e2af", text: "#1e1e2e" },
  { bg: "#94e2d5", text: "#1e1e2e" },
  { bg: "#89dceb", text: "#1e1e2e" },
  { bg: "#74c7ec", text: "#1e1e2e" },
  { bg: "#f5c2e7", text: "#1e1e2e" },
  { bg: "#b4befe", text: "#1e1e2e" },
  { bg: "#f2cdcd", text: "#1e1e2e" },
];

function getCategoryColor(category) {
  const idx = metadata.categories.indexOf(category);
  return CATEGORY_COLORS[idx >= 0 ? idx % CATEGORY_COLORS.length : 0];
}

function populateSelect(elementId, items, placeholder) {
  const sel = document.getElementById(elementId);
  sel.innerHTML = `<option value="">${placeholder}</option>`;
  items.forEach((item) => {
    const opt = document.createElement("option");
    opt.value = item;
    opt.textContent = item;
    sel.appendChild(opt);
  });
}

async function loadTouchpoints() {
  try {
    [allTouchpoints, metadata] = await Promise.all([
      api("/touchpoints"),
      api("/metadata"),
    ]);
  } catch (err) {
    console.error("Failed to load data:", err);
    allTouchpoints = [];
    metadata = { categories: [], tags: [] };
  }
  populateSelect("tp-category", metadata.categories, "Select category...");
  renderTagSelector();
  populateSelect("filter-category", metadata.categories, "All Categories");
  renderFilterTags();
  renderTouchpointList();
}

function renderTagSelector() {
  const container = document.getElementById("tp-tags");
  container.innerHTML = "";
  metadata.tags.forEach((tag) => {
    const chip = document.createElement("span");
    chip.className = "tag-chip";
    chip.textContent = tag;
    chip.dataset.tag = tag;
    chip.addEventListener("click", () => chip.classList.toggle("selected"));
    container.appendChild(chip);
  });
}

function renderFilterTags() {
  const container = document.getElementById("filter-tags");
  container.innerHTML = "";
  metadata.tags.forEach((tag) => {
    const chip = document.createElement("span");
    chip.className = "filter-tag-chip";
    chip.textContent = tag;
    chip.dataset.tag = tag;
    chip.addEventListener("click", () => {
      chip.classList.toggle("selected");
      renderTouchpointList();
      if (typeof initDashboard === "function") initDashboard();
    });
    container.appendChild(chip);
  });
}

function getFilterTags() {
  return Array.from(document.querySelectorAll("#filter-tags .filter-tag-chip.selected")).map(
    (el) => el.dataset.tag
  );
}

function getFilterCategory() {
  return document.getElementById("filter-category").value;
}

function getFilteredTouchpoints() {
  let tps = [...allTouchpoints];
  const cat = getFilterCategory();
  const tags = getFilterTags();
  if (cat) tps = tps.filter((tp) => tp.category === cat);
  if (tags.length) tps = tps.filter((tp) => tags.some((t) => (tp.tags || []).includes(t)));
  return tps;
}

function getSelectedTags() {
  return Array.from(document.querySelectorAll("#tp-tags .tag-chip.selected")).map(
    (el) => el.dataset.tag
  );
}

function setSelectedTags(tags) {
  document.querySelectorAll("#tp-tags .tag-chip").forEach((chip) => {
    chip.classList.toggle("selected", tags.includes(chip.dataset.tag));
  });
}

function renderTouchpointList() {
  const container = document.getElementById("touchpoint-list");
  const filtered = getFilteredTouchpoints();

  if (!filtered.length) {
    container.innerHTML = '<div class="empty-state">No touchpoints yet. Click "New Touchpoint" to add one.</div>';
    return;
  }

  const sorted = [...filtered].sort((a, b) => new Date(b.date) - new Date(a.date));

  container.innerHTML = sorted
    .map((tp) => {
      const cc = getCategoryColor(tp.category);
      const safeId = escapeHtml(tp.id);
      return `
      <div class="tp-row" data-id="${safeId}">
        <span class="text-xs text-overlay0 w-16 flex-shrink-0">${formatDateShort(tp.date)}</span>
        <span class="cat-badge flex-shrink-0" style="background:${cc.bg};color:${cc.text}">${escapeHtml(tp.category)}</span>
        <span class="text-sm text-text flex-1 truncate">${escapeHtml(tp.description)}</span>
        ${(tp.tags || []).map((t) => `<span class="tag-chip" style="cursor:default;font-size:0.6rem;padding:1px 7px">${escapeHtml(t)}</span>`).join("")}
        ${tp.people_involved && tp.people_involved.length ? `<span class="text-xs text-overlay0 flex-shrink-0">${tp.people_involved.map(escapeHtml).join(", ")}</span>` : ""}
        ${tp.url ? `<a href="${escapeHtml(tp.url)}" target="_blank" rel="noopener" class="text-blue text-xs hover:underline flex-shrink-0">Link</a>` : ""}
        <div class="tp-actions">
          <button onclick="startEdit('${safeId}')" class="text-xs text-blue hover:text-lavender transition cursor-pointer px-1">Edit</button>
          <button onclick="deleteTouchpoint('${safeId}')" class="text-xs text-red hover:brightness-125 transition cursor-pointer px-1">Delete</button>
        </div>
      </div>`;
    })
    .join("");
}

const modal = document.getElementById("tp-modal");
const modalDialog = document.getElementById("tp-modal-dialog");

function openModal() {
  modal.classList.add("show");
  requestAnimationFrame(() => {
    modalDialog.style.transform = "scale(1)";
    modalDialog.style.opacity = "1";
  });
}

function closeModal() {
  modalDialog.style.transform = "scale(0.95)";
  modalDialog.style.opacity = "0";
  setTimeout(() => {
    modal.classList.remove("show");
    cancelEdit();
  }, 200);
}

document.getElementById("btn-new-touchpoint").addEventListener("click", () => {
  cancelEdit();
  openModal();
});

document.getElementById("modal-close").addEventListener("click", closeModal);
document.getElementById("tp-modal-backdrop").addEventListener("click", closeModal);

document.addEventListener("keydown", (e) => {
  if (e.key === "Escape" && modal.classList.contains("show")) closeModal();
});

document.getElementById("touchpoint-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const input = {
    description: document.getElementById("tp-description").value.trim(),
    category: document.getElementById("tp-category").value,
    tags: getSelectedTags(),
    people_involved: document
      .getElementById("tp-people")
      .value.split(",")
      .map((s) => s.trim())
      .filter(Boolean),
    url: document.getElementById("tp-url").value.trim(),
  };

  try {
    if (editingId) {
      await api(`/touchpoints/${editingId}`, {
        method: "PUT",
        body: JSON.stringify(input),
      });
    } else {
      await api("/touchpoints", {
        method: "POST",
        body: JSON.stringify(input),
      });
    }
    document.getElementById("touchpoint-form").reset();
    setSelectedTags([]);
    closeModal();
    await loadTouchpoints();
    if (typeof initDashboard === "function") initDashboard();
  } catch (err) {
    alert("Error: " + err.message);
  }
});

function startEdit(id) {
  const tp = allTouchpoints.find((t) => t.id === id);
  if (!tp) return;

  editingId = id;
  document.getElementById("form-title").textContent = "Edit Touchpoint";
  document.getElementById("form-submit").textContent = "Update";
  document.getElementById("form-cancel").classList.remove("hidden");
  document.getElementById("edit-id").value = id;
  document.getElementById("tp-description").value = tp.description;
  document.getElementById("tp-category").value = tp.category;
  setSelectedTags(tp.tags || []);
  document.getElementById("tp-people").value = (tp.people_involved || []).join(", ");
  document.getElementById("tp-url").value = tp.url || "";

  openModal();
}

function cancelEdit() {
  editingId = null;
  document.getElementById("form-title").textContent = "New Touchpoint";
  document.getElementById("form-submit").textContent = "Add Touchpoint";
  document.getElementById("form-cancel").classList.add("hidden");
  document.getElementById("edit-id").value = "";
  document.getElementById("touchpoint-form").reset();
  setSelectedTags([]);
}

document.getElementById("form-cancel").addEventListener("click", () => {
  closeModal();
});

async function deleteTouchpoint(id) {
  if (!confirm("Delete this touchpoint?")) return;
  try {
    await api(`/touchpoints/${id}`, { method: "DELETE" });
    await loadTouchpoints();
    if (typeof initDashboard === "function") initDashboard();
  } catch (err) {
    alert("Error: " + err.message);
  }
}

document.getElementById("filter-category").addEventListener("change", () => {
  renderTouchpointList();
  if (typeof initDashboard === "function") initDashboard();
});
