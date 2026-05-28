export interface SalesSummary {
  order_count: number;
  total_sales: number;
  avg_ticket: number;
  total_tax: number;
  total_tips: number;
  total_discounts: number;
}

export interface SummaryResponse {
  restaurant_id: string;
  restaurant_name: string;
  from: string;
  to: string;
  sales: SalesSummary;
}

export interface SalesByEmployee {
  employee_id: string;
  employee_name: string;
  order_count: number;
  total_sales: number;
  avg_ticket: number;
}

export interface SalesByEmployeeResponse {
  from: string;
  to: string;
  employees: SalesByEmployee[];
}

export interface TopProduct {
  product_id: string;
  product_name: string;
  total_quantity: number;
  total_revenue: number;
}

export interface TopProductsResponse {
  from: string;
  to: string;
  products: TopProduct[];
}

export interface CategoryCoverage {
  employee_id: string;
  employee_name: string;
  categories_covered: number;
  total_categories: number;
  coverage_percent: number;
}

export interface CategoryCoverageResponse {
  from: string;
  to: string;
  employees: CategoryCoverage[];
}

export interface TicketIdealResult {
  total_orders: number;
  complete_orders: number;
  incomplete_orders: number;
  completion_rate: number;
}

export interface TicketIdealResponse {
  from: string;
  to: string;
  result: TicketIdealResult;
  rules?: Record<string, number>;
}

export interface DateRange {
  from: Date;
  to: Date;
}
