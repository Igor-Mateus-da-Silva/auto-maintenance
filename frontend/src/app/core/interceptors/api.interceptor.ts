import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const apiInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);

  const isAuthBypassEndpoint =
    req.url.endsWith('/login') ||
    req.url.endsWith('/logout') ||
    req.url.endsWith('/health') ||
    req.url.endsWith('/auth/verify');

  const request = req.clone({
    withCredentials: true,
    setHeaders: {
      Accept: 'application/json'
    }
  });

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401 && !isAuthBypassEndpoint) {
        authService.clearSession();
      }
      return throwError(() => error);
    })
  );
};
