
MATCH(r:Role {code: 'control-systems-edit'})
DETACH DELETE r;

MATCH(r:Role {code: 'control-systems-view'})
DETACH DELETE r;

