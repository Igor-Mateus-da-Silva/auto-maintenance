import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { LucideAngularModule, Car, LogOut } from 'lucide-angular';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, LucideAngularModule],
  template: `
    <header class="border-b border-zinc-800/80 bg-zinc-950/80 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-8">
          <a routerLink="/dashboard" class="flex items-center gap-2.5 group">
            <div class="w-8 h-8 rounded-lg bg-zinc-900 border border-zinc-800 flex items-center justify-center text-zinc-200 group-hover:border-zinc-700 transition">
              <lucide-icon [img]="CarIcon" class="w-4 h-4"></lucide-icon>
            </div>
            <span class="font-semibold text-zinc-100 tracking-tight text-sm sm:text-base">AutoCare</span>
          </a>

          <nav class="flex items-center gap-1">
            <a
              routerLink="/dashboard"
              routerLinkActive="bg-zinc-800/90 text-zinc-100"
              class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition"
            >
              Dashboard
            </a>
            <a
              routerLink="/vehicles"
              routerLinkActive="bg-zinc-800/90 text-zinc-100"
              class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition"
            >
              Veículos
            </a>
            <a
              routerLink="/maintenances"
              routerLinkActive="bg-zinc-800/90 text-zinc-100"
              class="px-3 py-1.5 rounded-md text-sm font-medium text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900 transition"
            >
              Manutenções
            </a>
          </nav>
        </div>

        <button
          type="button"
          (click)="logout()"
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-md text-xs font-medium text-zinc-400 hover:text-zinc-100 hover:bg-zinc-900 border border-transparent hover:border-zinc-800 transition cursor-pointer"
        >
          <lucide-icon [img]="LogOutIcon" class="w-3.5 h-3.5"></lucide-icon>
          <span>Sair</span>
        </button>
      </div>
    </header>
  `
})
export class NavbarComponent {
  private readonly authService = inject(AuthService);

  protected readonly CarIcon = Car;
  protected readonly LogOutIcon = LogOut;

  logout(): void {
    this.authService.logout();
  }
}
