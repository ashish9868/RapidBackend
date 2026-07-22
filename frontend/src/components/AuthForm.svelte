<script>
    import { goto, route } from "@mateothegreat/svelte5-router";
    import { LoginMode } from "../constants/Auth";
    import { ResourceApis, Resources } from "../api/ResourceApis";
    import {
        ArrowRight,
        EyeClosed,
        Key,
        KeyRound,
        Mail,
        User,
    } from "@lucide/svelte";
    import FormInput from "./form/FormInput.svelte";
    import { ToastsUtil } from "../utils/ToastsUtil";

    const { mode } = $props();
    let form = $state({
        email: "",
        password: "",
        confirm_password: "",
    });
    let loading = $state(false);
    let formErrors = $state({});
    let passwordToggle = $state(false);
    let isResetPassword = $derived(mode === LoginMode.RESET);
    let isLogin = $derived(mode === LoginMode.LOGIN);
    let isSetPass = $derived(mode === LoginMode.SET_PASSWORD);
</script>

<div>
    <div class="space-y-1">
        <h1 class="text-2xl font-bold text-white tracking-tight py-2">
            {#if mode === LoginMode.SET_PASSWORD}
                Set New Password
            {:else if mode === LoginMode.RESET}
                Reset Password
            {:else}
                Login to Dashboard
            {/if}
        </h1>
    </div>

    <!-- Security Inputs -->
    <form
        class="space-y-2"
        onsubmit={async (event) => {
            loading = true;
            event.preventDefault();
            let resource = Resources.LOGIN;
            if (isResetPassword) {
                resource = Resources.RESET_PASSWORD;
            }
            const { errors, success } = await ResourceApis.create(
                resource,
                form,
            );
            loading = false;
            if (!success) {
                formErrors = errors;
            } else {
                form.confirm_password = "";
                form.email = "";
                form.password = "";

                isLogin && goto("/dashboard");
                isResetPassword &&
                    ToastsUtil.showSuccess(
                        "You'll get a reset password instruction on your email, if Email is valid.",
                        15000,
                    );
            }
        }}
    >
        <!-- Username/Email -->
        {#if isLogin || isResetPassword}
            <div class="space-y-1.5">
                <FormInput
                    label="Email"
                    required={true}
                    type={"email"}
                    bind:value={form.email}
                    name={"email"}
                    error={formErrors?.email}
                    iconStart={Mail}
                    placeholder={"Enter your email"}
                />
            </div>
        {/if}

        {#if isLogin || isSetPass}
            <!-- Token/Password -->
            <div class="space-y-1.5">
                <FormInput
                    label="Password"
                    required={true}
                    type={"password"}
                    bind:value={form.password}
                    name={"email"}
                    error={formErrors?.password}
                    iconStart={Key}
                />
            </div>
        {/if}

        {#if isSetPass}
            <!-- Token/Password -->
            <div class="space-y-1.5">
                <FormInput
                    label="Confirm Password"
                    required={true}
                    type={"password"}
                    bind:value={form.confirm_password}
                    name={"confirm_password"}
                    error={formErrors?.confirm_password}
                    iconStart={Key}
                />
            </div>
        {/if}

        <br />
        <!-- Utility Preferences -->
        <!-- Authentication Execution -->
        <button
            type="submit"
            disabled={loading}
            class="w-full py-2 bg-{loading
                ? 'gray'
                : 'indigo'}-600 text-white rounded-lg text-sm font-semibold hover:bg-{loading
                ? 'red'
                : 'indigo'}-500 focus:outline-none focus:ring-2 focus:ring-{loading
                ? 'gray'
                : 'indigo'}-500 focus:ring-offset-2 focus:ring-offset-zinc-950 shadow-lg shadow-{loading
                ? 'gray'
                : 'indigo'}-600/20 transition-all flex items-center justify-center gap-2 group"
        >
            <span>{loading ? "Wait..." : "Proceed"}</span>
            {#if loading}<ArrowRight
                    class="w-4 h-4 group-hover:translate-x-0.5 transition-transform"
                />
            {/if}
        </button>
    </form>

    <!-- System Provision Link -->
    {#if isLogin || isResetPassword}
        <div class="flex justify-between text-gray-50 my-2">
            {isLogin ? "Forgot Password?" : "Remember Password?"}
            <a
                href={isLogin ? "/reset-password" : "/"}
                use:route
                class="font-bold text-indigo-400 hover:text-indigo-300 transition-colors flex flex-row"
            >
                {isLogin ? "Reset Password" : "Login"}
                <span class="htmx-indicator m-1"></span>
            </a>
        </div>
    {/if}
</div>
