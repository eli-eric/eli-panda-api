MERGE(r:Role {code: 'catalogue-service-view'})
SET 
r.name = 'Service Viewer',
r.uid = '1ee97cba-3734-404f-9eb0-8544a1ebf54c';

MERGE(r:Role {code: 'catalogue-service-edit'})
SET 
r.name = 'Service Editor',
r.uid = '3608a1dd-edc2-4721-99d2-bceb2fc473af';