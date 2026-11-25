import {
  Component,
  ElementRef,
  OnInit,
  computed,
  inject,
  signal,
  viewChild
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import {
  LucideAngularModule,
  Car,
  Gauge,
  CheckCircle2,
  AlertTriangle,
  AlertCircle,
  Clock,
  RefreshCw,
  Check,
  X,
  Plus,
  Loader2
} from 'lucide-angular';
import { DashboardService } from '../../core/services/dashboard.service';
import { VehicleService } from '../../core/services/vehicle.service';
import {
  DashboardResponse,
  MaintenanceStatus,
  VehicleMaintenanceDetail
} from '../../shared/models/maintenance.model';
import { Vehicle } from '../../shared/models/vehicle.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [FormsModule, RouterLink, LucideAngularModule],
  template: `
    <div class="space-y-8">
      <!-- A. Header & Vehicle Filter -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-zinc-800/70">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-zinc-100">Dashboard</h1>
          <p class="text-sm text-zinc-400 mt-1">
            Acompanhamento de vencimentos e quilometragem da frota.
          </p>
        </div>

        @if (dashboard()?.vehicles?.length) {
          <div class="flex items-center gap-2.5">
            <label for="vehicle-filter" class="text-xs font-medium text-zinc-400 whitespace-nowrap">
              Veículo:
            </label>
            <select
              id="vehicle-filter"
              [value]="selectedVehicleId() ?? ''"
              (change)="onVehicleFilterChange($event)"
              class="rounded-lg bg-zinc-900 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-600 transition min-w-48"
            >
              <option value="">Todos os veículos</option>
              @for (vehicle of dashboard()!.vehicles; track vehicle.id) {
                <option [value]="vehicle.id">{{ vehicle.name }}</option>
              }
            </select>
          </div>
        }
      </div>

      <!-- Loading Skeleton -->
      @if (loading() && !dashboard()) {
        <div class="space-y-6 animate-pulse">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="h-20 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
            <div class="h-20 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
            <div class="h-20 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
          </div>
          <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <div class="h-24 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
            <div class="h-24 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
            <div class="h-24 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
            <div class="h-24 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
          </div>
          <div class="h-48 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
        </div>
      } @else if (errorMessage()) {
        <div class="p-6 rounded-xl border border-rose-500/20 bg-rose-500/10 text-rose-300 text-sm flex items-center justify-between gap-4">
          <span>{{ errorMessage() }}</span>
          <button
            type="button"
            (click)="loadDashboard()"
            class="px-3 py-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-200 hover:bg-zinc-800 text-xs font-medium transition cursor-pointer"
          >
            Tentar novamente
          </button>
        </div>
      } @else if (dashboard(); as dash) {
        <!-- Empty State: No Vehicles -->
        @if (dash.vehicles.length === 0) {
          <div class="py-16 px-6 rounded-xl border border-zinc-800/80 bg-zinc-900/30 text-center max-w-lg mx-auto">
            <div class="w-12 h-12 rounded-xl bg-zinc-800/80 border border-zinc-700/60 flex items-center justify-center text-zinc-400 mx-auto mb-4">
              <lucide-icon [img]="CarIcon" class="w-6 h-6"></lucide-icon>
            </div>
            <h2 class="text-base font-semibold text-zinc-100">Nenhum veículo cadastrado</h2>
            <p class="text-sm text-zinc-400 mt-1.5 mb-6">
              Cadastre seu primeiro veículo para começar a acompanhar quilometragem e manutenções.
            </p>
            <a
              routerLink="/vehicles"
              class="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-sm font-medium transition"
            >
              <lucide-icon [img]="PlusIcon" class="w-4 h-4"></lucide-icon>
              <span>Cadastrar primeiro veículo</span>
            </a>
          </div>
        } @else {
          <!-- B. Quick Mileage Update (KM Rápido) -->
          <section aria-label="Atualização rápida de quilometragem" class="space-y-3">
            <h2 class="text-xs font-medium uppercase tracking-wider text-zinc-500">
              Quilometragem Atual (KM Rápido)
            </h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              @for (vehicle of dash.vehicles; track vehicle.id) {
                <div class="p-4 rounded-xl border border-zinc-800/80 bg-zinc-900/40 flex flex-col justify-between gap-3">
                  <div class="flex items-start justify-between gap-2">
                    <div class="min-w-0">
                      <div class="text-sm font-medium text-zinc-100 truncate">{{ vehicle.name }}</div>
                      @if (vehicle.plate) {
                        <div class="text-xs text-zinc-500 font-mono mt-0.5">{{ vehicle.plate }}</div>
                      }
                    </div>
                    @if (mileageFeedback() === vehicle.id) {
                      <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 shrink-0">
                        KM atualizado
                      </span>
                    }
                  </div>

                  @if (editingVehicleId() === vehicle.id) {
                    <form (ngSubmit)="saveQuickMileage(vehicle.id)" class="flex items-center gap-2">
                      <input
                        type="number"
                        min="0"
                        [ngModel]="editingMileageValue()"
                        (ngModelChange)="editingMileageValue.set($event)"
                        name="mileage-{{ vehicle.id }}"
                        required
                        class="w-full rounded-lg bg-zinc-950 border border-zinc-700 px-2.5 py-1.5 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                      />
                      <button
                        type="submit"
                        [disabled]="savingMileageId() === vehicle.id"
                        aria-label="Salvar quilometragem"
                        class="p-1.5 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 transition disabled:opacity-50 cursor-pointer shrink-0"
                      >
                        @if (savingMileageId() === vehicle.id) {
                          <lucide-icon [img]="LoaderIcon" class="w-4 h-4 animate-spin"></lucide-icon>
                        } @else {
                          <lucide-icon [img]="CheckIcon" class="w-4 h-4"></lucide-icon>
                        }
                      </button>
                      <button
                        type="button"
                        (click)="cancelQuickMileage()"
                        aria-label="Cancelar edição"
                        class="p-1.5 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-300 transition cursor-pointer shrink-0"
                      >
                        <lucide-icon [img]="XIcon" class="w-4 h-4"></lucide-icon>
                      </button>
                    </form>
                  } @else {
                    <div class="flex items-center justify-between gap-2 pt-1">
                      <div class="flex items-center gap-1.5 text-zinc-200 font-semibold text-base tabular-nums">
                        <lucide-icon [img]="GaugeIcon" class="w-4 h-4 text-zinc-500"></lucide-icon>
                        <span>{{ formatNumber(vehicle.current_mileage) }} km</span>
                      </div>
                      <button
                        type="button"
                        (click)="startQuickMileage(vehicle)"
                        class="text-xs font-medium px-2.5 py-1.5 rounded-lg bg-zinc-800/80 hover:bg-zinc-800 text-zinc-300 border border-zinc-700/60 transition cursor-pointer"
                      >
                        Atualizar KM
                      </button>
                    </div>
                  }
                </div>
              }
            </div>
          </section>

          <!-- C. Summary Filter Cards -->
          <section aria-label="Resumo de manutenções" class="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <!-- Total -->
            <button
              type="button"
              (click)="setStatusFilter('all')"
              class="text-left p-4 rounded-xl border transition cursor-pointer"
              [class.border-zinc-500]="statusFilter() === 'all'"
              [class.bg-zinc-900]="statusFilter() === 'all'"
              [class.border-zinc-800/80]="statusFilter() !== 'all'"
              [class.bg-zinc-900/30]="statusFilter() !== 'all'"
              [class.hover:border-zinc-700]="statusFilter() !== 'all'"
            >
              <span class="text-xs font-medium text-zinc-400">Total</span>
              <div class="mt-2 text-2xl font-semibold text-zinc-100 tabular-nums">
                {{ dash.summary.total }}
              </div>
            </button>

            <!-- Em dia -->
            <button
              type="button"
              (click)="setStatusFilter('ok')"
              class="text-left p-4 rounded-xl border transition cursor-pointer"
              [class.border-emerald-500/60]="statusFilter() === 'ok'"
              [class.bg-emerald-500/5]="statusFilter() === 'ok'"
              [class.border-zinc-800/80]="statusFilter() !== 'ok'"
              [class.bg-zinc-900/30]="statusFilter() !== 'ok'"
              [class.hover:border-zinc-700]="statusFilter() !== 'ok'"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-emerald-400">Em dia</span>
                <lucide-icon [img]="CheckCircleIcon" class="w-4 h-4 text-emerald-500/70"></lucide-icon>
              </div>
              <div class="mt-2 text-2xl font-semibold text-zinc-100 tabular-nums">
                {{ dash.summary.ok }}
              </div>
            </button>

            <!-- Próximas de vencer -->
            <button
              type="button"
              (click)="setStatusFilter('warning')"
              class="text-left p-4 rounded-xl border transition cursor-pointer"
              [class.border-amber-500/60]="statusFilter() === 'warning'"
              [class.bg-amber-500/5]="statusFilter() === 'warning'"
              [class.border-zinc-800/80]="statusFilter() !== 'warning'"
              [class.bg-zinc-900/30]="statusFilter() !== 'warning'"
              [class.hover:border-zinc-700]="statusFilter() !== 'warning'"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-amber-400">Próximas de vencer</span>
                <lucide-icon [img]="AlertTriangleIcon" class="w-4 h-4 text-amber-500/70"></lucide-icon>
              </div>
              <div class="mt-2 text-2xl font-semibold text-zinc-100 tabular-nums">
                {{ dash.summary.warning }}
              </div>
            </button>

            <!-- Vencidas -->
            <button
              type="button"
              (click)="setStatusFilter('overdue')"
              class="text-left p-4 rounded-xl border transition cursor-pointer"
              [class.border-rose-500/60]="statusFilter() === 'overdue'"
              [class.bg-rose-500/5]="statusFilter() === 'overdue'"
              [class.border-zinc-800/80]="statusFilter() !== 'overdue'"
              [class.bg-zinc-900/30]="statusFilter() !== 'overdue'"
              [class.hover:border-zinc-700]="statusFilter() !== 'overdue'"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-rose-400">Vencidas</span>
                <lucide-icon [img]="AlertCircleIcon" class="w-4 h-4 text-rose-500/70"></lucide-icon>
              </div>
              <div class="mt-2 text-2xl font-semibold text-zinc-100 tabular-nums">
                {{ dash.summary.overdue }}
              </div>
            </button>
          </section>

          <!-- D. Maintenances List (Sorted by Urgency) -->
          <section aria-label="Lista de manutenções" class="space-y-3">
            <div class="flex items-center justify-between">
              <h2 class="text-xs font-medium uppercase tracking-wider text-zinc-500">
                Manutenções ({{ filteredMaintenances().length }})
              </h2>
            </div>

            @if (filteredMaintenances().length === 0) {
              <div class="py-12 px-6 rounded-xl border border-zinc-800/80 bg-zinc-900/20 text-center">
                @if (dash.summary.total === 0) {
                  <p class="text-sm text-zinc-400">
                    Nenhuma manutenção vinculada ainda.
                  </p>
                  <a
                    routerLink="/vehicles"
                    class="inline-flex items-center gap-1.5 mt-3 text-xs font-medium text-zinc-200 hover:text-white underline underline-offset-4"
                  >
                    Ir para Veículos e vincular manutenções
                  </a>
                } @else {
                  <p class="text-sm text-zinc-400">
                    Nenhuma manutenção encontrada neste filtro de status.
                  </p>
                }
              </div>
            } @else {
              <div class="divide-y divide-zinc-800/70 rounded-xl border border-zinc-800/80 bg-zinc-900/30 overflow-hidden">
                @for (item of filteredMaintenances(); track item.id) {
                  <div class="p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-zinc-900/50 transition">
                    <div class="space-y-2 min-w-0">
                      <div class="flex flex-wrap items-center gap-2.5">
                        <!-- Status Badge -->
                        @switch (item.status) {
                          @case ('overdue') {
                            <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-rose-500/10 text-rose-400 border border-rose-500/20">
                              Vencida
                            </span>
                          }
                          @case ('warning') {
                            <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/20">
                              Próxima
                            </span>
                          }
                          @case ('ok') {
                            <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                              Em dia
                            </span>
                          }
                        }

                        <h3 class="text-sm font-semibold text-zinc-100 truncate">
                          {{ item.maintenance_type_name }}
                        </h3>

                        @if (!selectedVehicleId()) {
                          <span class="text-xs text-zinc-500 bg-zinc-800/60 px-2 py-0.5 rounded-md">
                            {{ item.vehicle_name }}
                          </span>
                        }
                      </div>

                      <!-- Expiration Indicators -->
                      <div class="flex flex-wrap items-center gap-x-5 gap-y-1.5 text-xs text-zinc-400">
                        @if (item.expiration_type === 'time' || item.expiration_type === 'both') {
                          <div class="flex items-center gap-1.5">
                            <lucide-icon [img]="ClockIcon" class="w-3.5 h-3.5 text-zinc-500"></lucide-icon>
                            <span>{{ formatTimeStatus(item) }}</span>
                          </div>
                        }

                        @if (item.expiration_type === 'km' || item.expiration_type === 'both') {
                          <div class="flex items-center gap-1.5">
                            <lucide-icon [img]="GaugeIcon" class="w-3.5 h-3.5 text-zinc-500"></lucide-icon>
                            <span>{{ formatKmStatus(item) }}</span>
                          </div>
                        }
                      </div>
                    </div>

                    <!-- Quick Renew Button -->
                    <div class="shrink-0">
                      <button
                        type="button"
                        (click)="openRenewModal(item)"
                        class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-zinc-800/80 hover:bg-zinc-800 text-zinc-200 border border-zinc-700/60 text-xs font-medium transition cursor-pointer"
                      >
                        <lucide-icon [img]="RefreshIcon" class="w-3.5 h-3.5"></lucide-icon>
                        <span>Renovar</span>
                      </button>
                    </div>
                  </div>
                }
              </div>
            }
          </section>
        }
      }

      <!-- Native <dialog> Modal for Renew Maintenance -->
      <dialog
        #renewDialog
        closedby="any"
        aria-labelledby="renew-dialog-title"
        (close)="onDialogClose()"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-md w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        @if (renewTarget(); as target) {
          <div class="space-y-5">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h1 id="renew-dialog-title" class="text-base font-semibold text-zinc-100">
                  Renovar Manutenção
                </h1>
                <p class="text-xs text-zinc-400 mt-1">
                  {{ target.maintenance_type_name }} &middot; {{ target.vehicle_name }}
                </p>
              </div>
              <button
                type="button"
                (click)="closeRenewModal()"
                aria-label="Fechar modal"
                class="p-1 rounded-md text-zinc-500 hover:text-zinc-200 hover:bg-zinc-800 transition cursor-pointer"
              >
                <lucide-icon [img]="XIcon" class="w-4 h-4"></lucide-icon>
              </button>
            </div>

            <form (ngSubmit)="submitRenew()" class="space-y-4">
              <div>
                <label for="renew-date" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Data de Realização
                </label>
                <input
                  id="renew-date"
                  type="date"
                  name="renewDate"
                  [ngModel]="renewDate()"
                  (ngModelChange)="renewDate.set($event)"
                  required
                  autofocus
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                />
              </div>

              <div>
                <label for="renew-mileage" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Quilometragem na Realização (km)
                </label>
                <input
                  id="renew-mileage"
                  type="number"
                  min="0"
                  name="renewMileage"
                  [ngModel]="renewMileage()"
                  (ngModelChange)="renewMileage.set($event)"
                  required
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                />
              </div>

              <div>
                <label for="renew-notes" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Observações (opcional)
                </label>
                <textarea
                  id="renew-notes"
                  name="renewNotes"
                  rows="2"
                  [ngModel]="renewNotes()"
                  (ngModelChange)="renewNotes.set($event)"
                  placeholder="Ex: Troca com filtro original"
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500 resize-none"
                ></textarea>
              </div>

              @if (renewError()) {
                <p role="alert" class="text-xs text-rose-400 bg-rose-500/10 border border-rose-500/20 rounded-lg px-3 py-2">
                  {{ renewError() }}
                </p>
              }

              <div class="flex items-center justify-end gap-2.5 pt-2">
                <button
                  type="button"
                  (click)="closeRenewModal()"
                  class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  [disabled]="savingRenew()"
                  class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-xs font-medium transition disabled:opacity-50 cursor-pointer"
                >
                  @if (savingRenew()) {
                    <lucide-icon [img]="LoaderIcon" class="w-3.5 h-3.5 animate-spin"></lucide-icon>
                    <span>Salvando...</span>
                  } @else {
                    <lucide-icon [img]="CheckIcon" class="w-3.5 h-3.5"></lucide-icon>
                    <span>Confirmar Renovação</span>
                  }
                </button>
              </div>
            </form>
          </div>
        }
      </dialog>
    </div>
  `
})
export class DashboardComponent implements OnInit {
  private readonly dashboardService = inject(DashboardService);
  private readonly vehicleService = inject(VehicleService);

