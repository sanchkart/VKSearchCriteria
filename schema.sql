--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

--
-- Name: request_type_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE request_type_enum AS ENUM (
    'INTERSECTION'
);


ALTER TYPE public.request_type_enum OWNER TO postgres;

--
-- Name: status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE status_enum AS ENUM (
    'PROCESSING',
    'DONE',
    'CANCELLED',
    'ERROR'
);


ALTER TYPE public.status_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: requests; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE requests (
    request_uuid uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_uuid uuid NOT NULL,
    request_type request_type_enum NOT NULL,
    created_at timestamp without time zone NOT NULL,
    status status_enum DEFAULT 'PROCESSING'::status_enum NOT NULL,
    params text NOT NULL
);


ALTER TABLE public.requests OWNER TO postgres;

--
-- Name: results; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE results (
    request_uuid uuid NOT NULL,
    id text NOT NULL,
    added_at timestamp without time zone NOT NULL,
    result_id bigint NOT NULL
);


ALTER TABLE public.results OWNER TO postgres;

--
-- Name: results_result_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE results_result_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.results_result_id_seq OWNER TO postgres;

--
-- Name: results_result_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE results_result_id_seq OWNED BY results.result_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users (
    user_uuid uuid DEFAULT uuid_generate_v4() NOT NULL,
    auth text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: result_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY results ALTER COLUMN result_id SET DEFAULT nextval('results_result_id_seq'::regclass);


--
-- Name: requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY requests
    ADD CONSTRAINT requests_pkey PRIMARY KEY (request_uuid);


--
-- Name: results_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY results
    ADD CONSTRAINT results_pkey PRIMARY KEY (result_id);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_uuid);


--
-- Name: requests_fk1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY requests
    ADD CONSTRAINT requests_fk1 FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE RESTRICT;


--
-- Name: results_fk1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY results
    ADD CONSTRAINT results_fk1 FOREIGN KEY (request_uuid) REFERENCES requests(request_uuid) ON DELETE RESTRICT;


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

