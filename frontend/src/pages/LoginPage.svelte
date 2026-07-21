<!-- LEFT PANEL: Distinct Cyberpunk-ish BaaS Terminal Interface -->
<script>
    import {
        ArrowRight,
        Eye,
        EyeClosed,
        KeyRound,
        TerminalIcon,
        User,
    } from "@lucide/svelte";
    import Terminal from "../components/Terminal.svelte";
    import { ResourceApis, Resources } from "../api/ResourceApis";
    import PublicLayout from "../layouts/PublicLayout.svelte";
    import { goto, route } from "@mateothegreat/svelte5-router";
    let AppName = "FastBackend";
    let Version = "0.0.1-prod";
    let DefaultPort = 7000;
    let form = $state({
        email: "",
        password: "",
    });
    let loading = $state(false);
    let formErrors = $state({});
    let passwordToggle = $state(false);
</script>

<PublicLayout>
    <div class="flex h-screen overflow-hidden">
        <div
            class="hidden lg:flex w-1/2 bg-zinc-900 border-r border-zinc-800 p-12 flex-col justify-between relative overflow-hidden"
        >
            <!-- Neon Background Accents -->
            <div
                class="absolute -top-40 -left-40 w-96 h-96 bg-indigo-600/20 rounded-full blur-[128px]"
            ></div>
            <div
                class="absolute -bottom-40 -right-40 w-96 h-96 bg-emerald-600/10 rounded-full blur-[128px]"
            ></div>

            <!-- Top Branding -->
            <div class="flex items-center gap-3 relative z-10">
                <div
                    class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white shadow-lg shadow-indigo-500/20"
                >
                    <TerminalIcon class="w-5 h-5" />
                </div>
                <div>
                    <span
                        class="text-white font-bold text-lg tracking-wide block"
                        >{AppName}</span
                    >
                    <span class="text-xs text-zinc-500 font-mono"
                        >{Version}</span
                    >
                </div>
            </div>

            <!-- The Simulated Live BaaS Service Console -->
            <div class="relative z-10 my-auto max-w-xl space-y-6">
                <div class="space-y-2">
                    <h2 class="text-3xl font-black tracking-tight text-white">
                        Create Hassle Free Backend for your Apps, Web Apps &
                        Internal Tools With <br />
                        <span
                            class="text-2xl text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-purple-400 to-emerald-400"
                        >
                            {AppName} BaaS Platform
                        </span>
                    </h2>
                    <p class="text-zinc-400 text-sm">
                        <i class="font-italic">{AppName}</i> is a platform that lets
                        developers build complete backend systems in minutes instead
                        of weeks.
                    </p>
                </div>

                <Terminal {AppName} {Version} {DefaultPort} />
            </div>

            <!-- Active Service Status Metrics -->
            <div
                class="hidden grid grid-cols-3 gap-4 border-t border-zinc-800/60 pt-6 relative z-10 font-mono"
            >
                <div>
                    <span
                        class="text-[10px] uppercase text-zinc-500 tracking-wider"
                        >Global Ping</span
                    >
                    <p class="text-sm font-bold text-emerald-400 mt-0.5">
                        14ms
                    </p>
                </div>
                <div>
                    <span
                        class="text-[10px] uppercase text-zinc-500 tracking-wider"
                        >Uptime</span
                    >
                    <p class="text-sm font-bold text-zinc-200 mt-0.5">
                        99.998%
                    </p>
                </div>
                <div>
                    <span
                        class="text-[10px] uppercase text-zinc-500 tracking-wider"
                        >Requests / sec</span
                    >
                    <p class="text-sm font-bold text-indigo-400 mt-0.5">
                        142.5k
                    </p>
                </div>
            </div>
        </div>

        <!-- RIGHT PANEL: Ultra Dark Login Form Container -->
        <div
            class="w-full lg:w-1/2 flex items-center justify-center p-8 md:p-16 bg-zinc-950"
        >
            <div class="w-full max-w-sm space-y-6">
                <!-- Mobile Top Header (Hidden on Desktop) -->
                <div class="flex lg:hidden items-center gap-2 mb-6">
                    <div
                        class="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white"
                    >
                        <i data-lucide="terminal" class="w-4 h-4"></i>
                    </div>
                    <span class="text-white font-bold tracking-wide text-sm"
                        >{AppName}</span
                    >
                </div>

                <div class="space-y-1">
                    <h1 class="text-2xl font-bold text-white tracking-tight">
                        Login to Dashboard
                    </h1>
                </div>

                <!-- Security Inputs -->
                <form
                    class="space-y-2"
                    onsubmit={async (event) => {
                        loading = true;
                        event.preventDefault();
                        const { errors, success } = await ResourceApis.create(
                            Resources.LOGIN,
                            form,
                        );
                        if (!success){
                            formErrors = errors
                        }else {
                            goto("/dashboard")
                        }
                        loading = false;
                    }}
                >
                    <!-- Username/Email -->
                    <div class="space-y-1.5">
                        <label
                            for="email"
                            class="text-xs font-medium text-zinc-400 font-mono"
                            >Email*</label
                        >
                        <div class="relative">
                            <span
                                class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-600"
                            >
                                <User class="w-4 h-4" />
                            </span>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                bind:value={form.email}
                                placeholder="username@example.com"
                                class="w-full pl-10 pr-4 py-2 text-sm bg-zinc-900 border border-zinc-800 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-zinc-100 placeholder-zinc-700 transition-colors"
                            />
                        </div>
                        <p class="text-red-500 py-1 first-letter:uppercase">
                            {formErrors?.email}
                        </p>
                    </div>

                    <!-- Token/Password -->
                    <div class="space-y-1.5">
                        <div class="flex items-center justify-between">
                            <label
                                for="password"
                                class="text-xs font-medium text-zinc-400 font-mono"
                                >Password*</label
                            >
                        </div>
                        <div class="relative">
                            <span
                                class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-600"
                            >
                                <KeyRound class="w-4 h-4" />
                            </span>
                            <input
                                type={passwordToggle ? "text" : "password"}
                                id="password"
                                name="password"
                                bind:value={form.password}
                                placeholder="••••••••"
                                class="w-full pl-10 pr-10 py-2 text-sm bg-zinc-900 border border-zinc-800 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-zinc-100 placeholder-zinc-700 transition-colors"
                            />
                            <button
                                onclick={() =>
                                    (passwordToggle = !passwordToggle)}
                                type="button"
                                aria-label="Login Button"
                                class="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-600 hover:text-zinc-400"
                            >
                                {#if passwordToggle}<Eye
                                        class="w-4 h-4"
                                    />{:else}
                                    <EyeClosed class="w-4 h-4" />
                                {/if}
                            </button>
                        </div>
                        <p class="text-red-500 py-1 first-letter:uppercase">
                            {formErrors?.password}
                        </p>
                    </div>

                    <!-- Utility Preferences -->
                    <!-- Authentication Execution -->
                    <button
                        type="submit"
                        disabled={loading}
                        class="w-full py-2 bg-{loading
                            ? 'gray'
                            : 'indigo'}-600 text-white rounded-lg text-sm font-semibold hover:bg-{loading
                            ? 'red'
                            : 'indigo'}-500 focus:outline-none focus:ring-2 focus:ring-{loading
                            ? 'gray'
                            : 'indigo'}-500 focus:ring-offset-2 focus:ring-offset-zinc-950 shadow-lg shadow-{loading
                            ? 'gray'
                            : 'indigo'}-600/20 transition-all flex items-center justify-center gap-2 group"
                    >
                        <span>{loading ? "Wait..." : "Proceed"}</span>
                        {#if loading}<ArrowRight
                                class="w-4 h-4 group-hover:translate-x-0.5 transition-transform"
                            />
                        {/if}
                    </button>
                </form>

                <!-- System Provision Link -->
                <div class="flex justify-between text-gray-50">
                    Forgot Password?
                    <a
                        href="/reset-password"
                        use:route
                        class="font-bold text-indigo-400 hover:text-indigo-300 transition-colors flex flex-row"
                    >
                        Reset Password <span class="htmx-indicator m-1"></span>
                    </a>
                </div>
            </div>
        </div>
    </div>
</PublicLayout>
