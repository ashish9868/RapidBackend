<script>
    import { Eye, EyeClosed } from '@lucide/svelte';

    /**
     * @param {{ label: string, required: boolean, type: "text" | "password" | "date" | "email" | "number" | "color", name: string, placeholder: string, error: string, value: string, iconStart: import('svelte').ComponentType }} props
     */
    let { label, required, type, name, value = $bindable(), placeholder, error, iconStart: IconStart, ...rest } = $props();
    let passwordToggle = $state(false);
    const isPassword = $derived(type === "password");
</script>

<div class="flex items-center justify-between">
    <label for="password" class="text-xs font-medium text-zinc-400">
        {label}
        {#if required}
        <i class="text-red-500">*</i>
        {/if}
    </label>
</div>
<div class="relative">
    <span
        class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-600"
    >
        <IconStart  size={12}/>
    </span>
    <input
        type={isPassword ? (passwordToggle ? "text" : "password") : type}
        id={name}
        {name}
        bind:value
        placeholder={isPassword ? "••••••••" : placeholder}
        class="w-full pl-10 pr-10 py-2 text-sm bg-zinc-900 border border-zinc-800 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-zinc-100 placeholder-zinc-700 transition-colors"
        {...rest}
    />
    {#if isPassword}
        <button
            onclick={() => (passwordToggle = !passwordToggle)}
            type="button"
            aria-label="Toggle Password"
            class="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-600 hover:text-zinc-400"
        >
            {#if passwordToggle}
                <Eye class="w-4 h-4" />
            {:else}
                <EyeClosed class="w-4 h-4" />
            {/if}
        </button>
    {/if}
</div>
{#if error}
    <p class="text-red-500 first-letter:uppercase">
        {error}
    </p>
{/if}
