export interface ApiDataResponse<T> {
  code: number;
  data: T;
  message: string;
} 