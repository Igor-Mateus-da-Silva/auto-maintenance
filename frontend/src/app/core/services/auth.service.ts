import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';

export interface LoginResponse {
  token: string;
}

export interface VerifyResponse {
  authenticated: boolean;
}

const TOKEN_KEY = 'car_maintenance_token';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly apiUrl = 'http://localhost:8080/api';

  readonly isAuthenticated = signal<boolean>(this.hasToken());

  private hasToken(): boolean {
    if (typeof localStorage === 'undefined') {
      return false;
    }
    return !!localStorage.getItem(TOKEN_KEY);
  }

  getToken(): string | null {
    if (typeof localStorage === 'undefined') {
      return null;
    }
    return localStorage.getItem(TOKEN_KEY);
  }

  login(accessKey: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, { access_key: accessKey }).pipe(
      tap((res) => {
        if (typeof localStorage !== 'undefined') {
          localStorage.setItem(TOKEN_KEY, res.token);
        }
        this.isAuthenticated.set(true);
      })
    );
  }

  logout(): void {
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(TOKEN_KEY);
    }
    this.isAuthenticated.set(false);
    void this.router.navigate(['/login']);
  }

  verifySession(): Observable<VerifyResponse> {
    return this.http.get<VerifyResponse>(`${this.apiUrl}/auth/verify`);
  }
}
