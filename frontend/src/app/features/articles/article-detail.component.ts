import { Component, OnInit, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { switchMap } from 'rxjs';
import { ApiService, Article } from '../../core/services/api.service';
import { SeoService } from '../../core/services/seo.service';

@Component({
    selector: 'app-article-detail',
    standalone: true,
    imports: [CommonModule, RouterLink],
    changeDetection: ChangeDetectionStrategy.OnPush,
    styles: [`
    .page { padding: 4rem 0 5rem; }
    .back { display: inline-flex; align-items: center; gap: 0.4rem; font-size: 0.875rem; color: var(--muted); text-decoration: none; margin-bottom: 2.5rem; transition: color 150ms; &:hover { color: var(--accent); } }
    .tags { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1rem; }
    h1 { margin-bottom: 1rem; }
    .meta { display: flex; align-items: center; gap: 1rem; flex-wrap: wrap; font-size: 0.8125rem; color: var(--muted); }
    .sep { width: 3px; height: 3px; border-radius: 50%; background: var(--muted); }
    hr { border: none; border-top: 1px solid var(--border); margin: 2rem 0; }
    .nf { text-align: center; padding: 5rem 0; }
    .skel { background: linear-gradient(90deg, var(--surface) 25%, var(--surface-2) 50%, var(--surface) 75%); background-size: 200%; animation: sh 1.5s infinite; border-radius: 4px; }
    @keyframes sh { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
  `],
    template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        @if (loading()) {
          <div class="page">
            <div class="skel" style="height:14px;width:80px;margin-bottom:2.5rem"></div>
            <div class="skel" style="height:40px;width:80%;margin-bottom:1rem"></div>
            <div class="skel" style="height:12px;width:40%;margin-bottom:2rem"></div>
            @for (i of [1,2,3,4,5]; track i) { <div class="skel" style="height:14px;margin-bottom:0.75rem"></div> }
          </div>
        } @else if (!article()) {
          <div class="nf">
            <h1>Article not found</h1>
            <p style="color:var(--text-secondary);margin-bottom:1.5rem">This article may have been moved or deleted.</p>
            <a routerLink="/articles" class="btn btn--ghost">Back to articles</a>
          </div>
        } @else {
          <article class="page fade-in-up">
            <a routerLink="/articles" class="back">
              <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>
              All articles
            </a>
            @if (article()!.tags.length) {
              <div class="tags">@for (t of article()!.tags; track t) { <span class="tag">{{ t }}</span> }</div>
            }
            <h1>{{ article()!.title }}</h1>
            <div class="meta">
              <time [attr.datetime]="article()!.published_at">{{ fmt(article()!.published_at) }}</time>
              <span class="sep"></span>
              <span>{{ mins() }} min read</span>
            </div>
            <hr>
            <div class="prose" [innerHTML]="html()"></div>
          </article>
        }
      </div>
    </main>
  `,
})
export class ArticleDetailComponent implements OnInit {
    private readonly api = inject(ApiService);
    private readonly route = inject(ActivatedRoute);
    private readonly seo = inject(SeoService);
    private readonly sanitizer = inject(DomSanitizer);
    readonly article = signal<Article | null>(null);
    readonly loading = signal(true);
    readonly html = signal<SafeHtml>('');

    ngOnInit() {
        this.route.params.pipe(switchMap(p => this.api.getArticle(p['slug']))).subscribe({
            next: a => {
                this.article.set(a);
                this.html.set(this.sanitizer.bypassSecurityTrustHtml(a.content));
                this.seo.set({ title: a.title, description: a.summary, type: 'article', keywords: a.tags.join(', ') });
                this.loading.set(false);
            },
            error: () => { this.article.set(null); this.loading.set(false); },
        });
    }

    mins() { return Math.max(1, Math.round(((this.article()?.content ?? '').split(/\s+/).length) / 200)); }
    fmt(iso: string | null) { return iso ? new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' }) : ''; }
}
