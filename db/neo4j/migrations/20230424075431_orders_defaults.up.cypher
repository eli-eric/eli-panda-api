
MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c51'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'NONE', os.name = 'None', os.sortOrder = 0;

MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c52'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'CANCELLED', os.name = 'Cancelled', os.sortOrder = 1;

MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c53'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'PLANNED', os.name = 'Planned', os.sortOrde = 2;

MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c54'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'REQUESTED', os.name = 'Requested', os.sortOrder = 3;

MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c55'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'ORDERED', os.name = 'Ordered', os.sortOrder = 4;