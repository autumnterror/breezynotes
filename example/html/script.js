// --------- Base URL ---------

let API_BASE = "http://127.0.0.1:8080";

async function refreshToken() {
  const url = API_BASE + "/api/auth/token";
  try {
    const res = await fetch(url, { method: "GET", credentials: "include" });
    if (!res.ok) {
      throw new Error(`Refresh failed with status ${res.status}`);
    }
  } catch (err) {
    console.error("Token refresh error:", err);
    throw err;
  }
}

async function apiRequest(path, { method = "GET", body, retryOnAuth = true } = {}) {
  const url = API_BASE + path;

  const options = {
    method,
    headers: {},
    credentials: "include",
  };

  if (body !== undefined) {
    options.headers["Content-Type"] = "application/json";
    options.body = JSON.stringify(body);
  }

  let res;
  try {
    res = await fetch(url, options);
  } catch (err) {
    console.error("Network error:", err);
    throw err;
  }

  if (res.status === 401 && retryOnAuth) {
    try {
      await refreshToken();
      return apiRequest(path, { method, body, retryOnAuth: false });
    } catch (_) {
      // fall through to error handling below
    }
  }

  const text = await res.text();
  let data = null;
  try {
    data = text ? JSON.parse(text) : null;
  } catch (_) {
    data = text;
  }

  if (!res.ok) {
    const errorPayload =
      data && typeof data === "object"
        ? data
        : { error: res.statusText || "Unknown error" };

    console.error("API error", {
      status: res.status,
      statusText: res.statusText,
      error: errorPayload,
      url,
      method,
    });

    throw new Error(
      (errorPayload && errorPayload.error) ||
        `Request failed with status ${res.status}`
    );
  }

  return data;
}

// --------- Навигация по "страницам" ---------

const views = document.querySelectorAll(".view");

function showView(name) {
  views.forEach((v) => {
    const isActive = v.id === `view-${name}`;
    if (isActive) v.classList.add("active");
    else v.classList.remove("active");
  });
}

// --------- Элементы UI ---------

const headerLoginBtn = document.getElementById("header-login-btn");
const headerRegisterBtn = document.getElementById("header-register-btn");
const headerLogoutBtn = document.getElementById("header-logout-btn");

const landingLoginBtn = document.getElementById("landing-login-btn");
const landingRegisterBtn = document.getElementById("landing-register-btn");

const authTabLogin = document.getElementById("auth-tab-login");
const authTabRegister = document.getElementById("auth-tab-register");
const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");
const authTitle = document.getElementById("auth-title");
const authSubtitle = document.getElementById("auth-subtitle");

// навигация-иконки
const navButtons = document.querySelectorAll(".nav-btn");
const navNotesBtn = document.getElementById("nav-notes-btn");
const navTrashBtn = document.getElementById("nav-trash-btn");
const navTagsBtn = document.getElementById("nav-tags-btn");
const navProfileBtn = document.getElementById("nav-profile-btn");

// заметки / корзина
const notesListEl = document.getElementById("notes-list");
const addNoteBtn = document.getElementById("add-note-btn");
const openTrashBtn = document.getElementById("open-trash-btn");
const notesTagFilterSelect = document.getElementById("notes-tag-filter");
const notesMainEl = document.querySelector("#view-notes .notes-main");

const trashListEl = document.getElementById("trash-list");
const backToNotesBtn = document.getElementById("back-to-notes-btn");
const clearTrashBtn = document.getElementById("clear-trash-btn");
const trashView = document.getElementById("view-trash");

// теги
const tagsListEl = document.getElementById("tags-list");
const addTagBtn = document.getElementById("add-tag-btn");

// профиль
const profilePhotoImg = document.getElementById("profile-photo");
const profileAvatarFallback = document.getElementById("profile-avatar-fallback");
const profileLoginSpan = document.getElementById("profile-login");
const profileEmailSpan = document.getElementById("profile-email");
const profileIdSpan = document.getElementById("profile-id");
const profileAboutP = document.getElementById("profile-about");
const profileEmailInline = document.getElementById("profile-email-inline");

const editAboutBtn = document.getElementById("edit-about-btn");
const editEmailBtn = document.getElementById("edit-email-btn");
const editPhotoBtn = document.getElementById("edit-photo-btn");
const changePwBtn = document.getElementById("change-pw-btn");

// --------- Auth UI состояние ---------

let currentAuthMode = "login";
let isAuthenticated = false;
let tagsCache = [];
let currentUser = null;
let selectedNoteId = null;
let selectedNoteData = null;
let restoredNoteFromStorage = false;

const SELECTED_NOTE_STORAGE_KEY = "bn_selected_note_id";

// --------- Вспомогательные функции ---------

function setNavActive(viewName) {
  navButtons.forEach((btn) => {
    const isActive = btn.dataset.view === viewName;
    if (isActive) btn.classList.add("active");
    else btn.classList.remove("active");
  });
}

function setAuthMode(mode) {
  currentAuthMode = mode;

  if (mode === "login") {
    authTabLogin.classList.add("active");
    authTabRegister.classList.remove("active");
    loginForm.classList.remove("hidden");
    registerForm.classList.add("hidden");

    authTitle.textContent = "Войти в аккаунт";
    authSubtitle.textContent = "Введи логин или email и пароль.";
  } else {
    authTabLogin.classList.remove("active");
    authTabRegister.classList.add("active");
    loginForm.classList.add("hidden");
    registerForm.classList.remove("hidden");

    authTitle.textContent = "Создать аккаунт";
    authSubtitle.textContent = "Заполни данные, и мы создадим новый профиль.";
  }

  showView("auth");
  setNavActive("");
}

