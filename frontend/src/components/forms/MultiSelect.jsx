import { useEffect, useRef, useState } from "preact/hooks"
import { InputBase } from "./InputBase";
import { CheckCheck } from "lucide-react";

export const MultiSelect = ({

}) => {
        const [open, setOpen] = useState(false)
    const [selected, setSelected] = useState([])

     const ref = useRef(null);

    useEffect(() => {
        function handleClick(e) {
            if (!ref.current?.contains(e.target)) {
                console.log('out')
                setOpen(false)
            }
        }

        document.addEventListener("pointerdown", handleClick);

        return () => {
            document.removeEventListener("pointerdown", handleClick);
        };
    }, [setOpen]);
    const options = [
        { label: 'Option A', value: 'A' },
        { label: 'Option B', value: 'B' },
        { label: 'Option C', value: 'C' },
        { label: 'Option D', value: 'D' },
        { label: 'Option E', value: 'E' },
    ]
    return (
        <InputBase label={'Multi Select'} name={'optionsxz'} Icon={CheckCheck}>
        <div ref={ref} className="w-full relative">
            <button type="button" onClick={() => setOpen(!open)} className="flex w-full items-center justify-between rounded-lg px-4 py-2.5 text-sm font-medium text-gray-700 text-white focus:outline-none">
                <span>{selected.length > 0 ? `${selected.length} Items Selected`: `No Items selected`}</span>
                <svg className="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M5.22 8.22a.75.75 0 011.06 0L10 11.94l3.72-3.72a.75.75 0 111.06 1.06l-4.25 4.25a.75.75 0 01-1.06 0L5.22 9.28a.75.75 0 010-1.06z" clip-rule="evenodd" />
                </svg>
            </button>

            <div className={`absolute z-10 mt-2 w-full rounded-lg bg-black  text-white ${!open ? 'hidden': ''}`}>
                <ul className="max-h-60 overflow-y-auto p-3 space-y-1 text-sm">

                    {options.map(x => {
                        return (
                            <li>
                                <label className="flex cursor-pointer items-center rounded px-2 py-1.5 text-white">
                                    <input type="checkbox" onChange={(e) => {
                                        setSelected(e.target.checked ? [...selected, x?.value] : [...selected].filter(s => s != x?.value))
                                    }} className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
                                    <span className="ml-2 font-medium text-gray-900">Tailwind CSS</span>
                                </label>
                            </li>

                        )
                    })}

                </ul>
            </div>
        </div>
        </InputBase>
    )
}