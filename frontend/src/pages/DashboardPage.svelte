<script>
    import PrivateLayout from "../layouts/PrivateLayout.svelte";
</script>

<PrivateLayout>
    <div class="bg-white border border-slate-200 rounded-xl shadow-sm lg:col-span-2 overflow-hidden flex flex-col">
  <div class="p-5 border-b border-slate-200 flex items-center justify-between">
    <h3 class="font-bold text-slate-900">Recent Transactions</h3>
    <a href="#" class="text-sm font-medium text-indigo-600 hover:text-indigo-700">View all</a>
  </div>
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50 text-xs font-semibold uppercase tracking-wider text-slate-500 border-b border-slate-200">
          <th class="p-4">Customer</th>
          <th class="p-4">Status</th>
          <th class="p-4">Date</th>
          <th class="p-4 text-right">Amount</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-200 text-sm">
        <tr class="hover:bg-slate-50/70 transition-colors">
          <td class="p-4 flex items-center gap-3">
            <div
              class="w-8 h-8 rounded-full bg-indigo-100 text-indigo-700 font-semibold flex items-center justify-center text-xs">
              OL</div>
            <div>
              <p class="font-medium text-slate-900">Olivia Logan</p>
              <p class="text-xs text-slate-400">olivia@example.com</p>
            </div>
          </td>
          <td class="p-4">
            <span
              class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full bg-emerald-50 text-emerald-700">Success</span>
          </td>
          <td class="p-4 text-slate-500">Jul 17, 2026</td>
          <td class="p-4 text-right font-medium text-slate-900">$350.00</td>
        </tr>
        <tr class="hover:bg-slate-50/70 transition-colors">
          <td class="p-4 flex items-center gap-3">
            <div
              class="w-8 h-8 rounded-full bg-amber-100 text-amber-700 font-semibold flex items-center justify-center text-xs">
              JM</div>
            <div>
              <p class="font-medium text-slate-900">James Miller</p>
              <p class="text-xs text-slate-400">james@example.com</p>
            </div>
          </td>
          <td class="p-4">
            <span
              class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full bg-amber-50 text-amber-700">Pending</span>
          </td>
          <td class="p-4 text-slate-500">Jul 16, 2026</td>
          <td class="p-4 text-right font-medium text-slate-900">$1,200.00</td>
        </tr>
        <tr class="hover:bg-slate-50/70 transition-colors">
          <td class="p-4 flex items-center gap-3">
            <div
              class="w-8 h-8 rounded-full bg-rose-100 text-rose-700 font-semibold flex items-center justify-center text-xs">
              SR</div>
            <div>
              <p class="font-medium text-slate-900">Sophia Reed</p>
              <p class="text-xs text-slate-400">sophia@example.com</p>
            </div>
          </td>
          <td class="p-4">
            <span
              class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full bg-rose-50 text-rose-700">Failed</span>
          </td>
          <td class="p-4 text-slate-500">Jul 15, 2026</td>
          <td class="p-4 text-right font-medium text-slate-900">$85.50</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>

<!-- Activity / Sidebar Content (Takes up 1/3 col) -->
<div class="bg-white border border-slate-200 rounded-xl shadow-sm p-5 space-y-4">
  <h3 class="font-bold text-slate-900">Recent Activity</h3>
  <div class="flow-root">
    <ul class="-mb-8">
      <!-- Activity Item 1 -->
      <li>
        <div class="relative pb-8">
          <span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-slate-200" aria-hidden="true"></span>
          <div class="relative flex space-x-3">
            <div>
              <span
                class="h-8 w-8 rounded-full bg-emerald-50 text-emerald-600 flex items-center justify-center ring-8 ring-white">
                <i data-lucide="check-circle-2" class="w-4 h-4"></i>
              </span>
            </div>
            <div class="flex-1 min-w-0 pt-1.5 flex justify-between space-x-4">
              <div>
                <p class="text-sm text-slate-600">Database backup <span class="font-medium text-slate-900">successfully
                    completed</span></p>
              </div>
              <div class="text-right text-xs whitespace-nowrap text-slate-400">
                <time>1h ago</time>
              </div>
            </div>
          </div>
        </div>
      </li>
      <!-- Activity Item 2 -->
      <li>
        <div class="relative pb-8">
          <span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-slate-200" aria-hidden="true"></span>
          <div class="relative flex space-x-3">
            <div>
              <span
                class="h-8 w-8 rounded-full bg-indigo-50 text-indigo-600 flex items-center justify-center ring-8 ring-white">
                <i data-lucide="user-plus" class="w-4 h-4"></i>
              </span>
            </div>
            <div class="flex-1 min-w-0 pt-1.5 flex justify-between space-x-4">
              <div>
                <p class="text-sm text-slate-600">New user registration: <span class="font-medium text-slate-900">Sophia
                    Reed</span></p>
              </div>
              <div class="text-right text-xs whitespace-nowrap text-slate-400">
                <time>3h ago</time>
              </div>
            </div>
          </div>
        </div>
      </li>
      <!-- Activity Item 3 -->
      <li>
        <div class="relative pb-8">
          <div class="relative flex space-x-3">
            <div>
              <span
                class="h-8 w-8 rounded-full bg-amber-50 text-amber-600 flex items-center justify-center ring-8 ring-white">
                <i data-lucide="message-square" class="w-4 h-4"></i>
              </span>
            </div>
            <div class="flex-1 min-w-0 pt-1.5 flex justify-between space-x-4">
              <div>
                <p class="text-sm text-slate-600">New comment received on <span
                    class="font-medium text-slate-900">Quarterly Report v2</span></p>
              </div>
              <div class="text-right text-xs whitespace-nowrap text-slate-400">
                <time>Yesterday</time>
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</div>
</PrivateLayout>