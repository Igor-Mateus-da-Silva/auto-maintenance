import { Component } from '@angular/core';

@Component({
  selector: 'app-vehicles',
  standalone: true,
  template: `
    <div class="space-y-6">
      <div class="flex items-center justify-between pb-6 border-b border-zinc-800/60">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-zinc-100">Veículos</h1>
          <p class="text-sm text-zinc-400 mt-1">Gerenciamento de veículos da frota.</p>
        </div>
      </div>
      <div class="p-8 rounded-xl border border-zinc-800/80 bg-zinc-900/20 text-center text-zinc-400 text-sm">
        Módulo de veículos pronto para receber listagem e cadastro.
      </div>
    </div>
  `
})
export class VehiclesComponent {}
