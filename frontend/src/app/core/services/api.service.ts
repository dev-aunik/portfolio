import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface Article {
    id: string; title: string; slug: string; summary: string;
    content: string; tags: string[]; published_at: string | null;
    created_at: string; updated_at: string;
}
export interface ArticleListResult {
    articles: Article[]; total: number; page: number; total_pages: number;
}

export interface ContactPayload {
    name: string; email: string; subject: string; message: string;
}

@Injectable({ providedIn: 'root' })
export class ApiService {
    private readonly http = inject(HttpClient);
    private readonly base = environment.apiUrl + '/api/v1';

    getArticles(page = 1, limit = 20): Observable<ArticleListResult> {
        const params = new HttpParams().set('page', page).set('limit', limit);
        return this.http.get<ArticleListResult>(`${this.base}/articles`, { params });
    }
    getArticle(slug: string): Observable<Article> {
        return this.http.get<Article>(`${this.base}/articles/${slug}`);
    }
    searchArticles(query: string, page = 1): Observable<ArticleListResult> {
        const params = new HttpParams().set('q', query).set('page', page);
        return this.http.get<ArticleListResult>(`${this.base}/articles/search`, { params });
    }
    submitContact(data: ContactPayload): Observable<{ message: string }> {
        return this.http.post<{ message: string }>(`${this.base}/contact`, data);
    }
}
