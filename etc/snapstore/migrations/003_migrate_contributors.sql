update publications t1
set data['author'] = 
    (select json_agg(
        case
        when el ?| array['person_id', 'external_person']
        then el::jsonb
        when  el->'id' is not null
        then jsonb_build_object('person_id', el->'id') ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        else jsonb_build_object('external_person',
        		jsonb_build_object('first_name', el->'first_name', 'last_name', el->'last_name', 'full_name', concat(replace(el->>'first_name', '"', ''), ' ', replace(el->>'last_name', '"', '')))) ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        end 
    )
    from publications t2, jsonb_array_elements(t2.data->'author') as el 
     where t1.snapshot_id = t2.snapshot_id)
where data ? 'author';

update publications t1
set data['editor'] = 
    (select json_agg(
        case
        when el ?| array['person_id', 'external_person']
        then el::jsonb
        when  el->'id' is not null
        then jsonb_build_object('person_id', el->'id') ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        else jsonb_build_object('external_person',
        		jsonb_build_object('first_name', el->'first_name', 'last_name', el->'last_name', 'full_name', concat(replace(el->>'first_name', '"', ''), ' ', replace(el->>'last_name', '"', '')))) ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        end 
    )
    from publications t2, jsonb_array_elements(t2.data->'editor') as el 
        where t1.snapshot_id = t2.snapshot_id)
where data ? 'editor';

update publications t1
set data['supervisor'] = 
    (select json_agg(
        case
        when el ?| array['person_id', 'external_person']
        then el::jsonb
        when  el->'id' is not null
        then jsonb_build_object('person_id', el->'id') ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        else jsonb_build_object('external_person',
        		jsonb_build_object('first_name', el->'first_name', 'last_name', el->'last_name', 'full_name', concat(replace(el->>'first_name', '"', ''), ' ', replace(el->>'last_name', '"', '')))) ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        end 
    )
    from publications t2, jsonb_array_elements(t2.data->'supervisor') as el 
     where t1.snapshot_id = t2.snapshot_id)
where data ? 'supervisor';

update datasets t1
set data['author'] = 
    (select json_agg(
        case
        when el ?| array['person_id', 'external_person']
        then el::jsonb
        when  el->'id' is not null
        then jsonb_build_object('person_id', el->'id') ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        else jsonb_build_object('external_person',
        		jsonb_build_object('first_name', el->'first_name', 'last_name', el->'last_name', 'full_name', concat(replace(el->>'first_name', '"', ''), ' ', replace(el->>'last_name', '"', '')))) ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        end 
    )
    from datasets t2, jsonb_array_elements(t2.data->'author') as el 
    where t1.snapshot_id = t2.snapshot_id)
where data ? 'author';

update datasets t1
set data['contributor'] = 
    (select json_agg(
        case
        when el ?| array['person_id', 'external_person']
        then el::jsonb
        when  el->'id' is not null
        then jsonb_build_object('person_id', el->'id') ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        else jsonb_build_object('external_person',
        		jsonb_build_object('first_name', el->'first_name', 'last_name', el->'last_name', 'full_name', concat(replace(el->>'first_name', '"', ''), ' ', replace(el->>'last_name', '"', '')))) ||
        	el::jsonb - 'id' - 'first_name' - 'last_name' - 'full_name' - 'orcid' - 'ugent_id' - 'department'
        end 
    )
    from datasets t2, jsonb_array_elements(t2.data->'contributor') as el 
    where t1.snapshot_id = t2.snapshot_id)
where data ? 'contributor';

---- create above / drop below ----
