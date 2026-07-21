import { ShipWheelIcon } from "@lucide/svelte";
import { toasts } from "svelte-toasts";


export const ToastsUtil = {
    showError: (description, duration = 5000) => {
        toasts.add({
            title: 'Please Note',
            description,
            duration: 3000, // 0 or negative to avoid auto-remove
            placement: 'top-center',
            type: 'info',
            theme: 'dark',
            showProgress: true,
            type: 'error',
            theme: 'dark',
            onClick: () => { },
            onRemove: () => { },
        })
    },
    showSuccess: (description, duration = 5000) => {
        toasts.add({
            title: 'Success',
            description,
            duration: 3000, // 0 or negative to avoid auto-remove
            placement: 'top-center',
            type: 'info',
            theme: 'dark',
            showProgress: true,
            type: 'success',
            theme: 'dark',
            onClick: () => { },
            onRemove: () => { },
        })
    }
}