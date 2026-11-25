import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { LucideAngularModule, KeyRound, Eye, EyeOff, Loader2 } from 'lucide-angular';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, LucideAngularModule],
  template: `
    <div class="min-h-screen bg-zinc-950 flex items-center justify-center p-4">
      <div class="bg-zinc-900/70 border border-zinc-800 rounded-xl p-8 max-w-sm w-full shadow-2xl shadow-black/50">
        <!-- Header -->
        <div class="flex flex-col items-center text-center mb-8">
          <div class="w-11 h-11 rounded-xl bg-zinc-800/80 border border-zinc-700/60 flex items-center justify-center text-zinc-200 mb-4">
            <lucide-icon [img]="KeyRoundIcon" class="w-5 h-5"></lucide-icon>
          </div>
          <h1 class="text-xl font-semibold tracking-tight text-zinc-100">Acesso ao Sistema</h1>
          <p class="text-xs text-zinc-400 mt-1.5">
            Controle de manutenção de veículos
          </p>
        </div>

        <!-- Form -->
        <form (ngSubmit)="onSubmit()" class="space-y-5">
          <div>
            <label for="access-key" class="block text-xs font-medium text-zinc-400 mb-2">
              Chave de Acesso
            </label>
            <div class="relative">
              <input
                id="access-key"
                name="accessKey"
                [type]="showPassword() ? 'text' : 'password'"
                [(ngModel)]="accessKey"
                autocomplete="current-password"
                required
                placeholder="Digite sua chave..."
                [disabled]="loading()"
                class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3.5 py-2.5 pr-10 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500 transition disabled:opacity-50"
              />
              <button
                type="button"
                (click)="togglePasswordVisibility()"
                [attr.aria-label]="showPassword() ? 'Ocultar chave de acesso' : 'Mostrar chave de acesso'"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-500 hover:text-zinc-300 transition"
              >
                @if (showPassword()) {
                  <lucide-icon [img]="EyeOffIcon" class="w-4 h-4"></lucide-icon>
                } @else {
                  <lucide-icon [img]="EyeIcon" class="w-4 h-4"></lucide-icon>
                }
              </button>
            </div>
          </div>

          @if (errorMessage()) {
            <p role="alert" class="text-xs text-rose-400 bg-rose-500/10 border border-rose-500/20 rounded-lg px-3 py-2">
              {{ errorMessage() }}
            </p>
          }

          <button
            type="submit"
            [disabled]="loading() || !accessKey.trim()"
            class="w-full flex items-center justify-center gap-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 font-medium text-sm py-2.5 px-4 transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
          >
            @if (loading()) {
              <lucide-icon [img]="LoaderIcon" class="w-4 h-4 animate-spin"></lucide-icon>
              <span>Entrando...</span>
            } @else {
              <span>Entrar</span>
            }
          </button>
        </form>
      </div>
    </div>
  `
})
export class LoginComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly KeyRoundIcon = KeyRound;
  protected readonly EyeIcon = Eye;
  protected readonly EyeOffIcon = EyeOff;
  protected readonly LoaderIcon = Loader2;

  protected accessKey = '';
  protected readonly showPassword = signal(false);
  protected readonly loading = signal(false);
  protected readonly errorMessage = signal<string | null>(null);

  togglePasswordVisibility(): void {
    this.showPassword.update((v) => !v);
  }

  onSubmit(): void {
    if (!this.accessKey.trim() || this.loading()) {
      return;
    }

    this.loading.set(true);
    this.errorMessage.set(null);

    this.authService.login(this.accessKey).subscribe({
      next: () => {
        this.loading.set(false);
        void this.router.navigate(['/dashboard']);
      },
      error: (err: HttpErrorResponse) => {
        this.loading.set(false);
        if (err.status === 401) {
          this.errorMessage.set('Chave de acesso inválida.');
        } else if (err.status === 0) {
          this.errorMessage.set('Servidor indisponível. Verifique se a API está em execução.');
        } else {
          this.errorMessage.set('Erro ao realizar login. Tente novamente.');
        }
      }
    });
  }
}