  protected readonly renewDialog = viewChild<ElementRef<HTMLDialogElement>>('renewDialog');

  // Icons
  protected readonly CarIcon = Car;
  protected readonly GaugeIcon = Gauge;
  protected readonly CheckCircleIcon = CheckCircle2;
  protected readonly AlertTriangleIcon = AlertTriangle;
  protected readonly AlertCircleIcon = AlertCircle;
  protected readonly ClockIcon = Clock;
  protected readonly RefreshIcon = RefreshCw;
  protected readonly CheckIcon = Check;
  protected readonly XIcon = X;
  protected readonly PlusIcon = Plus;
  protected readonly LoaderIcon = Loader2;

  // Dashboard state
  protected readonly dashboard = signal<DashboardResponse | null>(null);
  protected readonly loading = signal<boolean>(true);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly selectedVehicleId = signal<number | undefined>(undefined);
  protected readonly statusFilter = signal<'all' | MaintenanceStatus>('all');

  protected readonly filteredMaintenances = computed(() => {
    const items = this.dashboard()?.maintenances ?? [];
    const filter = this.statusFilter();
    if (filter === 'all') {
      return items;
    }
    return items.filter((item) => item.status === filter);
  });

  // Quick mileage state
  protected readonly editingVehicleId = signal<number | null>(null);
  protected readonly editingMileageValue = signal<number>(0);
  protected readonly savingMileageId = signal<number | null>(null);
  protected readonly mileageFeedback = signal<number | null>(null);

