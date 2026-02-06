const API_URL = 'https://localhost:8080';

async function handleResponse(response) {
    const responseText = await response.text();
    if (!response.ok) {
        try { return { success: false, message: JSON.parse(responseText).error }; } 
        catch { return { success: false, message: responseText || `Ошибка: ${response.status}` }; }
    }
    try { return { success: true, data: JSON.parse(responseText) }; } 
    catch { return { success: true, message: responseText || "Успешно (пустой ответ)" }; }
}

async function request(method, path, body = null) {
    try {
        const options = { method, credentials: 'include', headers: {} };
        if (body) {
            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(body);
        }
        const response = await fetch(API_URL + path, options);
        return await handleResponse(response);
    } catch (error) {
        return { success: false, message: "Сетевая ошибка или сервер недоступен." };
    }
}

export const loginUser = (credentials) => request("POST", "/api/auth", credentials);
export const registerUser = (details) => request("POST", "/api/auth/reg", details);
export async function checkToken() {
    try {
        const response = await fetch(`${API_URL}/api/auth/token`, { method: 'GET', credentials: 'include' });
        return (response.status === 200 || response.status === 201) 
            ? { success: true } 
            : { success: false, message: "Сессия недействительна." };
    } catch (error) {
        return { success: false, message: "Сервер недоступен." };
    }
}

async function authenticatedRequest(method, path, body = null) {
    const tokenResult = await checkToken();
    if (!tokenResult.success) {
        return { success: false, message: tokenResult.message, unauthorized: true };
    }
    return await request(method, path, body);
}

// --- User ---
export const getUserData = () => authenticatedRequest("GET", "/api/user/data");
export const deleteUser = () => authenticatedRequest("DELETE", "/api/user");
export const updateUserAbout = (new_about) => authenticatedRequest("PATCH", "/api/user/about", { new_about });
export const updateUserEmail = (new_email) => authenticatedRequest("PATCH", "/api/user/email", { new_email });
export const updateUserPhoto = (new_photo) => authenticatedRequest("PATCH", "/api/user/photo", { new_photo });
export const changePassword = (pw_old, pw_new_1, pw_new_2) => authenticatedRequest("PATCH", "/api/user/pw", { pw_old, pw_new_1, pw_new_2 });

// --- Note ---
export const getNotes = (start = 0, end = 100) => authenticatedRequest("GET", `/api/note/all?start=${start}&end=${end}`);
export const createNote = (title, id_user) => authenticatedRequest("POST", "/api/note", { title, id_user });
// *** ИСПРАВЛЕНО ***
// Тело запроса теперь { "id": "...", "title": "..." }
export const updateNoteTitle = (id, title) => authenticatedRequest("PATCH", "/api/note/title", { id, title });
export const deleteNote = (id_note) => authenticatedRequest("DELETE", `/api/note`, { id_note });
