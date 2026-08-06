
export const AppButton = ({
  children,
  loading = false,
  startIcon,
  endIcon,
  disabled,
  className,
  type = 'button',
  ...props
}) => {
  return (
    <button
      type={type || 'button'}
      disabled={disabled || loading}
      className={
        `w-full cursor-pointer inline-flex items-center justify-center gap-2 rounded-md px-4 py-2.5 text-sm font-medium transition-colors focus:outline-none focus:ring-1 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed 
        outline-none focus:outline-none
        bg-primary opacity-80 text-white hover:opacity-90 focus:ring-2 focus:ring-primary/80
        ${className || ''}`
      }
      {...props}
    >
      {loading ? (
        <Spinner />
      ) : (
        startIcon && <span className="flex-shrink-0">{startIcon}</span>
      )}

      <span>{children}</span>

      {!loading && endIcon && (
        <span className="flex-shrink-0">{endIcon}</span>
      )}
    </button>
  );
}

function Spinner() {
  return (
    <svg
      className="h-4 w-4 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
    >
      <circle
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
        className="opacity-25"
      />

      <path
        fill="currentColor"
        className="opacity-75"
        d="M12 2a10 10 0 0 1 10 10h-4a6 6 0 0 0-6-6V2z"
      />
    </svg>
  );
}