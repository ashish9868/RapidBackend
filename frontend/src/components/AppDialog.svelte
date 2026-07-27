<script>
    /** @typedef {"sm" | "md" | "lg" | "xl" | "full"} Size */

    /** @type {{
     * open?: boolean,
     * title?: string,
     * size?: Size,
     * showClose?: boolean,
     * closeOnBackdrop?: boolean,
     * closeOnEscape?: boolean,
     * children?: import("svelte").Snippet,
     * footer?: import("svelte").Snippet
     * }} */
    let {
        open = false,
        title = "",
        size = "md",
        showClose = true,
        closeOnBackdrop = true,
        closeOnEscape = true,
        children,
        footer
    } = $props();

    const sizes = {
        sm: "max-w-sm",
        md: "max-w-lg",
        lg: "max-w-2xl",
        xl: "max-w-4xl",
        full: "max-w-[95vw] h-[95vh]"
    };

    function close() {
        open = false;
    }

    function backdropClick(e) {
        if (closeOnBackdrop && e.target === e.currentTarget) {
            close();
        }
    }

    function keydown(e) {
        if (closeOnEscape && e.key === "Escape") {
            close();
        }
    }
</script>

    <svelte:window onkeydown={keydown} />

{#if open}

    <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
        onclick={backdropClick}>
        <div
            class={`bg-neutral-900 border border-neutral-800 rounded-xl shadow-2xl w-full ${sizes[size]}`}
        >
            <!-- Header -->
            <div class="flex items-center justify-between border-b border-neutral-800 px-5 py-4">
                <h2 class="text-lg font-semibold">
                    {title}
                </h2>

                {#if showClose}
                    <button
                        class="rounded p-1 hover:bg-neutral-800"
                        onclick={close}
                    >
                        ✕
                    </button>
                {/if}
            </div>

            <!-- Body -->
            <div class="p-5 overflow-auto">
                {@render children?.()}
            </div>

            <!-- Footer -->
            {#if footer}
                <div class="flex justify-end gap-2 border-t border-neutral-800 px-5 py-4">
                    {@render footer()}
                </div>
            {/if}
        </div>
    </div>
{/if}