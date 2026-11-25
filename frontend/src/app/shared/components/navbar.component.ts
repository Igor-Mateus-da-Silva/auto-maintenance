import { Component, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  template: `
    <header class="border-b border-zinc-800/80 bg-zinc-950/80 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-8">
          <a routerLink="/" class="flex items-center gap-3 group">
            <div class="w-9 h-9 rounded-lg bg-blue-600/10 border border-blue-500/20 flex items-center justify-center text-blue-400 group-hover:border-blue-500/40 transition">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
              </svg>
            </div>
            <div>
              <span class="font-semibold text-zinc-100 tracking-tight text-base">AutoCare</span>
              <span class="text-xs text-zinc-500 ml-1.5 font-normal">v1.0</span>
            </div>
          </a>

          <nav class="hidden md:flex items-center gap-1">
            <a routerLink="/dashboard" routerLinkActive="bg-zinc-800 text-zinc-100" class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition">
              Dashboard
            </a>
            <a routerLink="/vehicles" routerLinkActive="bg-zinc-800 text-zinc-100" class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition">
              Veículos
            </a>
            <a routerLink="/maintenances" routerLinkActive="bg-zinc-800 text-zinc-100" class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition">
              Manutenções
            </a>
          </nav>
        </div>

        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2 text-xs font-medium px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
            Sistema Ativo
          </div>
        </div>
      </div>
    </header>
  `
})
export class NavbarComponent {}
