import { List, Option } from "lucide-react"
import { InputBase } from "./InputBase"
import { useRef, useState } from "react"


export const AppSelect = ({
    multiselect = false,
    name,
    label,
    value = '',
    placeholder = '',
    required = false,
    options = [],
    ...props
}) => {
    const ref = useRef(null)
    return (
        <InputBase label={label} name={name} required={required} Icon={List}>
            {!multiselect && <select className="hide-scroll w-full border-0 outline-none focus-none text-white text-xs py-2 px-0 rounded-lg" id={`input_${name}`} name={name} {...props} >
                {options.map(x => {
                    return <option className="text-black" value={x?.value || x?.label}>{x?.label}</option>
                })}
            </select>}
            {multiselect && (
                <div className="flex flex-col max-h-[100px] overflow-y-auto hide-scroll">
                    <input ref={ref} type="hidden" name={name} />
                    {options.map((x, i) => {
                        return <span className="py-2 flex justify-between text-white text-xs">{x?.label} <input onChange={(e) => {
                            e.stopPropagation()
                            e.preventDefault()
                            const selectedVal = `${ref.current?.value}`.split(',').filter(x => typeof x !== undefined && typeof x!== null && x !== '')
                            ref.current.value = (e.target.checked ? [...selectedVal, x?.value] : [...selectedVal].filter(s => s != x?.value)).join(',')
                            ref.current.dispatchEvent(new Event('input', {bubbles: true}))
                        }} value={x?.value} type="checkbox" /></span>
                    })}
                </div>
            )}
        </InputBase>
    )
}