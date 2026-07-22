<script>
  import Counter from "./lib/Counter.svelte";
  import { route, Router } from "@mateothegreat/svelte5-router";
  import DashboardPage from "./pages/DashboardPage.svelte";
  import AuthPage from "./pages/AuthPage.svelte";
    import { LoginMode } from "./constants/Auth";

  const routes = [
    {
      path: "/",
      component: AuthPage,
      props: {
        mode: LoginMode.LOGIN
      }
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
  ];
</script>

<Router {routes} />
