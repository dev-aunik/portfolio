import { Component, OnInit, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SeoService } from '../../core/services/seo.service';

@Component({
  selector: 'app-talks',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    .page { padding: 4rem 0 5rem; }
    .eyebrow { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    h1 { margin-bottom: 0.75rem; }
    .sub { color: var(--text-secondary); margin-bottom: 2.5rem; }
    .empty-state {
      padding: 3rem 2rem; border: 1px dashed var(--border); border-radius: 12px;
      text-align: center; max-width: 480px;
    }
    .empty-icon { font-size: 2.5rem; margin-bottom: 1rem; }
    .empty-title { font-size: 1rem; font-weight: 500; color: var(--text); margin-bottom: 0.5rem; }
    .empty-desc { font-size: 0.9rem; color: var(--muted); line-height: 1.6; }
  `],
  template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <section class="page fade-in-up">
          <p class="eyebrow">Speaking</p>
          <h1>Talks</h1>
          <p class="sub">Conference talks and community presentations.</p>
          <div class="empty-state">
            <div class="empty-icon">🎤</div>
            <p class="empty-title">No talks yet.</p>
            <p class="empty-desc">I haven't spoken at conferences yet — but I'm actively sharing insights through articles and open source work. Stay tuned!</p>
          </div>
        </section>
      </div>
    </main>
  `,
})
export class TalksComponent implements OnInit {
  private readonly seo = inject(SeoService);
  ngOnInit() { this.seo.set({ title: 'Talks — Mehedi Hasan', description: 'Conference talks and community presentations by Mehedi Hasan.' }); }
}
