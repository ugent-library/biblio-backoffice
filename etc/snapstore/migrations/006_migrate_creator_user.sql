update publications t1
set data = jsonb_build_object('creator_id', data->'creator'->>'id') || data - 'creator'
where data ? 'creator';

update publications t1
set data = jsonb_build_object('user_id', data->'user'->>'id') || data - 'user'
where data ? 'user';

update publications t1
set data = jsonb_build_object('last_user_id', data->'last_user'->>'id') || data - 'last_user'
where data ? 'last_user';

update datasets t1
set data = jsonb_build_object('creator_id', data->'creator'->>'id') || data - 'creator'
where data ? 'creator';

update datasets t1
set data = jsonb_build_object('user_id', data->'user'->>'id') || data - 'user'
where data ? 'user';

update datasets t1
set data = jsonb_build_object('last_user_id', data->'last_user'->>'id') || data - 'last_user'
where data ? 'last_user';

---- create above / drop below ----
