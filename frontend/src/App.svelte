<script>
  import Counter from "./lib/Counter.svelte";
  import { goto, route, Router } from "@mateothegreat/svelte5-router";
  import DashboardPage from "./pages/DashboardPage.svelte";
  import AuthPage from "./pages/AuthPage.svelte";
  import { LoginMode } from "./constants/Auth";
  import { ResourceApis, Resources } from "./api/ResourceApis";
  import SplashPage from "./pages/SplashPage.svelte";
  import { ToastsUtil } from "./utils/ToastsUtil";
    import { Routes } from "./routes";
  let detecting = $state(true);
  let user = $state(null);
  let routes = $state([
    {
      path: "/",
      component: AuthPage,
      props: {
        mode: LoginMode.LOGIN,
      },
    },
    {
      path: "/reset-password",
      component: AuthPage,
      props: {
        mode: LoginMode.RESET,
      },
    },
    {
      path: "/set-password/:token",
      component: AuthPage,
      props: {
        mode: LoginMode.SET_PASSWORD,
      },
    },
    {
      path: "/dashboard",
      component: DashboardPage,
    },
  ]);
  $effect(async () => {
    const [data] = (await ResourceApis.getPaginated(Resources.Me)).results;
    user = data;
    detecting = false;
    const isDashboard = window.location.pathname.startsWith(Routes.DASHBOARD.path)
    if (!data?.ID) {
      isDashboard && ToastsUtil.showError("Session Expired!", 5000);
      goto(Routes.LOGIN.path);
    } else if (!isDashboard) {
      goto(Routes.DASHBOARD.path);
    }
  });
</script>

{#if detecting}
  <SplashPage />
{:else}
  <Router {routes} />
{/if}
