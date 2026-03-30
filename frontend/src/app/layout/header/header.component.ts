import { Component, ChangeDetectionStrategy, inject, computed } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { ThemeService, Theme } from '../../core/services/theme.service';

interface NavItem { label: string; path: string; }

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    header {
      position: sticky; top: 0; z-index: 50;
      backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px);
      border-bottom: 1px solid var(--border);
      background: var(--header-bg, rgba(10,10,10,0.8));
      transition: background 200ms ease, border-color 200ms ease;
    }
    nav {
      display: flex; align-items: center; justify-content: space-between;
      height: 64px; max-width: 768px; margin: 0 auto; padding-inline: 1.5rem;
    }
    .logo {
      font-size: 0.9375rem; font-weight: 600; color: var(--text);
      letter-spacing: -0.02em; text-decoration: none;
      transition: color 150ms ease;
      &:hover { color: var(--accent); }
    }
    .logo-dot { color: var(--accent); }
    .right { display: flex; align-items: center; gap: 0.25rem; }
    .social-links { display: flex; align-items: center; gap: 0.125rem; margin-right: 0.25rem; }
    .social-link {
      display: inline-flex; align-items: center; justify-content: center;
      width: 34px; height: 34px; border-radius: 8px; border: 1px solid transparent;
      background: transparent; color: var(--text-secondary); text-decoration: none;
      transition: color 150ms ease, background 150ms ease, border-color 150ms ease;
      &:hover { color: var(--text); background: var(--surface); border-color: var(--border); }
      svg { display: block; }
    }
    .nav-links { display: flex; align-items: center; gap: 0.25rem; list-style: none; }
    .nav-links a {
      font-size: 0.875rem; color: var(--text-secondary); padding: 0.4rem 0.75rem;
      border-radius: 8px; text-decoration: none;
      transition: color 150ms ease, background 150ms ease;
      &:hover { color: var(--text); background: var(--surface); }
    }
    .nav-links a.active {
      color: var(--accent); background: rgba(99,102,241,0.15);
    }
    /* Theme toggle button */
    .theme-btn {
      display: inline-flex; align-items: center; justify-content: center;
      width: 34px; height: 34px; border-radius: 8px; border: 1px solid var(--border);
      background: transparent; cursor: pointer; color: var(--text-secondary);
      transition: color 150ms ease, background 150ms ease, border-color 150ms ease;
      margin-left: 0.5rem; flex-shrink: 0;
      &:hover { color: var(--text); background: var(--surface); border-color: var(--border-hover); }
      svg { display: block; }
    }
    .sr-only { position: absolute; width: 1px; height: 1px; overflow: hidden; clip: rect(0,0,0,0); }
  `],
  template: `
    <header role="banner">
      <nav aria-label="Main navigation">
        <a routerLink="/" class="logo" aria-label="Home">mehedi<span class="logo-dot">.</span>dev</a>
        <div class="right">
          <ul class="nav-links" role="list">
            @for (item of navItems; track item.path) {
              <li>
                <a [routerLink]="item.path" routerLinkActive="active"
                   [routerLinkActiveOptions]="{ exact: item.path === '/' }">
                  {{ item.label }}
                </a>
              </li>
            }
          </ul>
          <!-- Social links -->
          <div class="social-links">
            <a class="social-link" href="#" aria-label="GitHub" target="_blank" rel="noopener noreferrer">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z"/>
              </svg>
            </a>
            <a class="social-link" href="#" aria-label="LinkedIn" target="_blank" rel="noopener noreferrer">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 0 1-2.063-2.065 2.064 2.064 0 1 1 2.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/>
              </svg>
            </a>
            <a class="social-link" href="#" aria-label="Twitter / X" target="_blank" rel="noopener noreferrer">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.744l7.73-8.835L1.254 2.25H8.08l4.259 5.63zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
              </svg>
            </a>
          </div>
          <!-- Theme toggle: cycles dark → light → auto -->
          <button class="theme-btn" (click)="themeService.cycle()"
                  [title]="themeLabel()">
            <span class="sr-only">{{ themeLabel() }}</span>
            @if (themeService.theme() === 'dark') {
              <!-- Moon icon -->
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
              </svg>
            } @else if (themeService.theme() === 'light') {
              <!-- Sun icon -->
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="5"/>
                <line x1="12" y1="1" x2="12" y2="3"/>
                <line x1="12" y1="21" x2="12" y2="23"/>
                <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
                <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
                <line x1="1" y1="12" x2="3" y2="12"/>
                <line x1="21" y1="12" x2="23" y2="12"/>
                <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
                <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
              </svg>
            } @else {
              <!-- Auto/Monitor icon -->
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
                <line x1="8" y1="21" x2="16" y2="21"/>
                <line x1="12" y1="17" x2="12" y2="21"/>
              </svg>
            }
          </button>
        </div>
      </nav>
    </header>
  `,
})
export class HeaderComponent {
  readonly themeService = inject(ThemeService);

  readonly themeLabel = computed(() => {
    const labels: Record<Theme, string> = {
      dark: 'Dark mode — click for light',
      light: 'Light mode — click for auto',
      auto: 'Auto mode — click for dark',
    };
    return labels[this.themeService.theme()];
  });

  navItems: NavItem[] = [
    { label: 'About', path: '/about' },
    { label: 'Articles', path: '/articles' },
    { label: 'Projects', path: '/projects' },
    { label: 'Contact', path: '/contact' },
  ];
}
