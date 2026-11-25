import { Component } from '@angular/core';

@Component({
  selector: 'app-login',
  standalone: true,
  template: `
    <div class="max-w-md mx-auto my-16 p-8 rounded-2xl border border-zinc-800/80 bg-zinc-900/40">
      <h2 class="text-xl font-bold text-zinc-100">Acesso ao Sistema</h2>
      <p class="text-sm text-zinc-400 mt-1">Autenticação simplificada.</p>
    </div>
  `
})
export class LoginComponent {}
