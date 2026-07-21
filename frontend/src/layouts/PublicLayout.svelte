<script>
    import { FlatToast, ToastContainer } from "svelte-toasts";
    import AppToastContainer from "../components/AppToastContainer.svelte";
        let { children } = $props();

</script>

<div>
                {@render children()}
                <AppToastContainer />
</div>
