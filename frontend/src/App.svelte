<script>
  import Counter from "./lib/Counter.svelte";
  import { goto, route, Router } from "@mateothegreat/svelte5-router";
  import DashboardPage from "./pages/DashboardPage.svelte";
  import AuthPage from "./pages/AuthPage.svelte";
  import { LoginMode } from "./constants/Auth";
  import { ResourceApis, Resources } from "./api/ResourceApis";
  import SplashPage from "./pages/SplashPage.svelte";
  import { ToastsUtil } from "./utils/ToastsUtil";
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
  ]);
  $effect(async () => {
    const [data] = (await ResourceApis.getPaginated(Resources.Me)).results;
    user = data;
    if (data?.ID) {
      routes = [
        {
          path: "/dashboard",
          component: DashboardPage,
        },
      ];
      console.log(routes, data?.ID);
    }
    detecting = false;
    !data?.ID && ToastsUtil.showError("Session Expired!", 5000);
    goto(routes[0].path);
  });
</script>

{#if detecting}
  <SplashPage />
{:else}
  <Router {routes} />
{/if}
