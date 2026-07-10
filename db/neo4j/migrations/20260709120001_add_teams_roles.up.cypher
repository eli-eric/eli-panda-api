MERGE (r:Role {name: 'Teams View', code: 'teams-view'})
  ON CREATE SET r.uid = apoc.create.uuid();
MERGE (r:Role {name: 'Teams Edit', code: 'teams-edit'})
  ON CREATE SET r.uid = apoc.create.uuid();