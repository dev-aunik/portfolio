import { Component, OnInit, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SeoService } from '../../core/services/seo.service';

@Component({
  selector: 'app-about',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    .page { padding: 4rem 0 5rem; }
    .eyebrow { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    h1 { margin-bottom: 2.5rem; }
    .bio { color: var(--text-secondary); font-size: 1.0625rem; line-height: 1.85; margin-bottom: 1.25rem; }
    .bio strong { color: var(--text); font-weight: 500; }
    .divider { border: none; border-top: 1px solid var(--border); margin: 3rem 0; }
    .skills-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 1.25rem; }
    .skill-card {
      background: var(--surface); border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem;
      transition: border-color 150ms ease;
      &:hover { border-color: var(--border-hover); }
    }
    .sk-cat { font-size: 0.75rem; letter-spacing: 0.06em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    .sk-items { display: flex; flex-wrap: wrap; gap: 0.4rem; }
    .sk-item { font-size: 0.8rem; color: var(--text-secondary); background: var(--surface-2); border: 1px solid var(--border); border-radius: 4px; padding: 0.2rem 0.55rem; }
    .tl-item { display: flex; gap: 1.25rem; padding: 1.25rem 0; border-bottom: 1px solid var(--border); &:last-child { border: none; } }
    .tl-year { font-size: 0.8125rem; color: var(--muted); min-width: 90px; padding-top: 0.15rem; white-space: nowrap; }
    .tl-role { font-size: 0.9375rem; font-weight: 500; color: var(--text); margin-bottom: 0.2rem; }
    .tl-company { font-size: 0.8125rem; color: var(--accent); margin-bottom: 0.35rem; }
    .tl-desc { font-size: 0.875rem; color: var(--text-secondary); }
    .edu-item { padding: 1rem 0; border-bottom: 1px solid var(--border); &:last-child { border: none; } }
    .edu-degree { font-size: 0.9375rem; font-weight: 500; color: var(--text); margin-bottom: 0.2rem; }
    .edu-inst { font-size: 0.875rem; color: var(--text-secondary); }
    .edu-date { font-size: 0.8125rem; color: var(--muted); margin-top: 0.2rem; }
    .cert-list { display: flex; flex-direction: column; gap: 0.75rem; }
    .cert-item { font-size: 0.875rem; color: var(--text-secondary); display: flex; align-items: center; gap: 0.75rem; }
    .cert-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--accent); flex-shrink: 0; }
  `],
  template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <article class="page fade-in-up">
          <p class="eyebrow">About</p>
          <h1>Building things that work.</h1>

          <p class="bio">I'm a <strong>Software Engineer</strong> with over 5 years of professional experience specializing in <strong>backend engineering</strong>, scalable architecture, and high-performance API development. I'm based in <strong>Dhaka, Bangladesh</strong>.</p>
          <p class="bio">My work spans large-scale <strong>eCommerce platforms</strong>, real-time <strong>auction systems</strong>, <strong>CRM solutions</strong>, and <strong>SaaS applications</strong>. I build primarily with <strong>Go</strong> and <strong>PHP/Laravel</strong>, backed by MySQL/PostgreSQL, Redis, RabbitMQ, and Docker — always following <strong>Clean Architecture</strong> principles and SOLID design.</p>
          <p class="bio">Currently at <strong>Daisen Technologies Ltd</strong>, where I architect and build Go-based APIs, CI/CD pipelines on GitHub Actions, and manage production deployments on <strong>Google Cloud Platform</strong>.</p>

          <hr class="divider">
          <h2 style="font-size:1rem;font-weight:600;margin-bottom:1.5rem">Technical Stack</h2>
          <div class="skills-grid">
            @for (s of skills; track s.category) {
              <div class="skill-card">
                <p class="sk-cat">{{ s.category }}</p>
                <div class="sk-items">
                  @for (i of s.items; track i) { <span class="sk-item">{{ i }}</span> }
                </div>
              </div>
            }
          </div>

          <hr class="divider">
          <h2 style="font-size:1rem;font-weight:600;margin-bottom:0.5rem">Experience</h2>
          @for (t of timeline; track t.role) {
            <div class="tl-item">
              <span class="tl-year">{{ t.year }}</span>
              <div>
                <p class="tl-role">{{ t.role }}</p>
                <p class="tl-company">{{ t.company }}</p>
                <p class="tl-desc">{{ t.desc }}</p>
              </div>
            </div>
          }

          <hr class="divider">
          <h2 style="font-size:1rem;font-weight:600;margin-bottom:1.25rem">Education</h2>
          @for (e of education; track e.degree) {
            <div class="edu-item">
              <p class="edu-degree">{{ e.degree }}</p>
              <p class="edu-inst">{{ e.institution }}</p>
              <p class="edu-date">{{ e.year }} @if (e.grade) { · {{ e.grade }} }</p>
            </div>
          }

          <hr class="divider">
          <h2 style="font-size:1rem;font-weight:600;margin-bottom:1.25rem">Certifications</h2>
          <div class="cert-list">
            @for (c of certifications; track c.name) {
              <div class="cert-item">
                <span class="cert-dot"></span>
                <span><strong>{{ c.name }}</strong> — {{ c.issuer }}</span>
              </div>
            }
          </div>
        </article>
      </div>
    </main>
  `,
})
export class AboutComponent implements OnInit {
  private readonly seo = inject(SeoService);
  ngOnInit() {
    this.seo.set({
      title: 'About — Mehedi Hasan',
      description: 'Software Engineer specializing in Go, PHP/Laravel, MySQL, Redis, Docker, and GCP. 5+ years building scalable backend systems.'
    });
  }

