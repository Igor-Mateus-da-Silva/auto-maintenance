import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { ApiService, HealthCheckResponse } from '../../core/services/api.service';

interface MetricCard {
  label: string;
  value: string;
  change: string;
  isPositive: boolean;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="space-y-8">
      <!-- Top banner / Welcome -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-zinc-800/60">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-zinc-100">Visão Geral</h1>
          <p class="text-sm text-zinc-400 mt-1">
            Controle simplificado e direto da frota e histórico de serviços.
          </p>
        </div>
        
        <div class="flex items-center gap-3">
          <a routerLink="/vehicles" class="inline-flex items-center justify-center px-4 py-2 text-sm font-medium rounded-lg bg-zinc-900 hover:bg-zinc-800 text-zinc-200 border border-zinc-800 transition">
            Gerenciar Veículos
          </a>
          <a routerLink="/maintenances" class="inline-flex items-center justify-center px-4 py-2 text-sm font-medium rounded-lg bg-blue-600 hover:bg-blue-500 text-white shadow-sm shadow-blue-900/30 transition">
            Nova Manutenção
          </a>
        </div>
      </div>

      <!-- Health check indicator with backend Go API -->
      <div class="p-4 rounded-xl border border-zinc-800/80 bg-zinc-900/40 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="w-2.5 h-2.5 rounded-full" [class.bg-emerald-400]="backendStatus() === 'online'" [class.bg-amber-400]="backendStatus() === 'checking'" [class.bg-red-400]="backendStatus() === 'offline'"></div>
          <div>
            <div class="text-sm font-medium text-zinc-200">
              Conexão com Backend Go
            </div>
            <div class="text-xs text-zinc-500">
              @if (backendStatus() === 'online') {
                API respondendo normalmente ({{ healthData()?.service }}).
              } @else if (backendStatus() === 'checking') {
                Verificando endpoint /api/health...
              } @else {
                Aguardando inicialização do servidor Go na porta 8080.
              }
            </div>
          </div>
        </div>

        <button (click)="checkBackend()" class="text-xs font-medium px-3 py-1.5 rounded-md bg-zinc-800 hover:bg-zinc-700 text-zinc-300 transition">
          Testar Conexão
        </button>
      </div>

      <!-- Metrics grid using @for -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        @for (metric of metrics(); track metric.label) {
          <div class="p-5 rounded-xl border border-zinc-800/80 bg-zinc-900/30 hover:border-zinc-700/80 transition flex flex-col justify-between">
            <span class="text-xs font-medium text-zinc-400 uppercase tracking-wider">{{ metric.label }}</span>
            <div class="mt-4 flex items-baseline justify-between">
              <span class="text-3xl font-semibold tracking-tight text-zinc-100">{{ metric.value }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full" [class.bg-emerald-500_10]="metric.isPositive" [class.text-emerald-400]="metric.isPositive" [class.bg-zinc-800]="!metric.isPositive" [class.text-zinc-400]="!metric.isPositive">
                {{ metric.change }}
              </span>
            </div>
          </div>
        }
      </div>

      <!-- Recent logs placeholder -->
      <div class="rounded-xl border border-zinc-800/80 bg-zinc-900/20 overflow-hidden">
        <div class="p-5 border-b border-zinc-800/80 flex items-center justify-between">
          <h2 class="text-base font-semibold text-zinc-200">Próximas Manutenções Programadas</h2>
          <span class="text-xs text-zinc-500">Task 1: Fundação Concluída</span>
        </div>
        <div class="p-8 text-center text-zinc-500 text-sm">
          A base de dados SQLite e o cadastro de veículos/manutenções serão integrados nas próximas etapas.
        </div>
      </div>
    </div>
  `
})
export class DashboardComponent implements OnInit {
  private readonly apiService = inject(ApiService);

  protected readonly backendStatus = signal<'checking' | 'online' | 'offline'>('checking');
  protected readonly healthData = signal<HealthCheckResponse | null>(null);

  protected readonly metrics = signal<MetricCard[]>([
    { label: 'Veículos Monitorados', value: '4', change: '+1 este mês', isPositive: true },
    { label: 'Manutenções no Mês', value: '2', change: 'Em dia', isPositive: true },
    { label: 'Gasto Total Estimado', value: 'R$ 1.450', change: 'Previsto', isPositive: false }
  ]);

  ngOnInit(): void {
    this.checkBackend();
  }

  checkBackend(): void {
    this.backendStatus.set('checking');
    this.apiService.checkHealth().subscribe({
      next: (res) => {
        this.healthData.set(res);
        this.backendStatus.set('online');
      },
      error: () => {
        this.backendStatus.set('offline');
      }
    });
  }
}
