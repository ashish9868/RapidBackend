<script>
  import Counter from "./lib/Counter.svelte";
  import { route, Router } from "@mateothegreat/svelte5-router";
    import DashboardPage from "./pages/DashboardPage.svelte";
    import LoginPage from "./pages/LoginPage.svelte";

  const routes = [
    {
      path: "/",
      component: LoginPage,
    },
    {
      path: "/dashboard",
      component: DashboardPage
    }
  ];
</script>

<Router {routes} />
