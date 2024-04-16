ALTER TABLE organizations ADD COLUMN ceased_on DATE;

---- create above / drop below ----

ALTER TABLE organizations DROP COLUMN ceased_on;
