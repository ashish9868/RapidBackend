<script>
    import { Edit, Edit2 } from "@lucide/svelte";
    import Dialog from "./AppDialog.svelte";
    import FormInput from "./form/FormInput.svelte";
    import AppDialog from "./AppDialog.svelte";
    let open = $state(false)
</script>

<div class="flex flex-row">
    <div>
        Selected Project  <br />
        <div class="flex flex-row gap-2">
    <small class="font-bold text-xs">Default </small>
    <button>
        <Edit2 onclick={() => open = true} size={12}/>

        </button>
    </div>
    </div>
    <AppDialog bind:open>
        <FormInput />
    </AppDialog>
</div>