function setAuthUI(authed) {
  isAuthenticated = authed;
  if (authed) {
    headerLoginBtn.classList.add("hidden");
    headerRegisterBtn.classList.add("hidden");
    headerLogoutBtn.classList.remove("hidden");
  } else {
    headerLoginBtn.classList.remove("hidden");
    headerRegisterBtn.classList.remove("hidden");
    headerLogoutBtn.classList.add("hidden");
  }
}

// --------- Logout ---------

function clearAllCookies() {
  const cookies = document.cookie ? document.cookie.split(";") : [];
  cookies.forEach((cookie) => {
    const eqPos = cookie.indexOf("=");
    const name = eqPos > -1 ? cookie.slice(0, eqPos).trim() : cookie.trim();
    if (!name) return;
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
  });
}

function logout() {
  clearAllCookies();
  setAuthUI(false);
  showView("landing");
  setNavActive("");
}

// --------- Навигация по кнопкам ---------

if (authTabLogin) {
  authTabLogin.addEventListener("click", () => setAuthMode("login"));
}
if (authTabRegister) {
  authTabRegister.addEventListener("click", () => setAuthMode("register"));
}

[headerLoginBtn, landingLoginBtn].forEach((btn) => {
  if (!btn) return;
  btn.addEventListener("click", () => setAuthMode("login"));
});

[headerRegisterBtn, landingRegisterBtn].forEach((btn) => {
  if (!btn) return;
  btn.addEventListener("click", () => setAuthMode("register"));
});

if (headerLogoutBtn) {
  headerLogoutBtn.addEventListener("click", logout);
}

// левое меню-иконки
if (navNotesBtn) {
  navNotesBtn.addEventListener("click", async () => {
    showView("notes");
    setNavActive("notes");
    await loadNotes();
  });
}

if (navTrashBtn) {
  navTrashBtn.addEventListener("click", async () => {
    showView("trash");
    setNavActive("trash");
    await loadTrashNotes();
  });
}

if (navTagsBtn) {
  navTagsBtn.addEventListener("click", async () => {
    showView("tags");
    setNavActive("tags");
    await loadTags();
  });
}

if (navProfileBtn) {
  navProfileBtn.addEventListener("click", async () => {
    showView("profile");
    setNavActive("profile");
    await loadUserData();
  });
}

// старые кнопки перехода между notes/trash
if (openTrashBtn) {
  openTrashBtn.addEventListener("click", async () => {
    showView("trash");
    setNavActive("trash");
    await loadTrashNotes();
  });
}

if (backToNotesBtn) {
  backToNotesBtn.addEventListener("click", async () => {
    showView("notes");
    setNavActive("notes");
    await loadNotes();
  });
}

// --------- Форматирование дат ---------

function formatDate(ts) {
  if (!ts) return "";
  try {
    const d = new Date(ts * 1000);
    return d.toLocaleString("ru-RU");
  } catch {
    return "";
  }
}

// --------- Helpers ---------

function toItemsArray(payload) {
  if (Array.isArray(payload)) return payload;
  if (payload && Array.isArray(payload.items)) return payload.items;
  return [];
}

function setSelectedNoteId(noteId) {
  selectedNoteId = noteId || null;
  if (noteId) {
    try {
      localStorage.setItem(SELECTED_NOTE_STORAGE_KEY, noteId);
    } catch (_) {}
  } else {
    try {
      localStorage.removeItem(SELECTED_NOTE_STORAGE_KEY);
    } catch (_) {}
  }
}

const NOTES_PAGE_SIZE = 20;
const notesPager = {
  currentTag: "",
  start: 0,
  end: NOTES_PAGE_SIZE - 1,
  hasMore: true,
  loading: false,
};

function resetNotesPager(tagId = "") {
  notesPager.currentTag = tagId || "";
  notesPager.start = 0;
  notesPager.end = NOTES_PAGE_SIZE - 1;
  notesPager.hasMore = true;
}

// --------- Фильтр по тегу (селект) ---------

function refreshTagFilterOptions() {
  if (!notesTagFilterSelect) return;

  const prevValue = notesTagFilterSelect.value;
  notesTagFilterSelect.innerHTML = "";

  const optAll = document.createElement("option");
  optAll.value = "";
  optAll.textContent = "Все теги";
  notesTagFilterSelect.appendChild(optAll);

  (tagsCache || []).forEach((tag) => {
    if (!tag || !tag.id) return;
    const opt = document.createElement("option");
    opt.value = tag.id;
    const label =
      ((tag.emoji || "") + " " + (tag.title || "")).trim() || tag.id;
    opt.textContent = label;
    notesTagFilterSelect.appendChild(opt);
  });

  if (
    prevValue &&
    Array.from(notesTagFilterSelect.options).some(
      (o) => o.value === prevValue
    )
  ) {
    notesTagFilterSelect.value = prevValue;
  } else {
    notesTagFilterSelect.value = "";
  }
}

if (notesTagFilterSelect) {
  notesTagFilterSelect.addEventListener("change", async () => {
    const id = notesTagFilterSelect.value;
    if (!id) {
      await loadNotes();
    } else {
      await loadNotesByTag(id);
    }
  });

  notesTagFilterSelect.addEventListener("change", () => {
    notesListEl?.scrollTo({ top: 0 });
  });
}

