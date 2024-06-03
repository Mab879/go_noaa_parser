CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS alerts
(
    id            serial                        NOT NULL  primary key,
    nws_url       character varying,
    alert_cap     xml,
    updated       timestamp with time zone,
    published     timestamp with time zone,
    author_name   character varying,
    title         character varying,
    summary       text,
    cap_event     character varying,
    cap_effective timestamp with time zone,
    cap_expires   timestamp with time zone,
    cap_status    integer,
    link          character varying,
    cap_msgtype   integer,
    cap_category  integer,
    cap_urgency   integer,
    cap_severity  integer,
    cap_certainty integer,
    cap_areadesc  text,
    cap_polygon   public.GEOMETRY(POLYGON, 4326),
    cap_geocode   json,
    cap_parameter json,
    created_at    timestamp(6) with time zone NOT NULL,
    updated_at    timestamp(6) with time zone NOT NULL
);

create unique index if not exists alerts_nws_url_uindex
    on alerts (nws_url);
