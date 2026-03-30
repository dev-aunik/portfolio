import { Routes } from '@angular/router';

export const routes: Routes = [
    {
        path: '',
        loadComponent: () => import('./features/home/home.component').then(m => m.HomeComponent),
        title: 'Mehedi Hasan — Software Engineer',
    },
    {
        path: 'about',
        loadComponent: () => import('./features/about/about.component').then(m => m.AboutComponent),
        title: 'About — Mehedi Hasan',
    },
    {
        path: 'articles',
        loadComponent: () => import('./features/articles/articles-list.component').then(m => m.ArticlesListComponent),
        title: 'Articles — Mehedi Hasan',
    },
    {
        path: 'articles/:slug',
        loadComponent: () => import('./features/articles/article-detail.component').then(m => m.ArticleDetailComponent),
    },
    {
        path: 'talks',
        loadComponent: () => import('./features/talks/talks.component').then(m => m.TalksComponent),
        title: 'Talks — Mehedi Hasan',
    },
    {
        path: 'projects',
        loadComponent: () => import('./features/projects/projects.component').then(m => m.ProjectsComponent),
        title: 'Projects — Mehedi Hasan',
    },
    {
        path: 'contact',
        loadComponent: () => import('./features/contact/contact.component').then(m => m.ContactComponent),
        title: 'Contact — Mehedi Hasan',
    },
    {
        path: '**',
        loadComponent: () => import('./features/not-found/not-found.component').then(m => m.NotFoundComponent),
        title: '404 — Mehedi Hasan',
    },
];
