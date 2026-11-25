import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, catchError, map, of, tap } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface VerifyResponse {
  authenticated: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly apiUrl = environment.apiUrl;

  readonly isAuthenticated = signal<boolean>(false);
  readonly sessionChecked = signal<boolean>(false);

  login(accessKey: string): Observable<VerifyResponse> {
    return this.http.post<VerifyResponse>(`${this.apiUrl}/login`, { access_key: accessKey }).pipe(
      tap((res) => {
        this.isAuthenticated.set(res.authenticated);
        this.sessionChecked.set(true);
      })
    );
  }

  logout(): void {
    this.http
      .post<VerifyResponse>(`${this.apiUrl}/logout`, {})
      .pipe(catchError(() => of(null)))
      .subscribe(() => {
        this.clearSession();
      });
  }

  clearSession(): void {
    this.isAuthenticated.set(false);
    this.sessionChecked.set(true);
    void this.router.navigate(['/login']);
  }

  verifySession(): Observable<VerifyResponse> {
    return this.http.get<VerifyResponse>(`${this.apiUrl}/auth/verify`);
  }

  checkSession(): Observable<boolean> {
    return this.verifySession().pipe(
      map((res) => {
        const valid = !!res.authenticated;
        this.isAuthenticated.set(valid);
        this.sessionChecked.set(true);
        return valid;
      }),
      catchError(() => {
        this.isAuthenticated.set(false);
        this.sessionChecked.set(true);
        return of(false);
      })
    );
  }
}
