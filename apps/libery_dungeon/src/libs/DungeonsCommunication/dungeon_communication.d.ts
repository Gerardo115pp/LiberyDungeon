export interface Stringable {
    toString(): string;
}

export type PaginatedResponse<T> = {
    page: number;
    page_size: number;
    total_pages: number;
    content: T[];
}