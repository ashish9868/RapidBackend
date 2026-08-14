import { useState } from 'preact/hooks'
import { AppButton } from './components/AppButton'
import { Save, ChevronRight, Trash2, ArrowRight } from "lucide-react";
import { AppInput } from './components/forms/AppInput';
import Logo from './assets/logo.svg'
import { useRef } from 'react';
import AppFileInput from './components/forms/AppFileInput';
import { AppSelect } from './components/forms/AppSelect';
import { MultiSelect } from './components/forms/MultiSelect';
export function App() {
  const ref = useRef(null)
  const [count, setCount] = useState(0)

  const onChange = (e) => {
    const data = Object.fromEntries(
      new FormData(ref.current).entries()
    );
    console.log(data?.options)
  }

  return (
    <div className='flex w-full max-w-lg justify-self-center items-center self-center min-h-screen'>
      <form ref={ref} onInput={onChange} className='flex w-full flex-col gap-8 px-4 pb-16 items-center items-center'>
        <div class="text-white font-bold flex flex-col text-center items-center">
          <img width={100} src={Logo} alt='Rapid Backend' />
          <p>Superadmin Login</p>
        </div>
        <AppInput name={'email'} label="Email" required />
        <AppInput type='password' name={'password'} label="Password" required />
        <AppInput name={'dt'} type={'datetime-local'} label="Password" required />
        <AppInput name={'tx1'} type={'textarea'} rows={5} label="Password" required />
        <AppInput name={'tx2'} type={'datetime-local'} rows={5} label="Password" required />
        <AppInput name={'tx3'} type={'date'} rows={5} label="Password" required />
        <AppFileInput name="file" max={2}  required/>
        <AppSelect name={'options'} multiselect label={'Choose Options'} required options={[
          {label: 'Option A', value: 'A'},
          {label: 'Option B', value: 'B'},
          {label: 'Option C', value: 'C'},
          {label: 'Option D', value: 'D'},
          {label: 'Option E', value: 'E'},
        ]} />
        <AppButton color='zinc' type='submit' endIcon={<ArrowRight />}>Login</AppButton>
        <MultiSelect />
      </form>
    </div>
  )
}
