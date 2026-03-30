import { Component, ChangeDetectionStrategy } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HeaderComponent } from './header/header.component';
import { FooterComponent } from './footer/footer.component';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [RouterOutlet, HeaderComponent, FooterComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <a href="#main-content" class="skip-link">Skip to main content</a>
    <app-header />
    <router-outlet />
    <app-footer />
  `,
  styles: [`
    .skip-link {
      position: absolute;
      top: -40px;
      left: 1rem;
      background: var(--accent);
      color: #fff;
      padding: 0.5rem 1rem;
      border-radius: 0 0 var(--radius) var(--radius);
      font-size: 0.875rem;
      font-weight: 500;
      z-index: 9999;
      text-decoration: none;
      &:focus { top: 0; }
    }
  `],
})
export class LayoutComponent { }
