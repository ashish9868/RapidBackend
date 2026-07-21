<!-- Simulated Terminal Window -->
<script>
    let { AppName, Version, DefaultPort } = $props();

</script>

<div
    class="bg-black/60 border border-zinc-800 rounded-xl font-mono text-xs overflow-hidden shadow-2xl backdrop-blur-md"
>
    <!-- Terminal Header -->
    <!-- Terminal Output -->

    <div
        class="w-full max-w-4xl rounded-xl overflow-hidden border border-zinc-800 shadow-2xl shadow-black/40"
    >
        <!-- Terminal Header -->
        <div
            class="bg-zinc-900 border-b border-zinc-800 px-4 py-3 flex items-center justify-between"
        >
            <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
            </div>

            <span class="text-zinc-400 text-sm font-medium">
                {AppName} CLI
            </span>

            <div class="w-12"></div>
        </div>

        <!-- Terminal -->
        <div class="bg-black text-sm font-mono p-6 text-zinc-200 leading-7">
            <div class="text-cyan-400">
                {AppName}
                {Version}
            </div>

            <div class="text-zinc-500 mb-5">
                High-performance backend framework
            </div>

            <div>
                <span class="text-green-400">$</span>
                <span class="text-white ml-2">fastbackend serve</span>
            </div>

            <div class="mt-3 flex items-start gap-2">
                <span class="text-emerald-400">✔</span>
                <span>
                    Server running on
                    <span class="text-cyan-400 font-semibold"
                        >http://localhost:{DefaultPort}</span
                    >
                </span>
            </div>

            <div class="mt-2 flex items-start gap-2">
                <span class="text-emerald-400">✔</span>
                <span>Ready to accept requests...</span>
            </div>
        </div>
    </div>
</div>
