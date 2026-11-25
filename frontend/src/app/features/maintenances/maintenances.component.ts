import {
  Component,
  ElementRef,
  OnInit,
  inject,
  signal,
  viewChild
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import {
  LucideWrench,
  LucidePlus,
  LucidePencil,
  LucideTrash,
  LucideClock,
  LucideGauge,
  LucideX,
  LucideCheck,
  LucideLoaderCircle
} from '@lucide/angular';
import { MaintenanceTypeService } from '../../core/services/maintenance-type.service';
import {
  ExpirationType,
  MaintenanceType
} from '../../shared/models/maintenance.model';

@Component({
  selector: 'app-maintenances',
  standalone: true,
  imports: [
    FormsModule,
    LucideWrench,
    LucidePlus,
    LucidePencil,
    LucideTrash,
    LucideClock,
    LucideGauge,
    LucideX,
    LucideCheck,
    LucideLoaderCircle
  ],
  template: `
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-zinc-800/70">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-zinc-100">Catálogo de Manutenções</h1>
          <p class="text-sm text-zinc-400 mt-1">
            Tipos de serviços reutilizáveis e regras de periodicidade.
          </p>
        </div>

        <button
          type="button"
          (click)="openCreateModal()"
          class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-sm font-medium transition cursor-pointer shrink-0"
        >
          <svg lucidePlus class="w-4 h-4"></svg>
          <span>Nova Manutenção</span>
        </button>
      </div>

      <!-- Loading Skeleton -->
      @if (loading() && items().length === 0) {
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 animate-pulse">
          <div class="h-28 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
          <div class="h-28 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
          <div class="h-28 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
        </div>
      } @else if (errorMessage()) {
        <div class="p-6 rounded-xl border border-rose-500/20 bg-rose-500/10 text-rose-300 text-sm flex items-center justify-between gap-4">
          <span>{{ errorMessage() }}</span>
          <button
            type="button"
            (click)="loadItems()"
            class="px-3 py-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-200 hover:bg-zinc-800 text-xs font-medium transition cursor-pointer"
          >
            Tentar novamente
          </button>
        </div>
      } @else if (items().length === 0) {
        <!-- Empty State -->
        <div class="py-16 px-6 rounded-xl border border-zinc-800/80 bg-zinc-900/30 text-center max-w-lg mx-auto">
          <div class="w-12 h-12 rounded-xl bg-zinc-800/80 border border-zinc-700/60 flex items-center justify-center text-zinc-400 mx-auto mb-4">
            <svg lucideWrench class="w-6 h-6"></svg>
          </div>
          <h2 class="text-base font-semibold text-zinc-100">Catálogo vazio</h2>
          <p class="text-sm text-zinc-400 mt-1.5 mb-6">
            Crie tipos de manutenção padrão (ex: Troca de Óleo, Filtro de Ar) para vincular aos seus veículos.
          </p>
          <button
            type="button"
            (click)="openCreateModal()"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-sm font-medium transition cursor-pointer"
          >
            <svg lucidePlus class="w-4 h-4"></svg>
            <span>Criar primeiro tipo</span>
          </button>
        </div>
      } @else {
        <!-- Catalog Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          @for (item of items(); track item.id) {
            <div class="p-5 rounded-xl border border-zinc-800/80 bg-zinc-900/30 hover:border-zinc-700/80 transition flex flex-col justify-between gap-4">
              <div class="space-y-2">
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <h2 class="text-base font-semibold text-zinc-100 truncate">{{ item.name }}</h2>
                    @if (item.description) {
                      <p class="text-xs text-zinc-400 mt-1 line-clamp-2">{{ item.description }}</p>
                    }
                  </div>
                  <span class="text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-zinc-800 text-zinc-300 border border-zinc-700/60 shrink-0">
                    {{ formatExpirationLabel(item.expiration_type) }}
                  </span>
                </div>

                <div class="flex flex-wrap items-center gap-4 pt-1 text-xs text-zinc-400">
                  @if (item.interval_days) {
                    <div class="flex items-center gap-1.5">
                      <svg lucideClock class="w-3.5 h-3.5 text-zinc-500"></svg>
                      <span>A cada {{ formatNumber(item.interval_days) }} dias</span>
                    </div>
                  }
                  @if (item.interval_km) {
                    <div class="flex items-center gap-1.5">
                      <svg lucideGauge class="w-3.5 h-3.5 text-zinc-500"></svg>
                      <span>A cada {{ formatNumber(item.interval_km) }} km</span>
                    </div>
                  }
                </div>
              </div>

              <div class="flex items-center justify-end gap-2 pt-3 border-t border-zinc-800/60">
                <button
                  type="button"
                  (click)="openEditModal(item)"
                  class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-zinc-800/70 hover:bg-zinc-800 text-zinc-300 text-xs font-medium transition cursor-pointer"
                >
                  <svg lucidePencil class="w-3.5 h-3.5"></svg>
                  <span>Editar</span>
                </button>
                <button
                  type="button"
                  (click)="openDeleteModal(item)"
                  class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 text-xs font-medium transition cursor-pointer"
                >
                  <svg lucideTrash class="w-3.5 h-3.5"></svg>
                  <span>Excluir</span>
                </button>
              </div>
            </div>
          }
        </div>
      }

      <!-- Create / Edit Modal (<dialog>) -->
      <dialog
        #formDialog
        closedby="any"
        aria-labelledby="mt-dialog-title"
        (close)="onFormDialogClose()"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-md w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        <div class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <h1 id="mt-dialog-title" class="text-base font-semibold text-zinc-100">
              {{ editingId() ? 'Editar Tipo de Manutenção' : 'Nova Manutenção' }}
            </h1>
            <button
              type="button"
              (click)="closeFormModal()"
              aria-label="Fechar modal"
              class="p-1 rounded-md text-zinc-500 hover:text-zinc-200 hover:bg-zinc-800 transition cursor-pointer"
            >
              <svg lucideX class="w-4 h-4"></svg>
            </button>
          </div>

          <form (ngSubmit)="submitForm()" class="space-y-4">
            <div>
              <label for="mt-name" class="block text-xs font-medium text-zinc-400 mb-1.5">
                Nome *
              </label>
              <input
                id="mt-name"
                type="text"
                name="name"
                [ngModel]="formName()"
                (ngModelChange)="formName.set($event)"
                required
                autofocus
                placeholder="Ex: Troca de Óleo"
                class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
              />
            </div>

            <div>
              <label for="mt-desc" class="block text-xs font-medium text-zinc-400 mb-1.5">
                Descrição (opcional)
              </label>
              <textarea
                id="mt-desc"
                name="description"
                rows="2"
                [ngModel]="formDescription()"
                (ngModelChange)="formDescription.set($event)"
                placeholder="Ex: Óleo sintético 5W30 + filtro"
                class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500 resize-none"
              ></textarea>
            </div>

            <div>
              <span class="block text-xs font-medium text-zinc-400 mb-1.5">
                Tipo de Expiração *
              </span>
              <div class="grid grid-cols-3 gap-1.5 p-1 rounded-lg bg-zinc-950 border border-zinc-800">
                <button
                  type="button"
                  (click)="formExpirationType.set('time')"
                  class="py-1.5 px-2 rounded-md text-xs font-medium transition cursor-pointer"
                  [class.bg-zinc-800]="formExpirationType() === 'time'"
                  [class.text-zinc-100]="formExpirationType() === 'time'"
                  [class.text-zinc-400]="formExpirationType() !== 'time'"
                >
                  Tempo
                </button>
                <button
                  type="button"
                  (click)="formExpirationType.set('km')"
                  class="py-1.5 px-2 rounded-md text-xs font-medium transition cursor-pointer"
                  [class.bg-zinc-800]="formExpirationType() === 'km'"
                  [class.text-zinc-100]="formExpirationType() === 'km'"
                  [class.text-zinc-400]="formExpirationType() !== 'km'"
                >
                  KM
                </button>
                <button
                  type="button"
                  (click)="formExpirationType.set('both')"
                  class="py-1.5 px-2 rounded-md text-xs font-medium transition cursor-pointer"
                  [class.bg-zinc-800]="formExpirationType() === 'both'"
                  [class.text-zinc-100]="formExpirationType() === 'both'"
                  [class.text-zinc-400]="formExpirationType() !== 'both'"
                >
                  Tempo e KM
                </button>
              </div>
            </div>

            <!-- Conditional Fields -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              @if (formExpirationType() === 'time' || formExpirationType() === 'both') {
                <div>
                  <label for="mt-days" class="block text-xs font-medium text-zinc-400 mb-1.5">
                    Intervalo em dias *
                  </label>
                  <input
                    id="mt-days"
                    type="number"
                    min="1"
                    name="intervalDays"
                    [ngModel]="formIntervalDays()"
                    (ngModelChange)="formIntervalDays.set($event)"
                    required
                    placeholder="Ex: 180"
                    class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
                  />
                </div>
              }

              @if (formExpirationType() === 'km' || formExpirationType() === 'both') {
                <div>
                  <label for="mt-km" class="block text-xs font-medium text-zinc-400 mb-1.5">
                    Intervalo em KM *
                  </label>
                  <input
                    id="mt-km"
                    type="number"
                    min="1"
                    name="intervalKm"
                    [ngModel]="formIntervalKm()"
                    (ngModelChange)="formIntervalKm.set($event)"
                    required
                    placeholder="Ex: 10000"
                    class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
                  />
                </div>
              }
            </div>

            @if (formError()) {
              <p role="alert" class="text-xs text-rose-400 bg-rose-500/10 border border-rose-500/20 rounded-lg px-3 py-2">
                {{ formError() }}
              </p>
            }

            <div class="flex items-center justify-end gap-2.5 pt-2">
              <button
                type="button"
                (click)="closeFormModal()"
                class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
              >
                Cancelar
              </button>
              <button
                type="submit"
                [disabled]="saving() || !isFormValid()"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-xs font-medium transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                @if (saving()) {
                  <svg lucideLoaderCircle class="w-3.5 h-3.5 animate-spin"></svg>
                  <span>Salvando...</span>
                } @else {
                  <svg lucideCheck class="w-3.5 h-3.5"></svg>
                  <span>{{ editingId() ? 'Salvar Alterações' : 'Cadastrar' }}</span>
                }
              </button>
            </div>
          </form>
        </div>
      </dialog>

      <!-- Delete Confirmation Modal (<dialog>) -->
      <dialog
        #deleteDialog
        closedby="any"
        aria-labelledby="mt-delete-title"
        (close)="deleteTarget.set(null)"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-sm w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        @if (deleteTarget(); as target) {
          <div class="space-y-4">
            <h1 id="mt-delete-title" class="text-base font-semibold text-zinc-100">
              Excluir tipo de manutenção?
            </h1>
            <p class="text-xs text-zinc-400 leading-relaxed">
              Tem certeza que deseja excluir <strong class="text-zinc-200">{{ target.name }}</strong>?
              Todos os vínculos com veículos que usam este tipo também serão removidos.
            </p>
            <div class="flex items-center justify-end gap-2.5 pt-2">
              <button
                type="button"
                autofocus
                (click)="closeDeleteModal()"
                class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
              >
                Cancelar
              </button>
              <button
                type="button"
                [disabled]="deleting()"
                (click)="confirmDelete()"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-rose-600 hover:bg-rose-500 text-white text-xs font-medium transition disabled:opacity-50 cursor-pointer"
              >
                @if (deleting()) {
                  <svg lucideLoaderCircle class="w-3.5 h-3.5 animate-spin"></svg>
                  <span>Excluindo...</span>
                } @else {
                  <span>Excluir</span>
                }
              </button>
            </div>
          </div>
        }
      </dialog>
    </div>
  `
})
export class MaintenancesComponent implements OnInit {
  private readonly maintenanceTypeService = inject(MaintenanceTypeService);

  protected readonly formDialog = viewChild<ElementRef<HTMLDialogElement>>('formDialog');
  protected readonly deleteDialog = viewChild<ElementRef<HTMLDialogElement>>('deleteDialog');

  // Icons

  // List state
  protected readonly items = signal<MaintenanceType[]>([]);
  protected readonly loading = signal<boolean>(true);
  protected readonly errorMessage = signal<string | null>(null);

  // Form modal state
  protected readonly editingId = signal<number | null>(null);
  protected readonly formName = signal<string>('');
  protected readonly formDescription = signal<string>('');
  protected readonly formExpirationType = signal<ExpirationType>('both');
  protected readonly formIntervalDays = signal<number | null>(180);
  protected readonly formIntervalKm = signal<number | null>(10000);
  protected readonly saving = signal<boolean>(false);
  protected readonly formError = signal<string | null>(null);

  // Delete modal state
  protected readonly deleteTarget = signal<MaintenanceType | null>(null);
  protected readonly deleting = signal<boolean>(false);

  ngOnInit(): void {
    this.loadItems();
  }

  loadItems(): void {
    this.loading.set(true);
    this.errorMessage.set(null);

    this.maintenanceTypeService.list().subscribe({
      next: (data) => {
        this.items.set(data);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.errorMessage.set('Não foi possível carregar o catálogo de manutenções.');
      }
    });
  }

  openCreateModal(): void {
    this.editingId.set(null);
    this.formName.set('');
    this.formDescription.set('');
    this.formExpirationType.set('both');
    this.formIntervalDays.set(180);
    this.formIntervalKm.set(10000);
    this.formError.set(null);

    const dialogEl = this.formDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  openEditModal(item: MaintenanceType): void {
    this.editingId.set(item.id);
    this.formName.set(item.name);
    this.formDescription.set(item.description ?? '');
    this.formExpirationType.set(item.expiration_type);
    this.formIntervalDays.set(item.interval_days ?? null);
    this.formIntervalKm.set(item.interval_km ?? null);
    this.formError.set(null);

    const dialogEl = this.formDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeFormModal(): void {
    const dialogEl = this.formDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  onFormDialogClose(): void {
    this.editingId.set(null);
    this.formError.set(null);
  }

  isFormValid(): boolean {
    if (!this.formName().trim()) {
      return false;
    }
    const type = this.formExpirationType();
    const days = Number(this.formIntervalDays());
    const km = Number(this.formIntervalKm());

    if ((type === 'time' || type === 'both') && (!days || days <= 0)) {
      return false;
    }
    if ((type === 'km' || type === 'both') && (!km || km <= 0)) {
      return false;
    }
    return true;
  }

  submitForm(): void {
    if (!this.isFormValid() || this.saving()) {
      return;
    }

    const type = this.formExpirationType();
    const desc = this.formDescription().trim();
    const days = type === 'time' || type === 'both' ? Number(this.formIntervalDays()) : null;
    const km = type === 'km' || type === 'both' ? Number(this.formIntervalKm()) : null;

    const payload = {
      name: this.formName().trim(),
      description: desc ? desc : null,
      expiration_type: type,
      interval_days: days,
      interval_km: km
    };

    this.saving.set(true);
    this.formError.set(null);

    const id = this.editingId();
    const request$ = id
      ? this.maintenanceTypeService.update(id, payload)
      : this.maintenanceTypeService.create(payload);

    request$.subscribe({
      next: () => {
        this.saving.set(false);
        this.closeFormModal();
        this.loadItems();
      },
      error: (err: HttpErrorResponse) => {
        this.saving.set(false);
        this.formError.set(err.error?.error || 'Erro ao salvar tipo de manutenção.');
      }
    });
  }

  openDeleteModal(item: MaintenanceType): void {
    this.deleteTarget.set(item);
    const dialogEl = this.deleteDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeDeleteModal(): void {
    const dialogEl = this.deleteDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  confirmDelete(): void {
    const target = this.deleteTarget();
    if (!target || this.deleting()) {
      return;
    }

    this.deleting.set(true);
    this.maintenanceTypeService.delete(target.id).subscribe({
      next: () => {
        this.deleting.set(false);
        this.closeDeleteModal();
        this.loadItems();
      },
      error: () => {
        this.deleting.set(false);
      }
    });
  }

  formatExpirationLabel(type: ExpirationType): string {
    switch (type) {
      case 'time':
        return 'Tempo';
      case 'km':
        return 'KM';
      case 'both':
        return 'Tempo e KM';
    }
  }

  formatNumber(value: number): string {
    return value.toLocaleString('pt-BR');
  }
}