if (notesListEl) {
  notesListEl.addEventListener("scroll", async () => {
    const nearBottom =
      notesListEl.scrollTop + notesListEl.clientHeight >=
      notesListEl.scrollHeight - 60;
    if (nearBottom) {
      await loadMoreNotes();
    }
  });
}

// --------- Рендер заметок / корзины ---------

function hideAllNoteMenus() {
  const menus = document.querySelectorAll(".note-menu");
  menus.forEach((m) => m.classList.add("hidden"));
}

function renderNotesList(notes, { append = false } = {}) {
  if (!notesListEl) return;

  if (!append) {
    notesListEl.innerHTML = "";
  }

  if (!notes || notes.length === 0) {
    if (!notesListEl.childElementCount) {
      const empty = document.createElement("div");
      empty.className = "notes-list-empty";
      empty.textContent = "Пока нет заметок";
      notesListEl.appendChild(empty);
    }
    return;
  }

  const emptyEl = notesListEl.querySelector(".notes-list-empty");
  if (emptyEl) emptyEl.remove();

  notes.forEach((note) => {
    const item = document.createElement("div");
    item.className = "note-item";

    const row = document.createElement("div");
    row.className = "note-row";

    const content = document.createElement("div");
    content.className = "note-content";
  
    const title = document.createElement("div");
    title.className = "note-title";
    title.textContent = note.title || "Без названия";
    content.appendChild(title);

    const fb = document.createElement("div");
    fb.className = "note-fb";
    fb.textContent = note.first_block || "Без названия";
    content.appendChild(fb);
    if (note.tag) {
      const tagPill = document.createElement("div");
      tagPill.className = "note-tag-pill";

      const emojiSpan = document.createElement("span");
      emojiSpan.textContent = note.tag.emoji || "🏷️";

      const textSpan = document.createElement("span");
      textSpan.textContent = note.tag.title || "тег";

      tagPill.appendChild(emojiSpan);
      tagPill.appendChild(textSpan);
      content.appendChild(tagPill);
    }

    const meta = document.createElement("div");
    meta.className = "note-meta";
    const updated = formatDate(note.updated_at || note.created_at);
    meta.textContent = updated ? `Обновлено: ${updated}` : "";
    content.appendChild(meta);

    const menuBtn = document.createElement("button");
    menuBtn.type = "button";
    menuBtn.className = "note-menu-btn";
    menuBtn.textContent = "⋯";

    row.appendChild(content);
    row.appendChild(menuBtn);

    const menu = document.createElement("div");
    menu.className = "note-menu hidden";

    const toTrashBtn = document.createElement("button");
    toTrashBtn.type = "button";
    toTrashBtn.className = "btn-ghost small";
    toTrashBtn.textContent = "В корзину";
    toTrashBtn.disabled = !note || !note.id;

    const addTagBtnInline = document.createElement("button");
    addTagBtnInline.type = "button";
    addTagBtnInline.className = "btn-outline small";
    addTagBtnInline.textContent = "Добавить тег";
    addTagBtnInline.disabled = !note || !note.id;

    toTrashBtn.addEventListener("click", async (event) => {
      event.stopPropagation();
      await moveNoteToTrash(note, toTrashBtn);
      hideAllNoteMenus();
    });

    addTagBtnInline.addEventListener("click", async (event) => {
      event.stopPropagation();
      await handleAddTagToNote(note);
      hideAllNoteMenus();
    });

    menu.appendChild(toTrashBtn);
    menu.appendChild(addTagBtnInline);

    const renameBtn = document.createElement("button");
    renameBtn.type = "button";
    renameBtn.className = "btn-ghost small";
    renameBtn.textContent = "Переименовать";
    renameBtn.disabled = !note || !note.id;
    renameBtn.addEventListener("click", async (event) => {
      event.stopPropagation();
      await renameNote(note);
      hideAllNoteMenus();
    });

    menu.appendChild(renameBtn);

    menuBtn.addEventListener("click", (event) => {
      event.stopPropagation();
      const isHidden = menu.classList.contains("hidden");
      hideAllNoteMenus();
      if (isHidden) menu.classList.remove("hidden");
    });

    item.addEventListener("click", () => {
      hideAllNoteMenus();
      openNote(note);
    });

    item.appendChild(row);
    item.appendChild(menu);

    notesListEl.appendChild(item);
  });
}

function showNoteLoading() {
  if (!notesMainEl) return;
  notesMainEl.innerHTML = "";
  const placeholder = document.createElement("div");
  placeholder.className = "notes-main-placeholder";
  placeholder.textContent = "Загрузка заметки...";
  notesMainEl.appendChild(placeholder);
}