  // Renew modal state
  protected readonly renewTarget = signal<VehicleMaintenanceDetail | null>(null);
  protected readonly renewDate = signal<string>('');
  protected readonly renewMileage = signal<number>(0);
  protected readonly renewNotes = signal<string>('');
  protected readonly savingRenew = signal<boolean>(false);
  protected readonly renewError = signal<string | null>(null);

  ngOnInit(): void {
    this.loadDashboard();
  }

  loadDashboard(): void {
    this.loading.set(true);
    this.errorMessage.set(null);

    this.dashboardService.getDashboard(this.selectedVehicleId()).subscribe({
      next: (res) => {
        this.dashboard.set(res);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.errorMessage.set('Não foi possível carregar os dados do dashboard.');
      }
    });
  }

  onVehicleFilterChange(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    this.selectedVehicleId.set(value ? Number(value) : undefined);
    this.loadDashboard();
  }

  setStatusFilter(status: 'all' | MaintenanceStatus): void {
    this.statusFilter.set(status);
  }

  // Quick Mileage Actions
  startQuickMileage(vehicle: Vehicle): void {
    this.editingVehicleId.set(vehicle.id);
    this.editingMileageValue.set(vehicle.current_mileage);
  }

  cancelQuickMileage(): void {
    this.editingVehicleId.set(null);
  }

