import { useMemo, useState } from "react";

/** 通用分页 hook，列表页共用。 */
export function usePagination<T>(rows: T[] = [], initialPageSize = 8) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(initialPageSize);
  const pageRows = useMemo(() => rows.slice((page - 1) * pageSize, page * pageSize), [rows, page, pageSize]);
  const pageCount = Math.max(1, Math.ceil(rows.length / pageSize));
  const safeSetPage = (next: number) => setPage(Math.min(Math.max(1, next), pageCount));
  return { page, setPage: safeSetPage, pageSize, setPageSize, pageRows, total: rows.length, pageCount };
}
