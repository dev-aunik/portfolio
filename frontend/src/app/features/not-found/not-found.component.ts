import { Component, ChangeDetectionStrategy } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
    selector: 'app-not-found',
    standalone: true,
    imports: [RouterLink],
    changeDetection: ChangeDetectionStrategy.OnPush,
    styles: [`
    .nf {
      display: flex; flex-direction: column; align-items: center; justify-content: center;
      text-align: center; padding: 8rem 1.5rem; min-height: 60vh;
    }
    .c {
      font-size: clamp(5rem, 15vw, 8rem); font-weight: 700; letter-spacing: -0.05em; line-height: 1;
      background: linear-gradient(135deg, #6366f1 0%, #818cf8 50%, rgba(99,102,241,0.3) 100%);
      -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
      margin-bottom: 1rem;
    }
    h1 { margin-bottom: 0.75rem; }
    p { color: var(--text-secondary); margin-bottom: 2rem; }
  `],
    template: `
    <main id="main-content" tabindex="-1">
      <div class="nf fade-in-up">
        <p class="c" aria-hidden="true">404</p>
        <h1>Page not found</h1>
        <p>The page you're looking for doesn't exist or has moved.</p>
        <a routerLink="/" class="btn btn--primary">Back home</a>
      </div>
    </main>
  `,
})
export class NotFoundComponent { }
