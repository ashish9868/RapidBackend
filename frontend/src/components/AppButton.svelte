<script>
    /**
     * @type {{
     *  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'error',
     *  loading?: boolean,
     *  disabled?: boolean,
     *  icon?: import("svelte").Component,
     *  children?: import("svelte").Snippet,
     *  class?: string
     * }}
     **/
    let {
        variant = "primary",
        loading = false,
        disabled = false,
        icon = null,
        class: className = "",
        children,
        ...rest
    } = $props();

    let Icon = icon

    const variants = {
        primary:
            "bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500",
        secondary:
            "bg-gray-700 text-white hover:bg-gray-600 focus:ring-gray-500",
        success:
            "bg-green-600 text-white hover:bg-green-700 focus:ring-green-500",
        warning:
            "bg-yellow-500 text-black hover:bg-yellow-600 focus:ring-yellow-400",
        error: "bg-red-600 text-white hover:bg-red-700 focus:ring-red-500",
    };
</script>

<button
    {...rest}
    disabled={disabled || loading}
    class={`inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-medium shadow-sm transition-colors
        focus:outline-none focus:ring-2 focus:ring-offset-2
        disabled:opacity-50 disabled:cursor-not-allowed
        ${variants[variant] ?? variants.primary}
        ${className}`}
>
    {#if loading}
        <svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
            <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
            />
            <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
            />
        </svg>
    {:else if Icon}
        <Icon class="w-4 h-4" />
    {/if}

    {@render children?.()}
</button>
