const appHeader = document.getElementById('app-header');
const pageContent = document.getElementById('page-content');
const authView = document.getElementById('auth-view');
const appView = document.getElementById('app-view');

export function showView(viewName) {
    authView.classList.toggle('hidden', viewName !== 'auth');
    appView.classList.toggle('hidden', viewName !== 'app');
}

export function renderHomePage(notes = []) {
    pageContent.innerHTML = `
        <div class="notes-header">
            <h2>Мои заметки</h2>
            <button class="create-note-btn" id="create-note">Создать заметку</button>
        </div>
        <div class="notes-list" id="notes-list">
            ${notes.length > 0 ? notes.map(note => `
                <div class="note-item" data-note-id="${note.id}">
                    <span class="note-item-title">${note.title}</span>
                    <button class="note-menu-btn">⋮</button>
                </div>
            `).join('') : '<p>У вас пока нет заметок.</p>'}
        </div>
        <div id="note-context-menu" class="context-menu hidden">
            <button class="context-menu-item" id="rename-note-btn">Переименовать</button>
            <button class="context-menu-item delete" id="delete-note-btn">Удалить</button>
        </div>
    `;
}

export function renderAppHeader(user) {
    appHeader.innerHTML = `
        <div class="profile-button" id="profile-btn">
            <span>${user.login}</span>
            <img src="${user.photo || 'placeholder.png'}" alt="Avatar">
        </div>
    `;
    document.getElementById('profile-btn').addEventListener('click', () => {
        window.location.hash = '#/profile';
    });
}

export function renderProfilePage(user) {
    pageContent.innerHTML = `
        <form class="profile-form" id="profile-details-form">
            <h2>Профиль пользователя</h2>
            <label for="profile-photo">URL Аватара</label>
            <input type="text" id="profile-photo" value="${user.photo || ''}">
            <label for="profile-email">Email</label>
            <input type="email" id="profile-email" value="${user.email}">
            <label for="profile-about">О себе</label>
            <textarea id="profile-about">${user.about || ''}</textarea>
            <button type="submit">Сохранить изменения</button>
        </form>
        <form class="profile-form" id="password-change-form">
            <h2>Смена пароля</h2>
            <input type="password" id="pw_old" placeholder="Старый пароль" required>
            <input type="password" id="pw_new_1" placeholder="Новый пароль" required>
            <input type="password" id="pw_new_2" placeholder="Повторите новый пароль" required>
            <button type="submit">Сменить пароль</button>
        </form>
        <div class="danger-zone">
            <h2>Опасная зона</h2>
            <button id="delete-account-btn">Удалить аккаунт</button>
        </div>
    `;
}

const toastContainer = document.getElementById('toast-container');
export function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    toastContainer.appendChild(toast);
    setTimeout(() => toast.classList.add('show'), 100);
    setTimeout(() => {
        toast.classList.remove('show');
        toast.addEventListener('transitionend', () => toast.remove());
    }, 5000);
}
