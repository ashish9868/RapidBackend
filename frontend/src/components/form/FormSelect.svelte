<script>
    import { Eye, EyeClosed } from "@lucide/svelte";

    /**
     * @param {{ label: string, required: boolean, type: "text" | "password" | "date" | "email" | "number" | "color", name: string, placeholder: string, error: string, value: string, options: [], iconStart: import('svelte').ComponentType }} props
     */
    let {
        label,
        options,
        required,
        type,
        name,
        value = $bindable(),
        placeholder,
        error,
        iconStart: IconStart,
        ...rest
    } = $props();
</script>

{#if label}
    
<div class="flex items-center justify-between">
    <label for="password" class="text-xs font-medium text-zinc-400">
        {label}
        {#if required}
            <i class="text-red-500">*</i>
        {/if}
    </label>
</div>
{/if}

<div class="relative">
{#if  IconStart}
    
    <span
        class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-600"
    >
        <IconStart size={12} />
    </span>
    {/if}

    <select
        id={name}
        {name}
        bind:value
        {placeholder}
        class="w-full pl-10 pr-10 py-2 text-sm bg-zinc-900 border border-zinc-800 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-zinc-100 placeholder-zinc-700 transition-colors"
        {...rest}
    >
        {#each options as option}
            <option value={option?.value}>{option?.label}</option>
        {/each}
    </select>
</div>
{#if error}
    <p class="text-red-500 first-letter:uppercase">
        {error}
    </p>
{/if}