  skills = [
    { category: 'Languages', items: ['Go', 'PHP 8', 'TypeScript', 'SQL', 'Bash'] },
    { category: 'Frameworks', items: ['Laravel', 'Symfony Components', 'EC-Cube'] },
    { category: 'Databases', items: ['MySQL', 'PostgreSQL', 'MS SQL Server', 'Redis'] },
    { category: 'Real-Time', items: ['WebSockets', 'Firebase Cloud Messaging'] },
    { category: 'DevOps', items: ['Docker', 'GitHub Actions', 'GCP', 'AWS', 'Nginx', 'Apache'] },
    { category: 'Frontend', items: ['JavaScript', 'jQuery', 'Bootstrap', 'Tailwind CSS'] },
    { category: 'Tools', items: ['Git', 'GitHub', 'Postman', 'RabbitMQ', 'Linux'] },
    { category: 'Principles', items: ['Clean Architecture', 'SOLID', 'OOP', 'Design Patterns'] },
  ];

  timeline = [
    {
      year: 'Apr 2023 – Now',
      role: 'Software Engineer',
      company: 'Daisen Technologies Ltd · Dhaka, Bangladesh',
      desc: 'Built a real-time Japanese auction platform. Developing scalable Go and Laravel APIs with clean architecture, CI/CD via GitHub Actions, GCP deployments, Redis caching, and queue workers. Key systems: Car Parts eCommerce (Go), Inventory & Accounting SaaS, EC-Cube custom modules.',
    },
    {
      year: 'Apr 2021 – Mar 2023',
      role: 'Software Engineer',
      company: 'MY Outsourcing Ltd · Dhaka, Bangladesh',
      desc: 'Built CRM systems for call centers and custom eCommerce platforms using Laravel and MySQL. Worked with WordPress, CS-Cart, and EC-Cube customizations. Used Docker for development and production. Delivered multiple standalone systems for international clients.',
    },
    {
      year: '2019 – 2021',
      role: 'Freelance Developer',
      company: 'Self-Employed',
      desc: 'Delivered dozens of Laravel/PHP applications for startups, businesses, and academic clients. Built custom CMS solutions, mini eCommerce, and student project systems. Experienced with cPanel, Webmin, DNS configuration, and VPS setups.',
    },
  ];

  education = [
    { degree: 'B.Sc. in Software Engineering', institution: 'Daffodil International University', year: '2019', grade: 'CGPA 3.10' },
    { degree: 'HSC (Science)', institution: 'Ghatail Cantonment Public School & College', year: '2013', grade: 'GPA 4.80' },
    { degree: 'SSC (Science)', institution: 'Ghatail Gano Pilot High School', year: '2011', grade: 'GPA 5.00' },
  ];

  certifications = [
    { name: '.NET Fundamentals', issuer: 'NIIT Dhaka' },
    { name: 'Professional English', issuer: 'WESL, New Zealand' },
    { name: 'Professional Soft Skills Services', issuer: 'NIIT Dhaka' },
  ];
}
