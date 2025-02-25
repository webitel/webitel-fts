select id,
       dc,
       description,
       coalesce(close_result, '')                                                 close_result,
       coalesce(rating_comment, '')                                               rating_comment,
       coalesce(subject, '') subject,
       contact_info,
       (extract(epoch from created_at) * 1000)::int8 created_at
from cases."case" c;