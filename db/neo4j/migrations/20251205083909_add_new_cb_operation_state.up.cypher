

// Create OperationalState nodes

MERGE (os:OperationalState{ uid:'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e' })
SET os.name = 'OS1: Standard Operation', os.code = 'OS1';

MERGE (os:OperationalState{ uid:'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f' })
SET os.name = 'OS2: Experimental technology standby', os.code = 'OS2';

MERGE (os:OperationalState{ uid:'d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a' })
SET os.name = 'OS3: Experimental technology safe state', os.code = 'OS3';

MERGE (os:OperationalState{ uid:'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b' })
SET os.name = 'OS4: All technology Shutdown', os.code = 'OS4';

MERGE (os:OperationalState{ uid:'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c' })
SET os.name = 'OS5: Power shutdown', os.code = 'OS5';