function renderNoteDetails(note) {
  if (!notesMainEl) return;
  notesMainEl.innerHTML = "";
  selectedNoteData = note || null;

  if (!note) {
    const placeholder = document.createElement("div");
    placeholder.className = "notes-main-placeholder";
    placeholder.innerHTML = "<p>Выбери заметку слева или создай новую.</p>";
    notesMainEl.appendChild(placeholder);
    return;
  }

  const card = document.createElement("div");
  card.className = "note-detail";

  const titleEl = document.createElement("div");
  titleEl.className = "note-detail-title";
  titleEl.textContent = note.title || "Без названия";
  card.appendChild(titleEl);

  const metaEl = document.createElement("div");
  metaEl.className = "note-detail-meta";
  const updated = formatDate(note.updated_at || note.created_at);
  metaEl.textContent = updated ? `Обновлено: ${updated}` : "";
  card.appendChild(metaEl);

  // Создание текстового блока
  const createBlockBox = document.createElement("div");
  createBlockBox.className = "note-detail-create";

  const createTitle = document.createElement("div");
  createTitle.className = "note-detail-heading";
  createTitle.textContent = "Создать текстовый блок";
  createBlockBox.appendChild(createTitle);

  const textInput = document.createElement("textarea");
  textInput.className = "note-detail-block-input";
  textInput.rows = 3;
  textInput.placeholder =
    "Текст нового блока (если пусто — вставлю пример с жирным словом)";
  createBlockBox.appendChild(textInput);

  const createBtn = document.createElement("button");
  createBtn.type = "button";
  createBtn.className = "btn-primary small";
  createBtn.textContent = "Добавить блок";
  createBtn.addEventListener("click", async () => {
    const text = textInput.value.trim();
    await createTextBlock(note.id, text);
  });
  createBlockBox.appendChild(createBtn);

  card.appendChild(createBlockBox);

  const heading = document.createElement("div");
  heading.className = "note-detail-heading";
  heading.textContent = "Blocks";
  card.appendChild(heading);

  const blocksContainer = document.createElement("div");
  blocksContainer.className = "note-detail-blocks";

  if (note.blocks === null || note.blocks === undefined) {
    const info = document.createElement("div");
    info.textContent = "null";
    blocksContainer.appendChild(info);
  } else if (Array.isArray(note.blocks) && note.blocks.length > 0) {
    note.blocks.forEach((block, idx) => {
      const blockBox = document.createElement("div");
      blockBox.className = "note-block-card";
      const order = idx;
      const blockTitle = document.createElement("div");
      blockTitle.className = "note-block-title";
      blockTitle.textContent = `Block #${idx + 1} (type: ${
        block.type || "?"
      }, order: ${order ?? "?"})`;
      blockBox.appendChild(blockTitle);

      const pre = document.createElement("pre");
      pre.className = "note-block-pre";
      pre.textContent = JSON.stringify(block, null, 2);
      blockBox.appendChild(pre);

      const actions = document.createElement("div");
      actions.className = "note-block-actions";

      const startInput = document.createElement("input");
      startInput.type = "number";
      startInput.min = "0";
      startInput.value = "0";
      startInput.className = "note-block-field";
      startInput.placeholder = "start";

      const endInput = document.createElement("input");
      endInput.type = "number";
      endInput.min = "0";
      endInput.value = "2";
      endInput.className = "note-block-field";
      endInput.placeholder = "end";

      const styleInput = document.createElement("input");
      styleInput.type = "text";
      styleInput.value = "bold";
      styleInput.className = "note-block-field";
      styleInput.placeholder = "style";

      const applyBtn = document.createElement("button");
      applyBtn.type = "button";
      applyBtn.className = "btn-outline small";
      applyBtn.textContent = "Применить стиль";
      applyBtn.addEventListener("click", async () => {
        await applyStyleToBlock(block.id, note.id, {
          start: Number(startInput.value) || 0,
          end: Number(endInput.value) || 0,
          style: styleInput.value || "bold",
        });
      });

      const orderBox = document.createElement("div");
      orderBox.className = "note-block-order";

      const upBtn = document.createElement("button");
      upBtn.type = "button";
      upBtn.className = "btn-ghost small";
      upBtn.textContent = "Выше";
      upBtn.disabled = order === undefined || order === null;
      upBtn.addEventListener("click", async () => {
        console.log("Change order up", note.id, order, order - 1);
        await changeBlockOrder(note.id, order, order - 1);
      });

      const downBtn = document.createElement("button");
      downBtn.type = "button";
      downBtn.className = "btn-ghost small";
      downBtn.textContent = "Ниже";
      downBtn.disabled = order === undefined || order === null;
      downBtn.addEventListener("click", async () => {
        console.log("Change order down", note.id, order, order + 1);
        await changeBlockOrder(note.id, order, order + 1);
      });

      orderBox.appendChild(upBtn);
      orderBox.appendChild(downBtn);

      const typeInput = document.createElement("input");
      typeInput.type = "text";
      typeInput.value = block.type || "";
      typeInput.className = "note-block-field";
      typeInput.placeholder = "type";

      const typeBtn = document.createElement("button");
      typeBtn.type = "button";
      typeBtn.className = "btn-ghost small";
      typeBtn.textContent = "Сменить тип";
      typeBtn.addEventListener("click", async () => {
        await changeBlockType(block.id, typeInput.value || "text");
      });

      const deleteBtn = document.createElement("button");
      deleteBtn.type = "button";
      deleteBtn.className = "btn-outline small";
      deleteBtn.textContent = "Удалить";
      deleteBtn.addEventListener("click", async () => {
        await deleteBlock(block.id, note.id);
      });

      actions.appendChild(startInput);
      actions.appendChild(endInput);
      actions.appendChild(styleInput);
      actions.appendChild(applyBtn);
      actions.appendChild(orderBox);
      actions.appendChild(typeInput);
      actions.appendChild(typeBtn);
      actions.appendChild(deleteBtn);

      blockBox.appendChild(actions);
      blocksContainer.appendChild(blockBox);
    });
  } else if (Array.isArray(note.blocks) && note.blocks.length === 0) {
    const info = document.createElement("div");
    info.textContent = "[]";
    blocksContainer.appendChild(info);
  } else {
    const info = document.createElement("div");
    info.textContent = String(note.blocks);
    blocksContainer.appendChild(info);
  }

  card.appendChild(blocksContainer);

  notesMainEl.appendChild(card);
}

