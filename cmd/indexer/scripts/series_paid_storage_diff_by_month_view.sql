create materialized view if not exists series_paid_storage_size_diff_by_month_%s AS
with f as (
        select generate_series(
        date_trunc('month', date '2018-06-25'),
        date_trunc('month', now()),
        '1 month'::interval
        ) as val
)
select
        extract(epoch from f.val),
        sum(paid_storage_size_diff) as value
from f
left join operations on date_trunc('month', operations.timestamp) = f.val where ((network = %d))
group by 1
order by date_part