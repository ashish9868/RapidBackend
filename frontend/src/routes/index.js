import { BookUser, Layers, LayoutDashboard, List, ProjectorIcon, Settings, Users } from "@lucide/svelte"
import { LoginMode } from "../constants/Auth"
import AuthPage from "../pages/AuthPage.svelte"
import DashboardPage from "../pages/DashboardPage.svelte"

export const Routes = {
    LOGIN: {
        title: "Login",
        path: "/",
        component: AuthPage,
        props: {
            mode: LoginMode.LOGIN
        }
    },
    RESET_PASSWORD: {
        title: "Reset Password",
        path: "/reset-password",
        component: AuthPage,
        props: {
            mode: LoginMode.RESET
        }
    },
    SET_PASSWORD: {
        title: "Set Password",
        path: "/set-password/:token",
        component: AuthPage,
        props: {
            mode: LoginMode.SET_PASSWORD,
        }
    },
    DASHBOARD: {
        showInMenu: true,
        title: "Dashboard",
        path: "/dashboard",
        component: DashboardPage,
        icon: LayoutDashboard,
    },
    Projects: {
        showInMenu: true,
        title: "Projects",
        path: "/projects",
        component: DashboardPage,
        icon: Layers,
    },
    COLLECTIONS: {
        showInMenu: true,
        title: "Collections",
        path: "/collections",
        component: DashboardPage,
        icon: List,
    },
    Users: {
        showInMenu: true,
        title: "Users",
        path: "/users",
        component: DashboardPage,
        icon: Users,
    },
    SETTINGS: {
        showInMenu: true,
        title: "Settings",
        path: "/settings",
        component: DashboardPage,
        icon: Settings
    }
}