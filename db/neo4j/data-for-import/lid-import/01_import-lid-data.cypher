LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/lid-old-data.csv' AS row
WITH  row.code AS code, row.description AS description
WITH code,description,
     apoc.text.regexGroups(code, '^[A-Z]{1,4}') AS typeMatch,
     apoc.text.regexGroups(code, '[0-9]{1,2}(?=[0-9]{2}$)') AS zoneMatch,
     apoc.text.regexGroups(code, '[0-9]{2}$') AS serialMatch
WITH code, description, 
     CASE 
         WHEN size(typeMatch) > 0 AND size(zoneMatch) > 0 AND size(serialMatch) > 0
         THEN {
             type: head(typeMatch[0]), 
             zone: head(zoneMatch[0]), 
             serial: head(serialMatch[0]), 
             result: "R" // Recognized
         }
         ELSE {
             type: null, 
             zone: null, 
             serial: null, 
             result: "NR" // Not Recognized
         }
     END AS data
with code, data, description
MERGE (lid:LidDevice{code: code})
SET
lid.pandaCode = data.type + apoc.text.lpad((data.zone), 2, '0') + "-" + apoc.text.lpad((data.serial), 3, '0'),
lid.type = data.type,
lid.zone = data.zone,
lid.serial = data.serial,
lid.description = description;