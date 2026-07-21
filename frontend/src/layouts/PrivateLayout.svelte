<script>
    import { FlatToast, ToastContainer } from "svelte-toasts";
    import AppToastContainer from "../components/AppToastContainer.svelte";

    let { children } = $props();
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
                <i data-lucide="layers" class="text-indigo-400 w-6 h-6"></i>
                <span class="text-white font-semibold text-lg tracking-wide"
                    >Rapid Backend 33</span
                >
            </div>

            <!-- Navigation Links -->
            <nav class="p-4 space-y-1">
                <a
                    href="#"
                    class="flex items-center gap-3 px-4 py-2.5 rounded-lg bg-indigo-600 text-white font-medium transition-colors"
                >
                    <i data-lucide="layout-dashboard" class="w-5 h-5"></i>
                    <span>Dashboard</span>
                </a>
                <a
                    href="#"
                    class="flex items-center gap-3 px-4 py-2.5 rounded-lg hover:bg-slate-800 hover:text-slate-100 transition-colors"
                >
                    <i data-lucide="bar-chart-3" class="w-5 h-5"></i>
                    <span>Analytics</span>
                </a>
                <a
                    href="#"
                    class="flex items-center gap-3 px-4 py-2.5 rounded-lg hover:bg-slate-800 hover:text-slate-100 transition-colors"
                >
                    <i data-lucide="users" class="w-5 h-5"></i>
                    <span>Customers</span>
                </a>
                <a
                    href="#"
                    class="flex items-center gap-3 px-4 py-2.5 rounded-lg hover:bg-slate-800 hover:text-slate-100 transition-colors"
                >
                    <i data-lucide="credit-card" class="w-5 h-5"></i>
                    <span>Transactions</span>
                </a>
                <a
                    href="#"
                    class="flex items-center gap-3 px-4 py-2.5 rounded-lg hover:bg-slate-800 hover:text-slate-100 transition-colors"
                >
                    <i data-lucide="settings" class="w-5 h-5"></i>
                    <span>Settings</span>
                </a>
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
            <button class="text-slate-500 hover:text-slate-300">
                <i data-lucide="log-out" class="w-5 h-5"></i>
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
            <div class="relative w-64 md:w-96">
                <span
                    class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400"
                >
                    <i data-lucide="search" class="w-4 h-4"></i>
                </span>
                <input
                    type="text"
                    placeholder="Search dashboard..."
                    class="w-full pl-10 pr-4 py-2 text-sm bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
                />
            </div>

            <!-- Right Side Actions -->
            <div class="flex items-center gap-4">
                <button
                    class="relative p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors"
                >
                    <i data-lucide="bell" class="w-5 h-5"></i>
                    <span
                        class="absolute top-1.5 right-1.5 w-2 h-2 bg-rose-500 rounded-full"
                    ></span>
                </button>
                <button
                    class="p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors"
                >
                    <i data-lucide="help-circle" class="w-5 h-5"></i>
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
                        <i data-lucide="download" class="w-4 h-4"></i> Export
                    </button>
                    <button
                        class="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 shadow-sm transition-colors"
                    >
                        <i data-lucide="plus" class="w-4 h-4"></i> New Report
                    </button>
                </div>
            </div>

            <!-- METRIC CARDS GRID -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
                <!-- Card 1 -->
                <div
                    class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex items-start justify-between"
                >
                    <div class="space-y-2">
                        <span
                            class="text-sm font-medium text-slate-500 uppercase tracking-wider"
                            >Total Revenue</span
                        >
                        <h3 class="text-3xl font-bold text-slate-900">
                            $48,259.50
                        </h3>
                        <span
                            class="inline-flex items-center gap-1 text-xs font-medium text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full"
                        >
                            <i data-lucide="arrow-up-right" class="w-3 h-3"></i>
                            +12.5%
                        </span>
                    </div>
                    <div class="p-3 bg-indigo-50 text-indigo-600 rounded-lg">
                        <i data-lucide="dollar-sign" class="w-5 h-5"></i>
                    </div>
                </div>

                <!-- Card 2 -->
                <div
                    class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex items-start justify-between"
                >
                    <div class="space-y-2">
                        <span
                            class="text-sm font-medium text-slate-500 uppercase tracking-wider"
                            >Active Users</span
                        >
                        <h3 class="text-3xl font-bold text-slate-900">
                            10,843
                        </h3>
                        <span
                            class="inline-flex items-center gap-1 text-xs font-medium text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full"
                        >
                            <i data-lucide="arrow-up-right" class="w-3 h-3"></i>
                            +8.2%
                        </span>
                    </div>
                    <div class="p-3 bg-sky-50 text-sky-600 rounded-lg">
                        <i data-lucide="users" class="w-5 h-5"></i>
                    </div>
                </div>

                <!-- Card 3 -->
                <div
                    class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex items-start justify-between"
                >
                    <div class="space-y-2">
                        <span
                            class="text-sm font-medium text-slate-500 uppercase tracking-wider"
                            >Conversion Rate</span
                        >
                        <h3 class="text-3xl font-bold text-slate-900">2.46%</h3>
                        <span
                            class="inline-flex items-center gap-1 text-xs font-medium text-rose-600 bg-rose-50 px-2 py-0.5 rounded-full"
                        >
                            <i data-lucide="arrow-down-right" class="w-3 h-3"
                            ></i> -1.4%
                        </span>
                    </div>
                    <div class="p-3 bg-amber-50 text-amber-600 rounded-lg">
                        <i data-lucide="activity" class="w-5 h-5"></i>
                    </div>
                </div>

                <!-- Card 4 -->
                <div
                    class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex items-start justify-between"
                >
                    <div class="space-y-2">
                        <span
                            class="text-sm font-medium text-slate-500 uppercase tracking-wider"
                            >Open Tickets</span
                        >
                        <h3 class="text-3xl font-bold text-slate-900">23</h3>
                        <span
                            class="inline-flex items-center gap-1 text-xs font-medium text-slate-600 bg-slate-100 px-2 py-0.5 rounded-full"
                        >
                            Static
                        </span>
                    </div>
                    <div class="p-3 bg-rose-50 text-rose-600 rounded-lg">
                        <i data-lucide="alert-circle" class="w-5 h-5"></i>
                    </div>
                </div>
            </div>
            <!-- MAIN DATA LAYOUT -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {@render children()}
            </div>
                <AppToastContainer />
    </div>
</div>
