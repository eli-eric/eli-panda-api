MERGE (r:Role{ name: 'Zones View', code: 'zones-view' }) ON CREATE SET r.uid = randomUUID();
MERGE (r:Role{ name: 'Zones Edit', code: 'zones-edit' }) ON CREATE SET r.uid = randomUUID();
