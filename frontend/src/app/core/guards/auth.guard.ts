import { CanActivateFn } from '@angular/router';

export const authGuard: CanActivateFn = (route, state) => {
  // Baseline guard placeholder - will validate auth session in auth task
  return true;
};
