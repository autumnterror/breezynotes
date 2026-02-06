import * as api from './api.js';
import * as ui from './ui.js';

let currentUser = null;
let activeNoteId = null;

function handleUnauthorized(message) {
    ui.showToast(message || "Ваша сессия истекла. Пожалуйста, войдите снова.", 'error');
    document.cookie = "access_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    document.cookie = "refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    setTimeout(() => window.location.reload(), 2000);
}

async function showNotesPage() {
    const result = await api.getNotes();
    if (result.unauthorized) return handleUnauthorized(result.message);
    if (result.success && result.data && Array.isArray(result.data)) {
        ui.renderHomePage(result.data);
    } else {
        ui.renderHomePage([]);
        ui.showToast(result.message || 'Не удалось загрузить заметки', 'error');
    }
}

const routes = {
    '': showNotesPage,
    '#/': showNotesPage,
    '#/profile': () => ui.renderProfilePage(currentUser),
};

async function handleSuccessfulAuth() {
    const userResult = await api.getUserData();
    if (userResult.unauthorized) return handleUnauthorized(userResult.message);
    if (userResult.success) {
        currentUser = userResult.data;
        ui.showView('app');
        ui.renderAppHeader(currentUser);
        await handleRouteChange();
    } else {
        ui.showToast("Не удалось загрузить данные пользователя.", 'error');
        ui.showView('auth');
    }
}

async function handleRouteChange() {
    const hash = window.location.hash;
    const render = routes[hash] || (() => pageContent.innerHTML = '<h2>Страница не найдена</h2>');
    await render();
    addPageListeners();
}

function addPageListeners() {
    if (window.location.hash.startsWith('#/profile')) {
        addProfileFormListeners();
    } else {
        addNotesPageListeners();
    }
}

function addNotesPageListeners() {
    document.getElementById('create-note')?.addEventListener('click', handleCreateNote);
    document.getElementById('notes-list')?.addEventListener('click', handleNotesListClick);
    document.getElementById('rename-note-btn')?.addEventListener('click', handleRenameNote);
    document.getElementById('delete-note-btn')?.addEventListener('click', handleDeleteNote);
}

async function handleCreateNote() {
    const title = prompt("Введите название новой заметки:", "Новая заметка");
    if (title) {
        const result = await api.createNote(title, currentUser.id);
        if (result.unauthorized) return handleUnauthorized(result.message);
        ui.showToast(result.message || 'Заметка создана', result.success ? 'success' : 'error');
        if (result.success) await showNotesPage();
    }
}

function handleNotesListClick(e) {
    if (e.target.classList.contains('note-menu-btn')) {
        const noteItem = e.target.closest('.note-item');
        activeNoteId = noteItem.dataset.noteId;
        const menu = document.getElementById('note-context-menu');
        menu.style.top = `${e.target.offsetTop + e.target.offsetHeight}px`;
        menu.style.left = `${e.target.offsetLeft - menu.offsetWidth + e.target.offsetWidth}px`;
        menu.classList.remove('hidden');
    }
}

async function handleRenameNote() {
    if (!activeNoteId) return;
    const oldTitle = document.querySelector(`.note-item[data-note-id="${activeNoteId}"] .note-item-title`).textContent;
    const newTitle = prompt("Введите новое название:", oldTitle);
    if (newTitle && newTitle !== oldTitle) {
        const result = await api.updateNoteTitle(activeNoteId, newTitle);
        if (result.unauthorized) return handleUnauthorized(result.message);
        ui.showToast(result.message || 'Заметка переименована', result.success ? 'success' : 'error');
        if (result.success) await showNotesPage();
    }
    activeNoteId = null;
}

async function handleDeleteNote() {
    if (!activeNoteId) return;
    if (confirm("Вы уверены, что хотите удалить эту заметку?")) {
        const result = await api.deleteNote(activeNoteId);
        if (result.unauthorized) return handleUnauthorized(result.message);
        ui.showToast(result.message || 'Заметка удалена', result.success ? 'success' : 'error');
        if (result.success) await showNotesPage();
    }
    activeNoteId = null;
}

