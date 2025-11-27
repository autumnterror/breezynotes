// --------- Base URL ---------

let API_BASE = "http://127.0.0.1:8080";

async function apiRequest(path, { method = "GET", body } = {}) {
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

const NOTES_PAGE_SIZE = 20;
const notesPager = {
  currentTag: "",
  start: 1,
  end: NOTES_PAGE_SIZE,
  hasMore: true,
  loading: false,
};

function resetNotesPager(tagId = "") {
  notesPager.currentTag = tagId || "";
  notesPager.start = 1;
  notesPager.end = NOTES_PAGE_SIZE;
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

    menuBtn.addEventListener("click", (event) => {
      event.stopPropagation();
      const isHidden = menu.classList.contains("hidden");
      hideAllNoteMenus();
      if (isHidden) menu.classList.remove("hidden");
    });

    item.addEventListener("click", () => {
      hideAllNoteMenus();
    });

    item.appendChild(row);
    item.appendChild(menu);

    notesListEl.appendChild(item);
  });
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

    content.appendChild(title);
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

async function loadNotesPage({ tagId = "", append = false } = {}) {
  const targetTag = tagId || "";
  const tagChanged = targetTag !== notesPager.currentTag;

  if (tagChanged || !append) {
    resetNotesPager(targetTag);
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
      ? `/api/notes/by-tag?${qs}`
      : `/api/notes/all?${qs}`;

    const data = await apiRequest(path, { method: "GET" });
    const items = toItemsArray(data);
    renderNotesList(items, { append: append && !tagChanged });

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

async function loadTrashNotes() {
  try {
    const data = await apiRequest("/api/trash", { method: "GET" });
    const items = toItemsArray(data);
    renderTrashList(items);
  } catch (err) {
    console.error("Failed to load trash:", err);
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
      await apiRequest("/api/notes", {
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
    const data = await apiRequest("/api/tags/by-user", { method: "GET" });
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
    await apiRequest("/api/tags", {
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
      await apiRequest("/api/tags/title", {
        method: "PUT",
        body: { id: tag.id, title: newTitle },
      });
    }

    if (newEmoji !== null && newEmoji !== "" && newEmoji !== tag.emoji) {
      await apiRequest("/api/tags/emoji", {
        method: "PUT",
        body: { id: tag.id, emoji: newEmoji },
      });
    }

    if (newColor !== null && newColor !== "" && newColor !== tag.color) {
      await apiRequest("/api/tags/color", {
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
    await apiRequest("/api/notes/tag", {
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
