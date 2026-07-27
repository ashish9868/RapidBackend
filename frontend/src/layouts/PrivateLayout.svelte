<script>
    import { FlatToast, ToastContainer } from "svelte-toasts";
    import AppToastContainer from "../components/AppToastContainer.svelte";
    import {
        Activity,
        AlertCircle,
        ArrowDownRight,
        ArrowUpRight,
        BarChart,
        BarChart2,
        Bell,
        CreditCard,
        DollarSign,
        Download,
        HelpCircle,
        Layers,
        LayoutDashboard,
        LogOut,
        Plus,
        PlusCircle,
        Search,
        Settings,
        Users,
    } from "@lucide/svelte";
    import { ResourceApis, Resources } from "../api/ResourceApis";
    import { Routes } from "../routes";
    import { active, route } from "@mateothegreat/svelte5-router";
    import { Global } from "../constants/Global";
    import FormInput from "../components/form/FormInput.svelte";
    import FormSelect from "../components/form/FormSelect.svelte";
    import ProjectSelect from "../components/ProjectSelect.svelte";
    import AppButton from "../components/AppButton.svelte";

    let { children } = $props();
    let projects = $state([
        {label: "Default", value: 'default'}
    ])
</script>

<div class="flex h-screen overflow-hidden">
    <!-- SIDEBAR -->
     
    <aside
        class="w-64 bg-slate-900 text-slate-300 flex flex-col justify-between hidden md:flex shrink-0"
    >
        <div>
            <!-- Logo / Brand -->
            <div
                class="h-16 flex items-center px-6 border-b border-slate-800 gap-2"
            >
            </div>

            <!-- Navigation Links -->
            <nav class="p-4 space-y-1">
                {#each Object.values(Routes) as r}
                    {#if r.showInMenu}
                        {#if r?.sectionTitle}
                            <div class="flex items-center my-4">
                                <div
                                    class="flex-grow border-t border-gray-300"
                                ></div>
                                <span class="px-3 text-gray-500 font-medium"
                                    >{r.sectionTitle}</span
                                >
                                <div
                                    class="flex-grow border-t border-gray-300"
                                ></div>
                            </div>
                        {/if}
                        <a
                            href={r.path}
                            use:route={{
                                active: {
                                    class: "bg-indigo-700",
                                },
                            }}
                            class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-white font-medium transition-colors"
                        >
                            <r.icon class="w-5 h-5" />
                            <span>{r.title}</span>
                        </a>
                    {/if}
                {/each}
            </nav>
        </div>

        <!-- User Profile Footer -->
        <div
            class="p-4 border-t border-slate-800 flex items-center justify-between"
        >
            <div class="flex items-center gap-3">
                <img
                    src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                    alt="Avatar"
                    class="w-9 h-9 rounded-full bg-slate-800"
                />
                <div>
                    <p class="text-sm font-medium text-slate-200">
                        Alex Morgan
                    </p>
                    <p class="text-xs text-slate-500">alex@acme.com</p>
                </div>
            </div>
            <button
                onclick={async () => {
                    await ResourceApis.getPaginated(Resources.LOGOUT);
                    window.location = "/";
                }}
                class="text-slate-500 cursor-pointer hover:text-slate-300"
            >
                <LogOut class="w-5 h-5" />
            </button>
        </div>
    </aside>

    <!-- MAIN CONTENT AREA -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <!-- TOP NAVBAR -->
        <header
            class="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-6 shrink-0"
        >
            <!-- Search Bar -->
            <ProjectSelect />
            <!-- Right Side Actions -->
            <div class="flex items-center gap-4">
                <button
                    class="relative p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors"
                >
                    <Bell class="w-5 h-5" />
                    <span
                        class="absolute top-1.5 right-1.5 w-2 h-2 bg-rose-500 rounded-full"
                    ></span>
                </button>
                <button
                    class="p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors"
                >
                    <HelpCircle class="w-5 h-5" />
                </button>
            </div>
        </header>

        <!-- DASHBOARD BODY CONTAINER -->
        <main class="flex-1 overflow-y-auto p-6 space-y-6">
            <!-- Welcome Header -->
            <div
                class="flex flex-col md:flex-row md:items-center md:justify-between gap-4"
            >
                <div>
                    <h1 class="text-2xl font-bold text-slate-900">
                        Welcome back, Alex
                    </h1>
                    <p class="text-sm text-slate-500">
                        Here's what's happening with your projects today.
                    </p>
                </div>
                <div class="flex items-center gap-3">
                    <button
                        class="inline-flex items-center gap-2 px-4 py-2 border border-slate-200 rounded-lg bg-white text-sm font-medium hover:bg-slate-50 transition-colors"
                    >
                        <Download class="w-4 h-4" /> Export
                    </button>
                    <AppButton icon={PlusCircle} variant={'error'} >Hello</AppButton>
                </div>
            </div>
            <!-- MAIN DATA LAYOUT -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {@render children()}
            </div>
            <AppToastContainer />
        </main>
    </div>
</div>
