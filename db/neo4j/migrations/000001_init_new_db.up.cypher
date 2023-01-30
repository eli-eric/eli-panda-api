// test
MATCH (n)
DETACH DELETE n;

create(jirka:Person{ name:"jirka" })
create(jindra:Person{ name:"jindra" })
create(jirka)-[:IS_FATHER_OF]->(jindra);
