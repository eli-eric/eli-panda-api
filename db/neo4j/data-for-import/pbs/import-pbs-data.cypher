:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/2024_08_27_PBS_import.csv' AS line
with line
where line.ImportInstr = "Yes"
  CALL {
   // TODO: Add the rest of the query here
  } IN TRANSACTIONS OF 200 ROWS;
  
  