<!-- LEFT PANEL: Distinct Cyberpunk-ish BaaS Terminal Interface -->
<script>
    import {
        ArrowRight,
        Eye,
        EyeClosed,
        KeyRound,
        TerminalIcon,
    } from "@lucide/svelte";
    import Terminal from "../components/Terminal.svelte";
    import PublicLayout from "../layouts/PublicLayout.svelte";
    import { Global } from "../constants/Global";
    import AuthForm from "../components/AuthForm.svelte";
    let { mode } = $props();
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
                        >{Global.AppName}</span
                    >
                    <span class="text-xs text-zinc-500 font-mono"
                        >{Global.Version}</span
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
                            {Global.AppName} BaaS Platform
                        </span>
                    </h2>
                    <p class="text-zinc-400 text-sm">
                        <i class="font-italic">{Global.AppName}</i> is a platform that lets
                        developers build complete backend systems in minutes instead
                        of weeks.
                    </p>
                </div>

                <Terminal />
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
                        >{Global.AppName}</span
                    >
                </div>

               <AuthForm mode={mode} />
            </div>
        </div>
    </div>
</PublicLayout>
