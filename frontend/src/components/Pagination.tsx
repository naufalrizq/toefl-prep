interface PaginationProps {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  size?: number;
  onSizeChange?: (size: number) => void;
}

export default function Pagination({ page, totalPages, onPageChange, size, onSizeChange }: PaginationProps) {
  const prevDisabled = page <= 1;
  const nextDisabled = page >= totalPages;

  return (
    <div className="mt-3 flex flex-wrap items-center justify-between gap-2 text-sm text-ink-muted">
      {size && onSizeChange ? (
        <label className="flex items-center gap-1.5">
          <span>Per page</span>
          <select
            value={size}
            onChange={(e) => onSizeChange(Number(e.target.value))}
            className="rounded-sm border border-border-strong bg-card px-2 py-1 text-sm"
            aria-label="Items per page"
          >
            {[5, 10, 25, 50].map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </select>
        </label>
      ) : (
        <span />
      )}
      <div className="flex items-center gap-3">
        <button
          type="button"
          disabled={prevDisabled}
          onClick={() => onPageChange(page - 1)}
          className="rounded-sm px-3 py-1.5 transition-colors duration-200 hover:bg-card-muted disabled:opacity-40"
        >
          ← Prev
        </button>
        <span>
          Page {page} of {totalPages}
        </span>
        <button
          type="button"
          disabled={nextDisabled}
          onClick={() => onPageChange(page + 1)}
          className="rounded-sm px-3 py-1.5 transition-colors duration-200 hover:bg-card-muted disabled:opacity-40"
        >
          Next →
        </button>
      </div>
    </div>
  );
}
