update publications t1
set data = jsonb_build_object('related_projects', 
    (select json_agg(jsonb_build_object('project_id', el->'id') || el::jsonb - 'id' - 'name')
    from publications t2, jsonb_array_elements(t2.data->'project') as el 
    where t1.id = t2.id)) || data - 'project'
where data ? 'project';

update datasets t1
set data = jsonb_build_object('related_projects', 
    (select json_agg(jsonb_build_object('project_id', el->'id') || el::jsonb - 'id' - 'name')
    from datasets t2, jsonb_array_elements(t2.data->'project') as el 
    where t1.id = t2.id)) || data - 'project'
where data ? 'project';

---- create above / drop below ----
