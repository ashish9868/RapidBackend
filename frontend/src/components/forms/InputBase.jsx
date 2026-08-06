

export const InputBase = ({
    Icon = null,
    name,
    label,
    required,
    children
}) => {
    return (
        <div className="flex flex-col bg-primary/50 hover:bg-primary/60 gap-y-1 px-4 py-2 rounded-lg w-full">
            <label className="flex gap-1 text-white font-bold text-xs pt-1" for={`input_${name}`}>{Icon && <Icon size={12} />}{label} {required && <i className="text-red-500 p-0">*</i>}</label>
            {children}
        </div>
    )
}