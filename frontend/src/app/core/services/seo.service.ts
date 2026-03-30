import { Injectable, inject } from '@angular/core';
import { Title, Meta } from '@angular/platform-browser';

@Injectable({ providedIn: 'root' })
export class SeoService {
    private readonly title = inject(Title);
    private readonly meta = inject(Meta);

    set(cfg: { title: string; description?: string; keywords?: string; type?: string }): void {
        const full = cfg.title.includes('Aunik') ? cfg.title : `${cfg.title} | Aunik`;
        const desc = cfg.description ?? 'Senior backend engineer — Go, distributed systems, Kafka, Redis.';
        this.title.setTitle(full);
        this.meta.updateTag({ name: 'description', content: desc });
        this.meta.updateTag({ property: 'og:title', content: full });
        this.meta.updateTag({ property: 'og:description', content: desc });
        this.meta.updateTag({ property: 'og:type', content: cfg.type ?? 'website' });
        if (cfg.keywords) this.meta.updateTag({ name: 'keywords', content: cfg.keywords });
    }
}
