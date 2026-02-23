document.getElementById("btn-add-category").addEventListener("click", addCategory);
document.getElementById("new-category-input").addEventListener("keydown", (e) => {
  if (e.key === "Enter") addCategory();
});

document.getElementById("btn-add-tag").addEventListener("click", addTag);
document.getElementById("new-tag-input").addEventListener("keydown", (e) => {
  if (e.key === "Enter") addTag();
});

async function loadMetadata() {
  try {
    const md = await api("/metadata");
    renderList("categories-list", md.categories || [], removeCategory);
    renderList("tags-list", md.tags || [], removeTag);
  } catch (err) {
    document.getElementById("categories-list").innerHTML = '<div class="empty-state">Failed to load.</div>';
    document.getElementById("tags-list").innerHTML = '<div class="empty-state">Failed to load.</div>';
  }
}

function renderList(containerId, items, onRemove) {
  const container = document.getElementById(containerId);
  if (!items.length) {
    container.innerHTML = '<div class="empty-state py-4">None yet.</div>';
    return;
  }
  container.innerHTML = items
    .map(
      (item) =>
        `<div class="flex items-center justify-between px-3 py-2 rounded-xl hover:bg-surface0 transition group">
          <span class="text-sm text-text">${escapeHtml(item)}</span>
          <button onclick="event.stopPropagation(); ${onRemove.name}('${escapeHtml(item)}')"
                  class="text-overlay0 hover:text-red transition cursor-pointer opacity-0 group-hover:opacity-100"
                  title="Remove">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>`
    )
    .join("");
}

async function addCategory() {
  const input = document.getElementById("new-category-input");
  const name = input.value.trim();
  if (!name) return;

  try {
    await api("/metadata/categories", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
    input.value = "";
    await loadMetadata();
  } catch (err) {
    alert("Failed to add category: " + err.message);
  }
}

async function removeCategory(name) {
  try {
    await api(`/metadata/categories/${encodeURIComponent(name)}`, { method: "DELETE" });
    await loadMetadata();
  } catch (err) {
    alert("Failed to remove category: " + err.message);
  }
}

async function addTag() {
  const input = document.getElementById("new-tag-input");
  const name = input.value.trim();
  if (!name) return;

  try {
    await api("/metadata/tags", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
    input.value = "";
    await loadMetadata();
  } catch (err) {
    alert("Failed to add tag: " + err.message);
  }
}

async function removeTag(name) {
  try {
    await api(`/metadata/tags/${encodeURIComponent(name)}`, { method: "DELETE" });
    await loadMetadata();
  } catch (err) {
    alert("Failed to remove tag: " + err.message);
  }
}
