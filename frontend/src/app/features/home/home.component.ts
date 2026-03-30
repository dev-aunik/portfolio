import { Component, OnInit, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ApiService, Article } from '../../core/services/api.service';
import { SeoService } from '../../core/services/seo.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink, CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    .hero { padding: 6rem 0 4rem; }
    .eyebrow {
      display: inline-flex; align-items: center; gap: 0.5rem;
      font-size: 0.8125rem; letter-spacing: 0.08em; text-transform: uppercase;
      color: var(--accent); font-weight: 500; margin-bottom: 1.5rem;
    }
    .eyebrow::before { content: ''; width: 20px; height: 1px; background: var(--accent); }
    h1 {
      font-size: clamp(2.5rem, 6vw, 4rem); font-weight: 700; line-height: 1.1;
      letter-spacing: -0.03em; margin-bottom: 1.25rem;
    }
    .gradient {
      background: linear-gradient(135deg, #6366f1 0%, #818cf8 50%, #a5b4fc 100%);
      -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
    }
    .bio { font-size: 1.125rem; color: var(--text-secondary); max-width: 540px; line-height: 1.8; margin-bottom: 2rem; }
    .bio strong { color: var(--text); font-weight: 500; }
    .actions { display: flex; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 3.5rem; }
    .divider { border: none; border-top: 1px solid var(--border); margin: 3.5rem 0 2.5rem; }
    .section-label { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--muted); font-weight: 500; margin-bottom: 1.25rem; }
    .article-item {
      display: block; padding: 1.25rem 0; border-bottom: 1px solid var(--border); text-decoration: none;
      transition: padding-left 300ms ease;
      &:first-child { border-top: 1px solid var(--border); }
      &:hover { padding-left: 0.5rem; }
      &:hover .title { color: var(--accent); }
    }
    .meta { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.4rem; }
    .date { font-size: 0.8125rem; color: var(--muted); }
    .title { font-size: 1rem; font-weight: 500; color: var(--text); margin-bottom: 0.35rem; transition: color 150ms ease; }
    .summary { font-size: 0.875rem; color: var(--text-secondary); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
    .view-all {
      display: inline-flex; align-items: center; gap: 0.4rem; font-size: 0.875rem; color: var(--accent);
      margin-top: 1.5rem; text-decoration: none; transition: gap 150ms ease;
      &:hover { gap: 0.65rem; }
    }
    .tech-pills { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-top: 2.5rem; }
    .pill {
      font-size: 0.75rem; font-weight: 500; color: var(--muted); background: var(--surface);
      border: 1px solid var(--border); border-radius: 8px; padding: 0.3rem 0.7rem;
      transition: color 150ms ease, border-color 150ms ease;
      &:hover { color: var(--text); border-color: var(--border-hover); }
    }
    .skel { background: linear-gradient(90deg, var(--surface) 25%, var(--surface-2) 50%, var(--surface) 75%); background-size: 200%; animation: sh 1.5s infinite; border-radius: 4px; }
    @keyframes sh { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
  `],
  template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <section class="hero fade-in-up">
          <p class="eyebrow">Software Engineer — Backend</p>
          <h1>Hi, I'm <span class="gradient">Mehedi</span>.</h1>
          <p class="bio">
            I design and build <strong>scalable backend systems</strong> that power real-world products.
            Specializing in <strong>Go</strong> and <strong>PHP/Laravel</strong>, with deep experience in
            eCommerce platforms, auction systems, CRM solutions, and SaaS applications —
            turning complex problems into clean, maintainable architecture.
          </p>
          <div class="actions">
            <a routerLink="/about" class="btn btn--primary">About me</a>
            <a routerLink="/contact" class="btn btn--ghost">Get in touch</a>
          </div>
          <div class="tech-pills">
            @for (t of techs; track t) { <span class="pill">{{ t }}</span> }
          </div>
        </section>

        <hr class="divider">

        <section>
          <p class="section-label">Recent Articles</p>
          @if (loading()) {
            @for (i of [1,2,3]; track i) {
              <div style="padding:1.25rem 0; border-bottom:1px solid var(--border)">
                <div class="skel" style="height:12px;width:35%;margin-bottom:0.5rem"></div>
                <div class="skel" style="height:17px;width:65%;margin-bottom:0.4rem"></div>
                <div class="skel" style="height:13px;width:50%"></div>
              </div>
            }
          } @else if (articles().length > 0) {
            @for (a of articles(); track a.slug) {
              <a [routerLink]="['/articles', a.slug]" class="article-item">
                <div class="meta">
                  <span class="date">{{ fmt(a.published_at) }}</span>
                  @for (tag of a.tags.slice(0,2); track tag) { <span class="tag">{{ tag }}</span> }
                </div>
                <p class="title">{{ a.title }}</p>
                <p class="summary">{{ a.summary }}</p>
              </a>
            }
            <a routerLink="/articles" class="view-all">
              All articles
              <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
            </a>
          } @else {
            <p style="color:var(--muted); font-size:0.9rem">Articles coming soon.</p>
          }
        </section>
      </div>
    </main>
  `,
})
export class HomeComponent implements OnInit {
  private readonly api = inject(ApiService);
  private readonly seo = inject(SeoService);
  readonly articles = signal<Article[]>([]);
  readonly loading = signal(true);
  techs = ['Go', 'PHP 8', 'Laravel', 'MySQL', 'PostgreSQL', 'Redis', 'RabbitMQ', 'Docker', 'GitHub Actions', 'GCP', 'WebSockets', 'Nginx'];

  ngOnInit() {
    this.seo.set({ title: 'Mehedi Hasan — Software Engineer' });
    this.api.getArticles(1, 5).subscribe({ next: r => { this.articles.set(r.articles ?? []); this.loading.set(false); }, error: () => this.loading.set(false) });
  }
  fmt(iso: string | null) { return iso ? new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }) : ''; }
}
