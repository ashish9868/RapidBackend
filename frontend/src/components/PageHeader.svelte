<script>
    let {
        title = "",
        description = "",
        actions,
    } = $props();
</script>

<div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-6">
    <div>
        <h1 class="text-2xl font-bold text-slate-900">{title}</h1>
        {#if description}
            <p class="text-sm text-slate-500 mt-1">{description}</p>
        {/if}
    </div>
    {#if actions}
        <div class="flex items-center gap-3">
            {@render actions()}
        </div>
    {/if}
</div>
