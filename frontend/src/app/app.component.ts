import { Component, OnInit, ChangeDetectionStrategy, signal, computed, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LayoutComponent } from './layout/layout.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, LayoutComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `<app-layout />`,
})
export class AppComponent { }