function renderNoteError() {
  if (!notesMainEl) return;
  notesMainEl.innerHTML = "";
  selectedNoteData = null;
  const errorEl = document.createElement("div");
  errorEl.className = "notes-main-placeholder";
  errorEl.textContent = "Не удалось загрузить заметку.";
  notesMainEl.appendChild(errorEl);
}

function renderTrashList(notes) {
  if (!trashListEl) return;

  trashListEl.innerHTML = "";

  if (!notes || notes.length === 0) {
    const empty = document.createElement("div");
    empty.className = "notes-list-empty";
    empty.textContent = "Корзина пуста";
    trashListEl.appendChild(empty);
    return;
  }

  notes.forEach((note) => {
    const item = document.createElement("div");
    item.className = "note-item";

    const content = document.createElement("div");
    content.className = "note-content";

    const title = document.createElement("div");
    title.className = "note-title";
    title.textContent = note.title || "Без названия";

    const meta = document.createElement("div");
    meta.className = "note-meta";
    const updated = formatDate(note.updated_at || note.created_at);
    meta.textContent = updated ? `Обновлено: ${updated}` : "";

    const firstBlock = document.createElement("div");
    firstBlock.className = "note-first-block";
    const fb = note.first_block;
    firstBlock.textContent =
      typeof fb === "string" && fb.trim()
        ? fb
        : fb === null
        ? "first_block: null"
        : "first_block: —";

    content.appendChild(title);
    content.appendChild(firstBlock);
    content.appendChild(meta);

    const actions = document.createElement("div");
    actions.className = "note-menu";

    const restoreBtn = document.createElement("button");
    restoreBtn.type = "button";
    restoreBtn.className = "btn-outline small";
    restoreBtn.textContent = "Вернуть";
    restoreBtn.disabled = !note || !note.id;

    restoreBtn.addEventListener("click", async (event) => {
      event.stopPropagation();
      await restoreNoteFromTrash(note, restoreBtn);
    });

    actions.appendChild(restoreBtn);

    item.appendChild(content);
    item.appendChild(actions);

    trashListEl.appendChild(item);
  });
}

// --------- Загрузка заметок / корзины ---------

async function loadNotesPage({ tagId = "", append = false, keepSelection = false } = {}) {
  const targetTag = tagId || "";
  const tagChanged = targetTag !== notesPager.currentTag;

  if (tagChanged || !append) {
    resetNotesPager(targetTag);
    if (!keepSelection) {
      setSelectedNoteId(null);
      renderNoteDetails(null);
    }
  } else {
    if (!notesPager.hasMore || notesPager.loading) return;
    const nextStart = notesPager.end + 1;
    notesPager.start = nextStart;
    notesPager.end = nextStart + NOTES_PAGE_SIZE - 1;
  }

  if (notesPager.loading) return;
  notesPager.loading = true;

  try {
    const qs = new URLSearchParams({
      start: String(notesPager.start),
      end: String(notesPager.end),
    });
    if (targetTag) qs.set("id", targetTag);

    const path = targetTag
      ? `/api/note/by-tag?${qs}`
      : `/api/note/all?${qs}`;

    const data = await apiRequest(path, { method: "GET" });
    const items = toItemsArray(data);
    renderNotesList(items, { append: append && !tagChanged });

    if (!append && !restoredNoteFromStorage) {
      let storedId = null;
      try {
        storedId = localStorage.getItem(SELECTED_NOTE_STORAGE_KEY);
      } catch (_) {}
      if (storedId) {
        restoredNoteFromStorage = true;
        await openNote({ id: storedId });
      }
    }

    if (items.length < NOTES_PAGE_SIZE) {
      notesPager.hasMore = false;
    }
  } catch (err) {
    console.error("Failed to load notes:", err);
  } finally {
    notesPager.loading = false;
  }
}

async function loadNotes() {
  await loadNotesPage({ tagId: "", append: false });
}

async function loadNotesByTag(tagId) {
  if (!tagId) {
    await loadNotes();
    return;
  }
  await loadNotesPage({ tagId, append: false });
}

async function loadMoreNotes() {
  const tagId = notesTagFilterSelect?.value || "";
  await loadNotesPage({ tagId, append: true });
}

async function refreshNotesList() {
  const tagId = notesTagFilterSelect?.value || "";
  await loadNotesPage({ tagId, append: false, keepSelection: true });
}

async function openNote(note) {
  if (!note || !note.id) return;
  setSelectedNoteId(note.id);
  showNoteLoading();

  try {
    const data = await apiRequest(`/api/note?id=${encodeURIComponent(note.id)}`, {
      method: "GET",
    });
    if (selectedNoteId !== note.id) return;
    renderNoteDetails(data);
  } catch (err) {
    console.error("Failed to load note detail:", err);
    if (selectedNoteId === note.id) {
      renderNoteError();
    }
  }
}

async function loadTrashNotes() {
  try {
    const data = await apiRequest("/api/trash", { method: "GET" });
    const items = toItemsArray(data);
    renderTrashList(items);
  } catch (err) {
    console.error("Failed to load trash:", err);
  }
}

