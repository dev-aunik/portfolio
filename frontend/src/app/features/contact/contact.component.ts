import { Component, OnInit, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { SeoService } from '../../core/services/seo.service';

type State = 'idle' | 'submitting' | 'success' | 'error';

@Component({
  selector: 'app-contact',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    .page { padding: 4rem 0 5rem; }
    .eyebrow { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    h1 { margin-bottom: 0.75rem; }
    .sub { color: var(--text-secondary); margin-bottom: 2.5rem; font-size: 1.0625rem; }
    form { display: flex; flex-direction: column; gap: 1.25rem; max-width: 520px; }
    .field { display: flex; flex-direction: column; gap: 0.4rem; }
    label { font-size: 0.8125rem; font-weight: 500; color: var(--text-secondary); }
    input, textarea {
      background: var(--surface); border: 1px solid var(--border); border-radius: 8px;
      color: var(--text); font-size: 0.9375rem; padding: 0.7rem 0.875rem; outline: none;
      transition: border-color 150ms, box-shadow 150ms; font-family: var(--font-sans);
      &::placeholder { color: var(--muted); }
      &:focus { border-color: var(--accent); box-shadow: 0 0 0 3px rgba(99,102,241,0.15); }
      &.invalid { border-color: var(--error); }
    }
    textarea { min-height: 140px; resize: vertical; }
    .err { font-size: 0.8125rem; color: var(--error); }
    .submit {
      align-self: flex-start; background: var(--accent); color: #fff; border: 1px solid var(--accent); border-radius: 8px;
      padding: 0.7rem 1.5rem; font-size: 0.9375rem; font-weight: 500; cursor: pointer;
      display: inline-flex; align-items: center; gap: 0.5rem;
      transition: all 150ms ease; font-family: var(--font-sans);
      &:hover:not(:disabled) { background: var(--accent-hover); border-color: var(--accent-hover); box-shadow: 0 0 20px rgba(99,102,241,0.25); }
      &:disabled { opacity: 0.6; cursor: not-allowed; }
    }
    .spin { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: sp 0.6s linear infinite; }
    @keyframes sp { to { transform: rotate(360deg); } }
    .toast { padding: 1rem 1.25rem; border-radius: 8px; font-size: 0.9rem; margin-bottom: 1.5rem; display: flex; align-items: flex-start; gap: 0.75rem; max-width: 520px; }
    .ok  { background: rgba(34,197,94,0.1); border: 1px solid rgba(34,197,94,0.3); color: #22c55e; }
    .bad { background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.3); color: #ef4444; }
  `],
  template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <section class="page fade-in-up">
          <p class="eyebrow">Contact</p>
          <h1>Get in touch</h1>
          <p class="sub">Based in Dhaka, Bangladesh — open to remote projects and collaborations. Feel free to reach out!</p>
          @if (state() === 'success') {
            <div class="toast ok" role="alert">
              <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg>
              Message sent! I'll get back to you soon.
            </div>
          }
          @if (state() === 'error') {
            <div class="toast bad" role="alert">
              <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
              Something went wrong. Please try again.
            </div>
          }
          @if (state() !== 'success') {
            <form [formGroup]="form" (ngSubmit)="submit()" novalidate>
              <div class="field">
                <label for="name">Name *</label>
                <input id="name" type="text" formControlName="name" placeholder="Your name" autocomplete="name"
                       [class.invalid]="f['name'].invalid && f['name'].touched">
                @if (f['name'].invalid && f['name'].touched) { <span class="err" role="alert">Name is required.</span> }
              </div>
              <div class="field">
                <label for="email">Email *</label>
                <input id="email" type="email" formControlName="email" placeholder="you@example.com" autocomplete="email"
                       [class.invalid]="f['email'].invalid && f['email'].touched">
                @if (f['email'].errors?.['required'] && f['email'].touched) { <span class="err" role="alert">Email is required.</span> }
                @else if (f['email'].errors?.['email'] && f['email'].touched) { <span class="err" role="alert">Enter a valid email.</span> }
              </div>
              <div class="field">
                <label for="subject">Subject</label>
                <input id="subject" type="text" formControlName="subject" placeholder="What's this about?">
              </div>
              <div class="field">
                <label for="message">Message *</label>
                <textarea id="message" formControlName="message" placeholder="Tell me about your project or idea…"
                          [class.invalid]="f['message'].invalid && f['message'].touched"></textarea>
                @if (f['message'].invalid && f['message'].touched) { <span class="err" role="alert">Message must be at least 10 characters.</span> }
              </div>
              <button class="submit" type="submit" [disabled]="state() === 'submitting'">
                @if (state() === 'submitting') { <span class="spin"></span> Sending… }
                @else { Send message
                  <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z"/></svg>
                }
              </button>
            </form>
          }
        </section>
      </div>
    </main>
  `,
})
export class ContactComponent implements OnInit {
  private readonly api = inject(ApiService);
  private readonly seo = inject(SeoService);
  private readonly fb = inject(FormBuilder);
  readonly state = signal<State>('idle');
  form = this.fb.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    email: ['', [Validators.required, Validators.email]],
    subject: [''],
    message: ['', [Validators.required, Validators.minLength(10)]],
  });
  get f() { return this.form.controls; }
  ngOnInit() { this.seo.set({ title: 'Contact — Mehedi Hasan', description: 'Get in touch with Mehedi Hasan. Based in Dhaka, Bangladesh — open to remote projects and collaborations.' }); }
  submit() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    this.state.set('submitting');
    const { name, email, subject, message } = this.form.value;
    this.api.submitContact({ name: name!, email: email!, subject: subject ?? '', message: message! }).subscribe({
      next: () => { this.state.set('success'); this.form.reset(); },
      error: () => this.state.set('error'),
    });
  }
}
