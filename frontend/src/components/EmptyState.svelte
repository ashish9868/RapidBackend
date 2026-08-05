<script>
    import { Inbox } from "@lucide/svelte";

    let {
        title = "Nothing here yet",
        description = "Get started by creating your first item.",
        action,
    } = $props();
</script>

<div class="flex flex-col items-center justify-center py-16 px-6 text-center bg-white border border-dashed border-slate-200 rounded-xl">
    <div class="w-14 h-14 rounded-full bg-slate-100 flex items-center justify-center mb-4">
        <Inbox class="w-7 h-7 text-slate-400" />
    </div>
    <h3 class="text-lg font-semibold text-slate-800">{title}</h3>
    <p class="text-sm text-slate-500 mt-1 max-w-sm">{description}</p>
    {#if action}
        <div class="mt-6">
            {@render action()}
        </div>
    {/if}
</div>
