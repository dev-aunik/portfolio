import { Injectable, signal, effect, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';

export type Theme = 'dark' | 'light' | 'auto';

@Injectable({ providedIn: 'root' })
export class ThemeService {
    private readonly platformId = inject(PLATFORM_ID);
    private readonly storageKey = 'portfolio-theme';

    readonly theme = signal<Theme>(this.#loadTheme());

    constructor() {
        effect(() => {
            const t = this.theme();
            if (!isPlatformBrowser(this.platformId)) return;
            document.documentElement.setAttribute('data-theme', t);
            localStorage.setItem(this.storageKey, t);
        });
    }

    setTheme(t: Theme) { this.theme.set(t); }

    cycle() {
        const order: Theme[] = ['auto', 'dark', 'light'];
        const next = order[(order.indexOf(this.theme()) + 1) % order.length];
        this.setTheme(next);
    }

    #loadTheme(): Theme {
        if (!isPlatformBrowser(inject(PLATFORM_ID))) return 'light';
        const saved = localStorage.getItem(this.storageKey) as Theme | null;
        return saved ?? 'light';
    }
}
