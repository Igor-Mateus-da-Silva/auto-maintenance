import {
  Component,
  ElementRef,
  OnInit,
  inject,
  signal,
  viewChild
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import {
  LucideCar,
  LucidePlus,
  LucidePencil,
  LucideTrash,
  LucideWrench,
  LucideGauge,
  LucideCalendar,
  LucideX,
  LucideCheck,
  LucideLoaderCircle,
  LucideCircleExclamationPoint
} from '@lucide/angular';
import { VehicleService } from '../../core/services/vehicle.service';
import { MaintenanceTypeService } from '../../core/services/maintenance-type.service';
import { Vehicle } from '../../shared/models/vehicle.model';
import {
  MaintenanceType,
  VehicleMaintenanceDetail
} from '../../shared/models/maintenance.model';

@Component({
  selector: 'app-vehicles',
  standalone: true,
  imports: [
    FormsModule,
    RouterLink,
    LucideCar,
    LucidePlus,
    LucidePencil,
    LucideTrash,
    LucideWrench,
    LucideGauge,
    LucideCalendar,
    LucideX,
    LucideCheck,
    LucideLoaderCircle,
    LucideCircleExclamationPoint
  ],
  template: `
    <div class="space-y-8">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-zinc-800/70">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-zinc-100">Veículos</h1>
          <p class="text-sm text-zinc-400 mt-1">
            Gerencie seus veículos e vincule as manutenções preventivas.
          </p>
        </div>

        <button
          type="button"
          (click)="openCreateVehicleModal()"
          class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-sm font-medium transition cursor-pointer shrink-0"
        >
          <svg lucidePlus class="w-4 h-4"></svg>
          <span>Novo Veículo</span>
        </button>
      </div>

      <!-- Loading Skeleton -->
      @if (loading() && vehicles().length === 0) {
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 animate-pulse">
          <div class="h-36 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
          <div class="h-36 rounded-xl bg-zinc-900/60 border border-zinc-800/60"></div>
        </div>
      } @else if (errorMessage()) {
        <div class="p-6 rounded-xl border border-rose-500/20 bg-rose-500/10 text-rose-300 text-sm flex items-center justify-between gap-4">
          <span>{{ errorMessage() }}</span>
          <button
            type="button"
            (click)="loadVehicles()"
            class="px-3 py-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-200 hover:bg-zinc-800 text-xs font-medium transition cursor-pointer"
          >
            Tentar novamente
          </button>
        </div>
      } @else if (vehicles().length === 0) {
        <!-- Empty State -->
        <div class="py-16 px-6 rounded-xl border border-zinc-800/80 bg-zinc-900/30 text-center max-w-lg mx-auto">
          <div class="w-12 h-12 rounded-xl bg-zinc-800/80 border border-zinc-700/60 flex items-center justify-center text-zinc-400 mx-auto mb-4">
            <svg lucideCar class="w-6 h-6"></svg>
          </div>
          <h2 class="text-base font-semibold text-zinc-100">Nenhum veículo cadastrado</h2>
          <p class="text-sm text-zinc-400 mt-1.5 mb-6">
            Adicione seu primeiro veículo para começar a gerenciar manutenções.
          </p>
          <button
            type="button"
            (click)="openCreateVehicleModal()"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-sm font-medium transition cursor-pointer"
          >
            <svg lucidePlus class="w-4 h-4"></svg>
            <span>Cadastrar primeiro veículo</span>
          </button>
        </div>
      } @else {
        <!-- Vehicles Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          @for (vehicle of vehicles(); track vehicle.id) {
            <div
              class="p-5 rounded-xl border transition flex flex-col justify-between gap-4"
              [class.border-zinc-500]="selectedVehicle()?.id === vehicle.id"
              [class.bg-zinc-900]="selectedVehicle()?.id === vehicle.id"
              [class.border-zinc-800/80]="selectedVehicle()?.id !== vehicle.id"
              [class.bg-zinc-900/30]="selectedVehicle()?.id !== vehicle.id"
            >
              <div class="space-y-3">
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <h2 class="text-base font-semibold text-zinc-100 truncate">{{ vehicle.name }}</h2>
                    <div class="flex items-center gap-2 text-xs text-zinc-500 mt-1">
                      @if (vehicle.plate) {
                        <span class="font-mono bg-zinc-800/70 text-zinc-300 px-1.5 py-0.5 rounded">
                          {{ vehicle.plate }}
                        </span>
                      }
                      @if (vehicle.year) {
                        <span>{{ vehicle.year }}</span>
                      }
                      @if (!vehicle.plate && !vehicle.year) {
                        <span>Sem placa/ano</span>
                      }
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-1.5 text-sm font-medium text-zinc-200 tabular-nums">
                  <svg lucideGauge class="w-4 h-4 text-zinc-500"></svg>
                  <span>{{ formatNumber(vehicle.current_mileage) }} km</span>
                </div>
              </div>

              <div class="flex items-center justify-between gap-2 pt-3 border-t border-zinc-800/60">
                <button
                  type="button"
                  (click)="selectVehicle(vehicle)"
                  class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition cursor-pointer"
                  [class.bg-zinc-100]="selectedVehicle()?.id === vehicle.id"
                  [class.text-zinc-950]="selectedVehicle()?.id === vehicle.id"
                  [class.bg-zinc-800/80]="selectedVehicle()?.id !== vehicle.id"
                  [class.hover:bg-zinc-800]="selectedVehicle()?.id !== vehicle.id"
                  [class.text-zinc-200]="selectedVehicle()?.id !== vehicle.id"
                >
                  <svg lucideWrench class="w-3.5 h-3.5"></svg>
                  <span>Manutenções</span>
                </button>

                <div class="flex items-center gap-1.5">
                  <button
                    type="button"
                    (click)="openEditVehicleModal(vehicle)"
                    aria-label="Editar veículo"
                    class="p-1.5 rounded-lg bg-zinc-800/70 hover:bg-zinc-800 text-zinc-300 transition cursor-pointer"
                  >
                    <svg lucidePencil class="w-3.5 h-3.5"></svg>
                  </button>
                  <button
                    type="button"
                    (click)="openDeleteVehicleModal(vehicle)"
                    aria-label="Excluir veículo"
                    class="p-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 transition cursor-pointer"
                  >
                    <svg lucideTrash class="w-3.5 h-3.5"></svg>
                  </button>
                </div>
              </div>
            </div>
          }
        </div>

        <!-- Selected Vehicle Maintenances Panel -->
        @if (selectedVehicle(); as vehicle) {
          <section class="rounded-xl border border-zinc-800 bg-zinc-900/40 overflow-hidden">
            <div class="p-5 border-b border-zinc-800/80 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div>
                <h2 class="text-base font-semibold text-zinc-100">
                  Manutenções de {{ vehicle.name }}
                </h2>
                <p class="text-xs text-zinc-400 mt-0.5">
                  Quilometragem atual: {{ formatNumber(vehicle.current_mileage) }} km
                </p>
              </div>

              <button
                type="button"
                (click)="openAttachModal(vehicle)"
                class="inline-flex items-center justify-center gap-1.5 px-3.5 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-xs font-medium transition cursor-pointer shrink-0"
              >
                <svg lucidePlus class="w-3.5 h-3.5"></svg>
                <span>Adicionar Manutenção</span>
              </button>
            </div>

            @if (loadingMaintenances()) {
              <div class="p-8 text-center text-xs text-zinc-400 flex items-center justify-center gap-2">
                <svg lucideLoaderCircle class="w-4 h-4 animate-spin"></svg>
                <span>Carregando manutenções...</span>
              </div>
            } @else if (vehicleMaintenances().length === 0) {
              <div class="p-10 text-center text-sm text-zinc-400">
                Nenhuma manutenção vinculada a este veículo ainda.
              </div>
            } @else {
              <div class="divide-y divide-zinc-800/70">
                @for (m of vehicleMaintenances(); track m.id) {
                  <div class="p-4 sm:px-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-zinc-900/50 transition">
                    <div class="space-y-1.5 min-w-0">
                      <div class="flex flex-wrap items-center gap-2.5">
                        @switch (m.status) {
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
                        <span class="text-sm font-medium text-zinc-100">{{ m.maintenance_type_name }}</span>
                      </div>

                      <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-zinc-400">
                        <span class="inline-flex items-center gap-1">
                          <svg lucideCalendar class="w-3.5 h-3.5 text-zinc-500"></svg>
                          Última vez: {{ formatDate(m.last_performed_date) }}
                        </span>
                        <span class="inline-flex items-center gap-1">
                          <svg lucideGauge class="w-3.5 h-3.5 text-zinc-500"></svg>
                          Aos {{ formatNumber(m.last_performed_mileage) }} km
                        </span>
                        @if (m.notes) {
                          <span class="text-zinc-500 italic">{{ m.notes }}</span>
                        }
                      </div>
                    </div>

                    <button
                      type="button"
                      (click)="removeAttachedMaintenance(m.id)"
                      [disabled]="removingMaintenanceId() === m.id"
                      class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 text-xs font-medium transition disabled:opacity-50 cursor-pointer shrink-0 self-start sm:self-center"
                    >
                      @if (removingMaintenanceId() === m.id) {
                        <svg lucideLoaderCircle class="w-3.5 h-3.5 animate-spin"></svg>
                      } @else {
                        <svg lucideTrash class="w-3.5 h-3.5"></svg>
                      }
                      <span>Remover</span>
                    </button>
                  </div>
                }
              </div>
            }
          </section>
        }
      }

      <!-- Create / Edit Vehicle Modal (<dialog>) -->
      <dialog
        #vehicleDialog
        closedby="any"
        aria-labelledby="v-dialog-title"
        (close)="onVehicleDialogClose()"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-md w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        <div class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <h1 id="v-dialog-title" class="text-base font-semibold text-zinc-100">
              {{ editingVehicleId() ? 'Editar Veículo' : 'Novo Veículo' }}
            </h1>
            <button
              type="button"
              (click)="closeVehicleModal()"
              aria-label="Fechar modal"
              class="p-1 rounded-md text-zinc-500 hover:text-zinc-200 hover:bg-zinc-800 transition cursor-pointer"
            >
              <svg lucideX class="w-4 h-4"></svg>
            </button>
          </div>

          <form (ngSubmit)="submitVehicleForm()" class="space-y-4">
            <div>
              <label for="v-name" class="block text-xs font-medium text-zinc-400 mb-1.5">
                Nome *
              </label>
              <input
                id="v-name"
                type="text"
                name="name"
                [ngModel]="vName()"
                (ngModelChange)="vName.set($event)"
                required
                autofocus
                placeholder="Ex: Onix 1.0 LS"
                class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
              />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="v-plate" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Placa (opcional)
                </label>
                <input
                  id="v-plate"
                  type="text"
                  name="plate"
                  [ngModel]="vPlate()"
                  (ngModelChange)="vPlate.set($event)"
                  placeholder="Ex: ABC-1234"
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
                />
              </div>

              <div>
                <label for="v-year" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Ano (opcional)
                </label>
                <input
                  id="v-year"
                  type="number"
                  min="1900"
                  max="2100"
                  name="year"
                  [ngModel]="vYear()"
                  (ngModelChange)="vYear.set($event)"
                  placeholder="Ex: 2022"
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
                />
              </div>
            </div>

            <div>
              <label for="v-mileage" class="block text-xs font-medium text-zinc-400 mb-1.5">
                Quilometragem Atual (km) *
              </label>
              <input
                id="v-mileage"
                type="number"
                min="0"
                name="mileage"
                [ngModel]="vMileage()"
                (ngModelChange)="vMileage.set($event)"
                required
                class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
              />
            </div>

            @if (vFormError()) {
              <p role="alert" class="text-xs text-rose-400 bg-rose-500/10 border border-rose-500/20 rounded-lg px-3 py-2">
                {{ vFormError() }}
              </p>
            }

            <div class="flex items-center justify-end gap-2.5 pt-2">
              <button
                type="button"
                (click)="closeVehicleModal()"
                class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
              >
                Cancelar
              </button>
              <button
                type="submit"
                [disabled]="savingVehicle() || !vName().trim() || vMileage() < 0"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-xs font-medium transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                @if (savingVehicle()) {
                  <svg lucideLoaderCircle class="w-3.5 h-3.5 animate-spin"></svg>
                  <span>Salvando...</span>
                } @else {
                  <svg lucideCheck class="w-3.5 h-3.5"></svg>
                  <span>{{ editingVehicleId() ? 'Salvar Alterações' : 'Cadastrar' }}</span>
                }
              </button>
            </div>
          </form>
        </div>
      </dialog>

      <!-- Attach Maintenance Modal (<dialog>) -->
      <dialog
        #attachDialog
        closedby="any"
        aria-labelledby="attach-dialog-title"
        (close)="attachError.set(null)"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-md w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        <div class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <div>
              <h1 id="attach-dialog-title" class="text-base font-semibold text-zinc-100">
                Adicionar Manutenção
              </h1>
              @if (selectedVehicle(); as v) {
                <p class="text-xs text-zinc-400 mt-1">{{ v.name }}</p>
              }
            </div>
            <button
              type="button"
              (click)="closeAttachModal()"
              aria-label="Fechar modal"
              class="p-1 rounded-md text-zinc-500 hover:text-zinc-200 hover:bg-zinc-800 transition cursor-pointer"
            >
              <svg lucideX class="w-4 h-4"></svg>
            </button>
          </div>

          @if (maintenanceTypes().length === 0) {
            <div class="p-4 rounded-lg border border-amber-500/20 bg-amber-500/10 text-amber-200 text-xs space-y-2">
              <div class="flex items-center gap-2 font-medium">
                <svg lucideCircleExclamationPoint class="w-4 h-4 text-amber-400"></svg>
                <span>Catálogo de manutenções vazio</span>
              </div>
              <p class="text-amber-300/80">
                Cadastre primeiro os tipos de manutenção no catálogo para poder vinculá-los a um veículo.
              </p>
              <a
                routerLink="/maintenances"
                (click)="closeAttachModal()"
                class="inline-block font-medium text-amber-300 underline underline-offset-4 hover:text-white pt-1"
              >
                Ir para Catálogo de Manutenções
              </a>
            </div>
          } @else {
            <form (ngSubmit)="submitAttach()" class="space-y-4">
              <div>
                <label for="attach-type" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Tipo de Manutenção *
                </label>
                <select
                  id="attach-type"
                  name="maintenanceTypeId"
                  [ngModel]="attachTypeId()"
                  (ngModelChange)="attachTypeId.set($event)"
                  required
                  autofocus
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                >
                  @for (mt of maintenanceTypes(); track mt.id) {
                    <option [value]="mt.id">{{ mt.name }}</option>
                  }
                </select>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label for="attach-date" class="block text-xs font-medium text-zinc-400 mb-1.5">
                    Última Realização *
                  </label>
                  <input
                    id="attach-date"
                    type="date"
                    name="lastPerformedDate"
                    [ngModel]="attachDate()"
                    (ngModelChange)="attachDate.set($event)"
                    required
                    class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                  />
                </div>

                <div>
                  <label for="attach-km" class="block text-xs font-medium text-zinc-400 mb-1.5">
                    KM na Realização *
                  </label>
                  <input
                    id="attach-km"
                    type="number"
                    min="0"
                    name="lastPerformedMileage"
                    [ngModel]="attachMileage()"
                    (ngModelChange)="attachMileage.set($event)"
                    required
                    class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:border-zinc-500"
                  />
                </div>
              </div>

              <div>
                <label for="attach-notes" class="block text-xs font-medium text-zinc-400 mb-1.5">
                  Observações (opcional)
                </label>
                <textarea
                  id="attach-notes"
                  name="notes"
                  rows="2"
                  [ngModel]="attachNotes()"
                  (ngModelChange)="attachNotes.set($event)"
                  placeholder="Ex: Feito na revisão dos 40 mil km"
                  class="w-full rounded-lg bg-zinc-950 border border-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-zinc-500 resize-none"
                ></textarea>
              </div>

              @if (attachError()) {
                <p role="alert" class="text-xs text-rose-400 bg-rose-500/10 border border-rose-500/20 rounded-lg px-3 py-2">
                  {{ attachError() }}
                </p>
              }

              <div class="flex items-center justify-end gap-2.5 pt-2">
                <button
                  type="button"
                  (click)="closeAttachModal()"
                  class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  [disabled]="savingAttach()"
                  class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 text-xs font-medium transition disabled:opacity-50 cursor-pointer"
                >
                  @if (savingAttach()) {
                    <svg lucideLoaderCircle class="w-3.5 h-3.5 animate-spin"></svg>
                    <span>Vinculando...</span>
                  } @else {
                    <svg lucideCheck class="w-3.5 h-3.5"></svg>
                    <span>Vincular Manutenção</span>
                  }
                </button>
              </div>
            </form>
          }
        </div>
      </dialog>

      <!-- Delete Vehicle Confirmation Modal (<dialog>) -->
      <dialog
        #deleteVehicleDialog
        closedby="any"
        aria-labelledby="v-delete-title"
        (close)="deleteVehicleTarget.set(null)"
        class="m-auto bg-zinc-900 border border-zinc-800 rounded-xl p-6 max-w-sm w-[calc(100%-2rem)] text-zinc-100 shadow-2xl"
      >
        @if (deleteVehicleTarget(); as target) {
          <div class="space-y-4">
            <h1 id="v-delete-title" class="text-base font-semibold text-zinc-100">
              Excluir veículo?
            </h1>
            <p class="text-xs text-zinc-400 leading-relaxed">
              Tem certeza que deseja excluir <strong class="text-zinc-200">{{ target.name }}</strong>?
              Todas as manutenções vinculadas a ele serão removidas.
            </p>
            <div class="flex items-center justify-end gap-2.5 pt-2">
              <button
                type="button"
                autofocus
                (click)="closeDeleteVehicleModal()"
                class="px-3.5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium transition cursor-pointer"
              >
                Cancelar
              </button>
              <button
                type="button"
                [disabled]="deletingVehicle()"
                (click)="confirmDeleteVehicle()"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-rose-600 hover:bg-rose-500 text-white text-xs font-medium transition disabled:opacity-50 cursor-pointer"
              >
                @if (deletingVehicle()) {
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
export class VehiclesComponent implements OnInit {
  private readonly vehicleService = inject(VehicleService);
  private readonly maintenanceTypeService = inject(MaintenanceTypeService);

  protected readonly vehicleDialog = viewChild<ElementRef<HTMLDialogElement>>('vehicleDialog');
  protected readonly attachDialog = viewChild<ElementRef<HTMLDialogElement>>('attachDialog');
  protected readonly deleteVehicleDialog = viewChild<ElementRef<HTMLDialogElement>>('deleteVehicleDialog');

  // Icons

  // Vehicles state
  protected readonly vehicles = signal<Vehicle[]>([]);
  protected readonly loading = signal<boolean>(true);
  protected readonly errorMessage = signal<string | null>(null);

  // Selected vehicle & its maintenances
  protected readonly selectedVehicle = signal<Vehicle | null>(null);
  protected readonly vehicleMaintenances = signal<VehicleMaintenanceDetail[]>([]);
  protected readonly loadingMaintenances = signal<boolean>(false);
  protected readonly removingMaintenanceId = signal<number | null>(null);

  // Vehicle form modal state
  protected readonly editingVehicleId = signal<number | null>(null);
  protected readonly vName = signal<string>('');
  protected readonly vPlate = signal<string>('');
  protected readonly vYear = signal<number | null>(null);
  protected readonly vMileage = signal<number>(0);
  protected readonly savingVehicle = signal<boolean>(false);
  protected readonly vFormError = signal<string | null>(null);

  // Delete vehicle modal state
  protected readonly deleteVehicleTarget = signal<Vehicle | null>(null);
  protected readonly deletingVehicle = signal<boolean>(false);

  // Attach maintenance modal state
  protected readonly maintenanceTypes = signal<MaintenanceType[]>([]);
  protected readonly attachTypeId = signal<number>(0);
  protected readonly attachDate = signal<string>('');
  protected readonly attachMileage = signal<number>(0);
  protected readonly attachNotes = signal<string>('');
  protected readonly savingAttach = signal<boolean>(false);
  protected readonly attachError = signal<string | null>(null);

  ngOnInit(): void {
    this.loadVehicles();
    this.loadMaintenanceTypes();
  }

  loadVehicles(): void {
    this.loading.set(true);
    this.errorMessage.set(null);

    this.vehicleService.list().subscribe({
      next: (data) => {
        this.vehicles.set(data);
        this.loading.set(false);

        const currentSelected = this.selectedVehicle();
        if (currentSelected) {
          const updatedSelected = data.find((v) => v.id === currentSelected.id) ?? null;
          this.selectedVehicle.set(updatedSelected);
        }
      },
      error: () => {
        this.loading.set(false);
        this.errorMessage.set('Não foi possível carregar os veículos.');
      }
    });
  }

  loadMaintenanceTypes(): void {
    this.maintenanceTypeService.list().subscribe({
      next: (types) => {
        this.maintenanceTypes.set(types);
        if (types.length > 0 && !this.attachTypeId()) {
          this.attachTypeId.set(types[0].id);
        }
      }
    });
  }

  selectVehicle(vehicle: Vehicle): void {
    this.selectedVehicle.set(vehicle);
    this.loadVehicleMaintenances(vehicle.id);
  }

  loadVehicleMaintenances(vehicleId: number): void {
    this.loadingMaintenances.set(true);
    this.vehicleService.listMaintenances(vehicleId).subscribe({
      next: (list) => {
        this.vehicleMaintenances.set(list);
        this.loadingMaintenances.set(false);
      },
      error: () => {
        this.loadingMaintenances.set(false);
      }
    });
  }

  // Vehicle CRUD Modals
  openCreateVehicleModal(): void {
    this.editingVehicleId.set(null);
    this.vName.set('');
    this.vPlate.set('');
    this.vYear.set(null);
    this.vMileage.set(0);
    this.vFormError.set(null);

    const dialogEl = this.vehicleDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  openEditVehicleModal(vehicle: Vehicle): void {
    this.editingVehicleId.set(vehicle.id);
    this.vName.set(vehicle.name);
    this.vPlate.set(vehicle.plate ?? '');
    this.vYear.set(vehicle.year ?? null);
    this.vMileage.set(vehicle.current_mileage);
    this.vFormError.set(null);

    const dialogEl = this.vehicleDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeVehicleModal(): void {
    const dialogEl = this.vehicleDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  onVehicleDialogClose(): void {
    this.editingVehicleId.set(null);
    this.vFormError.set(null);
  }

  submitVehicleForm(): void {
    const mileage = Number(this.vMileage());
    if (!this.vName().trim() || isNaN(mileage) || mileage < 0 || this.savingVehicle()) {
      return;
    }

    const plateTrimmed = this.vPlate().trim();
    const yearVal = this.vYear();

    const payload = {
      name: this.vName().trim(),
      plate: plateTrimmed ? plateTrimmed : null,
      year: yearVal ? Number(yearVal) : null,
      current_mileage: mileage
    };

    this.savingVehicle.set(true);
    this.vFormError.set(null);

    const id = this.editingVehicleId();
    const request$ = id
      ? this.vehicleService.update(id, payload)
      : this.vehicleService.create(payload);

    request$.subscribe({
      next: (saved) => {
        this.savingVehicle.set(false);
        this.closeVehicleModal();
        this.loadVehicles();
        if (this.selectedVehicle()?.id === saved.id) {
          this.loadVehicleMaintenances(saved.id);
        }
      },
      error: (err: HttpErrorResponse) => {
        this.savingVehicle.set(false);
        this.vFormError.set(err.error?.error || 'Erro ao salvar veículo.');
      }
    });
  }

  openDeleteVehicleModal(vehicle: Vehicle): void {
    this.deleteVehicleTarget.set(vehicle);
    const dialogEl = this.deleteVehicleDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeDeleteVehicleModal(): void {
    const dialogEl = this.deleteVehicleDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  confirmDeleteVehicle(): void {
    const target = this.deleteVehicleTarget();
    if (!target || this.deletingVehicle()) {
      return;
    }

    this.deletingVehicle.set(true);
    this.vehicleService.delete(target.id).subscribe({
      next: () => {
        this.deletingVehicle.set(false);
        if (this.selectedVehicle()?.id === target.id) {
          this.selectedVehicle.set(null);
          this.vehicleMaintenances.set([]);
        }
        this.closeDeleteVehicleModal();
        this.loadVehicles();
      },
      error: () => {
        this.deletingVehicle.set(false);
      }
    });
  }

  // Attach Maintenance Modal
  openAttachModal(vehicle: Vehicle): void {
    const today = new Date().toISOString().slice(0, 10);
    const types = this.maintenanceTypes();
    this.attachTypeId.set(types.length > 0 ? types[0].id : 0);
    this.attachDate.set(today);
    this.attachMileage.set(vehicle.current_mileage);
    this.attachNotes.set('');
    this.attachError.set(null);

    const dialogEl = this.attachDialog()?.nativeElement;
    if (dialogEl && !dialogEl.open) {
      dialogEl.showModal();
    }
  }

  closeAttachModal(): void {
    const dialogEl = this.attachDialog()?.nativeElement;
    if (dialogEl && dialogEl.open) {
      dialogEl.close();
    }
  }

  submitAttach(): void {
    const vehicle = this.selectedVehicle();
    if (!vehicle || this.savingAttach()) {
      return;
    }

    const typeId = Number(this.attachTypeId());
    const mileage = Number(this.attachMileage());
    if (!typeId || !this.attachDate() || isNaN(mileage) || mileage < 0) {
      this.attachError.set('Preencha todos os campos obrigatórios corretamente.');
      return;
    }

    const notesTrimmed = this.attachNotes().trim();

    this.savingAttach.set(true);
    this.attachError.set(null);

    this.vehicleService
      .attachMaintenance(vehicle.id, {
        maintenance_type_id: typeId,
        last_performed_date: this.attachDate(),
        last_performed_mileage: mileage,
        notes: notesTrimmed ? notesTrimmed : null
      })
      .subscribe({
        next: () => {
          this.savingAttach.set(false);
          this.closeAttachModal();
          this.loadVehicleMaintenances(vehicle.id);
        },
        error: (err: HttpErrorResponse) => {
          this.savingAttach.set(false);
          this.attachError.set(err.error?.error || 'Erro ao vincular manutenção.');
        }
      });
  }

  removeAttachedMaintenance(id: number): void {
    const vehicle = this.selectedVehicle();
    if (!vehicle || this.removingMaintenanceId() !== null) {
      return;
    }

    this.removingMaintenanceId.set(id);
    this.vehicleService.deleteMaintenance(id).subscribe({
      next: () => {
        this.removingMaintenanceId.set(null);
        this.loadVehicleMaintenances(vehicle.id);
      },
      error: () => {
        this.removingMaintenanceId.set(null);
      }
    });
  }

  formatNumber(value: number): string {
    return value.toLocaleString('pt-BR');
  }

  formatDate(isoDate: string): string {
    const [year, month, day] = isoDate.split('-');
    if (!year || !month || !day) {
      return isoDate;
    }
    return `${day}/${month}/${year}`;
  }
}
