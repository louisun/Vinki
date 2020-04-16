-- 创建数据库 Vinki 和对应表
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
CREATE DATABASE vinki WITH ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';
ALTER DATABASE vinki OWNER TO postgres;
\connect vinki
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET default_tablespace = '';
SET default_with_oids = false;
CREATE TABLE public."Article" (
    id bigint NOT NULL,
    title character varying(255),
    file_path character varying(255),
    tag character varying(255),
    html_content text
);
ALTER TABLE public."Article" OWNER TO postgres;
CREATE TABLE public."Tag" (
    id bigint NOT NULL,
    name character varying(255),
    html_content text
);
ALTER TABLE public."Tag" OWNER TO postgres;
ALTER TABLE ONLY public."Article"
    ADD CONSTRAINT "Artical_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Tag"
    ADD CONSTRAINT "Tag_pkey" PRIMARY KEY (id);
