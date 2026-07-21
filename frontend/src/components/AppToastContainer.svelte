<script>
    import { FlatToast, ToastContainer } from "svelte-toasts";
</script>

<ToastContainer let:data>
    <FlatToast {data} />
</ToastContainer>
