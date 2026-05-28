import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/components/ui/table";
import type { TopProductsResponse } from "../types";

interface TopProductsTableProps {
  data?: TopProductsResponse;
  isLoading?: boolean;
}

export default function TopProductsTable({ data, isLoading }: TopProductsTableProps) {
  if (isLoading) {
    return <div className="text-muted-foreground">Loading...</div>;
  }

  const products = data?.products || [];

  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Product</TableHead>
            <TableHead className="text-right">Qty</TableHead>
            <TableHead className="text-right">Revenue</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {products.length === 0 ? (
            <TableRow>
              <TableCell colSpan={3} className="text-center text-muted-foreground">
                No data
              </TableCell>
            </TableRow>
          ) : (
            products.map((product: { product_id: string; product_name: string; total_quantity: number; total_revenue: number }) => (
              <TableRow key={product.product_id}>
                <TableCell className="font-medium">{product.product_name}</TableCell>
                <TableCell className="text-right">{product.total_quantity}</TableCell>
                <TableCell className="text-right">
                  ${(product.total_revenue / 100).toFixed(2)}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
