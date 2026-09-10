MATCH (r:Role {code: 'teams-view'})
DETACH DELETE r;
MATCH (r:Role {code: 'teams-edit'})
DETACH DELETE r;