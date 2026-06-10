export interface ApiResponse<T> {
  data: T;
}

export interface ApiError {
  message: string;
  status?: number;
  details?: Record<string, string[]>;
}

export interface PaginationParams {
  page?: number;
  limit?: number;
}
