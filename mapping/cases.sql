select id,
       dc,
       description,
       coalesce(close_result, '')                                                 close_result,
       coalesce(rating_comment, '')                                               rating_comment,
       coalesce(subject, '') subject,
       contact_info,
       (extract(epoch from created_at) * 1000)::int8 created_at,
       (select json_agg(distinct subject)
        from cases.case_acl
        where access&4=4
          and object = c.id) _role_ids
from cases."case" c
order by id desc;