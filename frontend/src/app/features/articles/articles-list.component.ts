import { Component, OnInit, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject, debounceTime, switchMap, of } from 'rxjs';
import { ApiService, Article } from '../../core/services/api.service';
import { SeoService } from '../../core/services/seo.service';

@Component({
    selector: 'app-articles-list',
    standalone: true,
    imports: [RouterLink, CommonModule, FormsModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    styles: [`
    .page { padding: 4rem 0 5rem; }
    .eyebrow { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    h1 { margin-bottom: 2rem; }
    .search { position: relative; margin-bottom: 2.5rem; }
    .search svg { position: absolute; left: 0.875rem; top: 50%; transform: translateY(-50%); color: var(--muted); pointer-events: none; }
    .search input {
      width: 100%; background: var(--surface); border: 1px solid var(--border); border-radius: 8px;
      color: var(--text); font-size: 0.9375rem; padding: 0.75rem 1rem 0.75rem 2.75rem; outline: none;
      transition: border-color 150ms, box-shadow 150ms; font-family: var(--font-sans);
      &::placeholder { color: var(--muted); }
      &:focus { border-color: var(--accent); box-shadow: 0 0 0 3px rgba(99,102,241,0.15); }
    }
    .item {
      display: block; padding: 1.5rem 0; border-bottom: 1px solid var(--border); text-decoration: none;
      transition: padding-left 300ms ease;
      &:first-child { border-top: 1px solid var(--border); }
      &:hover { padding-left: 0.5rem; }
      &:hover .t { color: var(--accent); }
    }
    .meta { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.4rem; flex-wrap: wrap; }
    .dt { font-size: 0.8125rem; color: var(--muted); }
    .t { font-size: 1.0625rem; font-weight: 500; color: var(--text); margin-bottom: 0.35rem; transition: color 150ms; }
    .s { font-size: 0.875rem; color: var(--text-secondary); }
    .empty { color: var(--muted); font-size: 0.9rem; padding: 2rem 0; }
    .pager { display: flex; gap: 0.5rem; margin-top: 2.5rem; align-items: center; }
    .pb {
      padding: 0.4rem 0.75rem; font-size: 0.875rem; border-radius: 8px;
      background: var(--surface); border: 1px solid var(--border); color: var(--text-secondary); cursor: pointer;
      transition: all 150ms ease;
      &:hover:not(:disabled) { color: var(--text); border-color: var(--border-hover); }
      &.active { background: rgba(99,102,241,0.15); border-color: var(--accent); color: var(--accent); }
      &:disabled { opacity: 0.4; cursor: not-allowed; }
    }
    .skel { background: linear-gradient(90deg, var(--surface) 25%, var(--surface-2) 50%, var(--surface) 75%); background-size: 200%; animation: sh 1.5s infinite; border-radius: 4px; }
    @keyframes sh { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
  `],
    template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <section class="page fade-in-up">
          <p class="eyebrow">Writing</p>
          <h1>Articles</h1>
          <div class="search">
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" aria-hidden="true"><circle cx="11" cy="11" r="8"/><path stroke-linecap="round" d="m21 21-4.35-4.35"/></svg>
            <input type="search" placeholder="Search articles…" [(ngModel)]="q" (ngModelChange)="onQ($event)" aria-label="Search articles" id="article-search"/>
          </div>
          @if (loading()) {
            @for (i of [1,2,3,4]; track i) {
              <div style="padding:1.5rem 0; border-bottom:1px solid var(--border)">
                <div class="skel" style="height:12px;width:30%;margin-bottom:0.6rem"></div>
                <div class="skel" style="height:18px;width:65%;margin-bottom:0.4rem"></div>
                <div class="skel" style="height:13px;width:50%"></div>
              </div>
            }
          } @else if (articles().length === 0) {
            <p class="empty">{{ q ? 'No results for "' + q + '".' : 'No articles yet.' }}</p>
          } @else {
            @for (a of articles(); track a.slug) {
              <a [routerLink]="['/articles', a.slug]" class="item">
                <div class="meta">
                  <span class="dt">{{ fmt(a.published_at) }}</span>
                  @for (tag of a.tags.slice(0,3); track tag) { <span class="tag">{{ tag }}</span> }
                </div>
                <p class="t">{{ a.title }}</p>
                <p class="s">{{ a.summary }}</p>
              </a>
            }
            @if (!q && totalPages() > 1) {
              <div class="pager">
                <button class="pb" [disabled]="page() <= 1" (click)="go(page()-1)">←</button>
                @for (p of pages(); track p) {
                  <button class="pb" [class.active]="p === page()" (click)="go(p)">{{ p }}</button>
                }
                <button class="pb" [disabled]="page() >= totalPages()" (click)="go(page()+1)">→</button>
              </div>
            }
          }
        </section>
      </div>
    </main>
  `,
})
export class ArticlesListComponent implements OnInit {
    private readonly api = inject(ApiService);
    private readonly seo = inject(SeoService);
    private readonly s$ = new Subject<string>();
    readonly articles = signal<Article[]>([]);
    readonly loading = signal(true);
    readonly page = signal(1);
    readonly totalPages = signal(1);
    readonly pages = signal<number[]>([]);
    q = '';

    ngOnInit() {
        this.seo.set({ title: 'Articles', description: 'Technical writing on Go, distributed systems, Kafka, Redis, and backend architecture.' });
        this.load(1);
        this.s$.pipe(debounceTime(350), switchMap(q => q.length > 1 ? this.api.searchArticles(q) : of(null))).subscribe(res => {
            if (res) { this.articles.set(res.hits ?? []); } else if (!this.q) { this.load(1); }
            this.loading.set(false);
        });
    }

    load(pg: number) {
        this.loading.set(true);
        this.api.getArticles(pg, 10).subscribe({
            next: r => {
                this.articles.set(r.articles ?? []);
                this.totalPages.set(r.total_pages ?? 1);
                this.page.set(r.page ?? 1);
                this.pages.set(Array.from({ length: r.total_pages }, (_, i) => i + 1));
                this.loading.set(false);
            },
            error: () => this.loading.set(false),
        });
    }

    onQ(val: string) {
        this.q = val;
        if (!val) { this.load(1); return; }
        this.loading.set(true);
        this.s$.next(val);
    }

    go(pg: number) {
        if (pg < 1 || pg > this.totalPages()) return;
        this.load(pg);
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    fmt(iso: string | null) { return iso ? new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }) : ''; }
}
