<script>
	/**
	 * @typedef {'primary'|'secondary'|'danger'|'success'|'ghost'} ButtonVariant
	 */

	/**
	 * @typedef DialogButton
	 * @property {string=} id
	 * @property {string} label
	 * @property {string=} action
	 * @property {ButtonVariant=} variant
	 * @property {boolean=} close
	 * @property {boolean=} disabled
	 * @property {boolean=} hidden
	 * @property {boolean=} loading
	 * @property {boolean=} autofocus
	 * @property {string=} class
	 */

	/** @type {{
		open: boolean,
		title?: string,
		description?: string,
		size?: 'sm'|'md'|'lg'|'xl'|'2xl'|'full',
		width?: string,
		closable?: boolean,
		closeOnEscape?: boolean,
		closeOnBackdrop?: boolean,
		loading?: boolean,
		showHeader?: boolean,
		showFooter?: boolean,
		buttons?: DialogButton[],
		onAction?: (button: DialogButton) => void,
		onClose?: () => void,
		children?: import('svelte').Snippet
	}} */
	let {
		open = $bindable(false),

		title = '',
		description = '',

		size = 'md',
		width = '',

		closable = true,
		closeOnEscape = true,
		closeOnBackdrop = true,

		loading = false,

		showHeader = true,
		showFooter = true,

		buttons = [
			{
				label: 'Close',
				close: true,
				variant: 'secondary'
			}
		],

		onAction = () => {},
		onClose = () => {},

		children
	} = $props();

	/** @type {Record<string,string>} */
	const sizes = {
		sm: 'max-w-sm',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl',
		'2xl': 'max-w-6xl',
		full: 'max-w-[95vw]'
	};

	/** @type {Record<ButtonVariant,string>} */
	const variants = {
		primary: 'bg-blue-600 hover:bg-blue-500 text-white',
		secondary:
			'border border-zinc-700 bg-zinc-800 text-zinc-200 hover:bg-zinc-700',
		danger: 'bg-red-600 hover:bg-red-500 text-white',
		success: 'bg-green-600 hover:bg-green-500 text-white',
		ghost: 'text-zinc-300 hover:bg-zinc-800'
	};

	function close() {
		if (loading) return;

		open = false;
		onClose();
	}

	function backdropClick(e) {
		if (e.target !== e.currentTarget) return;

		if (closeOnBackdrop) {
			close();
		}
	}

	function keydown(e) {
		if (e.key === 'Escape' && closeOnEscape) {
			e.preventDefault();
			close();
		}
	}

	/**
	 * @param {DialogButton} button
	 */
	function click(button) {
		if (button.disabled || loading) return;

		onAction(button);

		if (button.close) {
			close();
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm"
		onclick={backdropClick}
		onkeydown={keydown}
		tabindex="0"
		role="dialog"
		aria-modal="true"
	>
		<div
        data-theme="primary"
			class={`relative w-full rounded-xl border bg-white shadow-2xl ${
				width || sizes[size]
			}`}
		>
			<!-- Header -->

			{#if showHeader && (title || description)}

				<div class="relative border-b border-zinc-800 p-5">

					<h2 class="text-lg font-semibold text-white">
						{title}
					</h2>

					{#if description}
						<p class="mt-1 text-sm text-zinc-400">
							{description}
						</p>
					{/if}

					{#if closable}
						<button
							type="button"
							class="absolute right-4 top-4 rounded p-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
							onclick={close}
						>
							✕
						</button>
					{/if}

				</div>

			{/if}

			<!-- Body -->

			<div class="max-h-[70vh] overflow-auto p-5">

				{#if children}
					{@render children()}
				{/if}

			</div>

			<!-- Footer -->

			{#if showFooter && buttons.length}

				<div class="flex justify-end gap-2 border-t border-zinc-800 p-4">

					{#each buttons as button}

						{#if !button.hidden}

							<button
								type="button"
								class={`rounded-lg px-4 py-2 text-sm transition disabled:cursor-not-allowed disabled:opacity-50 ${
									variants[button.variant ?? 'secondary']
								} ${button.class ?? ''}`}
								disabled={button.disabled || loading}
								autofocus={button.autofocus}
								onclick={() => click(button)}
							>

								{#if button.loading}
									<span class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"></span>
								{/if}

								{button.label}

							</button>

						{/if}

					{/each}

				</div>

			{/if}

			<!-- Loading Overlay -->

			{#if loading}

				<div class="absolute inset-0 flex items-center justify-center rounded-xl bg-black/40 backdrop-blur-sm">

					<div class="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>

				</div>

			{/if}

		</div>
	</div>
{/if}