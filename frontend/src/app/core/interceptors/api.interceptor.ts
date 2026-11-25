import { HttpInterceptorFn } from '@angular/common/http';

export const apiInterceptor: HttpInterceptorFn = (req, next) => {
  // Configured for future auth token attachment and base request handling
  const cloned = req.clone({
    setHeaders: {
      'Accept': 'application/json',
    }
  });

  return next(cloned);
};
