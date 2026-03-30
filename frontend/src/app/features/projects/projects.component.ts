import { Component, OnInit, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SeoService } from '../../core/services/seo.service';

interface Project {
    name: string;
    company: string;
    year: string;
    description: string;
    highlights: string[];
    tags: string[];
}

@Component({
    selector: 'app-projects',
    standalone: true,
    imports: [CommonModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    styles: [`
    .page { padding: 4rem 0 5rem; }
    .eyebrow { font-size: 0.75rem; letter-spacing: 0.1em; text-transform: uppercase; color: var(--accent); font-weight: 500; margin-bottom: 0.75rem; }
    h1 { margin-bottom: 0.75rem; }
    .sub { color: var(--text-secondary); margin-bottom: 2.5rem; font-size: 1.0625rem; max-width: 520px; line-height: 1.7; }
    .projects-grid { display: flex; flex-direction: column; gap: 0; }
    .project-card {
      border-bottom: 1px solid var(--border); padding: 2rem 0;
      &:first-of-type { border-top: 1px solid var(--border); }
    }
    .project-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; flex-wrap: wrap; margin-bottom: 0.65rem; }
    .project-name { font-size: 1.0625rem; font-weight: 600; color: var(--text); }
    .project-meta { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.75rem; flex-wrap: wrap; }
    .project-company { font-size: 0.8125rem; color: var(--accent); font-weight: 500; }
    .project-year { font-size: 0.8125rem; color: var(--muted); }
    .project-desc { font-size: 0.9rem; color: var(--text-secondary); line-height: 1.7; margin-bottom: 1rem; }
    .highlights { list-style: none; padding: 0; margin: 0 0 1rem; display: flex; flex-direction: column; gap: 0.35rem; }
    .highlights li {
      font-size: 0.875rem; color: var(--text-secondary); padding-left: 1rem; position: relative;
      &::before { content: '›'; position: absolute; left: 0; color: var(--accent); font-weight: 700; }
    }
    .tags { display: flex; flex-wrap: wrap; gap: 0.4rem; }
    .tag { font-size: 0.75rem; font-weight: 500; color: var(--muted); background: var(--surface); border: 1px solid var(--border); border-radius: 4px; padding: 0.2rem 0.55rem; }
  `],
    template: `
    <main id="main-content" tabindex="-1">
      <div class="container">
        <section class="page fade-in-up">
          <p class="eyebrow">Work</p>
          <h1>Projects</h1>
          <p class="sub">A selection of systems and products I've built professionally — spanning eCommerce, real-time platforms, SaaS, and CRM solutions.</p>
          <div class="projects-grid">
            @for (p of projects; track p.name) {
              <article class="project-card">
                <div class="project-header">
                  <h2 class="project-name">{{ p.name }}</h2>
                </div>
                <div class="project-meta">
                  <span class="project-company">{{ p.company }}</span>
                  <span aria-hidden="true">·</span>
                  <span class="project-year">{{ p.year }}</span>
                </div>
                <p class="project-desc">{{ p.description }}</p>
                <ul class="highlights">
                  @for (h of p.highlights; track h) { <li>{{ h }}</li> }
                </ul>
                <div class="tags">
                  @for (t of p.tags; track t) { <span class="tag">{{ t }}</span> }
                </div>
              </article>
            }
          </div>
        </section>
      </div>
    </main>
  `,
})
export class ProjectsComponent implements OnInit {
    private readonly seo = inject(SeoService);

    ngOnInit() {
        this.seo.set({
            title: 'Projects — Mehedi Hasan',
            description: 'Projects built by Mehedi Hasan — eCommerce platforms, real-time auction systems, SaaS inventory tools, CRM solutions, and more.'
        });
    }

    projects: Project[] = [
        {
            name: 'Real-Time Japanese Auction Platform',
            company: 'Daisen Technologies Ltd',
            year: '2023 – Present',
            description: 'A high-performance, real-time auction system for the Japanese market, enabling live bidding with instant price updates across multiple concurrent users.',
            highlights: [
                'Built WebSocket-based real-time bidding engine with instant price propagation',
                'Integrated Redis for session management, caching, and pub/sub messaging',
                'Implemented RabbitMQ queue workers for reliable async event processing',
                'Designed optimized SQL schemas with strategic indexing for high-throughput queries',
            ],
            tags: ['Go', 'WebSockets', 'Redis', 'RabbitMQ', 'MySQL', 'Docker', 'GCP'],
        },
        {
            name: 'Car Parts eCommerce Platform',
            company: 'Daisen Technologies Ltd',
            year: '2023 – Present',
            description: 'A scalable eCommerce platform for car parts, built with Go following clean architecture principles, supporting high product catalogue volumes and complex filtering.',
            highlights: [
                'Designed RESTful API in Go using clean/layered architecture (CSR pattern)',
                'Implemented advanced product search and filtering with SQL optimisation',
                'Built Redis-based caching layer for catalogue and pricing data',
                'Set up CI/CD pipelines via GitHub Actions with automated GCP deployments',
            ],
            tags: ['Go', 'Clean Architecture', 'MySQL', 'Redis', 'GitHub Actions', 'GCP', 'Docker'],
        },
        {
            name: 'Inventory & Accounting SaaS',
            company: 'Daisen Technologies Ltd',
            year: '2023 – Present',
            description: 'A multi-tenant SaaS application for inventory management and accounting, tailored for small-to-medium Japanese businesses.',
            highlights: [
                'Architected multi-tenant data isolation with per-tenant scoping',
                'Built financial reporting modules with real-time dashboard data',
                'Implemented role-based access control (RBAC) for different user tiers',
                'Optimized complex aggregate SQL queries with proper indexing strategies',
            ],
            tags: ['Go', 'Laravel', 'PostgreSQL', 'Redis', 'Docker', 'GCP'],
        },
        {
            name: 'EC-Cube Custom Module Development',
            company: 'Daisen Technologies Ltd',
            year: '2023 – Present',
            description: 'Custom plugin and module development for EC-Cube (Symfony-based eCommerce CMS), extending core functionality for Japanese client requirements.',
            highlights: [
                'Developed custom payment gateway integrations and shipping modules',
                'Built admin panel extensions using Symfony components',
                'Created custom CSV import/export tools for bulk product management',
            ],
            tags: ['PHP', 'Symfony', 'EC-Cube', 'MySQL', 'Twig', 'Docker'],
        },
        {
            name: 'CRM System for Call Centers',
            company: 'MY Outsourcing Ltd',
            year: '2021 – 2023',
            description: 'A comprehensive CRM platform built for call center operations, managing leads, customer interactions, agent performance tracking, and reporting.',
            highlights: [
                'Built call log management, lead assignment, and pipeline tracking features',
                'Implemented real-time agent monitoring dashboard',
                'Developed automated reporting with exportable PDF/Excel reports',
                'Designed normalized MySQL schema for high-volume transactional data',
            ],
            tags: ['PHP', 'Laravel', 'MySQL', 'jQuery', 'Bootstrap', 'Docker'],
        },
        {
            name: 'Custom eCommerce Systems',
            company: 'MY Outsourcing Ltd',
            year: '2021 – 2023',
            description: 'Multiple standalone eCommerce platforms built from scratch using Laravel for international clients, featuring product management, orders, payments, and logistics.',
            highlights: [
                'Delivered multiple independent eCommerce solutions for international clients',
                'Integrated payment gateways including Stripe and local providers',
                'Built inventory tracking and order management systems',
                'Customized CS-Cart and WordPress WooCommerce installations',
            ],
            tags: ['PHP', 'Laravel', 'MySQL', 'CS-Cart', 'WordPress', 'Docker', 'Nginx'],
        },
    ];
}
