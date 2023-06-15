update publications t1
set data = jsonb_build_object('related_organizations', 
    (select json_agg(jsonb_build_object('organization_id', el->'id') || el::jsonb - 'id' - 'tree')
    from publications t2, jsonb_array_elements(t2.data->'department') as el 
    where t1.id = t2.id)) || data - 'department'
where data ? 'department';

update datasets t1
set data = jsonb_build_object('related_organizations', 
    (select json_agg(jsonb_build_object('organization_id', el->'id') || el::jsonb - 'id' - 'tree')
    from datasets t2, jsonb_array_elements(t2.data->'department') as el 
    where t1.id = t2.id)) || data - 'department'
where data ? 'department';

---- create above / drop below ----
