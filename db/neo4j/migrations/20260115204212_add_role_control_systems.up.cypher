
MERGE(r:Role {code: 'control-systems-view'})
SET 
r.name = 'Control Systems Viewer',
r.uid = '3d928fe9-6faa-42d1-9a99-4a0ddc07d217';

MERGE(r:Role {code: 'control-systems-edit'})
SET 
r.name = 'Control Systems Editor',
r.uid = '6bfa4988-51ce-424a-85fb-bd7d2eded058';

