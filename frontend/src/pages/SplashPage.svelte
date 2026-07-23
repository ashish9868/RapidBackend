<script>
    import { Circle } from "@lucide/svelte";
    import PublicLayout from "../layouts/PublicLayout.svelte";
    import { Global } from "../constants/Global";
</script>

<PublicLayout>
    <div
        class="absolute inset-0 flex flex-col items-center justify-center bg-slate-900"
    >
        <span
            class="text-3xl text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-purple-400 to-emerald-400"
        >
            {Global.AppName}
        </span>
        <p class="py-2 text-white font-bold"></p>
        <div
            class="h-10 w-10 animate-spin text-white border rounded-full border-t-red-500"
        ></div>
        <p class="mt-4 text-white">Please wait...</p>
    </div>
</PublicLayout>
