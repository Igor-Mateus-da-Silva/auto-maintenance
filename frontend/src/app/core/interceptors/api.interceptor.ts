import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const apiInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);

  const isPublicEndpoint = req.url.endsWith('/login') || req.url.endsWith('/health');

  let request = req.clone({
    setHeaders: {
      Accept: 'application/json',
    },
  });

  if (!isPublicEndpoint) {
    const token = authService.getToken();
    if (token) {
      request = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`,
        },
      });
    }
  }

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401 && !isPublicEndpoint) {
        authService.logout();
      }
      return throwError(() => error);
    }),
  );
};