async function renameNote(note) {
  if (!note || !note.id) return;
  const currentTitle = note.title || "";
  const newTitle = prompt("Новое название заметки:", currentTitle);
  if (newTitle === null) return;
  const title = newTitle.trim();
  if (!title) {
    alert("Название не может быть пустым");
    return;
  }

  try {
    await apiRequest("/api/note/title", {
      method: "PATCH",
      body: { id: note.id, title },
    });

    await refreshNotesList();

    if (selectedNoteId === note.id) {
      await openNote({ id: note.id });
    }
  } catch (err) {
    console.error("Failed to rename note:", err);
  }
}

async function createTextBlock(noteId, text) {
  if (!noteId) return;
  const blocks = Array.isArray(selectedNoteData?.blocks)
    ? selectedNoteData.blocks
    : [];
  const pos = blocks.length;
  const baseText =
    text ||
    "Новый текстовый блок — здесь пример длинной строки, чтобы проверить разметку.";
  const payload = {
    type: "text",
    note_id: noteId,
    pos,
    data: {
      text: [
        { style: "default", text: baseText },
      ],
    },
  };

  try {
    await apiRequest("/api/block", { method: "POST", body: payload });
    await openNote({ id: noteId });
    await refreshNotesList();
  } catch (err) {
    console.error("Failed to create block:", err);
  }
}

async function applyStyleToBlock(blockId, noteId, styleData) {
  if (!blockId) return;
  const payload = {
    block_id: blockId,
    op: "apply_style",
    note_id: noteId,
    data: {
      start: styleData.start ?? 0,
      end: styleData.end ?? 0,
      style: styleData.style || "bold",
    },
  };

  try {
    await apiRequest("/api/block/op", { method: "POST", body: payload });
    if (selectedNoteId) {
      await openNote({ id: selectedNoteId });
    }
    await refreshNotesList();
  } catch (err) {
    console.error("Failed to apply style to block:", err);
  }
}

async function changeBlockOrder(noteId, oldOrder, newOrder) {
  const targetNoteId = noteId || selectedNoteId;
  if (
    !targetNoteId ||
    oldOrder === undefined ||
    oldOrder === null ||
    newOrder === undefined ||
    newOrder === null
  ) {
    return;
  }
  const oldInt = Number(oldOrder);
  const newInt = Math.max(0, Number(newOrder));
  if (Number.isNaN(oldInt) || Number.isNaN(newInt)) return;
  const payload = { note_id: targetNoteId, old_order: oldInt, new_order: newInt };

  try {
    await apiRequest("/api/block/order", { method: "PATCH", body: payload });
    if (selectedNoteId) {
      await openNote({ id: selectedNoteId });
    }
    await refreshNotesList();
  } catch (err) {
    console.error("Failed to change block order:", err);
  }
}

async function changeBlockType(blockId, newType) {
  if (!blockId || !newType) return;
  const payload = { id: blockId, new_type: newType };

  try {
    await apiRequest("/api/block/type", { method: "PATCH", body: payload });
    if (selectedNoteId) {
      await openNote({ id: selectedNoteId });
    }
    await refreshNotesList();
  } catch (err) {
    console.error("Failed to change block type:", err);
  }
}

async function deleteBlock(blockId, noteId) {
  if (!blockId) return;
  const targetNoteId = noteId || selectedNoteId;
  if (!targetNoteId) return;

  try {
    await apiRequest(
      `/api/block?block_id=${encodeURIComponent(
        blockId
      )}&note_id=${encodeURIComponent(targetNoteId)}`,
      { method: "DELETE" }
    );
    if (selectedNoteId === targetNoteId) {
      await openNote({ id: targetNoteId });
    }
    await refreshNotesList();
  } catch (err) {
    console.error("Failed to delete block:", err);
  }
}

async function moveNoteToTrash(note, triggerBtn) {
  if (!note || !note.id) return;
  if (triggerBtn) triggerBtn.disabled = true;

  try {
    await apiRequest(`/api/trash/to?id=${encodeURIComponent(note.id)}`, {
      method: "PUT",
    });
    const tagId = notesTagFilterSelect?.value;
    if (tagId) await loadNotesByTag(tagId);
    else await loadNotes();

    if (trashView && trashView.classList.contains("active")) {
      await loadTrashNotes();
    }
  } catch (err) {
    console.error("Failed to move note to trash:", err);
  } finally {
    if (triggerBtn) triggerBtn.disabled = false;
  }
}

async function restoreNoteFromTrash(note, triggerBtn) {
  if (!note || !note.id) return;
  if (triggerBtn) triggerBtn.disabled = true;

  try {
    await apiRequest(`/api/trash/from?id=${encodeURIComponent(note.id)}`, {
      method: "PUT",
    });
    await loadTrashNotes();
    const tagId = notesTagFilterSelect?.value;
    if (tagId) await loadNotesByTag(tagId);
    else await loadNotes();
  } catch (err) {
    console.error("Failed to restore note from trash:", err);
  } finally {
    if (triggerBtn) triggerBtn.disabled = false;
  }
}

// создание заметки
if (addNoteBtn) {
  addNoteBtn.addEventListener("click", async () => {
    addNoteBtn.disabled = true;
    try {
      await apiRequest("/api/note", {
        method: "POST",
        body: { title: "новая заметка" },
      });
      const tagId = notesTagFilterSelect?.value;
      if (tagId) await loadNotesByTag(tagId);
      else await loadNotes();
    } catch (err) {
      console.error("Failed to create note:", err);
    } finally {
      addNoteBtn.disabled = false;
    }
  });
}