function addProfileFormListeners() {
    const page = document.getElementById('page-content');
    page.querySelector("#profile-details-form")?.addEventListener('submit', handleProfileUpdate);
    page.querySelector("#password-change-form")?.addEventListener('submit', handlePasswordChange);
    page.querySelector("#delete-account-btn")?.addEventListener('click', handleDeleteAccount);
}

async function handleProfileUpdate(e) {
    e.preventDefault();
    ui.showToast("Обновление...");
    const tokenResult = await api.checkToken();
    if (!tokenResult.success) return handleUnauthorized(tokenResult.message);
    const photo = document.getElementById('profile-photo').value;
    const email = document.getElementById('profile-email').value;
    const about = document.getElementById('profile-about').value;
    if (photo !== currentUser.photo) await api.updateUserPhoto(photo).then(r => ui.showToast(r.message || "Фото обновлено", r.success ? 'success' : 'error'));
    if (email !== currentUser.email) await api.updateUserEmail(email).then(r => ui.showToast(r.message || "Email обновлен", r.success ? 'success' : 'error'));
    if (about !== currentUser.about) await api.updateUserAbout(about).then(r => ui.showToast(r.message || "Информация обновлена", r.success ? 'success' : 'error'));
}

async function handlePasswordChange(e) {
    e.preventDefault();
    const pw_old = document.getElementById('pw_old').value;
    const pw_new_1 = document.getElementById('pw_new_1').value;
    const pw_new_2 = document.getElementById('pw_new_2').value;
    if (pw_new_1 !== pw_new_2) {
        return ui.showToast("Новые пароли не совпадают", 'error');
    }
    const result = await api.changePassword(pw_old, pw_new_1, pw_new_2);
    if (result.unauthorized) return handleUnauthorized(result.message);
    ui.showToast(result.message || "Пароль изменен", result.success ? 'success' : 'error');
    if (result.success) e.target.reset();
}

async function handleDeleteAccount() {
    if (confirm("Вы уверены, что хотите удалить свой аккаунт? Это действие необратимо.")) {
        const result = await api.deleteUser();
        if (result.unauthorized) return handleUnauthorized(result.message);
        ui.showToast(result.message || "Аккаунт удален", result.success ? 'success' : 'error');
        if (result.success) {
            setTimeout(() => window.location.reload(), 2000);
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.addEventListener('click', (e) => {
        const menu = document.getElementById('note-context-menu');
        if (menu && !menu.contains(e.target) && !e.target.classList.contains('note-menu-btn')) {
            menu.classList.add('hidden');
            activeNoteId = null;
        }
    });
    document.getElementById('show-register').addEventListener('click', () => {
        document.getElementById('login-form').classList.add('hidden');
        document.getElementById('register-form').classList.remove('hidden');
    });
    document.getElementById('show-login').addEventListener('click', () => {
        document.getElementById('register-form').classList.add('hidden');
        document.getElementById('login-form').classList.remove('hidden');
    });
    document.getElementById('login-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const data = Object.fromEntries(new FormData(e.target));
        const result = await api.loginUser({ login: data.login, email: data.email, password: data.password });
        ui.showToast(result.message || (result.success ? "Успешный вход" : "Ошибка"), result.success ? 'success' : 'error');
        if (result.success) await handleSuccessfulAuth();
    });
    document.getElementById('register-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const data = Object.fromEntries(new FormData(e.target));
        if (data.pw1 !== data.pw2) {
            return ui.showToast("Пароли не совпадают!", 'error');
        }
        const result = await api.registerUser({ login: data.login, email: data.email, pw1: data.pw1, pw2: data.pw2 });
        ui.showToast(result.message || (result.success ? "Успешная регистрация" : "Ошибка"), result.success ? 'success' : 'error');
        if (result.success) await handleSuccessfulAuth();
    });
    api.checkToken().then(result => {
        if (result.success) {
            handleSuccessfulAuth();
        } else {
            ui.showView('auth');
        }
    });
    window.addEventListener('hashchange', handleRouteChange);
});
