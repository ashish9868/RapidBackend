const PROJECT_KEY = "rb_active_project";

export const authStore = $state({
    user: null,
    project: null,
    detecting: true,
});

export function loadStoredProject() {
    try {
        const raw = localStorage.getItem(PROJECT_KEY);
        return raw ? JSON.parse(raw) : null;
    } catch {
        return null;
    }
}

export function setActiveProject(project) {
    authStore.project = project;
    if (project) {
        localStorage.setItem(PROJECT_KEY, JSON.stringify(project));
    } else {
        localStorage.removeItem(PROJECT_KEY);
    }
}

export function setUser(user) {
    authStore.user = user;
}

export function clearSession() {
    authStore.user = null;
}