// очистка корзины
if (clearTrashBtn) {
  clearTrashBtn.addEventListener("click", async () => {
    clearTrashBtn.disabled = true;
    try {
      await apiRequest("/api/trash", { method: "DELETE" });
      await loadTrashNotes();
    } catch (err) {
      console.error("Failed to clear trash:", err);
    } finally {
      clearTrashBtn.disabled = false;
    }
  });
}

// --------- Теги ---------

function renderTagsList(tags) {
  if (!tagsListEl) return;

  tagsListEl.innerHTML = "";

  if (!tags || tags.length === 0) {
    const empty = document.createElement("div");
    empty.className = "notes-list-empty";
    empty.textContent = "Пока нет тегов";
    tagsListEl.appendChild(empty);
    return;
  }

  tags.forEach((tag) => {
    const item = document.createElement("div");
    item.className = "tag-item";

    const main = document.createElement("div");
    main.className = "tag-main";

    const colorDot = document.createElement("div");
    colorDot.className = "tag-color-dot";
    if (tag.color) {
      colorDot.style.background = tag.color;
    }

    const emojiSpan = document.createElement("div");
    emojiSpan.className = "tag-emoji";
    emojiSpan.textContent = tag.emoji || "🏷️";

    const texts = document.createElement("div");
    texts.className = "tag-texts";

    const titleDiv = document.createElement("div");
    titleDiv.className = "tag-title";
    titleDiv.textContent = tag.title || "Без названия";

    const metaDiv = document.createElement("div");
    metaDiv.className = "tag-meta";
    metaDiv.textContent = tag.id ? `ID: ${tag.id}` : "";

    texts.appendChild(titleDiv);
    texts.appendChild(metaDiv);

    main.appendChild(colorDot);
    main.appendChild(emojiSpan);
    main.appendChild(texts);

    const editBtn = document.createElement("button");
    editBtn.type = "button";
    editBtn.className = "btn-outline small";
    editBtn.textContent = "Изменить";

    editBtn.addEventListener("click", () => editTag(tag));

    item.appendChild(main);
    item.appendChild(editBtn);

    tagsListEl.appendChild(item);
  });
}

async function loadTags() {
  try {
    const data = await apiRequest("/api/tag/by-user", { method: "GET" });
    tagsCache = toItemsArray(data);
    renderTagsList(tagsCache);
    refreshTagFilterOptions();
  } catch (err) {
    console.error("Failed to load tags:", err);
  }
}

async function createTag() {
  const title = prompt("Название тега:");
  if (!title) return;

  const emoji = prompt("Emoji для тега (можно пусто):", "🏷️") || "";
  const color =
    prompt("Цвет тега (hex или имя, можно пусто):", "#ff9fd1") || "";

  try {
    await apiRequest("/api/tag", {
      method: "POST",
      body: { title, emoji, color },
    });
    await loadTags();
  } catch (err) {
    console.error("Failed to create tag:", err);
  }
}

async function editTag(tag) {
  if (!tag || !tag.id) return;

  const newTitle = prompt(
    "Новое название тега (оставь пустым, чтобы не менять):",
    tag.title || ""
  );
  const newEmoji = prompt(
    "Новый emoji (оставь пустым, чтобы не менять):",
    tag.emoji || ""
  );
  const newColor = prompt(
    "Новый цвет (оставь пустым, чтобы не менять):",
    tag.color || ""
  );

  try {
    if (newTitle !== null && newTitle !== "" && newTitle !== tag.title) {
      await apiRequest("/api/tag/title", {
        method: "PUT",
        body: { id: tag.id, title: newTitle },
      });
    }

    if (newEmoji !== null && newEmoji !== "" && newEmoji !== tag.emoji) {
      await apiRequest("/api/tag/emoji", {
        method: "PUT",
        body: { id: tag.id, emoji: newEmoji },
      });
    }

    if (newColor !== null && newColor !== "" && newColor !== tag.color) {
      await apiRequest("/api/tag/color", {
        method: "PUT",
        body: { id: tag.id, color: newColor },
      });
    }

    await loadTags();
    await loadNotes();
  } catch (err) {
    console.error("Failed to update tag:", err);
  }
}

if (addTagBtn) {
  addTagBtn.addEventListener("click", createTag);
}

// Добавление тега к заметке
async function handleAddTagToNote(note) {
  if (!note || !note.id) return;

  // если ещё не загружали теги — загрузим
  if (!tagsCache || tagsCache.length === 0) {
    await loadTags();
  }

  if (!tagsCache || tagsCache.length === 0) {
    alert("Пока нет тегов. Сначала создай их во вкладке тегов.");
    return;
  }

  const listText = tagsCache
    .map(
      (t, idx) => `${idx + 1}. ${t.title || "(без названия)"} ${t.emoji || ""}`
    )
    .join("\n");

  const input = prompt(
    "Выбери номер тега, который добавить к заметке:\n\n" + listText
  );

  if (!input) return;
  const index = Number(input);
  if (!Number.isInteger(index) || index < 1 || index > tagsCache.length) {
    alert("Неверный номер тега");
    return;
  }

  const tag = tagsCache[index - 1];
  if (!tag || !tag.id) {
    alert("Не удалось найти выбранный тег");
    return;
  }

  try {
    await apiRequest("/api/note/tag", {
      method: "POST",
      body: { noteId: note.id, tagId: tag.id },
    });
    const currentTagId = notesTagFilterSelect?.value;
    if (currentTagId) await loadNotesByTag(currentTagId);
    else await loadNotes();
  } catch (err) {
    console.error("Failed to add tag to note:", err);
  }
}

