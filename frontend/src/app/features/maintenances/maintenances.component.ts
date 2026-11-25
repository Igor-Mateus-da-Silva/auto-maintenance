import { Component } from '@angular/core';

@Component({
  selector: 'app-maintenances',
  standalone: true,
  template: `
    <div class="space-y-6">
      <div class="flex items-center justify-between pb-6 border-b border-zinc-800/60">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-zinc-100">Manutenções</h1>
          <p class="text-sm text-zinc-400 mt-1">Histórico e agendamento de serviços.</p>
        </div>
      </div>
      <div class="p-8 rounded-xl border border-zinc-800/80 bg-zinc-900/20 text-center text-zinc-400 text-sm">
        Módulo de manutenções pronto para receber os registros preventivos e corretivos.
      </div>
    </div>
  `
})
export class MaintenancesComponent {}
