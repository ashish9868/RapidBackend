import { Calendar, Calendar1, Clock, Clock10, Eye, EyeClosed, Key, Text } from "lucide-react"
import { useEffect, useMemo, useRef, useState } from "preact/hooks"
import flatpickr from 'flatpickr'
import "flatpickr/dist/flatpickr.min.css";
import { InputBase } from "./InputBase";

export const AppInput = ({
    type = 'text',
    name,
    label,
    value = '',
    placeholder = '',
    required = false,
    ...props
}) => {
    const ref = useRef(null)
    const [showPass, setShowPass] = useState(false)
    const isPasswordType = useMemo(() => { return `${type}` === 'password' }, [type])
    const isTextArea = useMemo(() => { return `${type}` === 'textarea' }, [type])
    const isDate = useMemo(() => { return `${type}` === 'date' }, [type])
    const isDateTime = useMemo(() => { return `${type}` === 'datetime-local' }, [type])
    const isTime = useMemo(() => { return `${type}` === 'time-local' }, [type])
    const hasTime = useMemo(() => { return isDateTime || isTime }, [isDateTime, isTime])
    const hasDate = useMemo(() => { return isDateTime || isDate }, [isDateTime, isDate])
    const {fieldType, Icon} = useMemo(() => {
        if (isDate || isDateTime || isTime) {
            return {
                fieldType: 'text',
                Icon: Calendar
            }
        }
        return {
           fieldType: isPasswordType && showPass ? 'text' : type,
           Icon: isPasswordType ? Key : Text
        }
    }, [isPasswordType, type, showPass, isDate, isDateTime, isTime])

    const fieldPlaceholder = useMemo(() => {
        if (hasDate && hasTime){
            placeholder = 'MM/DD/YYYY HH:MM'
        }else if (hasDate){
            placeholder = 'MM/DD/YYYY'
        } else if (hasTime){
            placeholder = 'HH:MM'
        } else {
            return placeholder || `Enter ${label}`
        }
    }, [placeholder, label, hasTime, hasDate])

    useEffect(() => {
        if (ref?.current) {
            ref.current.value = value || ''
        }
    }, [ref, value])

    useEffect(() => {
        if (hasDate || hasTime) {
            const fp = flatpickr(ref.current, {
                altInput: true,
                time_24hr: true,
                altFormat: `${hasDate ? 'm/d/Y' : ''} ${hasTime ? 'H:i' : ''}`.trim(),   // Display: 08/06/2026
                dateFormat: `${hasDate ? 'Y-m-d' : ''} ${hasTime ? 'H:i:S' : ''}`,  // Value: 2026-08-06
                noCalendar: isTime,
                enableTime: hasTime,
                allowInput: true,
                // onReady(_, __, fp) {
                //     return 
                //     if (!hasTime){
                //         return
                //     }
                //     const btn = document.createElement("button");
                //     btn.type = "button";
                //     btn.textContent = "Done";

                //     btn.className = "w-full px-3 py-1 rounded bg-blue-600 text-white text-sm";

                //     btn.addEventListener("click", () => {
                //         fp.close();
                //     });

                //     const footer = document.createElement("div");
                //     footer.style.padding = "8px";
                //     footer.style.textAlign = "right";
                //     footer.appendChild(btn);
                //     fp.calendarContainer.appendChild(footer);
                // }
            });
            return () => fp.destroy();
        }
    }, [hasTime, hasDate, isTime])

    return (
        <InputBase label={label} required={required} name={name} Icon={Icon}>
            <div className="flex justify-between py-1 relative">
                {isTextArea && <textarea ref={ref} autoComplete={name} placeholder={placeholder || `Enter ${label}`} className="w-full outline-none focus-none text-white text-xs py-1 overflow-visible" id={`input_${name}`} name={name} {...props}></textarea>}
                {!isTextArea && <input ref={ref} type={fieldType || 'text'} autoComplete={name} placeholder={placeholder || `Enter ${label}`} className="w-full outline-none focus-none text-white text-xs py-1" id={`input_${name}`} name={name} {...props} />}
                {isPasswordType && <button className="text-white cursor-pointer absolute right-0 top-3" onClick={() => setShowPass(!showPass)} type="button">
                    {showPass ? <EyeClosed size={12} /> : <Eye size={12} />}
                </button>}
                {(hasDate|| hasTime) && <span className="text-white absolute right-0 top-3">
                    {hasTime ? <Clock size={12} /> : <Calendar size={12} />}
                </span>}
            </div>
            </InputBase>
    )
}