// --------- Профиль ---------

function renderUserProfile() {
  const user = currentUser;
  if (!user) return;

  const login = user.login || "—";
  const email = user.email || "—";
  const about = user.about || "";
  const photo = user.photo || "";
  const id = user.id || "—";

  profileLoginSpan.textContent = login;
  profileEmailSpan.textContent = email;
  profileEmailInline.textContent = email;
  profileIdSpan.textContent = `ID: ${id}`;
  profileAboutP.textContent = about || "Описание пока пустое.";

  if (photo) {
    profilePhotoImg.src = photo;
    profilePhotoImg.classList.add("visible");
    profileAvatarFallback.style.display = "none";
  } else {
    profilePhotoImg.src = "";
    profilePhotoImg.classList.remove("visible");
    profileAvatarFallback.style.display = "flex";
  }
}

async function loadUserData() {
  try {
    const data = await apiRequest("/api/user/data", { method: "GET" });
    currentUser = data || null;
    renderUserProfile();
  } catch (err) {
    console.error("Failed to load user data:", err);
  }
}

if (editAboutBtn) {
  editAboutBtn.addEventListener("click", async () => {
    const current = currentUser?.about || "";
    const newAbout = prompt("Новое описание (about):", current);
    if (newAbout === null) return;

    try {
      await apiRequest("/api/user/about", {
        method: "PATCH",
        body: { new_about: newAbout },
      });
      await loadUserData();
    } catch (err) {
      console.error("Failed to update about:", err);
    }
  });
}

if (editEmailBtn) {
  editEmailBtn.addEventListener("click", async () => {
    const current = currentUser?.email || "";
    const newEmail = prompt("Новый email:", current);
    if (!newEmail) return;

    try {
      await apiRequest("/api/user/email", {
        method: "PATCH",
        body: { new_email: newEmail },
      });
      await loadUserData();
    } catch (err) {
      console.error("Failed to update email:", err);
    }
  });
}

if (editPhotoBtn) {
  editPhotoBtn.addEventListener("click", async () => {
    const current = currentUser?.photo || "";
    const newPhoto = prompt("Ссылка на новое фото (URL):", current);
    if (newPhoto === null) return;

    try {
      await apiRequest("/api/user/photo", {
        method: "PATCH",
        body: { new_photo: newPhoto },
      });
      await loadUserData();
    } catch (err) {
      console.error("Failed to update photo:", err);
    }
  });
}

if (changePwBtn) {
  changePwBtn.addEventListener("click", async () => {
    const pw1 = prompt("Новый пароль:");
    if (!pw1) return;
    const pw2 = prompt("Повтори новый пароль:");
    if (!pw2) return;

    try {
      await apiRequest("/api/user/pw", {
        method: "PATCH",
        body: { new_password: pw1, new_password_2: pw2 },
      });
      alert("Пароль обновлён (если бекэнд принял запрос).");
    } catch (err) {
      console.error("Failed to change password:", err);
    }
  });
}

// --------- Проверка токена и переход на notes ---------

async function verifyTokenAndGoToNotes() {
  try {
    await apiRequest("/api/auth/token", { method: "GET" });
    setAuthUI(true);
    showView("notes");
    setNavActive("notes");
    await Promise.all([loadNotes(), loadTags(), loadUserData()]);
  } catch (err) {
    console.error("Token check failed:", err);
  }
}

// --------- Обработка форм ---------

// Вход
if (loginForm) {
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("login-email").value.trim();
    const login = document.getElementById("login-login").value.trim();
    const password = document.getElementById("login-password").value;

    const submitBtn = loginForm.querySelector("button[type=submit]");
    submitBtn.disabled = true;

    const payload = { password };
    if (email) payload.email = email;
    if (login) payload.login = login;

    try {
      await apiRequest("/api/auth", {
        method: "POST",
        body: payload,
      });

      await verifyTokenAndGoToNotes();
    } catch (err) {
      console.error("Auth error:", err);
    } finally {
      submitBtn.disabled = false;
    }
  });
}

// Регистрация
if (registerForm) {
  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("reg-email").value.trim();
    const login = document.getElementById("reg-login").value.trim();
    const pw1 = document.getElementById("reg-pw1").value;
    const pw2 = document.getElementById("reg-pw2").value;

    const submitBtn = registerForm.querySelector("button[type=submit]");
    submitBtn.disabled = true;

    const payload = { email, login, pw1, pw2 };

    try {
      await apiRequest("/api/auth/reg", {
        method: "POST",
        body: payload,
      });

      await verifyTokenAndGoToNotes();
    } catch (err) {
      console.error("Registration error:", err);
    } finally {
      submitBtn.disabled = false;
    }
  });
}

// --------- Auto-check на загрузке ---------

(async function autoCheckAuthOnLoad() {
  try {
    await apiRequest("/api/auth/token", { method: "GET" });
    setAuthUI(true);
    showView("notes");
    setNavActive("notes");
    await Promise.all([loadNotes(), loadTags(), loadUserData()]);
  } catch {
    setAuthUI(false);
    showView("landing");
    setNavActive("");
  }
})();
