select comment, c.dc, c.case_id as parent_id, c.id,  (extract(epoch from created_at) * 1000)::int8 created_at
from "cases".case_comment c;