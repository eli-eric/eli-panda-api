MERGE (r:Role{ name: 'Zones Edit', code: 'zones-edit' }) ON CREATE SET r.uid = apoc.create.uuid();
