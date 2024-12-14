MERGE(r:Role {code: 'publications-view'})
SET 
r.name = 'Publications Viewer',
r.uid = 'a690c631-b648-4f70-a16b-7c104a12cf92';

MERGE(r:Role {code: 'publications-edit'})
SET 
r.name = 'Publications Editor',
r.uid = '65d6278d-5347-4802-a0a0-8e76a37f8193';