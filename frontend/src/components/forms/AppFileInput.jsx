import { forwardRef } from "preact/compat";
import { useEffect, useImperativeHandle, useRef, useState } from "preact/hooks";
import { InputBase } from "./InputBase";
import { Files } from "lucide-react";

const formatSize = bytes => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
};

const AppFileInput = forwardRef(function AppFileInput({
    name,
    max = 1,
    accept,
    label = "Choose files",
    required = false,
    helper = "Drop files here or click to browse",
    onChange
}, ref) {

    const inputRef = useRef();
    const filesRef = useRef([]);

    const [, refresh] = useState(0);

    const rerender = () => refresh(v => v + 1);

    const setFiles = files => {
        filesRef.current = max > 1 ? [...files].slice(0, max) : files.slice(0, 1);

        onChange?.(multiple ? filesRef.current : filesRef.current[0]);

        rerender();
    };

    const remove = index => {
        const files = [...filesRef.current];

        URL.revokeObjectURL(files[index].preview);

        files.splice(index, 1);

        setFiles(files);
    };

    const handleFiles = list => {

        const files = [...list].map(file => {

            file.preview = file.type.startsWith("image/")
                ? URL.createObjectURL(file)
                : null;

            return file;
        });

        setFiles(files);
    };

    useImperativeHandle(ref, () => ({
        getValue: () => multiple ? filesRef.current : filesRef.current[0],
        setValue: setFiles,
        clear() {
            filesRef.current.forEach(f => f.preview && URL.revokeObjectURL(f.preview));
            filesRef.current = [];
            inputRef.current.value = "";
            rerender();
        }
    }));

    useEffect(() => () => {
        filesRef.current.forEach(f => f.preview && URL.revokeObjectURL(f.preview));
    }, []);

    return (
        <InputBase label={label} name={name} required={required} Icon={Files}>

            {filesRef.current.length < max &&
                <>
                    <input
                        ref={inputRef}
                        hidden
                        type="file"
                        multiple={max > 1}
                        accept={accept}
                        name={name}
                        onChange={e => handleFiles(e.target.files)}
                    />

                    <div
                        onClick={() => inputRef.current.click()}
                        onDragOver={e => e.preventDefault()}
                        onDrop={e => {
                            e.preventDefault();
                            handleFiles(e.dataTransfer.files);
                        }}
                        className="cursor-pointer border-dotted border-2 rounded-lg bg-surface px-2 my-2 py-4  text-center hover:border-primary"
                    >
                        <div className="mt-1 text-xs text-white flex flex-col">
                            {helper}
                            <small>Maximum {max} files are allowed.</small>
                        </div>
                    </div>
                </>}

            {filesRef.current.map((file, i) => (

                <div
                    key={i}
                    className="flex items-center gap-3 rounded-lg bg-surface p-2"
                >

                    {file.preview ? (

                        <img
                            src={file.preview}
                            className="h-16 w-16 rounded object-cover"
                        />

                    ) : (

                        <div className="flex h-16 w-16 items-center justify-center rounded bg-surface-hover text-3xl">
                            📄
                        </div>

                    )}

                    <div className="flex-1 overflow-hidden">

                        <div className="truncate text-sm text-white">
                            {file.name}
                        </div>

                        <div className="text-xs text-white">
                            {formatSize(file.size)}
                        </div>

                    </div>

                    <button
                        type="button"
                        onClick={() => remove(i)}
                        className="rounded p-2 bg-red-500 hover:bg-red-600 text-white cursor-pointer"
                    >
                        ✕
                    </button>

                </div>

            ))}

        </InputBase>
    );
});

export default AppFileInput;