  saveQuickMileage(vehicleId: number): void {
    const mileage = Number(this.editingMileageValue());
    if (isNaN(mileage) || mileage < 0 || this.savingMileageId() !== null) {
      return;
    }

    this.savingMileageId.set(vehicleId);
    this.vehicleService.updateMileage(vehicleId, mileage).subscribe({
      next: () => {
        this.savingMileageId.set(null);
        this.editingVehicleId.set(null);
        this.mileageFeedback.set(vehicleId);
        setTimeout(() => {
          if (this.mileageFeedback() === vehicleId) {
            this.mileageFeedback.set(null);
          }
        }, 3000);
        this.loadDashboard();
      },
      error: () => {
        this.savingMileageId.set(null);
      }
    });
  }

  // Renew Modal Actions
  openRenewModal(item: VehicleMaintenanceDetail): void {
    const today = new Date().toISOString().slice(0, 10);
    this.renewTarget.set(item);
    this.renewDate.set(today);
    this.renewMileage.set(item.vehicle_current_mileage);
    this.renewNotes.set(item.notes ?? '');
    this.renewError.set(null);

    const dialogEl = this.renewDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeRenewModal(): void {
    const dialogEl = this.renewDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  onDialogClose(): void {
    this.renewTarget.set(null);
    this.renewError.set(null);
  }

  submitRenew(): void {
    const target = this.renewTarget();
    if (!target || this.savingRenew()) {
      return;
    }

    const mileage = Number(this.renewMileage());
    if (!this.renewDate() || isNaN(mileage) || mileage < 0) {
      this.renewError.set('Preencha uma data e quilometragem válidas.');
      return;
    }

    this.savingRenew.set(true);
    this.renewError.set(null);

    const notesTrimmed = this.renewNotes().trim();

    this.vehicleService
      .renewMaintenance(target.id, {
        last_performed_date: this.renewDate(),
        last_performed_mileage: mileage,
        notes: notesTrimmed ? notesTrimmed : null
      })
      .subscribe({
        next: () => {
          this.savingRenew.set(false);
          this.closeRenewModal();
          this.loadDashboard();
        },
        error: () => {
          this.savingRenew.set(false);
          this.renewError.set('Erro ao renovar manutenção. Verifique os dados e tente novamente.');
        }
      });
  }

  // Formatting Helpers
  formatNumber(value: number | null | undefined): string {
    if (value === null || value === undefined) {
      return '0';
    }
    return value.toLocaleString('pt-BR');
  }

  formatDate(isoDate: string | null | undefined): string {
    if (!isoDate) {
      return '';
    }
    const [year, month, day] = isoDate.split('-');
    if (!year || !month || !day) {
      return isoDate;
    }
    return `${day}/${month}/${year}`;
  }

  formatTimeStatus(item: VehicleMaintenanceDetail): string {
    const days = item.days_remaining;
    const dueDate = this.formatDate(item.next_due_date);

    if (days === null || days === undefined) {
      return '';
    }
    if (days < 0) {
      const absDays = Math.abs(days);
      return `Vencida há ${absDays} ${absDays === 1 ? 'dia' : 'dias'}`;
    }
    if (days === 0) {
      return `Vence hoje (${dueDate})`;
    }
    return `Vence em ${days} ${days === 1 ? 'dia' : 'dias'} (${dueDate})`;
  }

  formatKmStatus(item: VehicleMaintenanceDetail): string {
    const km = item.km_remaining;
    const nextKm = this.formatNumber(item.next_due_mileage);

    if (km === null || km === undefined) {
      return '';
    }
    if (km < 0) {
      return `Passou ${this.formatNumber(Math.abs(km))} km do limite`;
    }
    if (km === 0) {
      return `Limite atingido (${nextKm} km)`;
    }
    return `Faltam ${this.formatNumber(km)} km (aos ${nextKm} km)`;
  }
}
