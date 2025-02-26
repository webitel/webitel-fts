select comment,
       c.dc,
       c.case_id as parent_id,
       c.id,
       (extract(epoch from created_at) * 1000)::int8 created_at,
       (select json_agg(distinct subject)
        from cases.case_comment_acl
        where access&4=4
          and object = c.id) _role_ids
from "cases".case_comment c;